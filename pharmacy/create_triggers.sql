CREATE OR REPLACE FUNCTION decrease_medicine_and_substance_stock() RETURNS TRIGGER AS $$
DECLARE
    rec RECORD;
    local_medicine_id INT;
BEGIN
    SELECT id INTO local_medicine_id FROM local_medicine WHERE medicine_id = NEW.medicine_id;

    -- Уменьшение количества медикаментов на складе
    UPDATE medicine_warehouse
    SET total_amount = total_amount - NEW.quantity_used
    WHERE medicine_id = NEW.medicine_id;

    -- Проверка критического уровня медикаментов
    IF (SELECT total_amount FROM medicine_warehouse WHERE medicine_id = NEW.medicine_id) < (SELECT critical_limit FROM medicine_warehouse WHERE medicine_id = NEW.medicine_id) THEN
        RAISE NOTICE 'Critical limit reached for medicine_id %', NEW.medicine_id;
    END IF;

    -- Уменьшение количества ингредиентов на складе на основе состава медикамента
    FOR rec IN
        SELECT substance_id, required_quantity
        FROM medicine_composition
        WHERE medicine_id = local_medicine_id
        LOOP
            UPDATE substance_warehouse
            SET total_amount = total_amount - (rec.required_quantity * NEW.quantity_used)
            WHERE substance_id = rec.substance_id;

            -- Проверка критического уровня ингредиентов
            IF (SELECT total_amount FROM substance_warehouse WHERE substance_id = rec.substance_id) < (SELECT critical_limit FROM substance_warehouse WHERE substance_id = rec.substance_id) THEN
                RAISE NOTICE 'Critical limit reached for substance_id %', rec.substance_id;
            END IF;

            -- Логирование использования ингредиентов
            INSERT INTO substance_usage_statistics (substance_id, quantity_used, usage_time)
            VALUES (rec.substance_id, rec.required_quantity * NEW.quantity_used, CURRENT_TIMESTAMP);
        END LOOP;

    -- Логирование использования медикаментов
    INSERT INTO medicine_usage_statistics (medicine_id, quantity_used, usage_time)
    VALUES (NEW.medicine_id, NEW.quantity_used, CURRENT_TIMESTAMP);

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION increase_medicine_and_substance_stock() RETURNS TRIGGER AS $$
DECLARE
    rec RECORD;
    local_medicine_id INT;
BEGIN
    -- Найти соответствующий local_medicine для medicine_id
    SELECT id INTO local_medicine_id FROM local_medicine WHERE medicine_id = OLD.medicine_id;

    -- Возврат количества медикаментов на склад
    UPDATE medicine_warehouse
    SET total_amount = total_amount + OLD.quantity_used
    WHERE medicine_id = OLD.medicine_id;

    -- Логирование возврата медикаментов
    INSERT INTO medicine_usage_statistics (medicine_id, quantity_used, usage_time)
    VALUES (OLD.medicine_id, -OLD.quantity_used, CURRENT_TIMESTAMP);

    -- Возврат количества ингредиентов на склад на основе состава медикамента
    FOR rec IN
        SELECT substance_id, required_quantity
        FROM medicine_composition
        WHERE medicine_id = local_medicine_id
        LOOP
            UPDATE substance_warehouse
            SET total_amount = total_amount + (rec.required_quantity * OLD.quantity_used)
            WHERE substance_id = rec.substance_id;

            -- Логирование возврата ингредиентов
            INSERT INTO substance_usage_statistics (substance_id, quantity_used, usage_time)
            VALUES (rec.substance_id, -(rec.required_quantity * OLD.quantity_used), CURRENT_TIMESTAMP);
        END LOOP;

    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_decrease_medicine_and_substance_stock
    AFTER INSERT ON medicine_list
    FOR EACH ROW
EXECUTE FUNCTION decrease_medicine_and_substance_stock();


CREATE TRIGGER trg_increase_medicine_and_substance_stock
    BEFORE DELETE ON medicine_list
    FOR EACH ROW
EXECUTE FUNCTION increase_medicine_and_substance_stock();

