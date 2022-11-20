package db

import (
	"BalanceService/internal/balance"
	"BalanceService/internal/reservation"
	"BalanceService/internal/transaction"
	"BalanceService/pkg/client/postgresql"
	"BalanceService/pkg/logging"
	"context"
	"fmt"
	"github.com/jackc/pgconn"
)

type db struct {
	client postgresql.Client
	logger *logging.Logger
}

func (d db) GetTransactionHistory(ctx context.Context, balance balance.Balance, sort *transaction.SortingHistory) ([]transaction.History, error) {
	q := `SELECT id, balance_id, transaction_type_id, service_id, amount, created_at FROM transaction_history WHERE balance_id = $1`

	transactions := make([]transaction.History, 0)

	q = q + fmt.Sprintf(` ORDER BY %s %s`, sort.Sorting, sort.OrderBy)

	q = q + fmt.Sprintf(` LIMIT 10 offset %d`, sort.Pagination)

	d.logger.Info(q)

	rows, err := d.client.Query(ctx, q, balance.ID)
	if err != nil {
		d.logger.Fatal(err)
	}
	for rows.Next() {
		var transact transaction.History
		err = rows.Scan(&transact.ID, &transact.BalanceID, &transact.TransactionTypeID, &transact.ServiceID, &transact.Amount, &transact.CreatedAt)
		if err != nil {
			d.logger.Fatal(err)
		}
		transactions = append(transactions, transact)
	}

	return transactions, nil
}

func (d db) AcceptReservation(ctx context.Context, reserv reservation.Reservation) error {
	q := `UPDATE reservation_status
		  SET status_id = $1
		  WHERE reservation_id = $2`

	if err := d.client.QueryRow(ctx, q, reservation.Accept, reserv.ID).Scan(); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code))
			d.logger.Error(newErr)
			return newErr
		}
	}
	return nil
}

func (d db) GetReservation(ctx context.Context, orderID uint) (reservation.Reservation, error) {
	q := `SELECT * FROM reservation WHERE order_id = $1`

	var reserv reservation.Reservation

	if err := d.client.QueryRow(ctx, q, orderID).Scan(&reserv.ID, &reserv.UserID, &reserv.ServiceID, &reserv.OrderID, &reserv.Amount, &reserv.CreatedAt, &reserv.UpdatedAt); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code))
			d.logger.Error(newErr)
			return reserv, newErr
		}
	}
	return reserv, nil
}

func (d db) WithdrawFunds(ctx context.Context, userBalance balance.Balance) error {
	q := `UPDATE balances
		  SET amount = $1
		  WHERE user_id = $2`

	if err := d.client.QueryRow(ctx, q, userBalance.Amount, userBalance.UserID).Scan(); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code))
			d.logger.Error(newErr)
			return newErr
		}
	}
	return nil
}

func (d db) Reservation(ctx context.Context, reservationDTO reservation.ReservationDTO) error {
	d.logger.Infof("UserID: %d, ServiceID: %d, OrderID: %d, Amount: %d", reservationDTO.UserID, reservationDTO.ServiceID, reservationDTO.OrderID, reservationDTO.Amount)
	q := `INSERT INTO reservation (user_id, service_id, order_id, amount, created_at, updated_at)
		  VALUES ($1, $2, $3, $4, $5, $6)
		  RETURNING id`

	var id uint

	if err := d.client.QueryRow(ctx, q, reservationDTO.UserID, reservationDTO.ServiceID, reservationDTO.OrderID, reservationDTO.Amount, reservationDTO.CreatedAt, reservationDTO.UpdatedAt).Scan(&id); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code))
			d.logger.Error(newErr)
			return newErr
		}
	}

	query := `INSERT INTO reservation_status (reservation_id, status_id)
		  VALUES ($1, $2)`

	if err := d.client.QueryRow(ctx, query, id, reservation.InProgress).Scan(); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code))
			d.logger.Error(newErr)
			return newErr
		}
	}
	return nil
}

func (d db) AddTransactionHistory(ctx context.Context, transactionHistory *transaction.History) error {
	if transactionHistory.TransactionTypeID == transaction.AddingAmount {
		q := `INSERT INTO transaction_history (balance_id, transaction_type_id, amount, created_at)
		  VALUES ($1, $2, $3, $4)`

		if err := d.client.QueryRow(ctx, q, transactionHistory.BalanceID, transactionHistory.TransactionTypeID, transactionHistory.Amount, transactionHistory.CreatedAt).Scan(); err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok {
				newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code))
				d.logger.Error(newErr)
				return newErr
			}
		}
		return nil
	}

	q := `INSERT INTO transaction_history (balance_id, transaction_type_id, service_id, amount, created_at)
		  VALUES ($1, $2, $3, $4, $5)`

	if err := d.client.QueryRow(ctx, q, transactionHistory.BalanceID, transactionHistory.TransactionTypeID, transactionHistory.ServiceID, transactionHistory.Amount, transactionHistory.CreatedAt).Scan(); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code))
			d.logger.Error(newErr)
			return newErr
		}
	}

	return nil
}

func (d db) AddAmount(ctx context.Context, userBalance balance.Balance) (balance.Balance, error) {
	d.logger.Info(userBalance.Amount)
	q := `UPDATE balances
		  SET amount = $1
		  WHERE id = $2
		  RETURNING amount`

	if err := d.client.QueryRow(ctx, q, userBalance.Amount, userBalance.ID).Scan(&userBalance.Amount); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code))
			d.logger.Error(newErr)
			return userBalance, newErr
		}
	}

	return userBalance, nil
}

func (d db) CreateUserBalance(ctx context.Context, userBalance balance.BalanceDTO) (balance.Balance, error) {
	q := `INSERT INTO balances(user_id, amount)
		  VALUES ($1, $2)
		  RETURNING id, user_id, amount`

	var NewUserBalance balance.Balance

	if err := d.client.QueryRow(ctx, q, userBalance.UserID, userBalance.Amount).Scan(&NewUserBalance.ID, &NewUserBalance.UserID, &NewUserBalance.Amount); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code))
			d.logger.Error(newErr)
			return NewUserBalance, newErr
		}
	}

	return NewUserBalance, nil
}

func (d db) GetUserBalance(ctx context.Context, id uint) (balance.Balance, error) {
	q := `SELECT * FROM balances WHERE user_id = $1`

	var userBalance balance.Balance

	if err := d.client.QueryRow(ctx, q, id).Scan(&userBalance.ID, &userBalance.UserID, &userBalance.Amount); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code))
			d.logger.Error(newErr)
			return userBalance, newErr
		}
	}

	return userBalance, nil
}

func NewDB(client postgresql.Client, logger *logging.Logger) balance.Storage {
	return &db{
		client: client,
		logger: logger,
	}
}
