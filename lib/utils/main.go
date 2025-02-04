package utils

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetEnv(k string) (string, error) {
	v, ok := os.LookupEnv(k)
	if !ok {
		return "", fmt.Errorf("undefined envar %s", k)
	}
	return v, nil
}
func GetIntParam(c *fiber.Ctx, p string) (int64, error) {
	ids, ok := c.AllParams()[p]
	if !ok {
		return 0, errors.Join(fiber.ErrBadRequest, fmt.Errorf("param %s missing\n", p))
	}

	id, err := strconv.ParseInt(ids, 10, 0)
	if err != nil {
		return 0, errors.Join(fiber.ErrBadRequest, fmt.Errorf("param %s has bad value %v\n", p, ids))
	}
	return id, nil
}
