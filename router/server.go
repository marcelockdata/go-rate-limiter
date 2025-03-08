package router

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
)

func InitilizeServer(router *chi.Mux) {
	server_port := os.Getenv("HTTP_PORT")
	server_host := "0.0.0.0"
	server_addr := fmt.Sprintf("%s:%s", server_host, server_port)
	srv := &http.Server{
		Handler:      router,
		Addr:         server_addr,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	fmt.Println("Starting the server on port:", server_port)
	log.Fatal(srv.ListenAndServe(), router)
}
