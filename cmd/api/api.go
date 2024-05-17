package api

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/xhoang0509/ecom-api/services/product"
	"github.com/xhoang0509/ecom-api/services/user"
	"log"
	"net/http"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subRouter)

	productStore := product.NewStore(s.db)
	productHandler := product.NewHandler(productStore, userStore)
	productHandler.RegisterRoutes(subRouter)

	log.Println("Listening on ", s.addr)

	return http.ListenAndServe(s.addr, router)
}
