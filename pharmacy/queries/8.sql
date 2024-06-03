SELECT o.id,
       o.customer_id,
       o.receipt_id,
       o.order_date,
       o.production_date,
       o.status
FROM orders o
WHERE o.status = 'in_production';