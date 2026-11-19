package ticket

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository gestiona la persistencia de tickets en PostgreSQL
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository crea una nueva instancia del repositorio de tickets
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// Create inserta un nuevo ticket y devuelve la entidad creada con sus valores generados
func (r *Repository) Create(ctx context.Context, input CreateInput) (*Ticket, error) {
	// Define la consulta SQL para insertar un nuevo ticket
	const query = `
		INSERT INTO tickets (title, description, assignee_id)
		VALUES ($1, $2, $3)
		RETURNING id, title, description, status, assignee_id, created_at, updated_at
	`

	// Ejecuta la consulta y escanea los resultados en una nueva instancia de Ticket
	var t Ticket
	err := r.db.QueryRow(
		ctx,
		query,
		input.Title,
		input.Description,
		input.AssigneeID,
	).Scan(
		&t.ID,
		&t.Title,
		&t.Description,
		&t.Status,
		&t.AssigneeID,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	// Maneja errores de integridad referencial y otros
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return nil, ErrInvalidAssignee
		}

		return nil, fmt.Errorf("error creando ticket: %w", err)
	}

	return &t, nil
}

// GetByID busca un ticket por su identificador único
func (r *Repository) GetByID(ctx context.Context, id string) (*Ticket, error) {
	// Define la consulta SQL para obtener un ticket por su ID
	const query = `
		SELECT id, title, description, status, assignee_id, created_at, updated_at
		FROM tickets
		WHERE id = $1
	`

	// Ejecuta la consulta y escanea los resultados en una nueva instancia de Ticket
	var t Ticket
	err := r.db.QueryRow(ctx, query, id).Scan(
		&t.ID,
		&t.Title,
		&t.Description,
		&t.Status,
		&t.AssigneeID,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	// Maneja errores de no encontrado y otros
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("error consultando ticket por id: %w", err)
	}

	return &t, nil
}

// List devuelve los tickets ordenados por fecha de creación aplicando los filtros opcionales
func (r *Repository) List(ctx context.Context, filter ListFilter) ([]Ticket, error) {
	// Define la consulta SQL para listar tickets con filtros opcionales
	const query = `
		SELECT id, title, description, status, assignee_id, created_at, updated_at
		FROM tickets
		WHERE ($1::ticket_status IS NULL OR status = $1::ticket_status)
		  AND ($2::uuid IS NULL OR assignee_id = $2::uuid)
		ORDER BY created_at ASC
	`

	// Comprueba si se ha proporcionado un estado del ticket para filtrar.
	// Convierte a str si se proporciona.
	var status *string
	if filter.Status != nil {
		s := string(*filter.Status)
		status = &s
	}

	// Ejecuta la consulta y obtiene las filas
	rows, err := r.db.Query(ctx, query, status, filter.AssigneeID)
	if err != nil {
		return nil, fmt.Errorf("error listando tickets: %w", err)
	}
	defer rows.Close()

	// Itera sobre las filas y escanea cada ticket en la lista de resultados
	var tickets []Ticket
	for rows.Next() {
		var t Ticket

		if err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Description,
			&t.Status,
			&t.AssigneeID,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error escaneando fila de ticket: %w", err)
		}

		tickets = append(tickets, t)
	}

	// Comprueba si hubo errores durante la iteración de las filas
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando listado de tickets: %w", err)
	}

	return tickets, nil
}

// Update modifica los campos informados de un ticket existente y actualiza su fecha de modificación
func (r *Repository) Update(ctx context.Context, id string, input UpdateInput) (*Ticket, error) {
	// Define la consulta SQL para actualizar un ticket por su ID con los valores proporcionados 
    // usando COALESCE para mantener los valores existentes si no se proporcionan nuevos
	const query = `
		UPDATE tickets
		SET
			title = COALESCE($2, title),
			description = COALESCE($3, description),
			status = COALESCE($4::ticket_status, status),
			assignee_id = COALESCE($5::uuid, assignee_id),
			updated_at = NOW()
		WHERE id = $1
		RETURNING id, title, description, status, assignee_id, created_at, updated_at
	`

	// Comprueba si se ha proporcionado un estado del ticket para filtrar.
	// Convierte a str si se proporciona.
	var status *string
	if input.Status != nil {
		s := string(*input.Status)
		status = &s
	}

	// Ejecuta la consulta y escanea los resultados en una nueva instancia de Ticket
	var t Ticket
	err := r.db.QueryRow(
		ctx,
		query,
		id,
		input.Title,
		input.Description,
		status,
		input.AssigneeID,
	).Scan(
		&t.ID,
		&t.Title,
		&t.Description,
		&t.Status,
		&t.AssigneeID,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

    // Maneja errores de no encontrado, integridad referencial y otros
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return nil, ErrInvalidAssignee
		}
        
		return nil, fmt.Errorf("error actualizando ticket: %w", err)
	}

	return &t, nil
}

// Delete elimina un ticket por su identificador único
func (r *Repository) Delete(ctx context.Context, id string) error {
	// Define la consulta SQL para eliminar un ticket por su ID
	const query = `
		DELETE FROM tickets
		WHERE id = $1
	`

	// Ejecuta la consulta y obtiene el resultado de la operación
	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error eliminando ticket: %w", err)
	}

	// Comprueba si se eliminó algún registro; si no, devuelve un error de ticket no encontrado
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
