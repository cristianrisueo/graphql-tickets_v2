package user

import (
	"errors"
	"time"
)

// Errores del dominio de usuarios
var (
	ErrNotFound      = errors.New("usuario no encontrado")
	ErrEmailConflict = errors.New("el email ya está registrado")
	ErrInvalidInput  = errors.New("datos de entrada inválidos")
)

// User representa un usuario dentro del dominio
type User struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
}
