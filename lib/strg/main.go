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
	ID   int64  `db:"id"`
	Name string `db:"name"`
	IsXSS bool   `db:"is_xss"`
}
type Domain struct {
	Name  string `db:"host"`
	IsXSS bool   `db:"is_xss"`
}
type GetCookies func(ctx context.Context, domain string) ([]ShortCookie, error)
type GetCookie func(ctx context.Context, id int64) (Cookie, error)
type GetDomains func(ctx context.Context) ([]Domain, error)
