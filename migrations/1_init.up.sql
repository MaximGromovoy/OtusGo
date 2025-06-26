CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    from_currency VARCHAR(10) NOT NULL,
    to_currency VARCHAR(10) NOT NULL,
    from_amount DOUBLE PRECISION NOT NULL,
    to_amount DOUBLE PRECISION NOT NULL,
    exchange_rate DOUBLE PRECISION NOT NULL,
    commission DOUBLE PRECISION NOT NULL,
    status VARCHAR(32),
    timestamp TIMESTAMP
);