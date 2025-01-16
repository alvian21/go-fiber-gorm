package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go-fiber-gorm/database"
	"go-fiber-gorm/model/entity"
	"go-fiber-gorm/model/request"
	"go-fiber-gorm/utils"
	"log"
)

func PhotoHandlerCreate(ctx *fiber.Ctx) error {
	photo := new(request.PhotoCreateRequest)

	if err := ctx.BodyParser(photo); err != nil {
		log.Println(err)
		return err
	}

	validate := validator.New()
	errValidate := validate.Struct(photo)
	if errValidate != nil {
		return ctx.Status(400).JSON(fiber.Map{"message": "failed", "errors": errValidate.Error()})
	}

	filenames := ctx.Locals("filenames")

	if filenames == nil {
		return ctx.Status(400).JSON(fiber.Map{"message": "failed", "errors": "cover is required"})
	} else {
		filenamesData := filenames.([]string)

		for _, filename := range filenamesData {

			newPhoto := entity.Photo{
				CategoryID: photo.CategoryID,
				Image:      filename,
			}

			errCreate := database.DB.Create(&newPhoto).Error

			if errCreate != nil {
				return ctx.Status(500).JSON(fiber.Map{"message": errCreate.Error()})
			}
		}
	}

	return ctx.JSON(fiber.Map{"message": "photo created"})
}

func PhotoHandlerDelete(ctx *fiber.Ctx) error {
	photoId := ctx.Params("id")

	var photo entity.Photo

	result := database.DB.Debug().First(&photo, "id = ? ", photoId).Error
	if result != nil {
		return ctx.Status(404).JSON(fiber.Map{"message": "photo not found"})
	}

	errDelete := database.DB.Debug().Delete(&photo).Error

	if errDelete != nil {
		return ctx.Status(500).JSON(fiber.Map{"message": "internal server error"})
	}

	errDeleteFile := utils.HandleRemovefile(photo.Image)

	if errDeleteFile != nil {
		return ctx.Status(500).JSON(fiber.Map{"message": "err remove file"})
	}

	return ctx.JSON(fiber.Map{"message": "photo deleted"})
}
