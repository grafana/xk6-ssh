package xk6ssh

import "go.k6.io/k6/js/modules"

func init() {
	modules.Register("k6/x/ssh", &K6SSH{})
}
