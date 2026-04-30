package domain

import "time"

type Subscribe struct {
	ServiceName string
	ClientId    string
	Price       string
	Start       time.Time
	End         time.Time
}
