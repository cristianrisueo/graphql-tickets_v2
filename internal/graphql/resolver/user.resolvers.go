// user.resolvers.go actúa como handler de User conectando las peticiones GraphQL (queries y mutaciones)
// con la capa servicio del dominio de User
package resolver

import (
	"context"

	"graphql-tickets-v2/internal/graphql/generated"
	"graphql-tickets-v2/internal/user"
)

// CreateUser registra un nuevo usuario delegando las validaciones al servicio
func (r *mutationResolver) CreateUser(ctx context.Context, name string, email string) (*user.User, error) {
	return r.UserService.Create(ctx, name, email)
}

// User obtiene un usuario por su identificador único 
func (r *queryResolver) User(ctx context.Context, id string) (*user.User, error) {
	return r.UserService.GetByID(ctx, id)
}

// Users devuelve el listado completo de usuarios registrados
func (r *queryResolver) Users(ctx context.Context) ([]user.User, error) {
	return r.UserService.List(ctx)
}

// Mutation devuelve la implementación de los resolvers de mutación raíz
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query devuelve la implementación de los resolvers de consulta raíz
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
