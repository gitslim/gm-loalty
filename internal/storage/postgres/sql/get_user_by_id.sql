SELECT id, login, password_hash, balance, created_at
FROM users
WHERE id = $1