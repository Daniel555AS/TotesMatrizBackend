package app

import (
	"time"
	"totesbackend/config"
	"totesbackend/controllers"
	"totesbackend/controllers/utilities"
	"totesbackend/database"
	"totesbackend/repositories"
	routes "totesbackend/router"
	"totesbackend/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // swagger handler
	"gorm.io/gorm"

	_ "totesbackend/docs" // importante para registrar los docs generados
)

var db *gorm.DB
var router *gin.Engine
var authUtil *utilities.AuthorizationUtil
var logUtil *utilities.LogUtil

// @schemes   https

// SetupAndRunApp initializes and configures the entire application server,
// including environment variables, database connection, middleware, route handlers,
// CORS policies, and Swagger documentation.
// It runs the HTTPS server on port 443 using TLS certificates.
// If any initialization step fails, it returns an error.
//
// The function also performs the following steps:
// - Loads environment variables
// - Starts and defers closure of the PostgreSQL connection
// - Applies database migrations
// - Initializes repositories, services, and utilities
// - Registers all API route groups (users, roles, auth, billing, etc.)
// - Enables CORS with specific allowed origins
// - Mounts the Swagger UI at /swagger/index.html
// - Starts the HTTPS server

func SetupAndRunApp() error {

	// load env
	err := config.LoadENV()
	if err != nil {
		return err
	}

	// start database
	err = database.StartPostgres()
	if err != nil {
		return err
	}

	// defer closing database
	defer database.ClosePostgres()

	db = database.GetDB()
	userRepo := repositories.NewUserRepository(db)
	authUtil = utilities.NewAuthorizationUtil(services.NewAuthorizationService(repositories.NewAuthorizationRepository(db), userRepo))
	logUtil = utilities.NewLogUtil(services.NewUserLogService(repositories.NewUserLogRepository(db)))
	router = gin.Default()
	database.MigrateDB() // recordar descomentar para inicializar la base de datos

	// Configurar CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:5503", "http://127.0.0.1:5500", "http://127.0.0.1:5501"}, // Especifica los orígenes permitidos
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Username"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	setUpUserRouter()
	setUpItemTypeRouter()
	setUpItemRouter()
	setUpPermissionRouter()
	setUpRoleRouter()
	setUpUserTypeRouter()
	setUpIdentifierTypeRouter()
	setUpUserStateTypeRouter()
	setUpEmployeeRouter()
	setUpAdditionalExpenseRouter()
	setUpHistoricalItemPriceRouter()
	setUpCommentRouter()
	setUpAuthRouter()
	setUpAppointmentRouter()
	setUpCustomerRouter()
	setUpOrderStateTypeRouter()
	setUpPurchaseOrderRouter()
	setUpDiscountTypeRouter()
	setUpUserCredentialValidationRouter()
	setUpTaxTypeRouter()
	setUpBillingRouter()
	setUpInvoice()
	setUpExternalSaleRouter()
	setUpSalesReportRouter()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err = router.RunTLS(":443", "certs/cert.pem", "certs/key.pem")
	if err != nil {
		panic(err)
	}

	return nil
}

func setUpPermissionRouter() {
	permissionRepo := repositories.NewPermissionRepository(db)
	permissionService := services.NewPermissionService(permissionRepo)
	permissionController := controllers.NewPermissionController(permissionService, authUtil, logUtil)
	routes.RegisterPermissionRoutes(router, permissionController)
}

func setUpEmployeeRouter() {
	employeeRepo := repositories.NewEmployeeRepository(db)
	employeeService := services.NewEmployeeService(employeeRepo)
	employeeController := controllers.NewEmployeeController(employeeService, authUtil, logUtil)
	routes.RegisterEmployeeRoutes(router, employeeController)
}

func setUpRoleRouter() {
	roleRepo := repositories.NewRoleRepository(db)
	roleService := services.NewRoleService(roleRepo)
	roleController := controllers.NewRoleController(roleService, authUtil, logUtil)
	routes.RegisterRoleRoutes(router, roleController)
}

func setUpItemTypeRouter() {
	itemTypeRepo := repositories.NewItemTypeRepository(db)
	itemTypeService := services.NewItemTypeService(itemTypeRepo)
	itemTypeController := controllers.NewItemTypeController(itemTypeService, authUtil, logUtil)
	routes.RegisterItemTypeRoutes(router, itemTypeController)
}

func setUpUserTypeRouter() {
	userTypeRepo := repositories.NewUserTypeRepository(db)
	userTypeService := services.NewUserTypeService(userTypeRepo)
	userTypeController := controllers.NewUserTypeController(userTypeService, authUtil, logUtil)
	routes.RegisterUserTypeRoutes(router, userTypeController)
}

func setUpItemRouter() {
	itemRepo := repositories.NewItemRepository(db)
	itemService := services.NewItemService(itemRepo)
	itemController := controllers.NewItemController(itemService, authUtil, logUtil)
	routes.RegisterItemRoutes(router, itemController)
}

func setUpUserStateTypeRouter() {
	userStateTypeRepo := repositories.NewUserStateTypeRepository(db)
	userStateTypeService := services.NewUserStateTypeService(userStateTypeRepo)
	userStateTypeController := controllers.NewUserStateTypeController(userStateTypeService, authUtil, logUtil)
	routes.RegisterUserStateTypeRoutes(router, userStateTypeController)
}

func setUpIdentifierTypeRouter() {
	identifierTypeRepo := repositories.NewIdentifierTypeRepository(db)
	identifierTypeService := services.NewIdentifierTypeService(identifierTypeRepo)
	identifierTypeController := controllers.NewIdentifierTypeController(identifierTypeService, authUtil, logUtil)
	routes.RegisterIdentifierTypeRoutes(router, identifierTypeController)
}

func setUpUserRouter() {
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService, authUtil, logUtil)
	routes.RegisterUserRoutes(router, userController)
}

func setUpAdditionalExpenseRouter() {
	addRepo := repositories.NewAdditionalExpenseRepository(db)
	addService := services.NewAdditionalExpenseService(addRepo)
	addController := controllers.NewAdditionalExpenseController(addService, authUtil, logUtil)
	routes.RegisterAdditionalExpenseRoutes(router, addController)
}

func setUpHistoricalItemPriceRouter() {
	hisRepo := repositories.NewHistoricalItemPriceRepository(db)
	hisService := services.NewHistoricalItemPriceService(hisRepo)
	hisController := controllers.NewHistoricalItemPriceController(hisService, authUtil, logUtil)
	routes.RegisterHistoricalItemPriceRoutes(router, hisController)
}

func setUpCommentRouter() {
	commentRepo := repositories.NewCommentRepository(db)
	commentService := services.NewCommentService(commentRepo)
	commentController := controllers.NewCommentController(commentService, authUtil, logUtil)
	routes.RegisterCommentRoutes(router, commentController)
}

func setUpAuthRouter() {
	authRepo := repositories.NewAuthorizationRepository(db)
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthorizationService(authRepo, userRepo)
	authController := controllers.NewAuthorizationController(authService, logUtil)
	routes.RegisterAuthorizationRoutes(router, authController)
}

func setUpAppointmentRouter() {
	appointmentRepo := repositories.NewAppointmentRepository(db)
	appointmentService := services.NewAppointmentService(appointmentRepo)
	appointmentController := controllers.NewAppointmentController(appointmentService, authUtil, logUtil)
	routes.RegisterAppointmentRoutes(router, appointmentController)
}

func setUpCustomerRouter() {
	customerRepo := repositories.NewCustomerRepository(db)
	customerService := services.NewCustomerService(customerRepo)
	customerController := controllers.NewCustomerController(customerService, authUtil, logUtil)
	routes.RegisterCustomerRoutes(router, customerController)

}

func setUpOrderStateTypeRouter() {
	orderStateTypeRepo := repositories.NewOrderStateTypeRepository(db)
	orderStateTypeService := services.NewOrderStateTypeService(orderStateTypeRepo)
	orderStateTypeController := controllers.NewOrderStateTypeController(orderStateTypeService, authUtil, logUtil)
	routes.RegisterOrderStateTypeRoutes(router, orderStateTypeController)
}

func setUpPurchaseOrderRouter() {
	purchaseOrderRepo := repositories.NewPurchaseOrderRepository(db)
	itemRepo := repositories.NewItemRepository(db)
	billingRepo := repositories.NewItemRepository(db)
	discountRepo := repositories.NewDiscountTypeRepository(db)
	taxRepo := repositories.NewTaxTypeRepository(db)
	invoiceRepo := repositories.NewInvoiceRepository(db)

	billingService := services.NewBillingService(billingRepo, discountRepo, taxRepo)
	purchaseOrderService := services.NewPurchaseOrderService(purchaseOrderRepo, itemRepo, billingService, invoiceRepo)
	purchaseOrderController := controllers.NewPurchaseOrderController(purchaseOrderService, authUtil, logUtil)

	routes.RegisterPurchaseOrderRoutes(router, purchaseOrderController)
}

func setUpDiscountTypeRouter() {
	discountTypeRepo := repositories.NewDiscountTypeRepository(db)
	discountTypeService := services.NewDiscountTypeService(discountTypeRepo)
	discountTypeController := controllers.NewDiscountTypeController(discountTypeService, authUtil, logUtil)
	routes.RegisterDiscountTypeRoutes(router, discountTypeController)
}

func setUpUserCredentialValidationRouter() {
	userRepository := repositories.NewUserRepository(db)
	userCredentialValidationService := services.NewUserCredentialValidationService(userRepository)
	userCredentialValidationController := controllers.NewUserCredentialValidationController(userCredentialValidationService, authUtil, logUtil)
	routes.RegisterUserCredentialValidationRoutes(router, userCredentialValidationController)
}

func setUpTaxTypeRouter() {
	taxTypeRepo := repositories.NewTaxTypeRepository(db)
	taxTypeService := services.NewTaxTypeService(taxTypeRepo)
	taxTypeController := controllers.NewTaxTypeController(taxTypeService, authUtil, logUtil)
	routes.RegisterTaxTypeRoutes(router, taxTypeController)
}

func setUpBillingRouter() {
	billingRepo := repositories.NewItemRepository(db)
	discountRepo := repositories.NewDiscountTypeRepository(db)
	taxRepo := repositories.NewTaxTypeRepository(db)

	billingService := services.NewBillingService(billingRepo, discountRepo, taxRepo)
	billingController := controllers.NewBillingController(billingService, authUtil)

	routes.RegisterBillingRoutes(router, billingController)
}

func setUpInvoice() {
	invoiceRepo := repositories.NewInvoiceRepository(db)
	itemRepo := repositories.NewItemRepository(db)
	billingRepo := repositories.NewItemRepository(db)
	discountRepo := repositories.NewDiscountTypeRepository(db)
	taxRepo := repositories.NewTaxTypeRepository(db)

	billingService := services.NewBillingService(billingRepo, discountRepo, taxRepo)
	invoiceService := services.NewInvoiceService(invoiceRepo, itemRepo, billingService)
	invoiceController := controllers.NewInvoiceController(invoiceService, authUtil, logUtil)

	routes.RegisterInvoice(router, invoiceController)
}

func setUpExternalSaleRouter() {
	externalSaleRepo := repositories.NewExternalSaleRepository(db)
	customerRepo := repositories.NewCustomerRepository(db)
	externalSaleService := services.NewExternalSaleService(externalSaleRepo, customerRepo)
	externalSaleController := controllers.NewExternalSaleController(externalSaleService, authUtil, logUtil)
	routes.RegisterExternalSaleRoutes(router, externalSaleController)
}

func setUpSalesReportRouter() {
	invoiceRepo := repositories.NewInvoiceRepository(db)
	salesReportService := services.NewSalesReportService(invoiceRepo)
	salesReportController := controllers.NewSalesReportController(salesReportService, authUtil, logUtil)
	routes.RegisterSalesReportRoutes(router, salesReportController)
}
