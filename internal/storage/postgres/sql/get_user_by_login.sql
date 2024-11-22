SELECT id, login, password_hash, balance, created_at
FROM users
WHERE login = $1