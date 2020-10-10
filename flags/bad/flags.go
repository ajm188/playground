package bad

import (
	"flag"
)

var (
	x = flag.String("some-flag", "", "defined in playground/flags/bad/flags.go")
	y = flag.Int("some-other-flag", 0, "defined in playground/flags/bad/flags.go")
)
