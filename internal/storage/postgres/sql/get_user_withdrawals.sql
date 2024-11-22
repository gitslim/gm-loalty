SELECT id, user_id, order_number, sum, processed_at
FROM withdrawals
WHERE user_id = $1
ORDER BY processed_at DESC