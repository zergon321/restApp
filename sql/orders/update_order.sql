UPDATE orders
SET contract_date = $2
WHERE id = $1;