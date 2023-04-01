package entity

import "time"

type Wallet struct {
	tableName struct{} `pg:"mini_wallets"`
	ID string `json:"id" pg:"id,pk"`
	OwnedBy string `json:"-" pg:"owned_by"`
	Balance float64 `json:"-" pg:"balance"`
	IsEnabled bool `json:"-" pg:"is_enabled"`
	EnabledAt time.Time `json:"-" pg:"enabled_at"`
	DisabledAt time.Time `json:"-" pg:"disabled_at"`
}
