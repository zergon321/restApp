DELETE FROM orders_to_services
WHERE order_id = $1 AND service_id = $2;