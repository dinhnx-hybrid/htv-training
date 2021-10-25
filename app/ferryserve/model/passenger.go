package model

import "time"

type Passenger struct {
	Id         uint
	Name       string
	ComingTime time.Duration
}
