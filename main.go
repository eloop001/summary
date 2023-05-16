package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/go-mail/mail"
	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
)

// create function init
func init() {
	// load variable from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {

	var allowed_filetypes [4]string
	allowed_filetypes[0] = ".mp4"
	allowed_filetypes[1] = ".mp3"
	allowed_filetypes[2] = ".docx"
	allowed_filetypes[3] = ".pdf"

	// load template
	engine := html.New("./views", ".html")

	// create go fiber app
	app := fiber.New(fiber.Config{
		Views:     engine,
		BodyLimit: 1024 * 1024 * 1024,
	})
	// seave static assets
	app.Static("/static", "./public")
	// create route
	app.Get("/", func(c *fiber.Ctx) error {
		var current_user_key string

		queryValue := c.Query("key")
		if queryValue == "abc" {
			current_user_key = queryValue
			return c.Render("index", fiber.Map{
				"Title": "summary.run - GDPR compliant file upload",
			})
			//
		}
		if queryValue == "" && current_user_key == "abc" {
			return c.Render("index", url.Values{"key": {"abc"}})
		} else {
			return c.SendString("Invalid key")
		}

	})

	// create a post route to accept upload file
	app.Post("/", func(c *fiber.Ctx) error {
		// get file from form

		file, err := c.FormFile("upload")
		if err != nil {
			return err
		}

		var allowed bool
		allowed = false

		fileguid := uuid.New().String()
		file_suffix := file.Filename[strings.LastIndex(file.Filename, "."):]

		// check if file suffix is in a list of allowed file types
		for _, allowed_filetype := range allowed_filetypes {
			if file_suffix == allowed_filetype {
				err = c.SaveFile(file, "./uploads/"+fileguid+file_suffix)
				send_mail("bigdata@ilab.dk", 0, file.Filename)
				fmt.Print("Mail sent")
				allowed = true
				if err != nil {
					return err
				}
			}
		}
		// save file to public folder
		// 1st param is file that we get from form
		// 2nd param is path to save file

		//return c.Render("index", url.Values{"key": {"abc"}})
		// return success message
		if allowed == false {
			return c.Render("index", fiber.Map{
				"Title":  "Home",
				"Msg":    "❌ File type not allowed!",
				"File":   file.Filename,
				"Params": url.Values{"key": {c.Params("key")}},
			})
		} else {
			return c.Render("index", fiber.Map{
				"Title":  "Home",
				"Msg":    "✔ File uploaded successfully!",
				"File":   file.Filename,
				"Params": url.Values{"key": {c.Params("key")}},
			})
		}
	})

	// start app with port from .env file
	err := app.Listen(os.Getenv("HOST") + ":" + os.Getenv("PORT"))
	if err != nil {
		println(err.Error())
	}

}

func send_mail(To string, Type int, fname string) {
	if len(fname) > 10 {
		fname = fname[0:4] + "***" + fname[len(fname)-5:]
	} else {
		fname = fname[0:1] + "*" + fname[len(fname)-4:]
	}

	var Mail_texts [3]string
	var Subject_texts [3]string

	Mail_texts[0] = "Thank you for using summary.run. We have recieved your file " + fname + ", and summarization has started."
	Mail_texts[1] = "Your summary is ready. Please see attached file."
	Mail_texts[2] = ""

	Subject_texts[0] = "Summary.run - File recieved"
	Subject_texts[1] = "Summary.run - Summary ready"

	m := mail.NewMessage()

	m.SetHeader("From", "run@summary.run")
	m.SetHeader("To", To)
	//m.SetAddressHeader("Cc", "oliver.doe@example.com", "Oliver")

	m.SetHeader("Subject", Subject_texts[Type])
	m.SetBody("text/html", Mail_texts[Type])
	//m.Attach("lolcat.jpg")
	d := mail.NewDialer("send.one.com", 465, "run@summary.run", "gpOUy5T8gT2jiSnYvy8B")

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
