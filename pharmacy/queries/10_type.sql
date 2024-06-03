SELECT m.id AS medicine_id,
       m.name AS medicine_name,
       pt.method_of_production,
       pt.time_to_product
FROM medicine m
         JOIN local_medicine lm ON m.id = lm.medicine_id
         JOIN production_techonology pt ON lm.production_techology = pt.id
WHERE
    m.type IN ($1);