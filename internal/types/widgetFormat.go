package types

type PercentageSwitch struct {
	MinValue     int `json:"min_value"`
	MaxValue     int `json:"max_value"`
	CurrentValue int `json:"current_value"`
}
