package service

import (
	"errors"
	"math"
)

type ExchangeService interface {
	GetConvertedTotal(target_currency string, amount int) (total float64, err error)
}

type CurrencyRepo interface {
	FetchPrice(currency_name string) (price float64, found bool)
}

var (
	ErrInvalidTargetCurrency  = errors.New("Invalid Target Currency Requested")
	ErrTargetCurrencyNotFound = errors.New("Currency Not Found")
	ErrInvalidAmount          = errors.New("Invalid Amount Requested")
)

type convertService struct {
	repo CurrencyRepo
}

func NewExchangeService(cr CurrencyRepo) (cs *convertService) {
	cs = &convertService{
		repo: cr,
	}

	return cs
}

func (cs *convertService) GetConvertedTotal(target_currency string, amount int) (total float64, err error) {
	if target_currency == "" {
		return 0.0, ErrInvalidTargetCurrency
	}
	if amount <= 0 {
		return 0.0, ErrInvalidAmount
	}

	price, found := cs.repo.FetchPrice(target_currency)
	if !found {
		return 0.0, ErrTargetCurrencyNotFound
	}

	total = price * float64(amount)

	return math.Round(total*100) / 100, nil
}
