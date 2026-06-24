package config

import (
	"time"

	iikoConfig "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/config"
	iikoMenuService "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/menu/service"
	"github.com/scrumno/scrumno-api/infrastructure/integration-system/shared/interfaces"
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action"
	authAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/auth"
	cartAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/cart/cart"
	cartProductAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/cart/product"
	healthAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/health"
	iikoAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/iiko"
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action/menu"
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action/orders"
	queueAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/queue"
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action/user"
	appConfig "github.com/scrumno/scrumno-api/internal/app/entity/app-config"
	getWorkingTime "github.com/scrumno/scrumno-api/internal/app/query/get-working-time"
	checkOntimeCode "github.com/scrumno/scrumno-api/internal/authorize/command/check-ontime-code"
	createAuthorizeCode "github.com/scrumno/scrumno-api/internal/authorize/command/create-authorize-code"
	createAuthorizeTokens "github.com/scrumno/scrumno-api/internal/authorize/command/create-authorize-tokens"
	createUserAuth "github.com/scrumno/scrumno-api/internal/authorize/command/create-user"
	logout "github.com/scrumno/scrumno-api/internal/authorize/command/logout"
	authEntity "github.com/scrumno/scrumno-api/internal/authorize/entity"
	codes "github.com/scrumno/scrumno-api/internal/authorize/entity/codes"
	authorizeTokens "github.com/scrumno/scrumno-api/internal/authorize/entity/tokens"
	findUserByPhone "github.com/scrumno/scrumno-api/internal/authorize/query/find-user-by-phone"
	getRefreshTokensAvailable "github.com/scrumno/scrumno-api/internal/authorize/query/get-refresh-tokens-available"
	getSmsCode "github.com/scrumno/scrumno-api/internal/authorize/query/get-sms-code"
	getSmsCodeSendAvailable "github.com/scrumno/scrumno-api/internal/authorize/query/get-sms-code-send-available"
	createUniqueCode "github.com/scrumno/scrumno-api/internal/authorize/service/create-unique-code"
	"github.com/scrumno/scrumno-api/internal/health/entity/status"
	checkStatusConnectDb "github.com/scrumno/scrumno-api/internal/health/query/check-status-connect-db"
	refreshMenu "github.com/scrumno/scrumno-api/internal/menu/command/refresh-menu"
	saveMenu "github.com/scrumno/scrumno-api/internal/menu/command/save-menu"
	category "github.com/scrumno/scrumno-api/internal/menu/entity/category"
	section "github.com/scrumno/scrumno-api/internal/menu/entity/section"

	saveMenuListener "github.com/scrumno/scrumno-api/internal/menu/listener/save-menu"
	createOrder "github.com/scrumno/scrumno-api/internal/orders/command/create-order"
	createOrderDraft "github.com/scrumno/scrumno-api/internal/orders/command/create-order-draft"
	payOrderDraft "github.com/scrumno/scrumno-api/internal/orders/command/pay-order-draft"
	processProviderWebhook "github.com/scrumno/scrumno-api/internal/orders/command/process-provider-webhook"
	ordersEntity "github.com/scrumno/scrumno-api/internal/orders/entity"
	orderProviderCreatedListener "github.com/scrumno/scrumno-api/internal/orders/listener/order-provider-created"
	orderStatusChangedListener "github.com/scrumno/scrumno-api/internal/orders/listener/order-status-changed"
	ordersService "github.com/scrumno/scrumno-api/internal/orders/service"
	saveModifier "github.com/scrumno/scrumno-api/internal/products/command/save-modifier"
	saveProductCommand "github.com/scrumno/scrumno-api/internal/products/command/save-product"
	modifier "github.com/scrumno/scrumno-api/internal/products/entity/modifier"
	"github.com/scrumno/scrumno-api/internal/products/entity/product"
	saveModifierListener "github.com/scrumno/scrumno-api/internal/products/listener/save-modifier"
	saveProductListener "github.com/scrumno/scrumno-api/internal/products/listener/save-product"
	queueEntity "github.com/scrumno/scrumno-api/internal/queue/entity"
	queueOrderProviderCreatedListener "github.com/scrumno/scrumno-api/internal/queue/listener/order-provider-created"
	queueOrderStatusChangedListener "github.com/scrumno/scrumno-api/internal/queue/listener/order-status-changed"
	getQueue "github.com/scrumno/scrumno-api/internal/queue/query/get-queue"
	queueService "github.com/scrumno/scrumno-api/internal/queue/service"
	updateUserProfile "github.com/scrumno/scrumno-api/internal/users/command/update-user-profile"
	conditionsUpdateProfilePolicy "github.com/scrumno/scrumno-api/internal/users/service/conditions-update-profile"
	"github.com/scrumno/scrumno-api/shared/services/jwt"
	"github.com/scrumno/scrumno-api/shared/services/sms"

	// Cart
	addProduct "github.com/scrumno/scrumno-api/internal/cart/command/add-product-to-cart"
	clearCart "github.com/scrumno/scrumno-api/internal/cart/command/clear-cart"
	createCart "github.com/scrumno/scrumno-api/internal/cart/command/create-cart"
	removeProduct "github.com/scrumno/scrumno-api/internal/cart/command/remove-product"
	updateProduct "github.com/scrumno/scrumno-api/internal/cart/command/update-product-cart"
	cart "github.com/scrumno/scrumno-api/internal/cart/entity"
	getCart "github.com/scrumno/scrumno-api/internal/cart/query/get-cart-by-user-id"
	"github.com/scrumno/scrumno-api/shared/services/snapshot"
	fileStorage "github.com/scrumno/scrumno-api/shared/services/storage"

	// Customer
	// customer "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/customer/service"
	order "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/order-delivery/service"
	getCategories "github.com/scrumno/scrumno-api/internal/menu/query/get-categories"
	getSections "github.com/scrumno/scrumno-api/internal/menu/query/get-sections"
	getProducts "github.com/scrumno/scrumno-api/internal/products/query/get-products"
)

func DI() (*action.Actions, *action.Listeners) {
	cfg := Load()

	em := GetEventManager()

	/* INTEGRATION SYSTEMs */

	var (
		// service
		orderProvider   interfaces.OrderProvider
		orderBuilder    interfaces.OrderBodyBuilder
		snapshotService interfaces.SnapshotService
		snapshotStore   interfaces.SnapshotStore

		//customer
		// cBuilder  interfaces.CustomerBodyBuilder
		// cProvider interfaces.CustomerProvider
		// cSync     interfaces.CustomerSyncService

		menuProvider interfaces.MenuProvider

		// handlers
		getMenuHandler interfaces.GetMenuHandler

		// actions
		refreshMenuAction menu.RefreshMenuAction

		// config
	)
	var iikoCfg *iikoConfig.Config

	switch cfg.IntegrationSystem.IntegrationSystem {

	case "iiko":
		iikoCfg = iikoConfig.Load()

		// services
		menuProvider = iikoMenuService.NewMenuProvider(iikoCfg)
		snapshotStore = fileStorage.NewFileStore(iikoCfg.SnapshotFilePath)
		snapshotService = snapshot.NewSnapshotService(snapshotStore)

		//customer
		// cBuilder = customer.NewCustomerBodyBuilder(iikoCfg)
		// cProvider = customer.NewCustomerProvider(iikoCfg)
		// cSync = customer.NewCustomerSyncService(cBuilder, cProvider)

		//order
		orderBuilder = order.NewOrderBodyBuilder(iikoCfg)
		orderProvider = order.NewOrderProvider(iikoCfg)

		// handlers
		getMenuHandler = refreshMenu.NewHandler(menuProvider, em, snapshotService)

		// actions
		refreshMenuAction = menu.NewRefreshMenuAction(getMenuHandler)
	}

	/* INTEGRATION SYSTEMs END */

	smsService := sms.NewSmsService(sms.Config{
		ApiKey:         cfg.Sms.ApiKey,
		ApiPhoneNumber: cfg.Sms.ApiPhoneNumber,
	})

	// repository
	statusRepo := status.NewStatusRepository(DB)
	registrationRepo := authEntity.NewRegistrationRepository(DB)
	tokensRepo := authorizeTokens.NewTokensRepository(DB)
	codesRepo := codes.NewSmsCodesRepository(DB)
	productRepo := product.NewProductRepository(DB)
	cartRepo := cart.NewCartRepository(DB)
	modifierRepo := modifier.NewModifierRepository(DB)
	sectionRepo := section.NewSectionRepository(DB)
	categoryRepo := category.NewCategoryRepository(DB)
	appConfigRepo := appConfig.NewAppConfigRepository(DB)
	queueRepo := queueEntity.NewQueueRepository(DB)
	orderRepo := ordersEntity.NewOrderRepository(DB)

	jwtManager := jwt.NewManager(jwt.Config{
		AccessSecret:    string(cfg.JWT.SecretKey),
		RefreshSecret:   string(cfg.JWT.SecretKey),
		AccessTokenTtl:  15 * time.Minute,
		RefreshTokenTtl: 7 * 24 * time.Hour,
	})

	// service
	checkStatusFetcher := checkStatusConnectDb.NewFetcher(statusRepo)

	// service (нужен до createAuthorizeCodeHandler)
	createUniqueCodeSvc := createUniqueCode.NewCreateUniqueCodeService()

	// command
	conditionsProfilePolicy := conditionsUpdateProfilePolicy.NewHandler()
	// Регистрация: CustomerSync (iiko) отключён — SyncGet/Sync в create-user могут долго висеть
	// на внешнем API и приводить к невыполненному ответу / timeout на клиенте.
	updateUserProfileHandler := updateUserProfile.NewHandler(registrationRepo, conditionsProfilePolicy, nil)
	logoutHandler := logout.NewHandler(tokensRepo)
	checkOntimeCodeHandler := checkOntimeCode.NewHandler(codesRepo)
	// Регистрация: CustomerSync (iiko) отключён — SyncGet/Sync в create-user могут долго висеть
	// на внешнем API и приводить к невыполненному ответу / timeout на клиенте.
	createUserAuthHandler := createUserAuth.NewHandler(registrationRepo, nil)
	createAuthorizeTokensHandler := createAuthorizeTokens.NewHandler(tokensRepo, jwtManager)
	createAuthorizeCodeHandler := createAuthorizeCode.NewHandler(codesRepo, createUniqueCodeSvc)

	createOrderHandler := createOrder.NewHandler(orderProvider, orderBuilder)
	createOrderDraftHandler := createOrderDraft.NewHandler(orderRepo)
	paymentStubService := ordersService.NewPaymentStubService()
	var commandStatusProvider *order.CommandStatusProvider
	if iikoCfg != nil {
		commandStatusProvider = order.NewCommandStatusProvider(iikoCfg)
	}
	payOrderDraftHandler := payOrderDraft.NewHandler(orderRepo, paymentStubService, createOrderHandler, commandStatusProvider, em)
	processProviderWebhookHandler := processProviderWebhook.NewHandler(orderRepo, em)
	ordersWsHub := ordersService.NewOrdersWebSocketHub()

	saveProductHandler := saveProductCommand.NewHandler(productRepo)

	createCartHandler := createCart.NewHandler(cartRepo)
	clearCartHandler := clearCart.NewHandler(cartRepo)
	addProductHandler := addProduct.NewHandler(cartRepo)
	removeProductHandler := removeProduct.NewHandler(cartRepo)
	updateProductHandler := updateProduct.NewHandler(cartRepo)
	getCartFetcher := getCart.NewFetcher(cartRepo)
	getWorkingTimeFetcher := getWorkingTime.NewFetcher(appConfigRepo)
	getQueueFetcher := getQueue.NewFetcher(queueRepo)

	saveModifierHandler := saveModifier.NewHandler(modifierRepo)
	saveMenuHandler := saveMenu.NewHandler(sectionRepo, categoryRepo)

	queueCalculator := queueService.NewOrdersQueueService(&queueEntity.OrdersQueueConfigTable{
		KitchenParallelSlots:  1,
		QueueGrowthFactor:     0.15,
		OrderReserveMinutes:   2,
		RestaurantOpenAt:      "10:00",
		RestaurantCloseAt:     "22:00",
		EmptyQueueWaitMinMins: 10,
		EmptyQueueWaitMaxMins: 10,
		QueueTimeMinFactor:    0.90,
		QueueTimeMaxFactor:    1.25,
	}, DB)
	queueMapper := queueService.NewQueueOrderMapper(productRepo, modifierRepo)

	getQueueAction := queueAction.NewGetQueueAction(
		getWorkingTimeFetcher,
		getQueueFetcher,
		queueCalculator,
		queueMapper,
		nil,
		getCartFetcher,
	)
	getNearestRangeAction := queueAction.NewGetNearestRangeAction(getQueueAction)
	// query
	getRefreshTokensFetcher := getRefreshTokensAvailable.NewFetcher(tokensRepo, jwtManager)
	findUserByPhoneFetcher := findUserByPhone.NewFetcher(registrationRepo)
	getSmsCodeSendAvailableFetcher := getSmsCodeSendAvailable.NewFetcher(codesRepo)
	getSmsCodeFetcher := getSmsCode.NewFetcher(smsService)
	createOrderDraftAction := orders.NewCreateOrderDraftAction(createOrderDraftHandler, findUserByPhoneFetcher, getCartFetcher)
	payOrderDraftAction := orders.NewPayOrderDraftAction(payOrderDraftHandler)
	ordersWebSocketAction := orders.NewOrdersWebSocketAction(ordersWsHub, orderRepo)
	iikoOrderWebhookAction := iikoAction.NewOrderWebhookAction(processProviderWebhookHandler)

	getCategoriesFetcher := getCategories.NewFetcher(categoryRepo)
	getSectionsFetcher := getSections.NewFetcher(sectionRepo)
	getProductsFetcher := getProducts.NewFetcher(productRepo)

	// listeners
	saveProductListener := saveProductListener.NewListener(saveProductHandler)
	saveModifierListener := saveModifierListener.NewListener(saveModifierHandler)
	saveMenuListener := saveMenuListener.NewListener(saveMenuHandler)
	orderProviderCreatedOrderListener := orderProviderCreatedListener.NewListener(orderRepo, cartRepo)
	orderStatusChangedOrderListener := orderStatusChangedListener.NewListener(orderRepo, ordersWsHub)
	orderProviderCreatedQueueListener := queueOrderProviderCreatedListener.NewListener(queueRepo)
	orderStatusChangedQueueListener := queueOrderStatusChangedListener.NewListener(queueRepo)

	return &action.Actions{
			CheckStatusConnectDB: healthAction.NewCheckStatusConnectDBAction(checkStatusFetcher),

			// users
			UpdateUserProfile: user.NewUpdateUserProfileAction(updateUserProfileHandler),

			// auth
			Registration:  authAction.NewRegistrationAction(findUserByPhoneFetcher, checkOntimeCodeHandler, createUserAuthHandler, createAuthorizeTokensHandler),
			Authorize:     authAction.NewAuthorizeAction(findUserByPhoneFetcher, checkOntimeCodeHandler, createAuthorizeTokensHandler),
			Logout:        authAction.NewLogoutAction(logoutHandler, findUserByPhoneFetcher),
			RefreshTokens: authAction.NewRefreshTokensAction(getRefreshTokensFetcher, findUserByPhoneFetcher, createAuthorizeTokensHandler),

			JWTManager: jwtManager,
			SmsCode:    authAction.NewAuthCodeAction(getSmsCodeSendAvailableFetcher, getSmsCodeFetcher, createAuthorizeCodeHandler, findUserByPhoneFetcher),

			// orders
			CreateOrder:      orders.NewCreateOrderAction(createOrderHandler, findUserByPhoneFetcher, getCartFetcher),
			CreateOrderDraft: createOrderDraftAction,
			PayOrderDraft:    payOrderDraftAction,
			OrdersWebSocket:  ordersWebSocketAction,

			// cart
			CreateCart: cartAction.NewCreateAction(createCartHandler),
			ClearCart:  cartAction.NewClearAction(clearCartHandler),
			GetCart:    cartAction.NewGetCartAction(getCartFetcher),

			AddProductToCart:      cartProductAction.NewAddProductAction(addProductHandler),
			RemoveProductFromCart: cartProductAction.NewRemoveProductAction(removeProductHandler),
			UpdateProductFromCart: cartProductAction.NewUpdateAction(updateProductHandler),

			// общие экшены для всех интеграционных систем
			RefreshMenu:      &refreshMenuAction,
			IikoOrderWebhook: iikoOrderWebhookAction,
			GetMenu:          menu.NewGetMenuAction(getCategoriesFetcher, getSectionsFetcher, getProductsFetcher),

			// queue
			GetNearestRange: getNearestRangeAction,
		},
		&action.Listeners{
			SaveProduct:               saveProductListener,
			SaveModifier:              saveModifierListener,
			SaveMenu:                  saveMenuListener,
			OrderProviderCreated:      orderProviderCreatedOrderListener,
			OrderStatusChanged:        orderStatusChangedOrderListener,
			QueueOrderProviderCreated: orderProviderCreatedQueueListener,
			QueueOrderStatusChanged:   orderStatusChangedQueueListener,
		}
}
