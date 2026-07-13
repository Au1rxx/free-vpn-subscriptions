// Package migrations embeds the versioned MySQL schema files.
package migrations

import "embed"

// Files contains this directory. The migration loader selects only *.sql.
//
//go:embed *.sql
var Files embed.FS
