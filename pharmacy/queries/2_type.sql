SELECT DISTINCT c.surname,
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
         JOIN medicine_list ON receipt.id = medicine_list.receipt_id
         LEFT JOIN medicine m ON medicine_list.medicine_id = m.id
WHERE o.status = 'in_production' AND m.type = $1;