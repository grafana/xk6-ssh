package xk6ssh

import "github.com/loadimpact/k6/js/modules"

const version = "v0.0.1"

func init() {
	modules.Register("k6/x/ssh", &K6SSH{Version: version})
}
