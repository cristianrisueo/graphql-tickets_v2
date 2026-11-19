package ticket

import (
	"errors"
	"time"
)

// Errores del dominio de tickets
var (
	ErrNotFound          = errors.New("ticket no encontrado")
	ErrInvalidInput      = errors.New("datos de entrada inválidos")
	ErrInvalidStatus     = errors.New("estado de ticket no válido")
	ErrInvalidTransition = errors.New("transición de estado no permitida")
	ErrInvalidAssignee   = errors.New("el usuario asignado no existe")
)

// Status representa los posibles estados dentro del ciclo de vida de un ticket
type Status string

const (
	StatusOpen       Status = "OPEN"
	StatusInProgress Status = "IN_PROGRESS"
	StatusClosed     Status = "CLOSED"
)

// IsValid comprueba si el estado pertenece al conjunto de estados válidos. Con
// GraphQL es innecesario para validar la entrada, pero blinda frente a otras fuentes
func (s Status) IsValid() bool {
	switch s {
	case StatusOpen, StatusInProgress, StatusClosed:
		return true
	default:
		return false
	}
}

// CanTransitionTo comprueba si el cambio de estado está permitido
func (s Status) CanTransitionTo(next Status) bool {
	// Primero comprueba si el estado de destino es válido
	if !next.IsValid() {
		return false
	}

	// Segundo comprueba si el estado de destino es el mismo que el actual
	if s == next {
		return false
	}

	// Por último, define la máquina de estados permitidos según el estado actual
	// OPEN -> IN_PROGRESS o CLOSED, IN_PROGRESS -> CLOSED
	switch s {
	case StatusOpen:
		return next == StatusInProgress || next == StatusClosed
	case StatusInProgress:
		return next == StatusClosed
	default:
		return false
	}
}

// Ticket representa un ticket dentro del dominio
type Ticket struct {
	ID          string
	Title       string
	Description *string
	Status      Status
	AssigneeID  *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CreateInput define los datos necesarios para crear un nuevo ticket
type CreateInput struct {
	Title       string
	Description *string
	AssigneeID  *string
}

// UpdateInput define los datos opcionales para actualizar un ticket existente
type UpdateInput struct {
	Title       *string
	Description *string
	Status      *Status
	AssigneeID  *string
}

// ListFilter define los filtros opcionales para consultar tickets
type ListFilter struct {
	Status     *Status
	AssigneeID *string
}
