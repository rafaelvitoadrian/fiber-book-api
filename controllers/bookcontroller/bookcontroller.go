package bookcontroller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rafaelvitoadrian/fiber-book-api/models"
	"gorm.io/gorm"
)

func Index(c *fiber.Ctx) error {
	var books []models.Book
	models.DB.Find(&books)
	return c.Status(fiber.StatusOK).JSON(books)
}

func Show(c *fiber.Ctx) error {
	id := c.Params("id")
	var book []models.Book
	if err := models.DB.First(&book, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"Message": "Data Tidak Ditemukan",
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"Messages": "Data tidak Ditemukan",
		})
	}

	return c.JSON(book)
}

func Create(c *fiber.Ctx) error {

	var books models.Book
	if err := c.BodyParser(&books); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := models.DB.Create(&books).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(books)
}

func Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var book models.Book
	if err := c.BodyParser(&book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if models.DB.Where("id = ?", id).Updates(&book).RowsAffected == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "Tidak Dapat mengupdate data",
		})
	}

	return c.JSON(fiber.Map{
		"msg": "data berhasil di update",
	})
}

func Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	var book models.Book
	if models.DB.Delete(&book, id).RowsAffected == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "Tidak dapat menghapus data",
		})
	}

	return c.JSON(fiber.Map{
		"msg": "data berhasil dihapus",
	})
}
