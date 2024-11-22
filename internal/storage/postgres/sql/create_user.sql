INSERT INTO users (login, password_hash, balance, created_at)
VALUES ($1, $2, $3, $4)
RETURNING id