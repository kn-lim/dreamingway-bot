package pixelmon

import (
	"fmt"
	"os"
)

var (
	ServerURL = fmt.Sprintf("%v.%v", os.Getenv("PIXELMON_SUBDOMAIN"), os.Getenv("PIXELMON_DOMAIN"))
)
