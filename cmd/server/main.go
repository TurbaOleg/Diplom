package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/TurbaOleg/Diplom/lib/profile/firefox"
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
	var connstr string
	connstr, err := utils.GetEnv("DB_CONNECTION")
	if err != nil {
		connstr, err = firefox.GetWindowsDbConnectionPath()
		if err != nil {
			return errors.Join(err)
		}
	}

	db, err := initDB(connstr)
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
