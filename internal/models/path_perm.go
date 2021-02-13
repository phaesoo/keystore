package models

type PathPermission struct {
	ID          int    `json:"id"`
	PathPattern string `json:"path_pattern"`
}
