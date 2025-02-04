package www

import (
	"embed"
	"html/template"
	"net/http"

	"codeberg.org/shinyzero0/oleg-soul-2024/lib/strg"
	"codeberg.org/shinyzero0/oleg-soul-2024/lib/strg/sqlite"
	"codeberg.org/shinyzero0/oleg-soul-2024/lib/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/jmoiron/sqlx"
)

//go:embed view/*
var viewfs embed.FS

func MakeApp(db *sqlx.DB) (*fiber.App, error) {
	gcc, err := sqlite.MakeGetCookies(db)
	if err != nil {
		return nil, err
	}
	gd, err := sqlite.MakeGetDomains(db)
	if err != nil {
		return nil, err
	}
	gc, err := sqlite.MakeGetCookie(db)
	if err != nil {
		return nil, err
	}

	engine := html.NewFileSystem(http.FS(viewfs), ".tmpl")
	engine.AddFunc(
		// add unescape function
		"unescape", func(s string) template.HTML {
			return template.HTML(s)
		},
	)

	app := fiber.New(fiber.Config{Views: engine})
	app.Get("/", func(c *fiber.Ctx) error { return c.RedirectToRoute("cookies", fiber.Map{}) })
	app.Get("/cookies", MakeGetDomains(gd)).Name("cookies")
	app.Get("/cookies/:domain", MakeGetCookies(gcc, gd))
	app.Get("/cookie/:id", MakeGetCookie(gcc, gd, gc))
	return app, nil
}

func MakeGetCookie(gcc strg.GetCookies, gd strg.GetDomains, gc strg.GetCookie) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := utils.GetIntParam(c, "id")
		if err != nil {
			return err
		}
		cookie, err := gc(c.Context(), id)
		if err != nil {
			return err
		}
		cookies, err := gcc(c.Context(), cookie.Host)
		if err != nil {
			return err
		}
		domains, err := gd(c.Context())
		if err != nil {
			return err
		}
		return c.Render("view/cookie",
			fiber.Map{
				"cookies": cookies,
				"domains": domains,
				"cookie":  cookie,
				"ID":      id,
				"domain":  cookie.Host,
				"samesite": map[int64]string{
					0: "Нет", 1: "Lax", 2: "Strict"},
			})
	}
}
func MakeGetCookies(gc strg.GetCookies, gd strg.GetDomains) fiber.Handler {
	return func(c *fiber.Ctx) error {
		domain := c.Params("domain")
		cookies, err := gc(c.Context(), domain)
		if err != nil {
			return err
		}
		domains, err := gd(c.Context())
		if err != nil {
			return err
		}
		return c.Render("view/cookies",
			fiber.Map{
				"cookies": cookies,
				"domains": domains,
				"domain":  domain})
	}
}
func MakeGetDomains(gd strg.GetDomains) fiber.Handler {
	return func(c *fiber.Ctx) error {
		domains, err := gd(c.Context())
		if err != nil {
			return err
		}
		return c.Render("view/domains", fiber.Map{"domains": domains})
	}
}
