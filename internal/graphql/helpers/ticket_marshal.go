// Proporciona adaptadores de serialización para Ticket entre el dominio interno y el adaptador de GraphQL
package helpers

import (
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"

	"graphql-tickets-v2/internal/ticket"
)

// MarshalTicketStatus serializa ticket.Status de dominio hacia el enum TicketStatus de GraphQL
func MarshalTicketStatus(s ticket.Status) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		fmt.Fprint(w, strconv.Quote(string(s)))
	})
}

// UnmarshalTicketStatus valida y convierte un valor entrante de GraphQL hacia el valor de dominio ticket.Status
func UnmarshalTicketStatus(v any) (ticket.Status, error) {
	str, ok := v.(string)
	if !ok {
		return "", ticket.ErrInvalidStatus
	}

	status := ticket.Status(str)
	if !status.IsValid() {
		return "", ticket.ErrInvalidStatus
	}

	return status, nil
}
