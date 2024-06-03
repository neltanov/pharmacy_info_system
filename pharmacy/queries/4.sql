SELECT substance.*, SUM(substance_usage_statistics.quantity_used) AS total_used
FROM substance
         JOIN substance_usage_statistics ON substance.id = substance_usage_statistics.substance_id
WHERE substance.name IN ($1)
  AND substance_usage_statistics.usage_time BETWEEN $2 AND $3
GROUP BY substance.id
HAVING SUM(substance_usage_statistics.quantity_used) > 0;