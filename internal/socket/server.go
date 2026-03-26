package socket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	fibersocket "github.com/gofiber/contrib/socketio"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	mu      sync.RWMutex
	clients map[string]map[string]struct{}
}

type IncomingMessage struct {
	Event string          `json:"event"`
	To    string          `json:"to,omitempty"`
	Data  json.RawMessage `json:"data,omitempty"`
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

func (s *Server) Upgrade(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}

	return fiber.ErrUpgradeRequired
}

func (s *Server) Handler() fiber.Handler {
	return fibersocket.New(func(kws *fibersocket.Websocket) {
		clientID := kws.Query("clientId")
		if clientID == "" {
			clientID = kws.UUID
		}

		kws.SetAttribute("client_id", clientID)
		s.addClient(clientID, kws.UUID)

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

	if incoming.Event == "" {
		incoming.Event = "message"
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
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.clients[clientID]; !ok {
		s.clients[clientID] = make(map[string]struct{})
	}

	s.clients[clientID][socketID] = struct{}{}
}

func (s *Server) removeClient(clientID string, socketID string) bool {
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
	s.mu.RLock()
	defer s.mu.RUnlock()

	sockets := s.clients[clientID]
	result := make([]string, 0, len(sockets))
	for socketID := range sockets {
		result = append(result, socketID)
	}

	return result
}

func now() string {
	return time.Now().UTC().Format(time.RFC3339)
}
