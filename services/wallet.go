package services

import (
	"fmt"
	"math/rand"

	"github.com/greatdaveo/SendlyPay/models"
)

var walletStore = map[string]float64{}

// To fund the wallet
func SetInitialBalance(user string, amount float64) {
	walletStore[user] = amount
}

// To get the current wallet balance
func GetBalance(user string) float64 {
	return walletStore[user]
}

// To check and deduct funds from wallet amount & for insufficient funds
func DeductFromWallet(user string, amount float64) bool {
	current := walletStore[user]
	if current >= amount {
		walletStore[user] = current - amount
		return true
	}

	return false
}

// To generate receipt for each transactions
func GenerateReceipt(from string, info models.PaymentInfo) models.TransactionReceipt {
	return models.TransactionReceipt{
		Reference:     fmt.Sprintf("TXN-%d", rand.Intn(100000000)),
		FromUser:      from,
		ToName:        info.RecipientName,
		AccountNumber: info.AccountNumber,
		SortCode:      info.SortCode,
		Amount:        info.Amount,
		Currency:      info.Currency,
		Status:        "Success",
	}
}

// To view past transactions
var transactionStore = map[string][]models.TransactionReceipt{}

// To save transaction
func SaveTransaction(user string, receipt models.TransactionReceipt) {
	transactionStore[user] = append(transactionStore[user], receipt)
}

// To get the recent transactions
func GetRecentTransactions(user string, count int) []models.TransactionReceipt {
	allTransactions := transactionStore[user]
	if len(allTransactions) <= count {
		return allTransactions
	}

	return allTransactions[len(allTransactions)-count:]
}
