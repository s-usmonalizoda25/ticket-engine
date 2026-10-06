CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

-- Мероприятия концерты и прочее
CREATE TABLE events (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    total_tickets INT NOT NULL CHECK (total_tickets >= 0),
    tickets_left INT NOT NULL CHECK (tickets_left >= 0),
    price NUMERIC(10, 2) NOT NULL CHECK (price >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);


CREATE TABLE tickets (
    id BIGSERIAL PRIMARY KEY,
    event_id BIGINT NOT NULL REFERENCES events(id) ON DELETE RESTRICT,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    status VARCHAR(50) NOT NULL DEFAULT 'booked',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

-- изменения билетов логи
CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,
    ticket_id BIGINT NOT NULL,
    old_status VARCHAR(50),
    new_status VARCHAR(50),
    changed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);


CREATE INDEX idx_tickets_user_status ON tickets(user_id, status) WHERE deleted_at IS NULL;

CREATE INDEX idx_tickets_booked_status ON tickets(status) WHERE status = 'booked' AND deleted_at IS NULL;

-- для быстроты нахождения активных мероприятий
CREATE INDEX idx_events_active ON events(id) WHERE deleted_at IS NULL;

-- триггер и фунция для записи логов билетов
CREATE OR REPLACE FUNCTION log_ticket_status_change()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.status IS DISTINCT FROM NEW.status THEN
        INSERT INTO audit_logs (ticket_id, old_status, new_status)
        VALUES (OLD.id, OLD.status, NEW.status);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_ticket_status_audit
AFTER UPDATE ON tickets
FOR EACH ROW
EXECUTE FUNCTION log_ticket_status_change();