SELECT c.surname,
       c.name,
       c.middle_name,
       COUNT(o.id) AS order_count
FROM customer c
         JOIN orders o ON c.id = o.customer_id
         JOIN medicine_list ml ON o.receipt_id = ml.receipt_id
         JOIN medicine m ON ml.medicine_id = m.id
WHERE m.type = $1
GROUP BY c.id
ORDER BY order_count DESC
LIMIT 10;