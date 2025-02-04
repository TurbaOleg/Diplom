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
