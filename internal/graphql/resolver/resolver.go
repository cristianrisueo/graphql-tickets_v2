package resolver

import (
	"graphql-tickets-v2/internal/ticket"
	"graphql-tickets-v2/internal/user"
)

// Resolver actúa como punto de inyección de dependencias para queries, mutaciones y campos
type Resolver struct {
	UserService   *user.Service
	TicketService *ticket.Service
}
