package firefox

import (
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/gofiber/fiber/v2"
	"gopkg.in/ini.v1"
)

func GetWindowsDbConnectionPath() (string, error) {
	p, err := xdg.SearchDataFile("Mozilla/Firefox/profiles.ini")
	if err != nil {
		return "", err
	}
	tree, err := ini.Load(p)
	if err != nil {
		return "", err
	}
	// var dprofile string
	for _, sec := range tree.Sections() {
		if !strings.HasPrefix(sec.Name(), "Install") {
			continue
		}
		if sec.HasKey("Default") {
			return filepath.Join(
				filepath.Dir(p),
				"Profiles",
				sec.Key("Default").Value(),
				"cookies.sqlite"), nil
		}
	}
	return "", fiber.ErrNotFound
}
