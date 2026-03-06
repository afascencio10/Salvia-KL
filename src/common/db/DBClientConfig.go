package db

type DBClientConfig struct {
	Hostname     string `json:"hostname"`
	Port         string `json:"port"`
	DatabaseName string `json:"database"`
	UserName     string `json:"user"`
	Password     string `json:"password"`
}
