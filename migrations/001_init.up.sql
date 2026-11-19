-- Enum para los estados válidos de un ticket
CREATE TYPE ticket_status AS ENUM (
    'OPEN',
    'IN_PROGRESS',
    'CLOSED'
);

-- Tabla de usuarios del sistema: ID, nombre, email y fecha de creación
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Tabla de tickets: ID, título, descripción, estado, ID del usuario asignado y fechas
CREATE TABLE tickets (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    title TEXT NOT NULL,
    description TEXT,
    status ticket_status NOT NULL DEFAULT 'OPEN',
    assignee_id UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Índice en tickets por el ID del usuario asignado
CREATE INDEX idx_tickets_assignee_id ON tickets(assignee_id);