UPDATE services
SET title = $2, service_description = $3, price = $4
WHERE id = $1;