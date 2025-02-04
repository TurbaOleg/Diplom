package strg

import "context"

type Cookie struct {
	Value      string `db:"value" form:"value"`
	Name       string `db:"name" form:"name"`
	Host       string `db:"host" form:"host"`
	Path       string `db:"path" form:"path"`
	Expiry     int64  `db:"expiry" form:"expiry"`
	IsSecure   bool   `db:"isSecure" form:"is_secure"`
	IsHttpOnly bool   `db:"isHttpOnly" form:"is_http_only"`
	SameSite   int64  `db:"sameSite" form:"same_site"`
	IsXSS      bool   `db:"is_xss" form:"is_xss"`
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

type SetCookie func(ctx context.Context, id int64, c Cookie) error

type DeleteCookie func(ctx context.Context, id int64) error
