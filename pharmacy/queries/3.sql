SELECT m.name, m.type, SUM(mus.quantity_used) AS total_used
FROM medicine_usage_statistics mus
         JOIN medicine m ON mus.medicine_id = m.id
GROUP BY m.id
HAVING SUM(mus.quantity_used) > 0
ORDER BY total_used DESC
LIMIT 10;