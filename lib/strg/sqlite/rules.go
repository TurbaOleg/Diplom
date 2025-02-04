package sqlite

import (
	"context"

	"codeberg.org/shinyzero0/oleg-soul-2024/lib/strg"
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
