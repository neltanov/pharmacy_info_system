select o.id as "Order ID", o.order_date as "Order Date", o.production_date as "Production date", o.status as "Status",
       c.surname, c.name, c.middle_name, c.phone_number, c.address,
       d.id, d.surname, d.name, d.middle_name,
       p.surname, p.name, p.middle_name,
       r.id
from orders o
join b_neltanov.customer c on o.customer_id = c.id
join b_neltanov.receipt r on r.id = o.receipt_id
join b_neltanov.doctor d on d.id = r.doctor_id
join b_neltanov.patient p on p.id = r.patient_id
