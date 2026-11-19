// user.resolvers.go actúa como handler de Ticket conectando las peticiones GraphQL (queries y mutaciones)
// con la capa servicio del dominio de Ticket
package resolver

import (
	"context"

	"graphql-tickets-v2/internal/graphql/generated"
	"graphql-tickets-v2/internal/graphql/helpers"
	"graphql-tickets-v2/internal/ticket"
	"graphql-tickets-v2/internal/user"
)

// CreateTicket crea un nuevo ticket en estado inicial OPEN
func (r *mutationResolver) CreateTicket(ctx context.Context, input ticket.CreateInput) (*ticket.Ticket, error) {
	return r.TicketService.Create(ctx, input)
}

// UpdateTicket aplica actualizaciones parciales o transiciones de estado sobre un ticket
func (r *mutationResolver) UpdateTicket(ctx context.Context, id string, input ticket.UpdateInput) (*ticket.Ticket, error) {
	return r.TicketService.Update(ctx, id, input)
}

// DeleteTicket elimina un ticket por su ID y devuelve true si la operación tuvo éxito
func (r *mutationResolver) DeleteTicket(ctx context.Context, id string) (bool, error) {
	if err := r.TicketService.Delete(ctx, id); err != nil {
		return false, err
	}
	return true, nil
}

// Ticket recupera un ticket concreto por su identificador único
func (r *queryResolver) Ticket(ctx context.Context, id string) (*ticket.Ticket, error) {
	return r.TicketService.GetByID(ctx, id)
}

// Tickets lista los tickets aplicando los filtros opcionales recibidos
func (r *queryResolver) Tickets(ctx context.Context, filter *ticket.ListFilter) ([]ticket.Ticket, error) {
	var f ticket.ListFilter
	if filter != nil {
		f = *filter
	}
	return r.TicketService.List(ctx, f)
}

// Assignee resuelve el usuario asignado al ticket agrupando las consultas mediante DataLoader
func (r *ticketResolver) Assignee(ctx context.Context, obj *ticket.Ticket) (*user.User, error) {
	// Comprueba que el ID del usuario asignado no sea nulo ni esté vacío
	if obj.AssigneeID == nil || *obj.AssigneeID == "" {
		return nil, nil
	}

	// Recupera el DataLoader del contexto para cargar en lote y evitar el problema N+1
	if loaders := helpers.LoadersFor(ctx); loaders != nil {
		return loaders.Users.Load(ctx, *obj.AssigneeID)
	}

	// Respaldo directo: obtiene el usuario por su ID si no hay DataLoader en el contexto
	return r.UserService.GetByID(ctx, *obj.AssigneeID)
}

// Ticket devuelve la implementación de los resolvers de campos sobre el tipo Ticket
func (r *Resolver) Ticket() generated.TicketResolver { return &ticketResolver{r} }

type ticketResolver struct{ *Resolver }
