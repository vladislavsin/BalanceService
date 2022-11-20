package transaction

import "time"

type History struct {
	ID                uint
	BalanceID         uint
	TransactionTypeID uint
	ServiceID         uint
	Amount            uint
	CreatedAt         time.Time
}

type SortingHistory struct {
	Pagination uint
	Sorting    string
	OrderBy    string
}

const (
	AddingAmount = 1
	Reservation  = 2
	PaidService  = 3
)

func New(userBalanceID uint, transactionTypeID uint, amount uint, ServiceID uint) *History {
	return &History{
		BalanceID:         userBalanceID,
		TransactionTypeID: transactionTypeID,
		ServiceID:         ServiceID,
		Amount:            amount,
		CreatedAt:         time.Now(),
	}
}

func NewSortingHistory(pagination uint, sorting string, orderBy string) *SortingHistory {
	return &SortingHistory{
		Pagination: pagination,
		Sorting:    sorting,
		OrderBy:    orderBy,
	}
}
