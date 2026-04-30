package domain

// Subscribtion model info
// @Description The subscribtion
type Subscribtion struct {
	ServiceName string `json:"service_name"`
	UserId      string `json:"user_id"`
	Price       int    `json:"price"`
	Start       *Date  `json:"start_date"`
	End         *Date  `json:"end_date"`
}
