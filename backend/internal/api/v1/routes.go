package v1

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/scrumno/scrumno-api/config"
	"github.com/scrumno/scrumno-api/internal/api/v1/collector"
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action"
	"github.com/scrumno/scrumno-api/internal/api/v1/middleware"
)

// SetupRouter создаёт маршруты
func SetupRouter(cfg *config.Config, actions *action.Actions) *mux.Router {
	router := mux.NewRouter()

	router.Use(middleware.Logging)
	router.Use(middleware.CORS)
	router.PathPrefix("/").Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	api := router.PathPrefix("/api/v1").Subrouter()

	healthPrefix := "/health"
	health := api.PathPrefix(healthPrefix).Subrouter()

	collectorRoutes := collector.NewEndpointCollector()
	if actions.JWTManager != nil {
		health.Use(middleware.NewAuthMiddleware(actions.JWTManager).Authenticator)
	}

	collectorRoutes.HandleFuncWithPostman(
		health,
		healthPrefix,
		actions.CheckStatusConnectDB.Action,
		actions.CheckStatusConnectDB.GetInputType(),
		"GET",
		"/check-status-connect-db",
	)

	collectorRoutes.HandleFuncWithPostman(
		health,
		healthPrefix,
		actions.CheckStatusConnectDB.Action,
		actions.CheckStatusConnectDB.GetInputType(),
		"POST",
		"/check-1221",
	)

	userPrefix := "/users"
	user := api.PathPrefix(userPrefix).Subrouter()

	if actions.JWTManager != nil {
		user.Use(middleware.NewAuthMiddleware(actions.JWTManager).Authenticator)
	}

	authPrefix := "/auth"
	auth := api.PathPrefix(authPrefix).Subrouter()
	collectorRoutes.HandleFuncWithPostman(
		auth,
		authPrefix,
		actions.Registration.Action,
		actions.Registration.GetInputType(),
		"POST",
		"/registration",
	)

	collectorRoutes.HandleFuncWithPostman(
		auth,
		authPrefix,
		actions.Authorize.Action,
		actions.Authorize.GetInputType(),
		"POST",
		"/authorize",
	)

	collectorRoutes.HandleFuncWithPostman(
		auth,
		authPrefix,
		actions.Logout.Action,
		actions.Logout.GetInputType(),
		"POST",
		"/logout",
	)

	collectorRoutes.HandleFuncWithPostman(
		auth,
		authPrefix,
		actions.SmsCode.Action,
		actions.SmsCode.GetInputType(),
		"POST",
		"/sms-code",
	)

	collectorRoutes.HandleFuncWithPostman(
		auth,
		authPrefix,
		actions.RefreshTokens.Action,
		actions.RefreshTokens.GetInputType(),
		"POST",
		"/refresh-tokens",
	)

	collectorRoutes.HandleFuncWithPostman(
		user,
		userPrefix,
		actions.UpdateUserProfile.Action,
		actions.UpdateUserProfile.GetInputType(),
		"PUT",
		"/update-user-profile",
	)

	cartPrefix := "/cart"
	cartRouter := api.PathPrefix(cartPrefix).Subrouter()
	if actions.JWTManager != nil {
		cartRouter.Use(middleware.NewAuthMiddleware(actions.JWTManager).Authenticator)
	}

	collectorRoutes.HandleFuncWithPostman(
		cartRouter,
		cartPrefix,
		actions.CreateCart.Action,
		actions.CreateCart.GetInputType(),
		"POST",
		"/create",
	)

	collectorRoutes.HandleFuncWithPostman(
		cartRouter,
		cartPrefix,
		actions.AddProductToCart.Action,
		actions.AddProductToCart.GetInputType(),
		"POST",
		"/add-product",
	)

	collectorRoutes.HandleFuncWithPostman(
		cartRouter,
		cartPrefix,
		actions.RemoveProductFromCart.Action,
		actions.RemoveProductFromCart.GetInputType(),
		"POST",
		"/remove-product",
	)

	collectorRoutes.HandleFuncWithPostman(
		cartRouter,
		cartPrefix,
		actions.UpdateProductFromCart.Action,
		actions.UpdateProductFromCart.GetInputType(),
		"PUT",
		"/update-product",
	)

	collectorRoutes.HandleFuncWithPostman(
		cartRouter,
		cartPrefix,
		actions.ClearCart.Action,
		actions.ClearCart.GetInputType(),
		"POST",
		"/clear-cart",
	)

	collectorRoutes.HandleFuncWithPostman(
		cartRouter,
		cartPrefix,
		actions.GetCart.Action,
		actions.GetCart.GetInputType(),
		"GET",
		"",
	)

	menuPrefix := "/menu"
	menu := api.PathPrefix(menuPrefix).Subrouter()
	collectorRoutes.HandleFuncWithPostman(
		menu,
		menuPrefix,
		actions.GetMenu.Action,
		actions.GetMenu.GetInputType(),
		"GET",
		"/get-menu",
	)

	// INTEGRATION SYSTEMs
	ordersPrefix := "/orders"
	orders := api.PathPrefix(ordersPrefix).Subrouter()
	collectorRoutes.HandleFuncWithPostman(
		orders,
		ordersPrefix,
		actions.CreateOrder.Action,
		actions.CreateOrder.GetInputType(),
		"POST",
		"/create-order",
	)

	// iiko integration endpoints
	iikoPrefix := "/iiko"
	iiko := api.PathPrefix(iikoPrefix).Subrouter()

	if actions.RefreshMenu != nil {
		collectorRoutes.HandleFuncWithPostman(
			iiko,
			iikoPrefix,
			actions.RefreshMenu.Action,
			actions.RefreshMenu.GetInputType(),
			"POST",
			"/refresh-menu",
		)
	}
	// INTEGRATION SYSTEMs END

	err := collectorRoutes.GeneratePostmanCollections()
	if err != nil {
		log.Printf("Ошибка генерации Postman: %v", err)
	}

	return router
}
