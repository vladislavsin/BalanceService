package balance

import (
	"BalanceService/internal/reservation"
	"BalanceService/internal/transaction"
	"BalanceService/pkg/logging"
	"context"
	"fmt"
	"time"
)

type Service struct {
	db     Storage
	logger *logging.Logger
}

func NewBalanceService(db Storage, logger *logging.Logger) *Service {
	return &Service{
		db:     db,
		logger: logger,
	}
}

func (s *Service) GetTransactionHistory(balance Balance, sort *transaction.SortingHistory) []transaction.ResponseTransactionDTO {
	transactions, err := s.db.GetTransactionHistory(context.TODO(), balance, sort)
	if err != nil {
		s.logger.Fatal(err)
	}

	responseTransactions := make([]transaction.ResponseTransactionDTO, 0)

	for _, item := range transactions {
		var responseTransaction transaction.ResponseTransactionDTO
		responseTransaction.TransactionType = s.getTransactionTypeMessage(item.TransactionTypeID, item.ServiceID)
		responseTransaction.Amount = item.Amount
		responseTransaction.Date = item.CreatedAt.Format("2006-01-02")
		responseTransactions = append(responseTransactions, responseTransaction)
	}

	return responseTransactions
}

func (s *Service) AddAmount(dto BalanceDTO) string {
	userBalance, err := s.db.GetUserBalance(context.TODO(), dto.UserID)
	if err != nil {
		s.logger.Fatal(err)
	}

	if userBalance.UserID != 0 {
		userBalance.Amount = userBalance.Amount + dto.Amount
		userBalance, err := s.db.AddAmount(context.TODO(), userBalance)
		if err != nil {
			s.logger.Fatal(err)
		}
		transactionHistory := transaction.New(userBalance.ID, transaction.AddingAmount, dto.Amount, 0)
		if err := s.db.AddTransactionHistory(context.TODO(), transactionHistory); err != nil {
			s.logger.Fatal(err)
		}

		return fmt.Sprintf("Пользователю с id: %d успешно начислена сумма в размере: %d", userBalance.UserID, dto.Amount)
	}

	newUserBalance, err := s.db.CreateUserBalance(context.TODO(), dto)
	if err != nil {
		s.logger.Fatalf("%v", err)
	}

	transactionHistory := transaction.New(newUserBalance.ID, transaction.AddingAmount, dto.Amount, 0)
	if err := s.db.AddTransactionHistory(context.TODO(), transactionHistory); err != nil {
		s.logger.Fatal(err)
	}

	return fmt.Sprintf("Пользователю с id: %d заведен баланс и начислена сумма в размере: %d", newUserBalance.UserID, dto.Amount)
}

func (s *Service) MakeReservation(dto reservation.ReservationDTO) (string, error) {
	userBalance, err := s.db.GetUserBalance(context.TODO(), dto.UserID)
	if err != nil {
		s.logger.Fatal(err)
	}

	if userBalance.UserID == 0 {
		return fmt.Sprintf("У пользователя с id: %d - нет баланса", dto.UserID), nil
	}

	if dto.Amount > userBalance.Amount {
		return fmt.Sprintf("У пользователя с id: %d - недостаточно средств на балансе", dto.UserID), nil
	}

	userBalance.Amount = userBalance.Amount - dto.Amount
	if err := s.db.WithdrawFunds(context.TODO(), userBalance); err != nil {
		s.logger.Fatal(err)
	}

	date := time.Now().Format("2006-01-02 15:04")
	dto.CreatedAt = date
	dto.UpdatedAt = date

	if err := s.db.Reservation(context.TODO(), dto); err != nil {
		s.logger.Fatal(err)
	}

	transactionHistory := transaction.New(userBalance.ID, transaction.Reservation, dto.Amount, dto.ServiceID)
	if err := s.db.AddTransactionHistory(context.TODO(), transactionHistory); err != nil {
		s.logger.Fatal(err)
	}

	return "Резервация средств прошла успешно!", nil
}

func (s *Service) AcceptReservation(dto reservation.ReservationDTO) (string, error) {
	reserv, err := s.db.GetReservation(context.TODO(), dto.OrderID)
	if err != nil {
		s.logger.Fatal(err)
	}

	if reserv.UserID == 0 {
		return fmt.Sprintf("Резерва с таким order_id: %d - нет", dto.UserID), nil
	}

	if err := s.db.AcceptReservation(context.TODO(), reserv); err != nil {
		s.logger.Fatal(err)
	}

	userBalance, err := s.db.GetUserBalance(context.TODO(), reserv.UserID)
	if err != nil {
		s.logger.Fatal(err)
	}

	transactionHistory := transaction.New(userBalance.ID, transaction.PaidService, reserv.Amount, reserv.ServiceID)
	if err := s.db.AddTransactionHistory(context.TODO(), transactionHistory); err != nil {
		s.logger.Fatal(err)
	}

	return "Cписание средств из резерва прошло успешно!", nil

}

func (s *Service) GetUserBalance(dto BalanceDTO) (Balance, error) {
	userBalance, err := s.db.GetUserBalance(context.TODO(), dto.UserID)
	if err != nil {
		s.logger.Fatal(err)
	}

	return userBalance, nil
}

func (s *Service) getTransactionTypeMessage(transactionType uint, serviceID uint) string {
	switch transactionType {
	case transaction.AddingAmount:
		return "Зачисление средств на баланс"
	case transaction.Reservation:
		return fmt.Sprintf("Резервация средств за услугу с id: %d", serviceID)
	case transaction.PaidService:
		return fmt.Sprintf("Списание средств за услугу с id: %d", serviceID)
	default:
		return "неизвестный тип"
	}
}
