package handler

type ResponseWallet struct {
	ID        string  `json:"id"`
	OwnedBy   string  `json:"owned_by"`
	Status    string  `json:"status"`
	EnabledAt string  `json:"enabled_at"`
	Balance   float64 `json:"balance"`
}

type ResponseDepositWallet struct {
	ID          string  `json:"id"`
	DepositedBy string  `json:"deposited_by"`
	Status      string  `json:"status"`
	DepositAt   string  `json:"deposited_at"`
	Amount      float64 `json:"amount"`
	ReferenceId string  `json:"reference_id"`
}

type ResponseWithdrawWallet struct {
	ID          string  `json:"id"`
	WithdrawnBy string  `json:"withdrawn_by"`
	Status      string  `json:"status"`
	WithdrawnAt string  `json:"withdrawn_at"`
	Amount      float64 `json:"amount"`
	ReferenceId string  `json:"reference_id"`
}

type ResponseTransactions struct {
	ID          	string  `json:"id"`
	Status      	string  `json:"status"`
	TransactedAt	string  `json:"transacted_at"`
	Type 			string  `json:"type"`
	Amount      	float64 `json:"amount"`
	ReferenceId 	string  `json:"reference_id"`
}

type EmptyResponse struct {
}
