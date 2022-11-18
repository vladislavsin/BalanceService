package reservation

type Reservation struct {
	ID        uint `json:"id"`
	UserID    uint `json:"user_id"`
	ServiceID uint `json:"service_id"`
	OrderID   uint `json:"order_id"`
	Amount    uint `json:"amount"`
	CreatedAt string
	UpdatedAt string
}

type ReservationDTO struct {
	UserID    uint `json:"user_id"`
	ServiceID uint `json:"service_id"`
	OrderID   uint `json:"order_id"`
	Amount    uint `json:"amount"`
	CreatedAt string
	UpdatedAt string
}

const (
	InProgress = 1
	Cancel     = 2
	Accept     = 3
)
