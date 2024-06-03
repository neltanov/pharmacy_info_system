SELECT m.name AS medicine_name,
       mw.total_amount,
       mw.critical_limit,
       m.type AS medicine_type
FROM medicine_warehouse mw
         JOIN medicine m ON mw.medicine_id = m.id
WHERE mw.total_amount <= mw.critical_limit;