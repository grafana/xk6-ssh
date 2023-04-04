// Package xk6ssh registers the xk6-ssh javascript extension
package xk6ssh

import (
	"github.com/spf13/afero"
	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/ssh", &K6SSH{fs: afero.NewOsFs()})
}
