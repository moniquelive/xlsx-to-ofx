// Package handler defines the http handlers
package handler

import (
	"fmt"
	"mime/multipart"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/utils"

	"github.com/moniquelive/xlsx-to-ofx/generator"
	"github.com/moniquelive/xlsx-to-ofx/parser"
)

var csrfErrorHandler = func(c *fiber.Ctx, err error) error {
	// Log the error so we can track who is trying to perform CSRF attacks
	// customize this to your needs
	log.Debugf("CSRF Error: %v Request: %v From: %v\n", err, c.OriginalURL(), c.IP())

	// Return a 403 Forbidden response for all other requests
	return c.Status(fiber.StatusForbidden).SendString("403 Forbidden")
}

// CsrfProtection acts as a middleware for the http routes
var CsrfProtection = csrf.New(csrf.Config{
	KeyLookup:      "form:_csrf",
	CookieName:     "csrf_",
	CookieSameSite: "Strict",
	Expiration:     1 * time.Hour,
	KeyGenerator:   utils.UUIDv4,
	ContextKey:     "token",
	ErrorHandler:   csrfErrorHandler,
	SingleUseToken: true,
})

// Index GET handler
func Index(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"csrfToken": c.Locals("token"),
	}, "layouts/main")
}

// DoConvert POST handler
func DoConvert(c *fiber.Ctx) (err error) {
	var formFile *multipart.FileHeader
	if formFile, err = c.FormFile("xlsxfile"); err != nil {
		return err
	}
	var file multipart.File
	if file, err = formFile.Open(); err != nil {
		return err
	}
	var records []parser.Record
	if records, err = parser.ParseReader(file, formFile.Size); err != nil {
		return err
	}
	acct := c.FormValue("account")
	c.Attachment(fmt.Sprintf("converted-%s.ofx", acct))
	data := generator.OFXData{
		Now:     time.Now(),
		Agencia: "1053",
		Conta:   acct,
		Records: records[1:],
		Balance: records[len(records)-1].Balance,
	}
	if err = generator.Fill(data, c.Response().BodyWriter()); err != nil {
		return err
	}
	return nil
}
