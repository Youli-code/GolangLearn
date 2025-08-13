package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	db := initDB()
	store := NewSQLiteStore(db)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", getTaskHandler(store))
	mux.HandleFunc("POST /tasks", postTaskHandler(store))
	mux.HandleFunc("GET /tasks/{ID}", getTaskByIDHandler(store))
	mux.HandleFunc("PUT /tasks/{ID}", updateTaskByIDHandler(store))
	mux.HandleFunc("DELETE /tasks/{ID}", deleteTaskByIDHandler(store))

	handler := Chain(mux,
		Recover,
		CORSFromEnv(),
		Logging,
	)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		fmt.Println("Server running on http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server crashed: %v\n", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
