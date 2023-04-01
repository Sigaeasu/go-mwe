-- BEGIN;
CREATE TABLE IF NOT EXISTS mini_wallets 
(
    id uuid DEFAULT gen_random_uuid () PRIMARY KEY,
    owned_by uuid NOT NULL UNIQUE,
    balance FLOAT NOT NULL DEFAULT 0,
    is_enabled BOOLEAN NOT NULL DEFAULT true,
    enabled_at TIMESTAMP DEFAULT now(),
    disabled_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS transactions (
    id uuid DEFAULT gen_random_uuid () PRIMARY KEY,
    amount FLOAT NOT NULL DEFAULT 0,
    type VARCHAR NOT NULL,
    status VARCHAR NOT NULL,
    reference_id uuid NOT NULL,
    created_by uuid NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- COMMIT;