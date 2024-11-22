INSERT INTO orders (number, user_id, status, accrual, uploaded_at, processed_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id