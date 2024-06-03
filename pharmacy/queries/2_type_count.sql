SELECT COUNT(DISTINCT customer.id)
FROM customer
         JOIN orders ON customer.id = orders.customer_id
         JOIN receipt ON orders.receipt_id = receipt.id
         JOIN medicine_list ON receipt.id = medicine_list.receipt_id
         LEFT JOIN medicine ON medicine_list.medicine_id = medicine.id
WHERE orders.status = 'in_production' AND medicine.type = $1;