package v1

import "github.com/tetrago/motmot/api/.gen/motmot/public/model"

type TokenModel struct {
	Token string `json:"token"`
}

type GroupModel struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type UserModel struct {
	Identifier  string       `json:"identifier"`
	DisplayName string       `json:"display_name"`
	Groups      []GroupModel `json:"groups"`
}

func MapToGroupModel(x model.Room, _ int) GroupModel {
	return GroupModel{x.ID, x.Name}
}
