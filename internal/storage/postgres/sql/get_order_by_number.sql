SELECT id, number, user_id, status, accrual, uploaded_at, processed_at
FROM orders
WHERE number = $1