package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/rafaelvitoadrian/fiber-book-api/controllers/bookcontroller"
	"github.com/rafaelvitoadrian/fiber-book-api/models"
)

func main() {
	models.ConnectDatabase()

	engine := html.New("./", ".html")
	app := fiber.New(fiber.Config{Views: engine})
	app.Post("/upload", bookcontroller.UploadImage)

	api := app.Group("/api")
	upload := api.Group("/upload")
	upload.Get("/", bookcontroller.UploadPhoto)

	book := api.Group("/book")

	book.Get("/", bookcontroller.Index)
	book.Get("/:id", bookcontroller.Show)
	book.Post("/", bookcontroller.Create)
	book.Put("/:id", bookcontroller.Update)
	book.Delete("/:id", bookcontroller.Delete)

	app.Listen(":3000")
}
