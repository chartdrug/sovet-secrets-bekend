package entity

import (
	"time"
)

type Cryptocloud struct {
	Id             string    `json:"-"`
	Owner          string    `json:"-"`
	Dt             time.Time `json:"dt"`
	Shopid         string    `json:"-"`
	Amount         float64   `json:"amount"`
	Currency       string    `json:"currency"`
	Email          string    `json:"-"`
	Invoiceid      string    `json:"invoiceid"`
	Status         string    `json:"status"`
	Payurl         string    `json:"payurl"`
	Reshtpp        string    `json:"-"`
	Dtpaym         time.Time `json:"dtpaym"`
	Statusinvoice  string    `json:"Statusinvoice"`
	Resthttpstatus string    `json:"Resthttpstatus"`
}

type CryptocloudPostback struct {
	Status       string  `json:"status"`
	InvoiceId    string  `json:"invoice_id"`
	AmountCrypto float64 `json:"amount_crypto"`
	Currency     string  `json:"currency"`
	OrderId      string  `json:"order_id"`
}
