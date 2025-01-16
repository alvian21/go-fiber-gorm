package utils

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const DefaultPathAssetImage = "./public/covers/"

func HandleSingleFile(ctx *fiber.Ctx) error {
	// HANDLE FILE
	file, errFile := ctx.FormFile("cover")
	if errFile != nil {
		log.Println("errFile = ", errFile)
	}

	var filename *string

	var newFilename string

	if file != nil {

		errCheckContentType := checkContentType(file, "image/jpeg", "image/png")

		if errCheckContentType != nil {
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"message": errCheckContentType.Error(),
			})
		}

		filename = &file.Filename
		extFile := filepath.Ext(*filename)
		// Get current Unix timestamp and convert it to a string
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		// Generate the new filename
		newFilename = fmt.Sprintf("image%s%s", timestamp, extFile)
		fmt.Println(newFilename)

		errSaveFile := ctx.SaveFile(file, fmt.Sprintf("./public/covers/%s", newFilename))
		if errSaveFile != nil {
			log.Println("errSaveFile = ", errSaveFile)
		}

	} else {
		log.Println("nothing file to upload")
	}

	if newFilename != "" {
		ctx.Locals("filename", newFilename)
	} else {
		ctx.Locals("filename", nil)
	}

	return ctx.Next()
}

func HandleMultipleFile(ctx *fiber.Ctx) error {
	form, errForm := ctx.MultipartForm()
	if errForm != nil {
		log.Println("errForm = ", errForm)
	}

	files := form.File["photos"]

	var filenames []string

	for i, file := range files {
		var filename string

		if file != nil {

			extFile := filepath.Ext(file.Filename)
			// Get current Unix timestamp and convert it to a string
			timestamp := strconv.FormatInt(time.Now().Unix(), 10)
			// Generate the new filename

			filename = fmt.Sprintf("%d-%s%s%s", i, "gambar", timestamp, extFile)

			errSaveFile := ctx.SaveFile(file, fmt.Sprintf("./public/covers/%s", filename))
			if errSaveFile != nil {
				log.Println("errSaveFile = ", errSaveFile)
			}

		} else {
			log.Println("nothing file to upload")
		}

		if filename != "" {
			filenames = append(filenames, filename)

		}

	}

	ctx.Locals("filenames", filenames)

	return ctx.Next()
}

func HandleRemovefile(filename string, pathFile ...string) error {

	if len(pathFile) > 0 {
		err := os.Remove(pathFile[0] + filename)
		if err != nil {
			log.Println("err remove file")
			return err
		}
	} else {
		err := os.Remove(DefaultPathAssetImage + filename)
		if err != nil {
			log.Println("err remove file default")
			return err
		}

	}

	return nil
}

func checkContentType(file *multipart.FileHeader, contentTypes ...string) error {
	if len(contentTypes) > 0 {
		for _, contentType := range contentTypes {
			contentTypeFile := file.Header.Get("Content-Type")
			if contentTypeFile == contentType {
				return nil
			}
		}

		return errors.New("image type not supported")
	} else {
		return errors.New("not found content type to be check")
	}
}
