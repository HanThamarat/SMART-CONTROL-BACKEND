package socket

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	fibersocket "github.com/gofiber/contrib/socketio"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"

	"github.com/HanThamarat/SMART-CONTROL-BACKEND/pkg/response"
)

type Server struct {
	mu            sync.RWMutex
	clients       map[string]map[string]struct{}
	mqttPublisher MQTTPublisher
}

type IncomingMessage struct {
	Event string          `json:"event"`
	To    string          `json:"to,omitempty"`
	Data  json.RawMessage `json:"data,omitempty"`
}

type MQTTPublisher interface {
	Publish(topic string, qos byte, retained bool, payload interface{}) error
}

type MQTTPublishRequest struct {
	Topic    string          `json:"topic"`
	Payload  json.RawMessage `json:"payload"`
	QoS      byte            `json:"qos,omitempty"`
	Retained bool            `json:"retained,omitempty"`
}

type OutgoingMessage struct {
	Event     string          `json:"event"`
	Message   string          `json:"message,omitempty"`
	ClientID  string          `json:"clientId,omitempty"`
	SocketID  string          `json:"socketId,omitempty"`
	From      string          `json:"from,omitempty"`
	To        string          `json:"to,omitempty"`
	Timestamp string          `json:"timestamp"`
	Data      json.RawMessage `json:"data,omitempty"`
}

func NewServer() *Server {
	server := &Server{
		clients: make(map[string]map[string]struct{}),
	}

	server.registerEvents()

	return server
}

func (s *Server) SetMQTTPublisher(publisher MQTTPublisher) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.mqttPublisher = publisher
}

func (s *Server) BroadcastMQTTMessage(topic string, payload []byte) {
	if strings.TrimSpace(topic) == "" {
		return
	}

	data, err := json.Marshal(map[string]interface{}{
		"topic":   topic,
		"payload": mqttSocketPayload(payload),
	})
	if err != nil {
		log.Printf("socket mqtt broadcast marshal error topic=%s err=%v", topic, err)
		return
	}

	message, err := json.Marshal(OutgoingMessage{
		Event:     "mqtt:message",
		Message:   "MQTT message received.",
		Timestamp: now(),
		Data:      data,
	})
	if err != nil {
		log.Printf("socket mqtt envelope marshal error topic=%s err=%v", topic, err)
		return
	}

	fibersocket.Broadcast(message, fibersocket.TextMessage)
}

func (s *Server) Upgrade(c *fiber.Ctx) error {
	if !websocket.IsWebSocketUpgrade(c) {
		return fiber.ErrUpgradeRequired
	}

	token, err := authenticateSocketRequest(c)
	if err != nil {
		return response.SetErrResponse(
			c,
			fiber.StatusUnauthorized,
			"Socket authentication failed.",
			err.Error(),
		)
	}

	c.Locals("user", token)
	c.Locals("allowed", true)
	return c.Next()
}

func (s *Server) Handler() fiber.Handler {
	return fibersocket.New(func(kws *fibersocket.Websocket) {
		clientID := normalizeClientID(kws.Query("clientId"))
		if clientID == "" {
			clientID = clientIDFromToken(kws.Locals("user"))
		}
		if clientID == "" {
			clientID = kws.UUID
		}

		kws.SetAttribute("client_id", clientID)
		setAuthAttributes(kws, kws.Locals("user"))
		s.addClient(clientID, kws.UUID)
		log.Printf("socket registered client_id=%s socket_id=%s active_sockets=%d", clientID, kws.UUID, s.clientSocketCount(clientID))

		s.emit(kws, OutgoingMessage{
			Event:     "socket:connected",
			Message:   "Socket connection established.",
			ClientID:  clientID,
			SocketID:  kws.UUID,
			Timestamp: now(),
		})
	})
}

func (s *Server) registerEvents() {
	fibersocket.On(fibersocket.EventConnect, func(ep *fibersocket.EventPayload) {
		clientID := ep.Kws.GetStringAttribute("client_id")
		log.Printf("socket connected client_id=%s socket_id=%s", clientID, ep.Kws.UUID)

		s.broadcast(ep.Kws, OutgoingMessage{
			Event:     "socket:user_joined",
			Message:   "A client joined the realtime channel.",
			ClientID:  clientID,
			SocketID:  ep.Kws.UUID,
			Timestamp: now(),
		})
	})

	fibersocket.On(fibersocket.EventMessage, func(ep *fibersocket.EventPayload) {
		s.handleMessage(ep)
	})

	fibersocket.On(fibersocket.EventDisconnect, func(ep *fibersocket.EventPayload) {
		s.handleDisconnect(ep)
	})

	fibersocket.On(fibersocket.EventClose, func(ep *fibersocket.EventPayload) {
		s.handleDisconnect(ep)
	})

	fibersocket.On(fibersocket.EventError, func(ep *fibersocket.EventPayload) {
		log.Printf("socket error client_id=%s socket_id=%s err=%v", ep.Kws.GetStringAttribute("client_id"), ep.Kws.UUID, ep.Error)
	})
}

func (s *Server) handleMessage(ep *fibersocket.EventPayload) {
	clientID := ep.Kws.GetStringAttribute("client_id")

	var incoming IncomingMessage
	if err := json.Unmarshal(ep.Data, &incoming); err != nil {
		s.emit(ep.Kws, OutgoingMessage{
			Event:     "socket:error",
			Message:   "Invalid payload. Expected JSON with event and optional data/to.",
			ClientID:  clientID,
			SocketID:  ep.Kws.UUID,
			Timestamp: now(),
		})
		return
	}

	incoming.To = normalizeClientID(incoming.To)

	if incoming.Event == "" {
		incoming.Event = "message"
	}

	if incoming.Event == "mqtt:publish" {
		s.handleMQTTPublish(ep, clientID, incoming.Data)
		return
	}

	outgoing := OutgoingMessage{
		Event:     incoming.Event,
		From:      clientID,
		To:        incoming.To,
		ClientID:  clientID,
		SocketID:  ep.Kws.UUID,
		Timestamp: now(),
		Data:      incoming.Data,
	}

	payload, err := json.Marshal(outgoing)
	if err != nil {
		log.Printf("socket marshal error client_id=%s socket_id=%s err=%v", clientID, ep.Kws.UUID, err)
		return
	}

	if incoming.To != "" {
		targetSockets := s.getClientSockets(incoming.To)
		if len(targetSockets) == 0 {
			log.Printf(
				"socket target unavailable from=%s socket_id=%s to=%s active_clients=%v",
				clientID,
				ep.Kws.UUID,
				incoming.To,
				s.clientIDs(),
			)
			s.emit(ep.Kws, OutgoingMessage{
				Event:     "socket:error",
				Message:   "Target client is not connected.",
				ClientID:  clientID,
				SocketID:  ep.Kws.UUID,
				To:        incoming.To,
				Timestamp: now(),
			})
			return
		}

		ep.Kws.EmitToList(targetSockets, payload, fibersocket.TextMessage)
		s.emit(ep.Kws, OutgoingMessage{
			Event:     "socket:delivered",
			Message:   "Message delivered to target client.",
			ClientID:  clientID,
			SocketID:  ep.Kws.UUID,
			To:        incoming.To,
			Timestamp: now(),
		})
		return
	}

	s.broadcastRaw(ep.Kws, payload)
	s.emit(ep.Kws, OutgoingMessage{
		Event:     "socket:delivered",
		Message:   "Message broadcast to connected clients.",
		ClientID:  clientID,
		SocketID:  ep.Kws.UUID,
		Timestamp: now(),
	})
}

func (s *Server) handleMQTTPublish(ep *fibersocket.EventPayload, clientID string, data json.RawMessage) {
	publisher := s.getMQTTPublisher()
	if publisher == nil {
		s.emit(ep.Kws, OutgoingMessage{
			Event:     "socket:error",
			Message:   "MQTT publisher is not configured.",
			ClientID:  clientID,
			SocketID:  ep.Kws.UUID,
			Timestamp: now(),
		})
		return
	}

	var req MQTTPublishRequest
	if err := json.Unmarshal(data, &req); err != nil {
		s.emit(ep.Kws, OutgoingMessage{
			Event:     "socket:error",
			Message:   "Invalid mqtt:publish payload. Expected topic and payload.",
			ClientID:  clientID,
			SocketID:  ep.Kws.UUID,
			Timestamp: now(),
		})
		return
	}

	req.Topic = strings.TrimSpace(req.Topic)
	if req.Topic == "" || len(req.Payload) == 0 {
		s.emit(ep.Kws, OutgoingMessage{
			Event:     "socket:error",
			Message:   "MQTT topic and payload are required.",
			ClientID:  clientID,
			SocketID:  ep.Kws.UUID,
			Timestamp: now(),
		})
		return
	}

	payload, err := mqttPayloadBytes(req.Payload)
	if err != nil {
		s.emit(ep.Kws, OutgoingMessage{
			Event:     "socket:error",
			Message:   "Unable to serialize MQTT payload.",
			ClientID:  clientID,
			SocketID:  ep.Kws.UUID,
			Timestamp: now(),
		})
		return
	}

	if err := publisher.Publish(req.Topic, req.QoS, req.Retained, payload); err != nil {
		log.Printf("socket mqtt publish error client_id=%s topic=%s err=%v", clientID, req.Topic, err)
		s.emit(ep.Kws, OutgoingMessage{
			Event:     "socket:error",
			Message:   "MQTT publish failed.",
			ClientID:  clientID,
			SocketID:  ep.Kws.UUID,
			Timestamp: now(),
			To:        req.Topic,
		})
		return
	}

	ackData, err := json.Marshal(map[string]interface{}{
		"topic":    req.Topic,
		"qos":      req.QoS,
		"retained": req.Retained,
	})
	if err != nil {
		log.Printf("socket mqtt ack marshal error client_id=%s topic=%s err=%v", clientID, req.Topic, err)
		ackData = nil
	}

	s.emit(ep.Kws, OutgoingMessage{
		Event:     "mqtt:published",
		Message:   "MQTT message published.",
		ClientID:  clientID,
		SocketID:  ep.Kws.UUID,
		Timestamp: now(),
		Data:      ackData,
	})
}

func (s *Server) handleDisconnect(ep *fibersocket.EventPayload) {
	clientID := ep.Kws.GetStringAttribute("client_id")
	if clientID == "" {
		return
	}

	if !s.removeClient(clientID, ep.Kws.UUID) {
		return
	}

	log.Printf("socket disconnected client_id=%s socket_id=%s", clientID, ep.Kws.UUID)

	s.broadcast(ep.Kws, OutgoingMessage{
		Event:     "socket:user_left",
		Message:   "A client left the realtime channel.",
		ClientID:  clientID,
		SocketID:  ep.Kws.UUID,
		Timestamp: now(),
	})
}

func (s *Server) emit(kws *fibersocket.Websocket, message OutgoingMessage) {
	payload, err := json.Marshal(message)
	if err != nil {
		log.Printf("socket emit marshal error client_id=%s socket_id=%s err=%v", kws.GetStringAttribute("client_id"), kws.UUID, err)
		return
	}

	kws.Emit(payload, fibersocket.TextMessage)
}

func (s *Server) broadcast(kws *fibersocket.Websocket, message OutgoingMessage) {
	payload, err := json.Marshal(message)
	if err != nil {
		log.Printf("socket broadcast marshal error client_id=%s socket_id=%s err=%v", kws.GetStringAttribute("client_id"), kws.UUID, err)
		return
	}

	s.broadcastRaw(kws, payload)
}

func (s *Server) broadcastRaw(kws *fibersocket.Websocket, payload []byte) {
	kws.Broadcast(payload, true, fibersocket.TextMessage)
}

func (s *Server) addClient(clientID string, socketID string) {
	clientID = normalizeClientID(clientID)
	if clientID == "" || socketID == "" {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.clients[clientID]; !ok {
		s.clients[clientID] = make(map[string]struct{})
	}

	s.clients[clientID][socketID] = struct{}{}
}

func (s *Server) removeClient(clientID string, socketID string) bool {
	clientID = normalizeClientID(clientID)
	if clientID == "" || socketID == "" {
		return false
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	sockets, ok := s.clients[clientID]
	if !ok {
		return false
	}

	if _, ok := sockets[socketID]; !ok {
		return false
	}

	delete(sockets, socketID)
	if len(sockets) == 0 {
		delete(s.clients, clientID)
	}

	return true
}

func (s *Server) getClientSockets(clientID string) []string {
	clientID = normalizeClientID(clientID)
	if clientID == "" {
		return nil
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	sockets := s.clients[clientID]
	result := make([]string, 0, len(sockets))
	for socketID := range sockets {
		result = append(result, socketID)
	}

	return result
}

func (s *Server) getMQTTPublisher() MQTTPublisher {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.mqttPublisher
}

func now() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func authenticateSocketRequest(c *fiber.Ctx) (*jwt.Token, error) {
	tokenString, err := tokenFromRequest(c)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

func tokenFromRequest(c *fiber.Ctx) (string, error) {
	authHeader := strings.TrimSpace(c.Get("Authorization"))
	if authHeader != "" {
		const bearerPrefix = "Bearer "
		if len(authHeader) <= len(bearerPrefix) || !strings.EqualFold(authHeader[:len(bearerPrefix)], bearerPrefix) {
			return "", errors.New("authorization header format must be Bearer {token}")
		}

		return strings.TrimSpace(authHeader[len(bearerPrefix):]), nil
	}

	token := strings.TrimSpace(c.Query("token"))
	if token != "" {
		return token, nil
	}

	return "", errors.New("missing JWT in Authorization header or token query parameter")
}

func normalizeClientID(clientID string) string {
	return strings.TrimSpace(clientID)
}

func (s *Server) clientSocketCount(clientID string) int {
	clientID = normalizeClientID(clientID)
	if clientID == "" {
		return 0
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.clients[clientID])
}

func (s *Server) clientIDs() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ids := make([]string, 0, len(s.clients))
	for clientID := range s.clients {
		ids = append(ids, clientID)
	}

	sort.Strings(ids)
	return ids
}

func clientIDFromToken(user interface{}) string {
	claims := tokenClaims(user)

	for _, key := range []string{"clientId", "clientID", "email", "name", "userId"} {
		if value := claimString(claims, key); value != "" {
			return normalizeClientID(value)
		}
	}

	return ""
}

func setAuthAttributes(kws *fibersocket.Websocket, user interface{}) {
	claims := tokenClaims(user)
	if len(claims) == 0 {
		return
	}

	if userID := claimString(claims, "userId"); userID != "" {
		kws.SetAttribute("user_id", userID)
	}
	if email := claimString(claims, "email"); email != "" {
		kws.SetAttribute("user_email", email)
	}
	if name := claimString(claims, "name"); name != "" {
		kws.SetAttribute("user_name", name)
	}
}

func tokenClaims(user interface{}) jwt.MapClaims {
	token, ok := user.(*jwt.Token)
	if !ok || token == nil {
		return nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil
	}

	return claims
}

func claimString(claims jwt.MapClaims, key string) string {
	if len(claims) == 0 {
		return ""
	}

	value, ok := claims[key]
	if !ok || value == nil {
		return ""
	}

	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v)
	case float64:
		return strconv.FormatInt(int64(v), 10)
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case json.Number:
		return v.String()
	default:
		return strings.TrimSpace(fmt.Sprint(v))
	}
}

func mqttPayloadBytes(raw json.RawMessage) ([]byte, error) {
	var textPayload string
	if err := json.Unmarshal(raw, &textPayload); err == nil {
		return []byte(textPayload), nil
	}

	if !json.Valid(raw) {
		return nil, errors.New("invalid JSON payload")
	}

	return []byte(raw), nil
}

func mqttSocketPayload(payload []byte) interface{} {
	if json.Valid(payload) {
		return json.RawMessage(payload)
	}

	return string(payload)
}
