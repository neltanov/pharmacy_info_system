SELECT m.name, m.type, SUM(mus.quantity_used) AS total_used
FROM medicine m
         JOIN medicine_usage_statistics mus ON m.id = mus.medicine_id
WHERE m.type = $1
GROUP BY m.id, m.name, m.type
HAVING SUM(mus.quantity_used) > 0
ORDER BY total_used DESC
LIMIT 10;