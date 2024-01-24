package app

import (
	"net/http"

	"github.com/FreeJ1nG/backend-template/app/auth"
	"github.com/FreeJ1nG/backend-template/app/dto"
	"github.com/FreeJ1nG/backend-template/util"
)

func (s *Server) InjectDependencies() {
	s.router.Use(util.LoggerMiddleware)
	s.router.Use(util.DefaultOptionsMiddleware)

	// Utils
	authUtil := auth.NewUtil()

	// Repositories
	authRepository := auth.NewRepository(s.db, s.rdb)

	// Services
	authService := auth.NewService(authRepository, authUtil, s.rdb)

	// Route Protector Wrapper
	routeProtector := util.NewRouteProtector(authUtil, authRepository)

	// Controllers
	authHandler := auth.NewHandler(authService)

	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		util.EncodeSuccessResponse(w, dto.SuccessNoPayloadResponse, http.StatusOK)
	}).Methods("GET")

	authRouter := s.router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/sign-in", authHandler.SignInUser)
	authRouter.HandleFunc("/sign-up", authHandler.SignUpUser)
	authRouter.HandleFunc("/refresh-jwt", authHandler.RefreshJwt)
	authRouter.HandleFunc("/me", routeProtector.Wrapper(authHandler.GetCurrentUser))
	authRouter.HandleFunc("/invalidate-jwt", routeProtector.Wrapper(authHandler.InvalidateToken))
}
