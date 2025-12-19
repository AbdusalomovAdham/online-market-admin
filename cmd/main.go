package main

import (
	"main/internal/cache"
	auth_controller "main/internal/controllers/http/v1/auth"
	cart_controller "main/internal/controllers/http/v1/cart"
	category_controller "main/internal/controllers/http/v1/category"
	order_controller "main/internal/controllers/http/v1/order"
	product_controller "main/internal/controllers/http/v1/product"
	rating_controller "main/internal/controllers/http/v1/rating"
	user_controller "main/internal/controllers/http/v1/user"
	wishlist_controller "main/internal/controllers/http/v1/wishlist"

	auth_middleware "main/internal/middleware/auth"

	"main/internal/pkg/config"
	"main/internal/pkg/postgres"

	auth "main/internal/repository/postgres/auth"
	"main/internal/repository/postgres/cart"
	"main/internal/repository/postgres/category"
	"main/internal/repository/postgres/order"
	product "main/internal/repository/postgres/product"
	"main/internal/repository/postgres/rating"
	"main/internal/repository/postgres/user"
	wishlist "main/internal/repository/postgres/wishlist"

	auth_service "main/internal/services/auth"
	cart_service "main/internal/services/cart"
	category_service "main/internal/services/category"
	order_service "main/internal/services/order"
	product_service "main/internal/services/product"
	rating_service "main/internal/services/rating"
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
	ratingRepository := rating.NewRepository(postgresDB)
	userRepository := user.NewRepository(postgresDB)
	categoryRepository := category.NewRepository(postgresDB)

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
	ratingService := rating_service.NewService(ratingRepository, authUseCase)
	userService := user_service.NewService(userRepository, authUseCase, fileUseCase)
	categoryService := category_service.NewService(categoryRepository, authUseCase)

	//controller
	authController := auth_controller.NewController(authService)
	wishlistController := wishlist_controller.NewController(wishlistService)
	productController := product_controller.NewController(productService)
	orderController := order_controller.NewController(orderService)
	cartController := cart_controller.NewController(cartService)
	ratingController := rating_controller.NewController(ratingService)
	userController := user_controller.NewController(userService)
	categoryController := category_controller.NewController(categoryService)

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
		v1.GET("/admin/user/list", authMiddleware.AuthMiddleware(), userController.AdminGetList)
		// get by id
		v1.GET("/admin/user/:id", authMiddleware.AuthMiddleware(), userController.AdminGetById)
		// create
		v1.POST("/admin/user/create", authMiddleware.AuthMiddleware(), userController.AdminCreateUser)
		//update
		v1.PATCH("/admin/user/:id", authMiddleware.AuthMiddleware(), userController.AdminUpdateUser)
		//delete
		v1.DELETE("/admin/user/delete/:id", authMiddleware.AuthMiddleware(), userController.AdminDeleteUser)

		// #category
		// list
		v1.GET("/admin/category/list", authMiddleware.AuthMiddleware(), categoryController.AdminGetCategoryList)
		// get by id
		v1.GET("/admin/category/:id", authMiddleware.AuthMiddleware(), categoryController.AdminGetCategoryById)
		// create
		v1.POST("/admin/category/create", authMiddleware.AuthMiddleware(), categoryController.AdminCreateCategory)
		// update
		v1.PATCH("/admin/category/:id", authMiddleware.AuthMiddleware(), categoryController.AdminUpdateCategory)
		// delete
		v1.DELETE("/admin/category/delete/:id", authMiddleware.AuthMiddleware(), categoryController.AdminDeleteCategory)

		// #wishlist
		// list
		v1.GET("/wishlist", authMiddleware.AuthMiddleware(), wishlistController.WishList)
		// create
		v1.POST("/wishlist/create", authMiddleware.AuthMiddleware(), wishlistController.Create)
		// delete
		v1.DELETE("/wishlist/delete/:id", authMiddleware.AuthMiddleware(), wishlistController.Delete)

		//  #products
		// create
		v1.POST("/admin/product/create", authMiddleware.AuthMiddleware(), productController.CreateProduct)
		// get by id
		v1.GET("/admin/product/:id", authMiddleware.AuthMiddleware(), productController.GetById)
		// list
		v1.GET("/admin/products", productController.GetProductsList)
		// update
		v1.PATCH("/admin/product/update/:id", authMiddleware.AuthMiddleware(), productController.UpdateProduct)
		// // delete
		// v1.DELETE("/admin/product/delete/:id", authMiddleware.AuthMiddleware(), productController.DeleteProduct)

		// #orders
		// create
		v1.POST("/order/create", authMiddleware.AuthMiddleware(), orderController.CreateOrder)
		// list
		v1.GET("/order/list", authMiddleware.AuthMiddleware(), orderController.GetOrderList)
		// get by id
		v1.GET("/order/:id", authMiddleware.AuthMiddleware(), orderController.GetOrderById)
		// delete
		v1.DELETE("/order/delete/:id", authMiddleware.AuthMiddleware(), orderController.DeleteOrder)

		// #cart
		// create
		v1.POST("/cart/create", authMiddleware.AuthMiddleware(), cartController.CreateCart)
		// cart item total update
		v1.PATCH("/cart/item/:id/update", authMiddleware.AuthMiddleware(), cartController.UpdateCartItemTotal)
		// delete cart item
		v1.DELETE("/cart/item/delete/:id", authMiddleware.AuthMiddleware(), cartController.DeleteCartItem)
		// get cart list
		v1.GET("/cart/list", authMiddleware.AuthMiddleware(), cartController.GetCartList)

		// #rating
		// create
		v1.POST("/create/rating/:id", authMiddleware.AuthMiddleware(), ratingController.CreateRating)

	}

	r.Run(serverPost)

}
