package controller

import (
	"github.com/hth0919/migrationcontroller/pkg/controller/migration"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, migration.Add)
}
