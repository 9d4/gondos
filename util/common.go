package util

import (
	"os"
	"strings"
)

func IsDevel() bool {
	return strings.ToLower(os.Getenv("ENV")) == "development"

}
