package app

import (
	"encoding/json"
	"net/http"

	"github.com/FreeJ1nG/backend-template/app/auth"
	"github.com/FreeJ1nG/backend-template/app/cms"
	"github.com/FreeJ1nG/backend-template/app/pagination"
	"github.com/FreeJ1nG/backend-template/util"
)

func (s *Server) InjectDependencies() {
	s.router.Use(util.LoggerMiddleware)
	s.router.Use(util.DefaultFormatMiddleware)

	paginator := pagination.NewPaginator(s.db)

	// Utils
	authUtil := auth.NewUtil()
	cmsUtil := cms.NewUtil()

	// Repositories
	authRepository := auth.NewRepository(s.db)
	cmsRepository := cms.NewRepository(s.db, cmsUtil, paginator)

	// Services
	authService := auth.NewService(authRepository, authUtil)
	cmsService := cms.NewService(cmsRepository, cmsUtil)

	// Route Protector Wrapper
	routeProtector := util.NewRouteProtector(authUtil, authService)

	// Controllers
	authHandler := auth.NewHandler(authService)
	cmsHandler := cms.NewHandler(cmsService)

	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "ok"})
	}).Methods("GET")

	authRouter := s.router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/sign-in", authHandler.SignInUser)
	authRouter.HandleFunc("/sign-up", authHandler.SignUpUser)
	authRouter.HandleFunc("/refresh-jwt", authHandler.RefreshJwt)
	authRouter.HandleFunc("/me", routeProtector.Wrapper(authHandler.GetCurrentUser, false))

	cmsRouter := s.router.PathPrefix("/cms").Subrouter()
	cmsRouter.HandleFunc("/{tableName}/info", routeProtector.Wrapper(cmsHandler.GetTableInfo, true)).Methods("GET")
	cmsRouter.HandleFunc("/{tableName}", routeProtector.Wrapper(cmsHandler.GetTableData, true)).Methods("GET")
	cmsRouter.HandleFunc("/{tableName}", routeProtector.Wrapper(cmsHandler.CreateTableData, true)).Methods("POST")
	cmsRouter.HandleFunc("/{tableName}/{pk}", routeProtector.Wrapper(cmsHandler.UpdateTableData, true)).Methods("PATCH")
	cmsRouter.HandleFunc("/{tableName}/{pk}", routeProtector.Wrapper(cmsHandler.DeleteTableData, true)).Methods("DELETE")
}
