package v1

type GroupModel struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type UserModel struct {
	ID          int64        `json:"id"`
	DisplayName string       `json:"display_name"`
	Groups      []GroupModel `json:"groups"`
}
