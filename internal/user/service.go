package user

import (
	"context"
	"strings"
)

// Service encapsula la lógica de negocio y validaciones del dominio de usuarios
type Service struct {
	repo *Repository
}

// NewService crea una nueva instancia del servicio de usuarios
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Create normaliza y valida los datos antes de registrar un nuevo usuario
func (s *Service) Create(ctx context.Context, name, email string) (*User, error) {
	// Normaliza el nombre y el email, eliminando espacios y convirtiendo a minúsculas
	cleanName := strings.TrimSpace(name)
	cleanEmail := strings.ToLower(strings.TrimSpace(email))

	// Valida que el nombre y el email no estén vacíos y que el email tenga un formato básico válido
	if cleanName == "" || cleanEmail == "" || !strings.Contains(cleanEmail, "@") {
		return nil, ErrInvalidInput
	}

	// LLama al repositorio para crear el usuario
	return s.repo.Create(ctx, cleanName, cleanEmail)
}

// GetByID recupera un usuario verificando previamente que el identificador no esté vacío
func (s *Service) GetByID(ctx context.Context, id string) (*User, error) {
	// Normaliza el ID eliminando espacios en blanco y comprueba que no esté vacío
	cleanID := strings.TrimSpace(id)
	if cleanID == "" {
		return nil, ErrInvalidInput
	}

	// Llama al repositorio para obtener el usuario por ID
	return s.repo.GetByID(ctx, cleanID)
}

// GetByIDs recupera un lote de usuarios por sus IDs para el DataLoader
func (s *Service) GetByIDs(ctx context.Context, ids []string) ([]User, error) {
	// Comprueba que la lista de IDs no esté vacía
	if len(ids) == 0 {
		return nil, nil
	}

	// Llama al repositorio para obtener los usuarios por sus IDs
	return s.repo.GetByIDs(ctx, ids)
}

// List devuelve el listado completo de usuarios registrados
func (s *Service) List(ctx context.Context) ([]User, error) {
	return s.repo.List(ctx)
}
