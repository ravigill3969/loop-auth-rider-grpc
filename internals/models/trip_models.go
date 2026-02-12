package models

type CancelRideS struct {
	TripID   string `json:"trip_id"`
	DriverID string `json:"driver_id"`
	Reason   string `json:"reason"`
}
