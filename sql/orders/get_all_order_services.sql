SELECT s.id, s.title, s.service_description, s.price
FROM orders o
INNER JOIN orders_to_services os
ON o.id = os.order_id
INNER JOIN services s
ON os.service_id = s.id
WHERE o.id = $1;