package db

import "embed"

// Migrations contains the service-owned schema.
//
//go:embed migrations/*.sql
var Migrations embed.FS
