package main

import (
	"KeyValueDB/db"
	"KeyValueDB/handlers"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var Database db.IDatabase

func main() {
	ctx := context.Background()

	Database = db.NewDatabase()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	mux := http.ServeMux{}
	mux.HandleFunc("/", handlers.IndexHandler(Database))

	server := http.Server{
		Addr:    ":8080",
		Handler: &mux,
	}

	go func() {
		fmt.Printf("Server is running on port %s\n", server.Addr)
		err := server.ListenAndServe()
		if err != nil {
			fmt.Printf("Server error: %s\n", err)
		}
	}()

	<-exit

	fmt.Println("Shutting down server...")
	cancelCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := server.Shutdown(cancelCtx)
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Server error whilst shutting down: %s\n", err)
		}
	}

	fmt.Println("Server is shut down.")
}
