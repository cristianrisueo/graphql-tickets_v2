-- Eliminación del índice en tickets por usuario asignado
DROP INDEX IF EXISTS idx_tickets_assignee_id;

-- Eliminación de la tabla de tickets
DROP TABLE IF EXISTS tickets;

-- Eliminación de la tabla de usuarios
DROP TABLE IF EXISTS users;

-- Eliminación del enum de estados de ticket
DROP TYPE IF EXISTS ticket_status;