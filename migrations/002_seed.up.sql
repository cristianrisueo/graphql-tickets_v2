-- Inserción de usuarios iniciales con UUIDs fijos para pruebas y consultas directas
INSERT INTO users (id, name, email) VALUES
    ('01933f10-0000-7000-8000-000000000001', 'Ana García', 'ana.garcia@example.com'),
    ('01933f10-0000-7000-8000-000000000002', 'Carlos Ruiz', 'carlos.ruiz@example.com'),
    ('01933f10-0000-7000-8000-000000000003', 'Elena Torres', 'elena.torres@example.com');

-- Inserción de tickets de prueba cubriendo distintos estados, asignaciones repetidas y campos nulos
INSERT INTO tickets (id, title, description, status, assignee_id) VALUES
    ('01933f20-0000-7000-8000-000000000001', 'Configurar pipeline CI/CD', 'Añadir pasos de lint y tests automáticos en cada pull request', 'OPEN', '01933f10-0000-7000-8000-000000000001'),
    ('01933f20-0000-7000-8000-000000000002', 'Optimizar consultas de listado', 'Revisar índices y evitar consultas N+1 en el resolver de tickets', 'IN_PROGRESS', '01933f10-0000-7000-8000-000000000001'),
    ('01933f20-0000-7000-8000-000000000003', 'Implementar suscripciones WebSocket', 'Notificar en tiempo real cuando se crea un nuevo ticket', 'IN_PROGRESS', '01933f10-0000-7000-8000-000000000002'),
    ('01933f20-0000-7000-8000-000000000004', 'Documentar variables de entorno', NULL, 'CLOSED', '01933f10-0000-7000-8000-000000000002'),
    ('01933f20-0000-7000-8000-000000000005', 'Investigar fuga de memoria en workers', 'Analizar perfiles pprof bajo carga sostenida', 'OPEN', NULL);