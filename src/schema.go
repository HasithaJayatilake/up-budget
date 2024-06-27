package main

import "time"

type Links struct {
	Self string `json:"self" db:"self_link"`
}

type Amount struct {
	CurrencyCode     string `json:"currencyCode" db:"currency_code"`
	Value            string `json:"value" db:"value"`
	ValueInBaseUnits int    `json:"valueInBaseUnits" db:"value_in_base_units"`
}

type HoldInfo struct {
	Amount        Amount  `json:"amount" db:"hold_amount"`
	ForeignAmount *Amount `json:"foreignAmount,omitempty" db:"hold_foreign_amount"`
}

type Attributes struct {
	Status             string     `json:"status" db:"status"`
	RawText            string     `json:"rawText" db:"raw_text"`
	Description        string     `json:"description" db:"description"`
	Message            *string    `json:"message,omitempty" db:"message"`
	IsCategorizable    bool       `json:"isCategorizable" db:"is_categorizable"`
	HoldInfo           *HoldInfo  `json:"holdInfo,omitempty" db:"hold_info"`
	RoundUp            *Amount    `json:"roundUp,omitempty" db:"round_up"`
	Cashback           *Amount    `json:"cashback,omitempty" db:"cashback"`
	Amount             Amount     `json:"amount" db:"amount"`
	ForeignAmount      *Amount    `json:"foreignAmount,omitempty" db:"foreign_amount"`
	CardPurchaseMethod *string    `json:"cardPurchaseMethod,omitempty" db:"card_purchase_method"`
	SettledAt          *time.Time `json:"settledAt,omitempty" db:"settled_at"`
	CreatedAt          time.Time  `json:"createdAt" db:"created_at"`
}

type Transaction struct {
	Type       string     `json:"type" db:"type"`
	ID         string     `json:"id" db:"id"`
	Attributes Attributes `json:"attributes" db:"attributes"`
}

type TransactionAPIResponse struct {
	Data  []Transaction `json:"data"`
	Links Links         `json:"links"`
}
