SELECT DISTINCT c.id,
                c.surname,
                c.name,
                c.middle_name,
                c.phone_number,
                c.address,
                o.receipt_id,
                o.order_date,
                o.production_date,
                o.status
FROM customer c
         JOIN orders o ON c.id = o.customer_id
         JOIN receipt ON o.receipt_id = receipt.id
WHERE o.status = 'in_production';