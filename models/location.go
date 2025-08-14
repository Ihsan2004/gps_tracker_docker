package models

type Location struct {
	DeviceID  int     `json:"device_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
