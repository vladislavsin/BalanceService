package balance

import (
	"BalanceService/internal/reservation"
	"BalanceService/internal/transaction"
	"context"
)

type Storage interface {
	AddAmount(ctx context.Context, userBalance Balance) (Balance, error)
	WithdrawFunds(ctx context.Context, userBalance Balance) error
	Reservation(ctx context.Context, reservation reservation.ReservationDTO) error
	CreateUserBalance(ctx context.Context, balance BalanceDTO) (Balance, error)
	GetUserBalance(ctx context.Context, id uint) (Balance, error)
	GetReservation(ctx context.Context, orderID uint) (reservation.Reservation, error)
	AcceptReservation(ctx context.Context, reservation reservation.Reservation) error
	GetTransactionHistory(ctx context.Context, balance Balance, sort *transaction.SortingHistory) ([]transaction.History, error)
	AddTransactionHistory(ctx context.Context, transaction *transaction.History) error
}
