SELECT m.name,
       m.price,
       mw.total_amount,
       m.type AS medicine_type,
       pt.method_of_production,
       pt.time_to_product
FROM medicine m
         LEFT JOIN local_medicine lm ON m.id = lm.medicine_id
         LEFT JOIN imported_medicine im ON m.id = im.medicine_id
         LEFT JOIN production_techonology pt ON lm.production_techology = pt.id
         LEFT JOIN medicine_warehouse mw ON m.id = mw.medicine_id
WHERE m.name = $1;