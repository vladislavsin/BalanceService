package balance

import (
	"BalanceService/internal/handlers"
	"BalanceService/internal/reservation"
	"BalanceService/internal/transaction"
	"BalanceService/pkg/logging"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type handler struct {
	service *Service
	logger  *logging.Logger
}

func NewHandler(db Storage, logger *logging.Logger) handlers.Handler {
	service := NewBalanceService(db, logger)
	return &handler{
		service: service,
		logger:  logger,
	}
}

func (h *handler) Register(router *gin.Engine) {
	router.GET("/balance/user/:userId", h.GetUserBalance)
	router.GET("/transactions/user/:userId", h.GetTransactionHistory)
	router.POST("/balance/add", h.AddAmount)
	router.POST("/reservation", h.Reservation)
	router.POST("/reservation/accept", h.AcceptReservation)
}

func (h *handler) GetTransactionHistory(ctx *gin.Context) {
	var dto BalanceDTO

	param := ctx.Param("userId")

	u64, err := strconv.ParseUint(param, 10, 0)
	if err != nil {
		h.logger.Fatal(err)
	}

	dto.UserID = uint(u64)

	userBalance, err := h.service.GetUserBalance(dto)
	if err != nil {
		h.logger.Fatal(err)
	}

	if userBalance.UserID == 0 {
		ctx.JSON(404, gin.H{
			"message": fmt.Sprintf("Баланса пользователя с id: %d - не найдено", dto.UserID),
		})
		return
	}

	page := ctx.Query("page")
	sorting := ctx.Query("sort")
	orderBy := ctx.Query("order")

	pagination := uint(0)

	if page != "" && page != "0" {
		ui64, err := strconv.ParseUint(page, 10, 0)
		if err != nil {
			h.logger.Fatal(err)
		}
		pagination = (uint(ui64) - 1) * 10
	}

	if sorting != "created_at" && sorting != "amount" {
		sorting = "created_at"
	}

	if orderBy != "desc" && orderBy != "asc" {
		orderBy = "desc"
	}

	sort := transaction.NewSortingHistory(pagination, sorting, orderBy)
	response := h.service.GetTransactionHistory(userBalance, sort)
	ctx.JSON(200, gin.H{
		"message": response,
	})
}

func (h *handler) AddAmount(ctx *gin.Context) {
	var dto BalanceDTO

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := h.service.AddAmount(dto)
	ctx.JSON(200, gin.H{
		"message": res,
	})
}

func (h *handler) GetUserBalance(ctx *gin.Context) {
	var dto BalanceDTO

	param := ctx.Param("userId")

	u64, err := strconv.ParseUint(param, 10, 0)
	if err != nil {
		h.logger.Fatal(err)
	}

	dto.UserID = uint(u64)

	userBalance, err := h.service.GetUserBalance(dto)
	if err != nil {
		h.logger.Fatal(err)
	}

	if userBalance.UserID == 0 {
		ctx.JSON(404, gin.H{
			"message": fmt.Sprintf("Баланса пользователя с id: %d - не найдено", dto.UserID),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"UserID": userBalance.UserID,
		"Amount": userBalance.Amount,
	})
}

func (h *handler) AcceptReservation(ctx *gin.Context) {
	var dto reservation.ReservationDTO

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.AcceptReservation(dto)
	if err != nil {
		h.logger.Fatal(err)
	}
	ctx.JSON(200, gin.H{
		"message": resp,
	})
}

func (h *handler) Reservation(ctx *gin.Context) {
	var dto reservation.ReservationDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.MakeReservation(dto)
	if err != nil {
		h.logger.Fatal(err)
	}

	ctx.JSON(200, gin.H{
		"message": res,
	})
}
