-- 1. Получить сведения о покупателях, которые не пришли забрать свой заказ в назначенное им время и общее их число.
SELECT customer.*, orders.*
FROM customer
         JOIN orders ON customer.id = orders.customer_id
WHERE orders.status = 'done'
  AND orders.production_date < CURRENT_DATE;

SELECT COUNT(DISTINCT customer.id)
FROM customer
         JOIN orders ON customer.id = orders.customer_id
WHERE orders.status = 'done'
  AND orders.production_date < CURRENT_DATE;

-- 2. Получить перечень и общее число покупателей,
-- которые ждут прибытия на склад нужных им медикаментов в целом
-- и по указанной категории медикаментов.

-- В целом
SELECT DISTINCT c.id,
                c.surname,
                c.name,
                c.middle_name,
                c.phone_number,
                c.address,
                o.receipt_id,
                o.order_date,
                o.production_date,
                o.status
FROM customer c
         JOIN orders o ON c.id = o.customer_id
         JOIN receipt ON o.receipt_id = receipt.id
WHERE o.status = 'in_production';

SELECT COUNT(DISTINCT customer.id)
FROM customer
         JOIN orders ON customer.id = orders.customer_id
WHERE orders.status = 'in_production';

-- Для указанной категории медикаментов
SELECT DISTINCT c.surname,
                c.name,
                c.middle_name,
                c.phone_number,
                c.address,
                o.receipt_id,
                o.order_date,
                o.production_date,
                o.status
FROM customer c
         JOIN orders o ON c.id = o.customer_id
         JOIN receipt ON o.receipt_id = receipt.id
         JOIN medicine_list ON receipt.id = medicine_list.receipt_id
         LEFT JOIN medicine m ON medicine_list.medicine_id = m.id
WHERE o.status = 'in_production' AND m.type = 'pill';

SELECT COUNT(DISTINCT customer.id)
FROM customer
         JOIN orders ON customer.id = orders.customer_id
         JOIN receipt ON orders.receipt_id = receipt.id
         JOIN medicine_list ON receipt.id = medicine_list.receipt_id
         LEFT JOIN medicine ON medicine_list.medicine_id = medicine.id
WHERE orders.status = 'in_production' AND medicine.type = 'pill';

-- 3. Получить перечень десяти наиболее часто используемых медикаментов в целом и указанной категории медикаментов.

-- В целом
SELECT m.name, SUM(mus.quantity_used) AS total_used
FROM medicine_usage_statistics mus
         JOIN medicine m ON mus.medicine_id = m.id
GROUP BY m.id
HAVING SUM(mus.quantity_used) > 0
ORDER BY total_used DESC
LIMIT 10;

-- Перечень десяти наиболее часто используемых медикаментов для указанной категории
-- Для локальных медикаментов
SELECT m.name, m.type, SUM(mus.quantity_used) AS total_used
FROM medicine m
         JOIN medicine_usage_statistics mus ON m.id = mus.medicine_id
WHERE m.type = 'pill' -- нужный тип локального медикамента
GROUP BY m.id, m.name, m.type
HAVING SUM(mus.quantity_used) > 0
ORDER BY total_used DESC
LIMIT 10;

-- 4. Получить какой объем указанных веществ использован за указанный период.
SELECT substance.*, SUM(substance_usage_statistics.quantity_used) AS total_used
FROM substance
         JOIN substance_usage_statistics ON substance.id = substance_usage_statistics.substance_id
WHERE substance.id IN (2, 14)                                                     -- нужные вещества
  AND substance_usage_statistics.usage_time BETWEEN '2024-01-01' AND '2024-08-01' -- нужный период
GROUP BY substance.id
HAVING SUM(substance_usage_statistics.quantity_used) > 0;

-- 5. Получить перечень и общее число покупателей, заказывавших определенное лекарство или определенные типы лекарств за данный период.
SELECT DISTINCT c.surname,
                c.name,
                c.middle_name,
                c.phone_number,
                c.address,
                o.order_date,
                o.status
FROM customer c
         JOIN orders o ON c.id = o.customer_id
         JOIN receipt ON o.receipt_id = receipt.id
         JOIN medicine_list ON receipt.id = medicine_list.receipt_id
WHERE medicine_list.medicine_id = 1                                 -- нужное лекарство
  AND o.order_date BETWEEN '2024-01-01' AND '2024-08-01';      -- нужный период

SELECT COUNT(DISTINCT customer.id)
FROM customer
         JOIN orders ON customer.id = orders.customer_id
         JOIN receipt ON orders.receipt_id = receipt.id
         JOIN medicine_list ON receipt.id = medicine_list.receipt_id
WHERE medicine_list.medicine_id = 1                                  -- нужное лекарство
  AND orders.order_date BETWEEN '2024-01-01' AND '2024-08-01';       -- нужный период

-- для определенных типов лекарств
SELECT DISTINCT c.surname,
                c.name,
                c.middle_name,
                c.phone_number,
                c.address,
                o.order_date,
                o.status
FROM customer c
         JOIN orders o ON c.id = o.customer_id
         JOIN receipt ON o.receipt_id = receipt.id
         JOIN medicine_list ON receipt.id = medicine_list.receipt_id
         JOIN medicine ON medicine_list.medicine_id = medicine.id
WHERE medicine.type = 'pill' AND o.order_date BETWEEN '2024-01-01' AND '2024-08-01';

SELECT COUNT(DISTINCT customer.id)
FROM customer
         JOIN orders ON customer.id = orders.customer_id
         JOIN receipt ON orders.receipt_id = receipt.id
         JOIN medicine_list ON receipt.id = medicine_list.receipt_id
         JOIN medicine ON medicine_list.medicine_id = medicine.id
WHERE medicine.type = 'pill' AND orders.order_date BETWEEN '2024-01-01' AND '2024-08-01';

-- 6. Получить перечень и типы лекарств, достигших своей критической нормы или закончившихся.
SELECT m.name AS medicine_name,
       mw.total_amount,
       mw.critical_limit,
       m.type AS medicine_type
FROM medicine_warehouse mw
         JOIN medicine m ON mw.medicine_id = m.id
WHERE mw.total_amount <= mw.critical_limit;

-- 7. Получить перечень лекарств с минимальным запасом на складе в целом и по указанной категории медикаментов.
-- в целом
SELECT m.name,
       mw.total_amount,
       m.type AS medicine_type
FROM medicine_warehouse mw
         JOIN medicine m ON mw.medicine_id = m.id
ORDER BY mw.total_amount, m.id;

-- по указанной категории медикаментов
SELECT
    m.name,
    mw.total_amount,
    m.type AS medicine_type
FROM medicine_warehouse mw
         JOIN medicine m ON mw.medicine_id = m.id
WHERE m.type = 'ointment'
ORDER BY mw.total_amount, m.id;

-- 8. Получить полный перечень и общее число заказов, находящихся в производстве.
-- Полный перечень заказов находящихся в производстве
SELECT o.id,
       o.customer_id,
       o.receipt_id,
       o.order_date,
       o.production_date,
       o.status
FROM orders o
WHERE o.status = 'in_production';

-- Общее число заказов находящихся в производстве
SELECT COUNT(*) FROM orders
WHERE status = 'in_production';


-- 9. Полный перечень и общее число препаратов, требующихся для заказов, находящихся в производстве
SELECT m.name,
       SUM(ml.quantity_used) AS total_quantity
FROM medicine_list ml
         JOIN orders o ON ml.receipt_id = o.receipt_id
         JOIN medicine m ON ml.medicine_id = m.id
WHERE o.status = 'in_production'
GROUP BY ml.medicine_id, m.name;

-- 10. Все технологии приготовления лекарств указанных типов,
-- конкретных лекарств, лекарств, находящихся в справочнике заказов в производстве

-- Технологии приготовления конкретных лекарств
SELECT
    m.id AS medicine_id,
    m.name AS medicine_name,
    pt.method_of_production,
    pt.time_to_product
FROM medicine m
         JOIN local_medicine lm ON m.id = lm.medicine_id
         JOIN production_techonology pt ON lm.production_techology = pt.id
WHERE m.id IN (11, 12);

-- Технологии приготовления лекарств указанных типов
SELECT m.id AS medicine_id,
       m.name AS medicine_name,
       pt.method_of_production,
       pt.time_to_product
FROM medicine m
         JOIN local_medicine lm ON m.id = lm.medicine_id
         JOIN production_techonology pt ON lm.production_techology = pt.id
WHERE
    m.type IN ('powder');

-- Технологии приготовления лекарств, находящихся в заказах в производстве
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

-- 11. Сведения о ценах на указанное лекарство в готовом виде, об объеме и ценах на все компоненты,
-- требующиеся для этого лекарства
-- Сведения о цене готового лекарства
SELECT m.id,
       m.name,
       m.price
FROM medicine m
WHERE m.id = 11;

-- Объем и цены на компоненты, требующиеся для лекарства
SELECT s.id AS id,
       s.name AS name,
       mc.required_quantity,
       s.price AS price
FROM medicine_composition mc
         JOIN substance s ON mc.substance_id = s.id
         JOIN local_medicine lm ON mc.medicine_id = lm.id
WHERE lm.medicine_id = 11;

-- 12. Сведения о наиболее часто делающих заказы клиентах на медикаменты определенного типа,
-- на конкретные медикаменты
-- Сведения о клиентах, часто заказывающих медикаменты определенного типа
SELECT c.surname,
       c.name,
       c.middle_name,
       COUNT(o.id) AS order_count
FROM customer c
         JOIN orders o ON c.id = o.customer_id
         JOIN medicine_list ml ON o.receipt_id = ml.receipt_id
         JOIN medicine m ON ml.medicine_id = m.id
WHERE m.type = 'pill'
GROUP BY c.id
ORDER BY order_count DESC
LIMIT 10;

-- Сведения о клиентах, часто заказывающих конкретные медикаменты
SELECT c.surname,
       c.name,
       c.middle_name,
       COUNT(o.id) AS order_count
FROM customer c
         JOIN orders o ON c.id = o.customer_id
         JOIN medicine_list ml ON o.receipt_id = ml.receipt_id
WHERE ml.medicine_id = 11
GROUP BY c.id
ORDER BY order_count DESC
LIMIT 10;


-- 13. Сведения о конкретном лекарстве (его тип, способ приготовления,
-- названия всех компонентов, цены, его количество на складе)
-- Сведения о конкретном лекарстве
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
WHERE m.id = 1;

-- Названия всех компонентов и их цены
SELECT s.name AS substance_name,
       mc.required_quantity,
       s.price AS substance_price
FROM medicine_composition mc
         JOIN substance s ON mc.substance_id = s.id
         JOIN local_medicine lm ON mc.medicine_id = lm.id
WHERE lm.medicine_id = 11;
