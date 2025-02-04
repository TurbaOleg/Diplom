Модуль: lib/utils
утилиты, которые используются в различных частях программы
```go
package utils

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetEnv(k string) (string, error) {
	v, ok := os.LookupEnv(k)
	if !ok {
		return "", fmt.Errorf("undefined envar %s", k)
	}
	return v, nil
}
func GetIntParam(c *fiber.Ctx, p string) (int64, error) {
	ids, ok := c.AllParams()[p]
	if !ok {
		return 0, errors.Join(fiber.ErrBadRequest, fmt.Errorf("param %s missing\n", p))
	}

	id, err := strconv.ParseInt(ids, 10, 0)
	if err != nil {
		return 0, errors.Join(fiber.ErrBadRequest, fmt.Errorf("param %s has bad value %v\n", p, ids))
	}
	return id, nil
}
```

Модуль: lib/strg
определяет интерфейс абстрактного хранилища файлов cookie и правил

```go
package strg

import "context"

type Cookie struct {
	Value      string `db:"value"`
	Name       string `db:"name"`
	Host       string `db:"host"`
	Path       string `db:"path"`
	Expiry     int64  `db:"expiry"`
	IsSecure   bool   `db:"isSecure"`
	IsHttpOnly bool   `db:"isHttpOnly"`
	SameSite   int64  `db:"sameSite"`
	IsXSS      bool   `db:"is_xss"`
}
type ShortCookie struct {
	ID    int64  `db:"id"`
	Name  string `db:"name"`
	IsXSS bool   `db:"is_xss"`
}
type Domain struct {
	Name  string `db:"host"`
	IsXSS bool   `db:"is_xss"`
}

// все домены
type GetDomains func(ctx context.Context) ([]Domain, error)

// все куки
type GetCookies func(ctx context.Context, domain string) ([]ShortCookie, error)

// один куки в подробностях
type GetCookie func(ctx context.Context, id int64) (Cookie, error)
package strg

import (
	"context"
)

type Rule struct {
	// Priority      int
	ID            int64  `db:"id"`
	DomainPattern string `db:"domain_pattern" form:"host"`
	IsSecure      bool   `db:"is_secure" form:"is_secure"`
	IsHttpOnly    bool   `db:"is_http_only" form:"is_http_only"`
	SameSite      int64  `db:"same_site" form:"same_site"`
}

// получить все правила
type GetRules func(ctx context.Context) ([]Rule, error)
// получить одно правило
type GetRule func(ctx context.Context, id int64) (Rule, error)
// удалить правило
type DeleteRule func(ctx context.Context, id int64) error
// добавить правило
type NewRule func(ctx context.Context, rule Rule) (int64, error)
// обновить правило
type SetRule func(ctx context.Context, id int64, rule Rule) error
```
Модуль: lib/rules
содержит в себе логику работы правил
```go
package rules

import (
	"context"

	"github.com/TurbaOleg/Diplom/lib/strg"
	"github.com/jmoiron/sqlx"
)

type ApplyRules func(ctx context.Context, rules []strg.Rule) error

// TODO: use some cool JOIN to do it in one expr
func MakeApplyRules(db *sqlx.DB) ApplyRules {
	return func(ctx context.Context, rules []strg.Rule) error {
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		defer tx.Rollback() 
		for _, rule := range rules {
			_, err := tx.ExecContext(ctx,
				`update moz_cookies
				set isSecure = ?,
					isHttpOnly = ?,
					sameSite = ?
				where host like ?`, rule.IsSecure, rule.IsHttpOnly, rule.SameSite, rule.DomainPattern)
			if err != nil {
				return err
			}
		}
		return tx.Commit()
	}
}
```
Модуль: lib/strg/sqlite

содержит имплементацию интерфейса strg для БД SQLite, используемой
для хранения cookies в браузере Firefox
```go
package sqlite

import (
	"context"
	"fmt"

	"github.com/TurbaOleg/Diplom/lib/strg"
	"github.com/jmoiron/sqlx"

	_ "modernc.org/sqlite"
)

func MakeGetCookies(db *sqlx.DB) (strg.GetCookies, error) {
	stmt, err := db.Preparex(`
		select
		id, name, value like '%x3c%' or value like '%<script>%' as is_xss
		from moz_cookies where host = ?
		order by is_xss desc`)
	if err != nil {
		return nil, err
	}
	return func(ctx context.Context, domain string) (out []strg.ShortCookie, err error) {
		if false {
			fmt.Println(domain)
		}
		err = stmt.SelectContext(ctx, &out, domain)
		return
	}, nil
}
func MakeGetDomains(db *sqlx.DB) (strg.GetDomains, error) {
	stmt, err := db.Preparex(`
		select distinct (host) host, max(value like 'true') OVER (partition by host) as is_xss
		from moz_cookies
		order by is_xss desc`)
	if err != nil {
		return nil, err
	}
	return func(ctx context.Context) (out []strg.Domain, err error) {
		err = stmt.SelectContext(ctx, &out)
		return
	}, nil
}
func MakeGetCookie(db *sqlx.DB) (strg.GetCookie, error) {
	stmt, err := db.Preparex(`
		select
		host, value, name, path, expiry, isSecure, isHttpOnly, sameSite
		from moz_cookies where id = ? limit 1`)
	if err != nil {
		return nil, err
	}
	return func(ctx context.Context, id int64) (out strg.Cookie, err error) {
		err = stmt.GetContext(ctx, &out, id)
		return
	}, nil
}
package sqlite

import (
	"context"

	"github.com/TurbaOleg/Diplom/lib/strg"
	"github.com/jmoiron/sqlx"
)

func MakeGetRules(db *sqlx.DB) (strg.GetRules, error) {
	stmt, err := db.Preparex(`
		select id, domain_pattern, is_secure, is_http_only, same_site from rules
	`)
	if err != nil {
		return nil, err
	}
	return func(ctx context.Context) (out []strg.Rule, err error) {
		err = stmt.SelectContext(ctx, &out)
		return
	}, nil
}
func MakeGetRule(db *sqlx.DB) (strg.GetRule, error) {
	stmt, err := db.Preparex(
		`select id, domain_pattern, is_secure, is_http_only, same_site from rules where id = ?`)
	if err != nil {
		return nil, err
	}
	return func(ctx context.Context, id int64) (out strg.Rule, err error) {
		err = stmt.GetContext(ctx, &out, id)
		return
	}, nil
}
func MakeNewRule(db *sqlx.DB) (strg.NewRule, error) {
	stmt, err := db.PrepareNamed(`
		insert into rules (domain_pattern, is_secure, is_http_only, same_site)
		values(:domain_pattern, :is_secure, :is_http_only, :same_site)
	`)
	if err != nil {
		return nil, err
	}
	return func(ctx context.Context, rule strg.Rule) (int64, error) {
		res, err := stmt.ExecContext(ctx, rule)
		if err != nil {
			return 0, err
		}
		return res.LastInsertId()
	}, nil

}
func MakeSetRule(db *sqlx.DB) (strg.SetRule, error) {
	stmt, err := db.PrepareNamed(`
		update rules
		set
			domain_pattern = :domain_pattern,
			is_secure = :is_secure,
			is_http_only = :is_http_only,
			same_site = :same_site
		where id = :id
	`)
	if err != nil {
		return nil, err
	}
	type RuleWithId struct {
		strg.Rule
		ID int64 `db:"id"`
	}
	return func(ctx context.Context, id int64, rule strg.Rule) error {
		_, err := stmt.ExecContext(ctx, RuleWithId{rule, id})
		return err
	}, nil

}
func MakeDeleteRule(db *sqlx.DB) (strg.DeleteRule, error) {
	stmt, err := db.Prepare(`
		delete from rules
		where id = ?
	`)
	if err != nil {
		return nil, err
	}
	return func(ctx context.Context, id int64) error {
		_, err := stmt.ExecContext(ctx, id)
		return err
	}, nil

}
```
Модуль: lib/www

определяет пользовательский интерфейс приложения
на веб-технологиях
```go
package www

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"

	"github.com/TurbaOleg/Diplom/lib/rules"
	"github.com/TurbaOleg/Diplom/lib/strg"
	"github.com/TurbaOleg/Diplom/lib/strg/sqlite"
	"github.com/TurbaOleg/Diplom/lib/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/jmoiron/sqlx"
)

//go:embed view/*
var viewfs embed.FS

var SAMESITE = map[int64]string{
	0: "Нет", 1: "Lax", 2: "Strict"}

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
	grs, err := sqlite.MakeGetRules(db)
	if err != nil {
		return nil, err
	}
	nr, err := sqlite.MakeNewRule(db)
	if err != nil {
		return nil, err

	}
	gr, err := sqlite.MakeGetRule(db)
	if err != nil {
		return nil, err
	}
	dr, err := sqlite.MakeDeleteRule(db)
	if err != nil {
		return nil, err
	}
	sr, err := sqlite.MakeSetRule(db)
	if err != nil {
		return nil, err
	}
	ar := rules.MakeApplyRules(db)

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

	app.Get("/rules", MakeGetRules(grs))
	app.Post("/rules", MakePostRules(nr))
	app.Post("/rules/apply", MakeApplyRules(ar, grs))
	app.Get("/rule/:id", MakeGetRule(gr))
	app.Post("/rule/:id", MakePostRule(sr))
	app.Post("/rule/:id/delete", MakeDeleteRule(dr))
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
				"cookies":  cookies,
				"domains":  domains,
				"cookie":   cookie,
				"ID":       id,
				"domain":   cookie.Host,
				"samesite": SAMESITE,
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
func MakeGetRules(lsr strg.GetRules) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rules, err := lsr(c.Context())
		if err != nil {
			return err
		}
		return c.Render("view/rules", fiber.Map{"rules": rules})
	}
}
func MakeGetRule(gr strg.GetRule) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := utils.GetIntParam(c, "id")
		if err != nil {
			return err
		}
		rule, err := gr(c.Context(), id)
		if err != nil {
			return err
		}
		return c.Render("view/rule", fiber.Map{
			"rule":     rule,
			"samesite": SAMESITE,
		})
	}
}
func MakePostRule(sr strg.SetRule) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := utils.GetIntParam(c, "id")
		if err != nil {
			return err
		}
		var in strg.Rule
		if err = c.BodyParser(&in); err != nil {
			return err
		}
		err = sr(c.Context(), id, in)
		if err != nil {
			return err
		}
		return c.Redirect(fmt.Sprintf("/rule/%d", id))
	}
}
func MakeApplyRules(ar rules.ApplyRules, grs strg.GetRules) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rs, err := grs(c.Context())
		if err != nil {
			return err
		}
		if err := ar(c.Context(), rs); err != nil {
			return err
		}
		return c.Redirect("/rules", fiber.StatusSeeOther)
	}
}
func MakePostRules(nr strg.NewRule) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var in strg.Rule
		if err := c.BodyParser(&in); err != nil {
			return err
		}
		id, err := nr(c.Context(), in)
		if err != nil {
			return err
		}
		return c.Redirect(fmt.Sprintf("/rule/%d", id))
	}
}
func MakeDeleteRule(dr strg.DeleteRule) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := utils.GetIntParam(c, "id")
		if err != nil {
			return err
		}
		var in strg.Rule
		if err = c.BodyParser(&in); err != nil {
			return err
		}
		err = dr(c.Context(), id)
		if err != nil {
			return err
		}
		return c.Redirect("/rules", fiber.StatusSeeOther)
	}
}
```
Модуль: cmd/server

Входная точка программы, которая собирает все модули
для совместной работы. Этот модуль компилируется в исполняемый файл
```go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/TurbaOleg/Diplom/lib/utils"
	"github.com/TurbaOleg/Diplom/lib/www"
	"github.com/jmoiron/sqlx"
	"github.com/skratchdot/open-golang/open"
)

func main() {
	if err := f(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func f() error {
	connstr, err := utils.GetEnv("DB_CONNECTION")
	if err != nil {
		return err
	}

	db, err := initDB(connstr)
	// db.MapperFunc(strings.ToLower)
	if err != nil {
		return err
	}
	app, err := www.MakeApp(db)
	if err != nil {
		return err
	}
	go func() {
		time.Sleep(time.Second)
		open.Start("http://localhost:8000")
	}()
	return app.Listen(":8000")

}
func initDB(connstr string) (*sqlx.DB, error) {
	return sqlx.Connect("sqlite", connstr)
}
```
