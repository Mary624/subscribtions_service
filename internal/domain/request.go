package domain

import "time"

type RepositoryRequest struct {
	ServiceName string    `json:"service_name,omitempty"`
	ClientId    string    `json:"user_id,omitempty"`
	Start       time.Time `json:"start_date"`
	End         time.Time `json:"end_date"` // TODO: can be nil
}
