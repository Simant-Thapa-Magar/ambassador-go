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

	authetication.Get("user", controllers.User)
	authetication.Post("logout", controllers.Logout)
}
