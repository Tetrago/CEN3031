package user

import (
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/tetrago/motmot/api/internal/crypt"
	"github.com/tetrago/motmot/api/internal/globals"

	"github.com/tetrago/motmot/api/.gen/motmot/public/model"
	. "github.com/tetrago/motmot/api/.gen/motmot/public/table"
)

func MakeIdentifier() (string, error) {
	var dest model.UserAccount

generate:
	ident, err := crypt.GenerateBase64(16)
	if err != nil {
		return "", err
	}

	stmt := SELECT(UserAccount.Identifier).FROM(UserAccount).WHERE(UserAccount.Identifier.EQ(String(ident)))

	if err := stmt.Query(globals.Database, &dest); err == qrm.ErrNoRows {
		return ident, nil
	} else if err != nil {
		return "", err
	} else {
		goto generate
	}
}
