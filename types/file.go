package types

type NTPConfiguration struct {
	Server string `json:"server"`
	IBurst bool   `json:"iburst"`
}
