package entity

import "time"

type Transaction struct {
	tableName	struct{} 	`pg:"transactions"`
	ID 			string 		`json:"id" pg:"id,pk"`
	Amount 		float64 	`json:"amount" pg:"amount"`
	Type 		string 		`json:"type" pg:"type"`
	ReferenceID string 		`json:"reference_id" pg:"reference_id"`
	Status 		string 		`json:"status" pg:"status"`
	CreatedBy 	string 		`json:"-" pg:"created_by"`
	CreatedAt 	time.Time 	`json:"transacted_at" pg:"created_at"`
}
