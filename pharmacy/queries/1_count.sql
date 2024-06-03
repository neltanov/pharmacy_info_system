SELECT COUNT(DISTINCT customer.id)
FROM customer
         JOIN orders ON customer.id = orders.customer_id
WHERE orders.status = 'done'
  AND orders.production_date < CURRENT_DATE;