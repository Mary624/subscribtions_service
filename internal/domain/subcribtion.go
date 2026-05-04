package domain

// SubscribtionRequest model info
// @Description The subscribtion
type SubscribtionRequest struct {
	ServiceName string `json:"service_name"`
	UserId      string `json:"user_id"`
	Price       int    `json:"price"`
	Start       string `json:"start_date"`
	End         string `json:"end_date"`
}

// Subscribtion model info
// @Description The subscribtion
type Subscribtion struct {
	ServiceName string `json:"service_name"`
	UserId      string `json:"user_id"`
	Price       int    `json:"price"`
	Start       *Date  `json:"start_date"`
	End         *Date  `json:"end_date"`
}
