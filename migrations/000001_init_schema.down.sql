DROP TRIGGER IF EXISTS trg_ticket_status_audit ON tickets;
DROP FUNCTION IF EXISTS log_ticket_status_change();
DROP INDEX IF EXISTS idx_tickets_booked_status;
DROP INDEX IF EXISTS idx_tickets_user_status;
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS tickets;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS users;