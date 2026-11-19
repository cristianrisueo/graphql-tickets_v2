// main.go es el punto de entrada de la API GraphQL: inicializa la conexión con PostgreSQL,
// inyecta las dependencias de los dominios y levanta el servidor HTTP con GraphQL Playground.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jackc/pgx/v5/pgxpool"

	"graphql-tickets-v2/internal/graphql/generated"
	"graphql-tickets-v2/internal/graphql/helpers"
	"graphql-tickets-v2/internal/graphql/resolver"
	"graphql-tickets-v2/internal/ticket"
	"graphql-tickets-v2/internal/user"
)

func main() {
	// Carga la configuración desde variables de entorno con valores por defecto
	port := getEnv("PORT", "8080")
	databaseURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/tickets?sslmode=disable")

	// Crea un contexto ligado a las señales del sistema (Ctrl+C / SIGTERM) para apagado limpio
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Inicializa el pool de conexiones a PostgreSQL
	pool, err := newDBPool(ctx, databaseURL)
	if err != nil {
		log.Fatalf("error inicializando la base de datos: %v", err)
	}
	defer pool.Close()

	// Instancia los repositorios y servicios de cada dominio
	userRepo := user.NewRepository(pool)
	userService := user.NewService(userRepo)

	ticketRepo := ticket.NewRepository(pool)
	ticketService := ticket.NewService(ticketRepo)

	// Configura el servidor de gqlgen inyectando los servicios en el Resolver raíz
	gqlServer := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: &resolver.Resolver{
					UserService:   userService,
					TicketService: ticketService,
				},
			},
		),
	)

	// Registra las rutas HTTP: Playground en "/" y el endpoint GraphQL envuelto con el DataLoader en "/query"
	mux := http.NewServeMux()
	mux.Handle("/", playground.Handler("GraphQL Tickets Playground", "/query"))
	mux.Handle("/query", helpers.LoaderMiddleware(userService, gqlServer))

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Arranca el servidor HTTP en segundo plano para escuchar señales de cierre en el hilo principal
	go func() {
		log.Printf("Servidor GraphQL escuchando en http://localhost:%s (Playground en /)", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error inesperado en el servidor HTTP: %v", err)
		}
	}()

	// Espera la señal de interrupción y da un margen de 10 segundos para cerrar peticiones activas
	<-ctx.Done()
	log.Println("Apagando servidor...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("error durante el apagado del servidor: %v", err)
	}
}

// newDBPool crea un pool de conexiones pgxpool y verifica la conectividad real mediante Ping
func newDBPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	// Realiza la conexión a la BD
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("creando el pool de conexiones: %w", err)
	}

	// Realiza un ping a la conexión y en caso de error la cierra
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("verificando conexión con PostgreSQL: %w", err)
	}

	return pool, nil
}

// getEnv devuelve el valor de una variable de entorno o un valor alternativo si no está definida
func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
