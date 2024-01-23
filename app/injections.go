package app

import (
	"encoding/json"
	"net/http"

	"github.com/FreeJ1nG/backend-template/app/auth"
	"github.com/FreeJ1nG/backend-template/util"
)

func (s *Server) InjectDependencies() {
	s.router.Use(util.LoggerMiddleware)

	// Utils
	authUtil := auth.NewUtil()

	// Repositories
	authRepository := auth.NewRepository(s.db)

	// Services
	authService := auth.NewService(authRepository, authUtil)

	// Route Protector Wrapper
	routeProtector := util.NewRouteProtector(authUtil, authService)

	// Controllers
	authHandler := auth.NewHandler(authService)

	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "ok"})
	}).Methods("GET")

	authRouter := s.router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/sign-in", authHandler.SignInUser)
	authRouter.HandleFunc("/sign-up", authHandler.SignUpUser)
	authRouter.HandleFunc("/refresh-jwt", authHandler.RefreshJwt)
	authRouter.HandleFunc("/me", routeProtector.Wrapper(authHandler.GetCurrentUser))
}
