package envfile

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

func Load(paths ...string) error {
	if len(paths) == 0 {
		paths = []string{".env"}
	}

	for _, path := range paths {
		if err := godotenv.Overload(path); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	return nil
}
