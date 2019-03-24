package main

import (
	"github.com/shopspring/decimal"
)

// Payment is the top-level struct of a payment record containing transaction
// details.
type Payment struct {
	Amount               decimal.Decimal    `json:"amount"`
	BeneficiaryParty     BeneficiaryParty   `json:"beneficiary_party"`
	ChargesInformation   ChargesInformation `json:"charges_information"`
	Currency             string             `json:"currency"`
	DebtorParty          DebtorParty        `json:"debtor_party"`
	EndToEndReference    string             `json:"end_to_end_reference"`
	Fx                   Fx                 `json:"fx"`
	NumericReference     string             `json:"numeric_reference"`
	PaymentID            string             `json:"payment_id"`
	PaymentPurpose       string             `json:"payment_purpose"`
	PaymentScheme        string             `json:"payment_scheme"`
	PaymentType          string             `json:"payment_type"`
	ProcessingDate       string             `json:"processing_date"`
	Reference            string             `json:"reference"`
	SchemePaymentSubType string             `json:"scheme_payment_sub_type"`
	SchemePaymentType    string             `json:"scheme_payment_type"`
	SponsorParty         SponsorParty       `json:"sponsor_party"`
}

// BeneficiaryParty details the party being credited from the payment
// transaction.
type BeneficiaryParty struct {
	AccountName       string `json:"account_name"`
	AccountNumber     string `json:"account_number"`
	AccountNumberCode string `json:"account_number_code"`
	AccountType       int    `json:"account_type"`
	Address           string `json:"address"`
	BankID            string `json:"bank_id"`
	BankIDCode        string `json:"bank_id_code"`
	Name              string `json:"name"`
}

// ChargesInformation describes the charges involved with processing the
// payment transaction.
type ChargesInformation struct {
	BearerCode              string          `json:"bearer_code"`
	ReceiverChargesAmount   decimal.Decimal `json:"receiver_charges_amount"`
	ReceiverChargesCurrency string          `json:"receiver_charges_currency"`
	SenderCharges           []SenderCharges `json:"sender_charges"`
}

// SenderCharges stores the currency and amount charged to the payment sender.
type SenderCharges struct {
	Amount   decimal.Decimal `json:"amount"`
	Currency string          `json:"currency"`
}

// DebtorParty details the party being debited for the payment transaction.
type DebtorParty struct {
	AccountName       string `json:"account_name"`
	AccountNumber     string `json:"account_number"`
	AccountNumberCode string `json:"account_number_code"`
	Address           string `json:"address"`
	BankID            string `json:"bank_id"`
	BankIDCode        string `json:"bank_id_code"`
	Name              string `json:"name"`
}

// Fx describes the foreign exchange details salient to the payment
// transaction.
type Fx struct {
	ContractReference string          `json:"contract_reference"`
	ExchangeRate      string          `json:"exchange_rate"`
	OriginalAmount    decimal.Decimal `json:"original_amount"`
	OriginalCurrency  string          `json:"original_currency"`
}

// SponsorParty details the party sponsoring the payment transaction.
type SponsorParty struct {
	AccountNumber string `json:"account_number"`
	BankID        string `json:"bank_id"`
	BankIDCode    string `json:"bank_id_code"`
}
