package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jgcaceres97/go-auth-jwt/src/controllers"
)

func Setup(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)
	api.Post("/logout", controllers.Logout)

	api.Get("/users", controllers.GetUsers)

	user := api.Group("/user")
	user.Get("/:id?", controllers.GetUser)
	user.Put("/:id?", controllers.UpdateUser)
	user.Delete("/:id?", controllers.DeleteUser)
}
