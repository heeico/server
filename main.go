package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/iamrahultanwar/heeico/api"
	"github.com/iamrahultanwar/heeico/database"
	"github.com/iamrahultanwar/heeico/model"
	"github.com/iamrahultanwar/heeico/types"
)

func init() {
	database.ConnectDB()
	database.Init()
}

func main() {
	app := fiber.New()

	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	app.Use(func(c *fiber.Ctx) error {
		paths := strings.Split(c.Path(), "/")
		formatPaths := []string{}
		for _, p := range paths {
			if len(p) > 0 {
				formatPaths = append(formatPaths, p)
			}
		}
		if len(formatPaths) == 0 {
			return c.Redirect("/admin")
		}
		if formatPaths[0] == "api" {
			return c.Next()
		}

		if formatPaths[0] != "admin" {
			return c.Next()
		}
		lastPath := formatPaths[len(formatPaths)-1]
		dir, _ := os.Getwd()
		fmt.Println(dir)
		if strings.Contains(lastPath, ".") {
			typeOfFile := paths[len(paths)-2]
			file := paths[len(paths)-1]
			if typeOfFile == "img" {
				return c.SendFile(filepath.Join(dir, "client", "img", file))
			}
			if typeOfFile == "mainfest.json" {
				return c.SendFile(filepath.Join(dir, "client", file))
			}
			return c.SendFile(filepath.Join(dir, "client", strings.Join(formatPaths[1:], "/")))
		}
		return c.SendFile(filepath.Join(dir, "client", "index.html"))
	})

	app.Use(cors.New())

	// redirect
	app.Get("/:alias", func(c *fiber.Ctx) error {
		alias := c.Params("alias")
		var urlShortener model.URLShortener
		result := database.DB.Where("alias = ?", alias).First(&urlShortener)
		if result.Error != nil {
			c.SendStatus(http.StatusBadRequest)
			if result.Error == sql.ErrNoRows {
				return c.JSON(types.FailResponse{Error: "No alias found", Status: false})
			}
			return c.JSON(types.FailResponse{Error: result.Error.Error(), Status: false})
		}
		if result.RowsAffected == 0 {
			return c.JSON(types.SuccessResponse{Message: "Alias not found", Status: false})
		}
		count := urlShortener.VisitCount + 1
		go database.DB.Model(&urlShortener).Where("id = ?", &urlShortener.ID).Update("visit_count", count)
		go database.DB.Create(&model.URLTrack{URLShortenerID: urlShortener.ID, IPAddress: c.IP()})
		return c.Redirect(urlShortener.RedirectURL)
	})

	app.Use(AuthMiddleware)

	api.Api(app)

	app.Use(func(c *fiber.Ctx) error {
		return c.Redirect("/admin")
	})

	app.Listen(":8080")

}

func AuthMiddleware(c *fiber.Ctx) error {

	skipPaths := []string{"/api/v1/admin/auth/register", "/api/v1/admin/auth/login"}

	for _, val := range skipPaths {
		if val == c.Path() {
			return c.Next()
		}
	}

	headers := c.GetReqHeaders()
	authToken := headers["Authorization"]

	if len(authToken) == 0 {
		c.SendStatus(http.StatusUnauthorized)
		return c.JSON(types.FailResponse{Status: false, Error: "Invalid auth token", Data: map[string]interface{}{"message": "Unauthorized access", "error": "Auth token not provided"}})
	}
	userId, err := model.ValidateAuthToken(authToken)
	c.Locals("role", "guest")
	if err == nil {
		c.Locals("userId", userId)
		c.Locals("role", "auth")
	}
	return c.Next()
}
