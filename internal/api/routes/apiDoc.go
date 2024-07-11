package routes

import (
	_ "ospm/docs/api"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
)

func SetupAPIDocs(rg fiber.Router) {
	rg.Get("/*", swagger.New(swagger.Config{
		DefaultModelsExpandDepth: -1, // Removes models from the bottom of the swagger page

	}))
}
