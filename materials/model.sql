-- Создание таблиц
CREATE TABLE segments (
    segment_id SERIAL PRIMARY KEY,
    segment_name VARCHAR(100) UNIQUE
);

CREATE TABLE users (
    user_id SERIAL PRIMARY KEY
);

CREATE TABLE user_segments (
    user_id BIGINT NOT NULL,
    segment_id BIGINT NOT NULL,
    PRIMARY KEY (user_id, segment_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (segment_id) REFERENCES segments(segment_id)
);

-- Создание триггера для проверки валидности юзера и сегмента
CREATE OR REPLACE FUNCTION validate_user_segment()
RETURNS TRIGGER AS $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM users WHERE user_id = NEW.user_id) THEN
        RAISE EXCEPTION 'user_id % does not exist in users table', NEW.user_id;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM segments WHERE segment_id = NEW.segment_id) THEN
        RAISE EXCEPTION 'segment_id % does not exist in segments table', NEW.segment_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER check_user_segment
BEFORE INSERT ON user_segments
FOR EACH ROW
EXECUTE FUNCTION validate_user_segment();

-- Создание таблицы аудита
CREATE TABLE segment_audit
(  user_id BIGINT NOT NULL,
   segment_id BIGINT NOT NULL,
   operation VARCHAR DEFAULT 'ADD' CONSTRAINT ch_type_event CHECK (operation IN('ADD', 'DELETE')),
   created timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
   FOREIGN KEY (user_id) REFERENCES users(user_id),
   FOREIGN KEY (segment_id) REFERENCES segments(segment_id)
);

-- Триггер для автоматизации добавления записей в аудит
CREATE FUNCTION fnc_trg_segment_audit() RETURNS TRIGGER AS $$
	BEGIN
		IF (TG_OP = 'INSERT') THEN
			INSERT INTO segment_audit(user_id, segment_id, operation) VALUES (NEW.user_id, NEW.segment_id, 'ADD');
		ELSIF (TG_OP = 'DELETE') THEN 
			INSERT INTO segment_audit(user_id, segment_id, operation) VALUES (OLD.user_id, OLD.segment_id, 'DELETE');
		END IF;
		RETURN NULL;
	END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_segment_audit AFTER INSERT OR DELETE ON user_segments
    FOR EACH ROW EXECUTE FUNCTION fnc_trg_segment_audit();

-- Изменение таблицы user_segments с добавлением дедлайна сегмента у пользователя
ALTER TABLE user_segments
ADD time_to_limit timestamp with time zone DEFAULT NULL;
