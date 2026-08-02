// Package migrations exposes the Authorizer audit migrations as an immutable
// filesystem embedded in the migrator binary.
package migrations

import "embed"

// Files contains every versioned SQL migration in this directory.
//
//go:embed files/*.sql
var Files embed.FS
