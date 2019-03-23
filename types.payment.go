package main

import (
	"github.com/shopspring/decimal"
)

// Payment is the top-level struct of a payment record containing transaction
// details.
type Payment struct {
	Amount               decimal.Decimal    `json:"amount"`
	BeneficiaryParty     beneficiaryParty   `json:"beneficiary_party"`
	ChargesInformation   chargesInformation `json:"charges_information"`
	Currency             string             `json:"currency"`
	DebtorParty          debtorParty        `json:"debtor_party"`
	EndToEndReference    string             `json:"end_to_end_reference"`
	Fx                   fx                 `json:"fx"`
	NumericReference     string             `json:"numeric_reference"`
	PaymentID            string             `json:"payment_id"`
	PaymentPurpose       string             `json:"payment_purpose"`
	PaymentScheme        string             `json:"payment_scheme"`
	PaymentType          string             `json:"payment_type"`
	ProcessingDate       string             `json:"processing_date"`
	Reference            string             `json:"reference"`
	SchemePaymentSubType string             `json:"scheme_payment_sub_type"`
	SchemePaymentType    string             `json:"scheme_payment_type"`
	SponsorParty         sponsorParty       `json:"sponsor_party"`
}

// beneficiaryParty details the party being credited from the payment
// transaction.
type beneficiaryParty struct {
	AccountName       string `json:"account_name"`
	AccountNumber     string `json:"account_number"`
	AccountNumberCode string `json:"account_number_code"`
	AccountType       int    `json:"account_type"`
	Address           string `json:"address"`
	BankID            string `json:"bank_id"`
	BankIDCode        string `json:"bank_id_code"`
	Name              string `json:"name"`
}

// chargesInformation describes the charges involved with processing the
// payment transaction.
type chargesInformation struct {
	BearerCode              string          `json:"bearer_code"`
	ReceiverChargesAmount   decimal.Decimal `json:"receiver_charges_amount"`
	ReceiverChargesCurrency string          `json:"receiver_charges_currency"`
	SenderCharges           []senderCharges `json:"sender_charges"`
}

// senderCharges stores the currency and amount charged to the payment sender.
type senderCharges struct {
	Amount   decimal.Decimal `json:"amount"`
	Currency string          `json:"currency"`
}

// debtorParty details the party being debited for the payment transaction.
type debtorParty struct {
	AccountName       string `json:"account_name"`
	AccountNumber     string `json:"account_number"`
	AccountNumberCode string `json:"account_number_code"`
	Address           string `json:"address"`
	BankID            string `json:"bank_id"`
	BankIDCode        string `json:"bank_id_code"`
	Name              string `json:"name"`
}

// fx describes the foreign exchange details salient to the payment
// transaction.
type fx struct {
	ContractReference string          `json:"contract_reference"`
	ExchangeRate      string          `json:"exchange_rate"`
	OriginalAmount    decimal.Decimal `json:"original_amount"`
	OriginalCurrency  string          `json:"original_currency"`
}

// sponsorParty details the party sponsoring the payment transaction.
type sponsorParty struct {
	AccountNumber string `json:"account_number"`
	BankID        string `json:"bank_id"`
	BankIDCode    string `json:"bank_id_code"`
}
