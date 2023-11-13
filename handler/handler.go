package handler

import (
	"bytes"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/utils"

	"github.com/moniquelive/xlsx-to-ofx/generator"
	"github.com/moniquelive/xlsx-to-ofx/parser"
)

var CsrfProtection = csrf.New(csrf.Config{
	KeyLookup:      "form:_csrf",
	CookieName:     "csrf_",
	CookieSameSite: "Strict",
	Expiration:     1 * time.Hour,
	KeyGenerator:   utils.UUIDv4,
	ContextKey:     "token",
})

// Index handler. HTTP method GET.
func Index(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"csrfToken": c.Locals("token"),
	}, "layouts/main")
}

// DoConvert handler. HTTP method POST.
func DoConvert(c *fiber.Ctx) error {
	formFile, err := c.FormFile("xlsxfile")
	if err != nil {
		return err
	}
	file, err := formFile.Open()
	if err != nil {
		return err
	}
	records, err := parser.ParseReader(file, formFile.Size)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	acct := c.FormValue("account")
	data := generator.OFXData{
		Now:     time.Now(),
		Agencia: "1053",
		Conta:   acct,
		Records: records[1:],
		Balance: records[len(records)-1].Balance,
	}
	if err := generator.Fill(data, buffer); err != nil {
		return err
	}
	c.Attachment(fmt.Sprintf("converted-%s.ofx", acct))
	return c.SendStream(buffer, buffer.Len())
}
