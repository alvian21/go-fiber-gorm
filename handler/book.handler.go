package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go-fiber-gorm/database"
	"go-fiber-gorm/model/entity"
	"go-fiber-gorm/model/request"
	"log"
)

func BookHandlerCreate(ctx *fiber.Ctx) error {
	book := new(request.BookCreateRequest)

	if err := ctx.BodyParser(book); err != nil {
		log.Println(err)
		return err
	}

	validate := validator.New()
	errValidate := validate.Struct(book)
	if errValidate != nil {
		return ctx.Status(400).JSON(fiber.Map{"message": "failed", "errors": errValidate.Error()})
	}

	filename := ctx.Locals("filename")

	var filenameString string

	if filename == "" {
		return ctx.Status(400).JSON(fiber.Map{"message": "failed", "errors": "cover is required"})
	} else {
		filenameString = filename.(string)
	}

	newBook := entity.Book{
		Title:  book.Title,
		Author: book.Author,
		Cover:  filenameString,
	}

	errCreate := database.DB.Create(&newBook).Error

	if errCreate != nil {
		return ctx.Status(500).JSON(fiber.Map{"message": errCreate.Error()})
	}

	return ctx.JSON(fiber.Map{"data": newBook, "message": "user created"})
}
