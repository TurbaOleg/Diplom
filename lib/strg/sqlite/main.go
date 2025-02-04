package sqlite

import (
	"context"
	"fmt"

	"github.com/TurbaOleg/Diplom/lib/strg"
	"github.com/jmoiron/sqlx"

	_ "embed"
	_ "modernc.org/sqlite"
)

//go:embed init.sql
var initstmt string

func InitDB(db *sqlx.DB) (err error) {
	var cnt int
	err = db.Get(&cnt,
		`SELECT count(*) FROM sqlite_master WHERE type='table' AND name='rules';`,
	)
	if err != nil {
		return err
	}
	if cnt == 0 {
		_, err = db.Exec(initstmt)
		return err
	} else {
		return nil
	}
}
func MakeGetCookies(db *sqlx.DB) (strg.GetCookies, error) {
	stmt, err := db.Preparex(`
		select
		id, name, value like 'true' as is_xss
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
func MakeSetCookie(db *sqlx.DB) (strg.SetCookie, error) {
	stmt, err := db.PrepareNamed(`
		update moz_cookies
		set
			host = :host,
			isSecure = :isSecure,
			isHttpOnly = :isHttpOnly,
			sameSite = :sameSite,
			name = :name,
			value = :value,
			expiry = :expiry,
			path = :path,
			host = :host
		where id = :id
	`)
	if err != nil {
		return nil, err
	}
	type CookieWithId struct {
		strg.Cookie
		ID int64 `db:"id"`
	}
	return func(ctx context.Context, id int64, rule strg.Cookie) error {
		_, err := stmt.ExecContext(ctx, CookieWithId{rule, id})
		return err
	}, nil

}
func MakeDeleteCookie(db *sqlx.DB) (strg.DeleteCookie, error) {
	stmt, err := db.Prepare(`
		delete from moz_cookies
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
