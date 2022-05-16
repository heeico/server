package api

import (
	"github.com/gofiber/fiber/v2"
	api "github.com/iamrahultanwar/heeico/api/admin"
)

func Api(app *fiber.App) {

	api.AdminApi(app)
	api.AuthApi(app)
	api.TeamApi(app)
}
