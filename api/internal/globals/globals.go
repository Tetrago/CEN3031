package globals

import (
	"database/sql"

	"github.com/tetrago/motmot/api/internal/options"
)

var Opts = options.LoadFromEnvironment()
var Database *sql.DB
