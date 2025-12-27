//go:build !sqlite3
// +build !sqlite3

package database

import (
	"github.com/glebarez/sqlite" // 纯 Go SQLite 驱动，无需 CGO
)

// 使用纯 Go 的 SQLite 驱动
func init() {
	opens["sqlite3"] = sqlite.Open
}
