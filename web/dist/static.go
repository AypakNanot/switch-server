//go:build !exclude_webdist
// +build !exclude_webdist

package dist

import "embed"

//go:embed *
var WebFS embed.FS
