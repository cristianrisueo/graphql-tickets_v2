package ticket

import (
	"context"
	"strings"
)

// Service encapsula la lógica de negocio y reglas de estado del dominio de tickets
type Service struct {
	repo *Repository
}

// NewService crea una nueva instancia del servicio de tickets
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Create valida y normaliza los datos antes de registrar un nuevo ticket
func (s *Service) Create(ctx context.Context, input CreateInput) (*Ticket, error) {
	// Normaliza el títutlo y comprueba que no sea una cadena vacía
	cleanTitle := strings.TrimSpace(input.Title)
	if cleanTitle == "" {
		return nil, ErrInvalidInput
	}

	cleanInput := CreateInput{
		Title: cleanTitle,
	}

	// Normaliza la descripción si viene proporcionada (es opcional)
	if input.Description != nil {
		desc := strings.TrimSpace(*input.Description)
		cleanInput.Description = &desc
	}

	// comprueba que el ID del usuario asignado no sea una cadena vacía tras eliminar espacios
	if input.AssigneeID != nil {
		assigneeID := strings.TrimSpace(*input.AssigneeID)
		if assigneeID == "" {
			return nil, ErrInvalidInput
		}

		cleanInput.AssigneeID = &assigneeID
	}

	// Llama al repositorio para crear el ticket con los datos validados y normalizados
	return s.repo.Create(ctx, cleanInput)
}

// GetByID recupera un ticket validando previamente que su identificador no esté vacío
func (s *Service) GetByID(ctx context.Context, id string) (*Ticket, error) {
	// Comprueba que el ID proporcionado no sea una cadena vacía después de eliminar espacios
	cleanID := strings.TrimSpace(id)
	if cleanID == "" {
		return nil, ErrInvalidInput
	}

	// Llama al repositorio para obtener el ticket por su ID
	return s.repo.GetByID(ctx, cleanID)
}

// List valida los filtros opcionales y devuelve los tickets coincidentes
func (s *Service) List(ctx context.Context, filter ListFilter) ([]Ticket, error) {
	// Comprueba que si se proporciona un estado del ticket, sea un estado válido
	if filter.Status != nil && !filter.Status.IsValid() {
		return nil, ErrInvalidStatus
	}

	cleanFilter := ListFilter{
		Status: filter.Status,
	}

	// Comprueba que si se proporciona un ID de usuario para filtrar, no sea una cadena vacía
	if filter.AssigneeID != nil {
		assigneeID := strings.TrimSpace(*filter.AssigneeID)
		if assigneeID == "" {
			return nil, ErrInvalidInput
		}

		cleanFilter.AssigneeID = &assigneeID
	}

	// Llama al repositorio para obtener la lista de tickets con los filtros validados y normalizados
	return s.repo.List(ctx, cleanFilter)
}

// Update valida los campos modificados, comprueba la transición de estado y actualiza el ticket
func (s *Service) Update(ctx context.Context, id string, input UpdateInput) (*Ticket, error) {
	// Comrpueba que el ID no sea una cadena vacía después de eliminar espacios en blanco
	cleanID := strings.TrimSpace(id)
	if cleanID == "" {
		return nil, ErrInvalidInput
	}

	// Rechaza peticiones sin ningún campo para actualizar
	if input.Title == nil && input.Description == nil && input.Status == nil && input.AssigneeID == nil {
		return nil, ErrInvalidInput
	}

	// Crea una nueva estructura de entrada para almacenar los campos normalizados
	var cleanInput UpdateInput

	// Normaliza y comprueba los campos opcionales de entrada antes de enviarlos al repositorio
	if input.Title != nil {
		title := strings.TrimSpace(*input.Title)
		if title == "" {
			return nil, ErrInvalidInput
		}

		cleanInput.Title = &title
	}

	if input.Description != nil {
		desc := strings.TrimSpace(*input.Description)
		cleanInput.Description = &desc
	}

	if input.AssigneeID != nil {
		assigneeID := strings.TrimSpace(*input.AssigneeID)
		if assigneeID == "" {
			return nil, ErrInvalidInput
		}

		cleanInput.AssigneeID = &assigneeID
	}

	// Para el estado del ticket valida que sea un estado válido, y luego que pueda transicionar desde el
	// estado actual al nuevo solicitado
	if input.Status != nil {
		if !input.Status.IsValid() {
			return nil, ErrInvalidStatus
		}

		existing, err := s.repo.GetByID(ctx, cleanID)
		if err != nil {
			return nil, err
		}

		if *input.Status != existing.Status && !existing.Status.CanTransitionTo(*input.Status) {
			return nil, ErrInvalidTransition
		}

		cleanInput.Status = input.Status
	}

	// LLama al repositorio para actualizar el ticket con los campos validados y normalizados
	return s.repo.Update(ctx, cleanID, cleanInput)
}

// Delete elimina un ticket validando previamente que su identificador no esté vacío
func (s *Service) Delete(ctx context.Context, id string) error {
	// Comrpueba que el ID no sea una cadena vacía después de eliminar espacios en blanco
	cleanID := strings.TrimSpace(id)
	if cleanID == "" {
		return ErrInvalidInput
	}

	// Llama al repositorio para eliminar el ticket por su ID
	return s.repo.Delete(ctx, cleanID)
}
