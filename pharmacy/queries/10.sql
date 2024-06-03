SELECT DISTINCT m.id AS medicine_id,
                m.name AS medicine_name,
                pt.method_of_production,
                pt.time_to_product
FROM medicine m
         JOIN local_medicine lm ON m.id = lm.medicine_id
         JOIN production_techonology pt ON lm.production_techology = pt.id
         JOIN medicine_list ml ON m.id = ml.medicine_id
         JOIN orders o ON ml.receipt_id = o.receipt_id
WHERE o.status = 'in_production';