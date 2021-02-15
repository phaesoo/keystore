package models

// AuthKey contains API Key base user auth info
type AuthKey struct {
	ID        int    `json:"id"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	UserUUID  string `json:"user_uuid"`
}
