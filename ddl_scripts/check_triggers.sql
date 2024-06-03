-- Проверка триггеров

select * from medicine_warehouse;
select * from substance_warehouse;
update substance_warehouse set total_amount=500 where total_amount!=500;
update medicine_warehouse set total_amount=1000 where total_amount!=1000;

-- Добавление в список готового медикамента (Парацетамол)
insert into medicine_list(receipt_id, medicine_id, quantity_used) VALUES (6, 1, 900);
-- Добавление в список медикамента, изготавливаемого аптекой (Микстура от кашля)
insert into medicine_list(receipt_id, medicine_id, quantity_used) VALUES (6, 11, 1);
-- Удаление лекарства из списка
delete from medicine_list where receipt_id=6;

insert into orders(customer_id, receipt_id, order_date, production_date, status)
VALUES (5, 6, NOW(), NOW(), 'in_production');

select * from medicine_usage_statistics;
select * from substance_usage_statistics;
