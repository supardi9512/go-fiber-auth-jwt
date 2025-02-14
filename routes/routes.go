package routes

import (
	c "go-fiber-auth-jwt/controllers"

	"github.com/gofiber/fiber/v2"

	jwtware "github.com/gofiber/contrib/jwt"
)

func Setup(app *fiber.App) {

	app.Post("/api/register", c.Register)
	app.Post("/api/login", c.Login)

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte("secret")},
	}))

	app.Post("/api/logout", c.Logout)
	app.Get("/api/user", c.User)
}
