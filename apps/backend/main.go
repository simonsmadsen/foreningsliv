package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"foreningsliv/backend/db"
	"foreningsliv/backend/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Connect to PostgreSQL and run migrations
	if err := db.Setup(); err != nil {
		log.Fatalf("Database setup failed: %v", err)
	}
	defer db.Close()

	staticDir := os.Getenv("STATIC_DIR")

	mux := http.NewServeMux()

	// GraphQL server
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{},
	}))

	// GraphQL endpoint with CORS
	mux.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		srv.ServeHTTP(w, r)
	})

	// GraphQL Playground at /playground
	mux.Handle("/playground", playground.Handler("GraphQL Playground", "/graphql"))

	// Serve static files if STATIC_DIR is set
	if staticDir != "" {
		fs := http.FileServer(http.Dir(staticDir))
		mux.Handle("/", fs)
		fmt.Printf("Serving static files from %s\n", staticDir)
	}

	fmt.Printf("GraphQL endpoint: http://localhost:%s/graphql\n", port)
	fmt.Printf("Playground:       http://localhost:%s/playground\n", port)
	fmt.Printf("Backend listening on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
