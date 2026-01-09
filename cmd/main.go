package main

import (
	"main/internal/cache"
	auth_controller "main/internal/controllers/http/v1/auth"
	cart_controller "main/internal/controllers/http/v1/cart"
	category_controller "main/internal/controllers/http/v1/category"
	order_controller "main/internal/controllers/http/v1/order"
	order_status_controller "main/internal/controllers/http/v1/order_status"
	param_controller "main/internal/controllers/http/v1/param"
	param_value_controller "main/internal/controllers/http/v1/param_value"
	payment_controller "main/internal/controllers/http/v1/payment"
	product_controller "main/internal/controllers/http/v1/product"
	role_controller "main/internal/controllers/http/v1/role"
	user_controller "main/internal/controllers/http/v1/user"
	wishlist_controller "main/internal/controllers/http/v1/wishlist"

	auth_middleware "main/internal/middleware/auth"

	"main/internal/pkg/config"
	"main/internal/pkg/postgres"

	"main/internal/repository/postgres/auth"
	"main/internal/repository/postgres/cart"
	"main/internal/repository/postgres/category"
	"main/internal/repository/postgres/order"
	"main/internal/repository/postgres/order_status"
	"main/internal/repository/postgres/param"
	"main/internal/repository/postgres/param_value"
	"main/internal/repository/postgres/payment"
	"main/internal/repository/postgres/product"
	"main/internal/repository/postgres/role"
	"main/internal/repository/postgres/user"
	"main/internal/repository/postgres/wishlist"

	auth_service "main/internal/services/auth"
	cart_service "main/internal/services/cart"
	category_service "main/internal/services/category"
	order_service "main/internal/services/order"
	order_status_service "main/internal/services/order_status"
	param_service "main/internal/services/param"
	param_value_service "main/internal/services/param_value"
	payment_service "main/internal/services/payment"
	product_service "main/internal/services/product"
	role_service "main/internal/services/role"
	user_service "main/internal/services/user"
	wishlist_service "main/internal/services/wishlist"

	auth_use_case "main/internal/usecase/auth"
	"main/internal/usecase/file"
	send_otp_use_case "main/internal/usecase/send_otp"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	serverPost := ":" + config.GetConfig().Port

	r := gin.Default()

	//databases
	postgresDB := postgres.NewDB()

	r.Static("/media", "../media")

	//cache
	newCache := cache.NewCache(config.GetConfig().RedisHost, config.GetConfig().RedisDB, time.Duration(config.GetConfig().RedisExpires)*time.Second)

	//repositories
	authRepository := auth.NewRepository(postgresDB)
	wishlistRepository := wishlist.NewRepository(postgresDB)
	productRepository := product.NewRepository(postgresDB)
	orderRepository := order.NewRepository(postgresDB)
	cartRepository := cart.NewRepository(postgresDB)
	userRepository := user.NewRepository(postgresDB)
	categoryRepository := category.NewRepository(postgresDB)
	orderStatusRepository := order_status.NewRepository(postgresDB)
	paymentStatusRepository := payment.NewRepository(postgresDB)
	roleRepository := role.NewRepository(postgresDB)
	paramRepository := param.NewRepository(postgresDB)
	paramValueRepository := param_value.NewRepository(postgresDB)

	//usecase
	authUseCase := auth_use_case.NewUseCase(authRepository)
	sendSMSUseCase := send_otp_use_case.NewUseCase()
	fileUseCase := file.NewUseCase()

	//services
	authService := auth_service.NewService(authRepository, newCache, sendSMSUseCase, authUseCase)
	wishlistService := wishlist_service.NewService(wishlistRepository, authUseCase)
	productService := product_service.NewService(productRepository, authUseCase, fileUseCase)
	orderService := order_service.NewService(orderRepository, authUseCase)
	cartService := cart_service.NewService(cartRepository, authUseCase)
	userService := user_service.NewService(userRepository, authUseCase, fileUseCase)
	categoryService := category_service.NewService(categoryRepository, authUseCase)
	orderStatusService := order_status_service.NewService(orderStatusRepository, authUseCase)
	paymenStatusService := payment_service.NewService(paymentStatusRepository, authUseCase)
	roleStatusService := role_service.NewService(roleRepository)
	paramService := param_service.NewService(paramRepository, authUseCase)
	paramValueService := param_value_service.NewService(paramValueRepository, authUseCase)

	//controller
	authController := auth_controller.NewController(authService)
	wishlistController := wishlist_controller.NewController(wishlistService)
	productController := product_controller.NewController(productService)
	orderController := order_controller.NewController(orderService)
	cartController := cart_controller.NewController(cartService)
	userController := user_controller.NewController(userService)
	categoryController := category_controller.NewController(categoryService)
	orderStatusController := order_status_controller.NewController(orderStatusService)
	paymentStatusController := payment_controller.NewController(paymenStatusService)
	roleController := role_controller.NewController(roleStatusService)
	paramController := param_controller.NewController(paramService)
	paramValueController := param_value_controller.NewController(*paramValueService)

	//middleware
	authMiddleware := auth_middleware.NewMiddleware(authUseCase)

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := r.Group("/api")
	{
		v1 := api.Group("v1")

		// #auth
		// send otp
		v1.POST("/admin/auth/sign-in", authController.SignIn)

		// #user
		// list
		v1.GET("/admin/user/list", authMiddleware.AuthMiddleware(), userController.AdminUserGetList)
		// get by id
		v1.GET("/admin/user/:id", authMiddleware.AuthMiddleware(), userController.AdminUserGetById)
		// create
		v1.POST("/admin/user/create", authMiddleware.AuthMiddleware(), userController.AdminCreateUser)
		//update
		v1.PATCH("/admin/user/:id", authMiddleware.AuthMiddleware(), userController.AdminUpdateUser)
		//delete
		v1.DELETE("/admin/user/delete/:id", authMiddleware.AuthMiddleware(), userController.AdminDeleteUser)

		// #category
		// list
		v1.GET("/admin/category/list", authMiddleware.AuthMiddleware(), categoryController.AdminCategoryGetList)
		// get by id
		v1.GET("/admin/category/:id", authMiddleware.AuthMiddleware(), categoryController.AdminCategoryGetById)
		// create
		v1.POST("/admin/category/create", authMiddleware.AuthMiddleware(), categoryController.AdminCategoryCreate)
		// update
		v1.PATCH("/admin/category/:id", authMiddleware.AuthMiddleware(), categoryController.AdminCategoryUpdate)
		// delete
		v1.DELETE("/admin/category/delete/:id", authMiddleware.AuthMiddleware(), categoryController.AdminCategoryDelete)
		// list by parent id
		v1.POST("/admin/category/:id", authMiddleware.AuthMiddleware(), categoryController.AdminGetByParentId)

		// #param
		// list
		v1.GET("/admin/param/list", authMiddleware.AuthMiddleware(), paramController.AdminParamGetList)
		// // get by id
		v1.GET("/admin/param/:id", authMiddleware.AuthMiddleware(), paramController.AdminParamGetById)
		// create
		v1.POST("/admin/param/create", authMiddleware.AuthMiddleware(), paramController.AdminParamCreate)
		// update
		v1.PATCH("/admin/param/:id", authMiddleware.AuthMiddleware(), paramController.AdminParamUpdate)
		// // delete
		v1.DELETE("/admin/param/delete/:id", authMiddleware.AuthMiddleware(), paramController.AdminParamDelete)

		// #param_value
		// list
		v1.GET("/admin/param-value/list", authMiddleware.AuthMiddleware(), paramValueController.AdminParamValueGetList)
		// get by id
		v1.GET("/admin/param-value/:id", authMiddleware.AuthMiddleware(), paramValueController.AdminParamValueGetById)
		// create
		v1.POST("/admin/param-value/create", authMiddleware.AuthMiddleware(), paramValueController.AdminParamValueCreate)
		// update
		v1.PATCH("/admin/param-value/:id", authMiddleware.AuthMiddleware(), paramValueController.AdminParamValueUpdate)
		//  delete
		v1.DELETE("/admin/param-value/delete/:id", authMiddleware.AuthMiddleware(), paramValueController.AdminParamValueDelete)
		// list by param id
		v1.GET("/admin/param-value/list/:id", authMiddleware.AuthMiddleware(), paramValueController.AdminParamValueGetListByParamId)

		// #wishlist
		// list
		v1.GET("/admin/wishlist", authMiddleware.AuthMiddleware(), wishlistController.AdminWishistGetList)
		// create
		v1.POST("/admin/wishlist/create", authMiddleware.AuthMiddleware(), wishlistController.AdminWishlistCreate)
		// delete
		v1.DELETE("/admin/wishlist/delete/:id", authMiddleware.AuthMiddleware(), wishlistController.AdminWishlistDelete)

		//  #products
		// create
		v1.POST("/admin/product/create", authMiddleware.AuthMiddleware(), productController.CreateProduct)
		// get by id
		v1.GET("/admin/product/:id", authMiddleware.AuthMiddleware(), productController.GetById)
		// list
		v1.GET("/admin/products", authMiddleware.AuthMiddleware(), productController.GetProductsList)
		// update
		v1.PATCH("/admin/product/update/:id", authMiddleware.AuthMiddleware(), productController.UpdateProduct)
		// // delete
		v1.DELETE("/admin/product/delete/:id", authMiddleware.AuthMiddleware(), productController.AdminDeleteProduct)

		// #orders
		// create
		// v1.POST("/order/create", authMiddleware.AuthMiddleware(), orderController.CreateOrder)
		// list
		v1.GET("/order/list", authMiddleware.AuthMiddleware(), orderController.AdminOrderGetList)
		// get by id
		v1.GET("/order/:id", authMiddleware.AuthMiddleware(), orderController.AdminOrderGetById)
		// update
		// v1.PATCH("/order/update/:id", authMiddleware.AuthMiddleware(), orderController.UpdateOrder)
		// delete
		v1.DELETE("/order/delete/:id", authMiddleware.AuthMiddleware(), orderController.AdminOrderDelete)

		// #cart
		// create
		v1.POST("/cart/create", authMiddleware.AuthMiddleware(), cartController.AdminCartCreate)
		// cart item total update
		v1.PATCH("/cart/item/:id/update", authMiddleware.AuthMiddleware(), cartController.AdminUpdateCartItemTotal)
		// delete cart item
		v1.DELETE("/cart/item/delete/:id", authMiddleware.AuthMiddleware(), cartController.AdminDeleteCartItem)
		// get cart list
		v1.GET("/cart/list", authMiddleware.AuthMiddleware(), cartController.AdminGetCartList)

		// #order status
		// list
		v1.GET("/admin/order-status/list", authMiddleware.AuthMiddleware(), orderStatusController.AdminOrderStatusGetList)
		// get by id
		v1.GET("/admin/order-status/:id", authMiddleware.AuthMiddleware(), orderStatusController.AdminOrderStatusGetById)
		// create
		v1.POST("/admin/order-status/create", authMiddleware.AuthMiddleware(), orderStatusController.AdminOrderStatusCreate)
		// update
		v1.PATCH("/admin/order-status/:id", authMiddleware.AuthMiddleware(), orderStatusController.AdminOrderStatusUpdate)
		// delete
		v1.DELETE("/admin/order-status/delete/:id", authMiddleware.AuthMiddleware(), orderStatusController.AdminOrderStatusDelete)

		// #payment status
		// list
		v1.GET("/admin/payment-status/list", authMiddleware.AuthMiddleware(), paymentStatusController.AdminPaymentGetList)
		// get by id
		v1.GET("/admin/payment-status/:id", authMiddleware.AuthMiddleware(), paymentStatusController.AdminPaymentGetById)
		// create
		v1.POST("/admin/payment-status/create", authMiddleware.AuthMiddleware(), paymentStatusController.AdminPaymentCreate)
		// update
		v1.PATCH("/admin/payment-status/:id", authMiddleware.AuthMiddleware(), paymentStatusController.AmdinPaymentUpdate)
		// delete
		v1.DELETE("/admin/payment-status/delete/:id", authMiddleware.AuthMiddleware(), paymentStatusController.AdminPaymentDelete)

		// #role
		// list
		v1.GET("/admin/role/list", authMiddleware.AuthMiddleware(), roleController.AdminGetList)
	}

	r.Run(serverPost)

}
