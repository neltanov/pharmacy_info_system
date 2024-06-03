SELECT m.name,
       SUM(ml.quantity_used) AS total_quantity
FROM medicine_list ml
         JOIN orders o ON ml.receipt_id = o.receipt_id
         JOIN medicine m ON ml.medicine_id = m.id
WHERE o.status = 'in_production'
GROUP BY ml.medicine_id, m.name;