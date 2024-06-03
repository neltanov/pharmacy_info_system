CREATE TYPE "medicine_type" AS ENUM (
  'pill',
  'ointment',
  'tincture',
  'mixture',
  'solution',
  'powder'
);

CREATE TYPE "order_status" AS ENUM (
  'in_production',
  'done'
);

CREATE TABLE "medicine" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar NOT NULL,
  "type" medicine_type NOT NULL,
  "price" float NOT NULL,
  "expiration_date" date NOT NULL
);

CREATE TABLE "substance" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar NOT NULL,
  "price" float NOT NULL
);

CREATE TABLE "medicine_list" (
  "id" SERIAL PRIMARY KEY,
  "receipt_id" int NOT NULL,
  "medicine_id" int NOT NULL,
  "quantity_used" float NOT NULL
);

CREATE TABLE "receipt" (
  "id" SERIAL PRIMARY KEY,
  "doctor_id" int NOT NULL,
  "patient_id" int NOT NULL
);

CREATE TABLE "doctor" (
  "id" SERIAL PRIMARY KEY,
  "surname" varchar NOT NULL,
  "name" varchar NOT NULL,
  "middle_name" varchar
);

CREATE TABLE "patient" (
  "id" SERIAL PRIMARY KEY,
  "surname" varchar NOT NULL,
  "name" varchar NOT NULL,
  "middle_name" varchar,
  "age" int,
  "diagnosis" varchar
);

CREATE TABLE "medicine_warehouse" (
  "id" SERIAL PRIMARY KEY,
  "total_amount" int NOT NULL,
  "critical_limit" int NOT NULL,
  "medicine_id" int NOT NULL
);

CREATE TABLE "substance_warehouse" (
  "id" SERIAL PRIMARY KEY,
  "total_amount" int NOT NULL,
  "critical_limit" int NOT NULL,
  "substance_id" int NOT NULL
);

CREATE TABLE "production_techonology" (
  "id" SERIAL PRIMARY KEY,
  "method_of_production" varchar NOT NULL,
  "time_to_product" varchar NOT NULL
);

CREATE TABLE "customer" (
  "id" SERIAL PRIMARY KEY,
  "surname" varchar NOT NULL,
  "name" varchar NOT NULL,
  "middle_name" varchar,
  "phone_number" varchar,
  "address" varchar
);

CREATE TABLE "orders" (
  "id" SERIAL PRIMARY KEY,
  "customer_id" int NOT NULL,
  "receipt_id" int NOT NULL,
  "order_date" date NOT NULL,
  "production_date" timestamp NOT NULL,
  "status" order_status NOT NULL
);

CREATE TABLE "medicine_composition" (
  "id" SERIAL PRIMARY KEY,
  "substance_id" int NOT NULL,
  "medicine_id" int NOT NULL,
  "required_quantity" float NOT NULL
);

CREATE TABLE "imported_medicine" (
  "id" SERIAL PRIMARY KEY,
  "medicine_id" int NOT NULL
);

CREATE TABLE "local_medicine" (
  "id" SERIAL PRIMARY KEY,
  "medicine_id" int NOT NULL,
  "production_techology" int NOT NULL
);

CREATE TABLE "medicine_usage_statistics" (
  "id" SERIAL PRIMARY KEY,
  "medicine_id" int NOT NULL,
  "quantity_used" float NOT NULL,
  "usage_time" timestamp NOT NULL
);

CREATE TABLE "substance_usage_statistics" (
  "id" SERIAL PRIMARY KEY,
  "substance_id" int NOT NULL,
  "quantity_used" float NOT NULL,
  "usage_time" timestamp NOT NULL
);

ALTER TABLE "medicine_list" ADD FOREIGN KEY ("receipt_id") REFERENCES "receipt" ("id");

ALTER TABLE "medicine_list" ADD FOREIGN KEY ("medicine_id") REFERENCES "medicine" ("id");

ALTER TABLE "receipt" ADD FOREIGN KEY ("doctor_id") REFERENCES "doctor" ("id");

ALTER TABLE "receipt" ADD FOREIGN KEY ("patient_id") REFERENCES "patient" ("id");

ALTER TABLE "medicine_warehouse" ADD FOREIGN KEY ("medicine_id") REFERENCES "medicine" ("id");

ALTER TABLE "substance_warehouse" ADD FOREIGN KEY ("substance_id") REFERENCES "substance" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("customer_id") REFERENCES "customer" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("receipt_id") REFERENCES "receipt" ("id");

ALTER TABLE "medicine_composition" ADD FOREIGN KEY ("substance_id") REFERENCES "substance" ("id");

ALTER TABLE "medicine_composition" ADD FOREIGN KEY ("medicine_id") REFERENCES "local_medicine" ("id");

ALTER TABLE "imported_medicine" ADD FOREIGN KEY ("medicine_id") REFERENCES "medicine" ("id");

ALTER TABLE "local_medicine" ADD FOREIGN KEY ("medicine_id") REFERENCES "medicine" ("id");

ALTER TABLE "local_medicine" ADD FOREIGN KEY ("production_techology") REFERENCES "production_techonology" ("id");

ALTER TABLE "medicine_usage_statistics" ADD FOREIGN KEY ("medicine_id") REFERENCES "medicine" ("id");

ALTER TABLE "substance_usage_statistics" ADD FOREIGN KEY ("substance_id") REFERENCES "substance" ("id");
