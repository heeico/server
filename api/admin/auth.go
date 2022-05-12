package api

import (
	"database/sql"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/iamrahultanwar/heeico/database"
	"github.com/iamrahultanwar/heeico/model"
	"github.com/iamrahultanwar/heeico/types"
)

func AuthApi(app *fiber.App) {

	// auth operations
	auth := app.Group("/api/v1/admin/auth")
	auth.Post("register", func(c *fiber.Ctx) error {
		user := model.User{}
		err := c.BodyParser(&user)
		if err != nil {
			c.SendStatus(http.StatusBadRequest)
			return c.JSON(types.FailResponse{Status: false, Error: "Valid data not supplied"})
		}
		user.HashPassword()
		result := database.DB.Create(&user)
		if result.Error != nil {
			c.SendStatus(http.StatusInternalServerError)
			return c.JSON(types.FailResponse{Status: false, Error: result.Error.Error()})
		}
		token := user.GetAuthToken()
		return c.JSON(types.SuccessResponse{Status: true, Message: "Register Successfull", Data: map[string]interface{}{"token": token}})
	})

	auth.Post("login", func(c *fiber.Ctx) error {
		user := model.User{}
		err := c.BodyParser(&user)
		if err != nil {
			c.SendStatus(http.StatusBadRequest)
			return c.JSON(types.FailResponse{Status: false, Error: "Valid data not supplied"})
		}
		u := model.User{}
		result := database.DB.Where("email", user.Email).First(&u)
		if result.Error != nil {
			if result.Error == sql.ErrNoRows {
				c.SendStatus(http.StatusBadRequest)
				return c.JSON(types.FailResponse{Status: false, Error: "email not found"})
			}
			c.SendStatus(http.StatusInternalServerError)
			return c.JSON(types.FailResponse{Status: false, Error: result.Error.Error()})
		}
		e := user.CheckPasswordHash(u.Password)
		if e != nil {
			c.SendStatus(http.StatusBadRequest)
			return c.JSON(types.FailResponse{Status: false, Error: "password does not match"})
		}
		token := user.GetAuthToken()

		return c.JSON(types.SuccessResponse{Status: true, Message: "Login Successfull", Data: map[string]interface{}{"token": token}})
	})
}
