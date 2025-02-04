package main

import (
	"fmt"
	"os"
	"time"

	"codeberg.org/shinyzero0/oleg-soul-2024/lib/utils"
	"codeberg.org/shinyzero0/oleg-soul-2024/lib/www"
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
