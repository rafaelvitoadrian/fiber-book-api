package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rafaelvitoadrian/fiber-book-api/controllers/bookcontroller"
	"github.com/rafaelvitoadrian/fiber-book-api/models"
)

func main() {
	models.ConnectDatabase()

	app := fiber.New()
	api := app.Group("/api")
	book := api.Group("/book")

	book.Get("/", bookcontroller.Index)
	book.Get("/:id", bookcontroller.Show)
	book.Post("/", bookcontroller.Create)
	book.Put("/:id", bookcontroller.Update)
	book.Delete("/:id", bookcontroller.Delete)

	app.Listen(":3000")
}
