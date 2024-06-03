SELECT m.name,
       mw.total_amount,
       m.type AS medicine_type
FROM medicine_warehouse mw
         JOIN medicine m ON mw.medicine_id = m.id
ORDER BY mw.total_amount, m.id;