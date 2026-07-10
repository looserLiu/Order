package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/ledger/backend/pkg/response"
)

type CurrencyHandler struct {
	db *gorm.DB
}

func NewCurrencyHandler(db *gorm.DB) *CurrencyHandler {
	return &CurrencyHandler{db: db}
}

func (h *CurrencyHandler) GetRates(c *gin.Context) {
	response.Success(c, defaultExchangeRates)
}

func (h *CurrencyHandler) Convert(c *gin.Context) {
	var req struct {
		From   string  `json:"from" binding:"required"`
		To     string  `json:"to" binding:"required"`
		Amount float64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	fromRate, ok := defaultExchangeRates[req.From]
	if !ok {
		response.Error(c, http.StatusBadRequest, "Unsupported currency: "+req.From)
		return
	}

	toRate, ok := defaultExchangeRates[req.To]
	if !ok {
		response.Error(c, http.StatusBadRequest, "Unsupported currency: "+req.To)
		return
	}

	converted := (req.Amount / fromRate) * toRate

	response.Success(c, gin.H{
		"from":       req.From,
		"to":         req.To,
		"amount":     req.Amount,
		"converted": converted,
		"rate":       toRate / fromRate,
	})
}

func (h *CurrencyHandler) ListCurrencies(c *gin.Context) {
	currencies := []Currency{
		{Code: "CNY", Name: "人民币", Symbol: "¥", Rate: defaultExchangeRates["CNY"]},
		{Code: "USD", Name: "美元", Symbol: "$", Rate: defaultExchangeRates["USD"]},
		{Code: "EUR", Name: "欧元", Symbol: "€", Rate: defaultExchangeRates["EUR"]},
		{Code: "JPY", Name: "日元", Symbol: "¥", Rate: defaultExchangeRates["JPY"]},
		{Code: "GBP", Name: "英镑", Symbol: "£", Rate: defaultExchangeRates["GBP"]},
		{Code: "HKD", Name: "港币", Symbol: "HK$", Rate: defaultExchangeRates["HKD"]},
		{Code: "KRW", Name: "韩元", Symbol: "₩", Rate: defaultExchangeRates["KRW"]},
		{Code: "AUD", Name: "澳元", Symbol: "A$", Rate: defaultExchangeRates["AUD"]},
		{Code: "CAD", Name: "加元", Symbol: "C$", Rate: defaultExchangeRates["CAD"]},
		{Code: "SGD", Name: "新加坡元", Symbol: "S$", Rate: defaultExchangeRates["SGD"]},
	}

	response.Success(c, currencies)
}

// CSV Import Handler
