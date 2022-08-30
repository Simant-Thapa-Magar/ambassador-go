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

	userRoutes := authetication.Group("user")
	userRoutes.Get("/", controllers.User)
	userRoutes.Put("update", controllers.UpdateUser)
	userRoutes.Put("update-password", controllers.UpdatePassword)

	authetication.Get("ambassadors", controllers.GetAmbassadors)

	productRoutes := authetication.Group("product")
	productRoutes.Post("/create", controllers.CreateProduct)
	productRoutes.Get("/:id", controllers.GetProduct)
	productRoutes.Put("/update/:id", controllers.UpdateProduct)
	productRoutes.Delete("/delete/:id", controllers.DeleteProduct)
}
