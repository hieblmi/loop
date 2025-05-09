CREATE TABLE IF NOT EXISTS withdrawals (
    -- id is the auto-incrementing primary key for a withdrawal.
    id INTEGER PRIMARY KEY,

    withdrawal_tx_id TEXT NOT NULL UNIQUE,

    deposit_outpoints TEXT NOT NULL,

    total_deposit_amount BIGINT NOT NULL,

    withdrawn_amount BIGINT NOT NULL,

    change_amount BIGINT NOT NULL,

    confirmation_height BIGINT NOT NULL
);