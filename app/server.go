package app

import (
	"log"
	"net/http"

	"github.com/FreeJ1nG/backend-template/util"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/cors"
)

type Server struct {
	config     util.Config
	router     *mux.Router
	db         *pgxpool.Pool
	httpServer *http.Server
}

func NewServer(config util.Config, db *pgxpool.Pool) *Server {
	r := mux.NewRouter().PathPrefix("/v1").Subrouter()
	server := &http.Server{
		Addr:    ":" + config.ServerPort,
		Handler: cors.AllowAll().Handler(r),
	}
	return &Server{
		config:     config,
		router:     r,
		db:         db,
		httpServer: server,
	}
}

func (s *Server) ListenAndServe() (err error) {
	if err = s.httpServer.ListenAndServe(); err != nil {
		log.Fatal("unable to start server: ", err.Error())
		return
	}
	return
}
