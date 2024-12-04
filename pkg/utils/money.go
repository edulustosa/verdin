package utils

import "github.com/Rhymond/go-money"

func SumMoney(current float64, transaction *money.Money) float64 {
	currentMoney := money.NewFromFloat(current, money.BRL)

	currentMoney, _ = currentMoney.Add(transaction)
	return currentMoney.AsMajorUnits()
}
