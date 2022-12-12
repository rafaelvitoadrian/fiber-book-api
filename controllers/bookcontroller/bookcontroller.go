package bookcontroller

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

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

func UploadPhoto(c *fiber.Ctx) error {
	return c.Render("index", nil)
}

func UploadImage(c *fiber.Ctx) error {
	var input struct {
		Nama_gambar string
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	gambar, err := c.FormFile("gambar")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	//ambil nama file
	fmt.Printf("Nama File: %s \n", gambar.Filename)

	//ambil ukuran file
	fmt.Printf("Ukuran File (byte): %d \n", gambar.Size)

	//kb
	fmt.Printf("Ukuran File (kb): %d \n", gambar.Size/1024)

	//mb
	fmt.Printf("Ukuran File (mb): %f \n", (float64(gambar.Size)/1024)/1024)

	//mengimbil mime type
	fmt.Printf("Mime type: %s \n", gambar.Header.Get("Content-Type"))

	// mengambil nama extensi
	splitDots := strings.Split(gambar.Filename, ".")
	ext := splitDots[len(splitDots)-1]
	fmt.Println(ext)

	namaFileBaru := fmt.Sprintf("%s.%s", time.Now().Format("2006-01-02-15-04-05"), ext)
	fmt.Println(namaFileBaru)

	//ambil ukuran gambar
	fileHeader, _ := gambar.Open()
	defer fileHeader.Close()

	imageConfig, _, err := image.DecodeConfig(fileHeader)
	if err != nil {
		log.Print(err)
	}

	width := imageConfig.Width
	height := imageConfig.Height
	fmt.Printf("Widht %d \n", width)
	fmt.Printf("Height %d \n", height)

	//membuat folder upload
	folderUpload := filepath.Join(".", "uploads")

	//mkdirALL => proses pembuat folder all agar tidak err ketika sudah ada
	if err := os.MkdirAll(folderUpload, 0770); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	//simpan ke directori
	if err := c.SaveFile(gambar, "./uploads/"+namaFileBaru); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"title":       input.Nama_gambar,
		"nama_gamabr": namaFileBaru,
		"messages":    "Berhasil Menyimpan Gambar",
	})
}
