package ssh

import (
	"github.com/loadimpact/k6/js/modules"
)

const version = "v0.0.1"

func init() {
	modules.Register("k6/x/ssh", &SSH{
		Version: version,
	})
}

// SSH is the main export of k6 docker extension
type SSH struct {
	Version string
}
