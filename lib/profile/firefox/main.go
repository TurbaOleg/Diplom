package firefox

import (
	"fmt"
	"path"
	"strings"

	"github.com/adrg/xdg"
	"github.com/gofiber/fiber/v2"
	"gopkg.in/ini.v1"
)

func GetWindowsDbConnectionPath() (string, error) {
	p, err := xdg.SearchDataFile("Firefox/profiles.ini")
	if err != nil {
		return "", err
	}
	tree, err := ini.Load(p)
	if err != nil {
		return "", err
	}
	// var dprofile string
	for _, sec := range tree.Sections() {
		if !strings.HasPrefix(sec.Name(), "Profile") {
			continue
		}
		if sec.Key("Default").Value() == "1" {
			return fmt.Sprintf("%s/Profiles/%s", path.Dir(p)+sec.Key("Path").String()), nil
		}
	}
	return "", fiber.ErrNotFound
}
