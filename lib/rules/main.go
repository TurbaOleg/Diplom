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
