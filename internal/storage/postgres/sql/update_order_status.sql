UPDATE orders
SET status = $2, accrual = $3, processed_at = CURRENT_TIMESTAMP
WHERE id = $1