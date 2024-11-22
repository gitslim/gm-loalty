INSERT INTO withdrawals (user_id, order_number, sum, processed_at)
VALUES ($1, $2, $3, $4)
RETURNING id