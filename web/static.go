//go:build !exclude_webdist
// +build !exclude_webdist

package web

import "embed"

//go:embed dist/*
var WebFS embed.FS
