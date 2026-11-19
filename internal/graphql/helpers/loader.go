// loader.go implementa el patrón DataLoader por petición HTTP para agrupar
// consultas de usuarios por lotes y evitar el problema de consultas N+1 en GraphQL.
package helpers

import (
	"context"
	"net/http"
	"time"

	"github.com/vikstrous/dataloadgen"

	"graphql-tickets-v2/internal/user"
)

type ctxKey string

const loadersKey ctxKey = "dataloaders"

// Loaders agrupa los cargadores por lotes disponibles durante una petición HTTP
type Loaders struct {
	Users *dataloadgen.Loader[string, *user.User]
}

// NewLoaders inicializa los DataLoaders con una ventana de agrupación de 2 ms
func NewLoaders(userService *user.Service) *Loaders {
	return &Loaders{
		Users: dataloadgen.NewLoader(func(ctx context.Context, ids []string) ([]*user.User, []error) {
			users, err := userService.GetByIDs(ctx, ids)
			if err != nil {
				return nil, []error{err}
			}

			// Indexa los usuarios recuperados por ID para reconstruir el orden exacto de entrada
			usersByID := make(map[string]*user.User, len(users))
			for i := range users {
				usersByID[users[i].ID] = &users[i]
			}

			// DataLoader exige que el slice devuelto tenga la misma longitud y orden que ids
			results := make([]*user.User, len(ids))
			for i, id := range ids {
				results[i] = usersByID[id]
			}

			return results, nil
		}, dataloadgen.WithWait(2*time.Millisecond)),
	}
}

// LoaderMiddleware crea una instancia limpia de Loaders por cada petición HTTP y la inyecta en el contexto
func LoaderMiddleware(userService *user.Service, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loaders := NewLoaders(userService)
		ctx := context.WithValue(r.Context(), loadersKey, loaders)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// LoadersFor extrae los DataLoaders asociados a la petición actual desde el contexto
func LoadersFor(ctx context.Context) *Loaders {
	loaders, _ := ctx.Value(loadersKey).(*Loaders)
	return loaders
}