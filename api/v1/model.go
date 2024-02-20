package v1

type GroupModel struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type UserModel struct {
	Identifier  string       `json:"identifier"`
	DisplayName string       `json:"display_name"`
	Groups      []GroupModel `json:"groups"`
}
