package balance

type BalanceDTO struct {
	UserID uint `json:"user_id" query:"userId"`
	Amount uint `json:"amount"`
}
