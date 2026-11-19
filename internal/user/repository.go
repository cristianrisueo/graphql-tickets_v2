package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository gestiona la persistencia de usuarios en PostgreSQL
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository crea una nueva instancia del repositorio de usuarios
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// Create inserta un nuevo usuario y devuelve la entidad creada con su UUID generado
func (r *Repository) Create(ctx context.Context, name, email string) (*User, error) {
	// Crea la consulta SQL para insertar un nuevo usuario
	const query = `
		INSERT INTO users (name, email)
		VALUES ($1, $2)
		RETURNING id, name, email, created_at
	`

	// Ejecuta la consulta y escanea los resultados en una nueva instancia de User
	// Comprueba si hay errores de conflicto de email u otros
	var u User
	err := r.db.QueryRow(ctx, query, name, email).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrEmailConflict
		}

		return nil, fmt.Errorf("error creando usuario: %w", err)
	}

	return &u, nil
}

// GetByID busca un usuario por su identificador único
func (r *Repository) GetByID(ctx context.Context, id string) (*User, error) {
	// Crea la consulta SQL para buscar un usuario por ID
	const query = `
		SELECT id, name, email, created_at
		FROM users
		WHERE id = $1
	`

	// Ejecuta la consulta y escanea los resultados en una nueva instancia de User
	// Maneja el caso en que no se encuentra el usuario y otros errores
	var u User
	err := r.db.QueryRow(ctx, query, id).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("error consultando usuario por id: %w", err)
	}

	return &u, nil
}

// GetByIDs recupera múltiples usuarios en una sola consulta para evitar problemas N+1
func (r *Repository) GetByIDs(ctx context.Context, ids []string) ([]User, error) {
	// Si no hay IDs proporcionados, devuelve un slice vacío
	if len(ids) == 0 {
		return nil, nil
	}

	// Crea la consulta SQL para buscar múltiples usuarios por sus IDs
	const query = `
		SELECT id, name, email, created_at
		FROM users
		WHERE id = ANY($1::uuid[])
	`

	// Ejecuta la consulta y maneja errores
	rows, err := r.db.Query(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("error consultando usuarios en bloque: %w", err)
	}
	defer rows.Close()

	// Itera sobre las filas devueltas y escanea cada usuario en un slice
	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("error escaneando fila de usuario: %w", err)
		}

		users = append(users, u)
	}

	// Comprueba si hubo errores durante la iteración de las filas
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando filas de usuarios: %w", err)
	}

	return users, nil
}

// List devuelve todos los usuarios ordenados por fecha de creación
func (r *Repository) List(ctx context.Context) ([]User, error) {
	// Crea la consulta SQL para listar todos los usuarios
	const query = `
		SELECT id, name, email, created_at
		FROM users
		ORDER BY created_at ASC
	`

	// Ejecuta la consulta y maneja errores
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error listando usuarios: %w", err)
	}
	defer rows.Close()

	// Itera sobre las filas devueltas y escanea cada usuario en un slice
	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("error escaneando fila de usuario: %w", err)
		}

		users = append(users, u)
	}

	// Comprueba si hubo errores durante la iteración de las filas
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando listado de usuarios: %w", err)
	}

	return users, nil
}
