package ddb

import (
	"fmt"
	"golang-migrator-example/src/migrator"
)

// Provide returns proper migration definitions or an error when not found.
func Provide(o *migrator.DefaultDDBProviderOptions) ([]migrator.Definition, error) {

	switch o.MigrationSet {
	case "example":
		return defsExample(o.Ctx, o.Conf, o.DB), nil
	default:
		return nil, fmt.Errorf("unknown migration set: %s", o.MigrationSet)
	}
}
