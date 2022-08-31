package route

import (
	"ambassador/src/controllers"
	"ambassador/src/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetUp(app *fiber.App) {
	api := app.Group("api")
	admin := api.Group("admin")

	admin.Post("register", controllers.Register)
	admin.Post("login", controllers.Login)

	authetication := admin.Use(middlewares.IsAuthenticatedUser)
	authetication.Post("logout", controllers.Logout)

	userRoutes := authetication.Group("users")
	userRoutes.Get("/", controllers.User)
	userRoutes.Put("update", controllers.UpdateUser)
	userRoutes.Put("update-password", controllers.UpdatePassword)
	userRoutes.Get("/:id/links", controllers.GetUserLinks)

	authetication.Get("ambassadors", controllers.GetAmbassadors)

	productRoutes := authetication.Group("products")
	productRoutes.Post("/create", controllers.CreateProduct)
	productRoutes.Get("/:id", controllers.GetProduct)
	productRoutes.Put("/update/:id", controllers.UpdateProduct)
	productRoutes.Delete("/delete/:id", controllers.DeleteProduct)

	authetication.Get("/orders", controllers.GetOrders)

	ambassador := api.Group("ambassador")
	ambassador.Post("register", controllers.Register)
	ambassador.Post("login", controllers.Login)

	ambassadorAuthentication := ambassador.Use(middlewares.IsAuthenticatedUser)
	ambassadorAuthentication.Post("logout", controllers.Logout)

	ambassadorUserRoutes := ambassadorAuthentication.Group("users")
	ambassadorUserRoutes.Get("/", controllers.User)
	ambassadorUserRoutes.Put("update", controllers.UpdateUser)
	ambassadorUserRoutes.Put("update-password", controllers.UpdatePassword)

	ambassadorProductRoutes := ambassadorAuthentication.Group("products")
	ambassadorProductRoutes.Get("/frontend", controllers.ProductFrontend)
	ambassadorProductRoutes.Get("/backend", controllers.ProductBackend)

	ambassador.Post("link", controllers.CreateLink)
	ambassador.Get("stats", controllers.Stats)
	ambassador.Get("rankings", controllers.GetRanking)

	checkoutRoute := api.Group("checkout")
	checkoutRoute.Get("links/:code", controllers.GetLink)
	checkoutRoute.Post("orders", controllers.CreateOrder)
}
