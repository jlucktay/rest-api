package storage

// Payment is the top-level struct of a payment record containing transaction details.
type Payment struct {
	Amount               float64            `bson:"amount" json:"amount"`
	BeneficiaryParty     BeneficiaryParty   `bson:"beneficiary_party" json:"beneficiary_party"`
	ChargesInformation   ChargesInformation `bson:"charges_information" json:"charges_information"`
	Currency             string             `bson:"currency" json:"currency"`
	DebtorParty          DebtorParty        `bson:"debtor_party" json:"debtor_party"`
	EndToEndReference    string             `bson:"end_to_end_reference" json:"end_to_end_reference"`
	Fx                   Fx                 `bson:"fx" json:"fx"`
	NumericReference     string             `bson:"numeric_reference" json:"numeric_reference"`
	PaymentID            string             `bson:"payment_id" json:"payment_id"`
	PaymentPurpose       string             `bson:"payment_purpose" json:"payment_purpose"`
	PaymentScheme        string             `bson:"payment_scheme" json:"payment_scheme"`
	PaymentType          string             `bson:"payment_type" json:"payment_type"`
	ProcessingDate       string             `bson:"processing_date" json:"processing_date"`
	Reference            string             `bson:"reference" json:"reference"`
	SchemePaymentSubType string             `bson:"scheme_payment_sub_type" json:"scheme_payment_sub_type"`
	SchemePaymentType    string             `bson:"scheme_payment_type" json:"scheme_payment_type"`
	SponsorParty         SponsorParty       `bson:"sponsor_party" json:"sponsor_party"`
}

// BeneficiaryParty details the party being credited from the payment transaction.
type BeneficiaryParty struct {
	AccountName       string `bson:"account_name" json:"account_name"`
	AccountNumber     string `bson:"account_number" json:"account_number"`
	AccountNumberCode string `bson:"account_number_code" json:"account_number_code"`
	AccountType       int    `bson:"account_type" json:"account_type"`
	Address           string `bson:"address" json:"address"`
	BankID            string `bson:"bank_id" json:"bank_id"`
	BankIDCode        string `bson:"bank_id_code" json:"bank_id_code"`
	Name              string `bson:"name" json:"name"`
}

// ChargesInformation describes the charges involved with processing the payment transaction.
type ChargesInformation struct {
	BearerCode              string          `bson:"bearer_code" json:"bearer_code"`
	ReceiverChargesAmount   float64         `bson:"receiver_charges_amount" json:"receiver_charges_amount"`
	ReceiverChargesCurrency string          `bson:"receiver_charges_currency" json:"receiver_charges_currency"`
	SenderCharges           []SenderCharges `bson:"sender_charges" json:"sender_charges"`
}

// SenderCharges stores the currency and amount charged to the payment sender.
type SenderCharges struct {
	Amount   float64 `bson:"amount" json:"amount"`
	Currency string  `bson:"currency" json:"currency"`
}

// DebtorParty details the party being debited for the payment transaction.
type DebtorParty struct {
	AccountName       string `bson:"account_name" json:"account_name"`
	AccountNumber     string `bson:"account_number" json:"account_number"`
	AccountNumberCode string `bson:"account_number_code" json:"account_number_code"`
	Address           string `bson:"address" json:"address"`
	BankID            string `bson:"bank_id" json:"bank_id"`
	BankIDCode        string `bson:"bank_id_code" json:"bank_id_code"`
	Name              string `bson:"name" json:"name"`
}

// Fx describes the foreign exchange details salient to the payment transaction.
type Fx struct {
	ContractReference string  `bson:"contract_reference" json:"contract_reference"`
	ExchangeRate      float64 `bson:"exchange_rate" json:"exchange_rate"`
	OriginalAmount    float64 `bson:"original_amount" json:"original_amount"`
	OriginalCurrency  string  `bson:"original_currency" json:"original_currency"`
}

// SponsorParty details the party sponsoring the payment transaction.
type SponsorParty struct {
	AccountNumber string `bson:"account_number" json:"account_number"`
	BankID        string `bson:"bank_id" json:"bank_id"`
	BankIDCode    string `bson:"bank_id_code" json:"bank_id_code"`
}
