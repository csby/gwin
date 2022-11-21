package config

type MsAd struct {
	Host     string `json:"host" note:"主机地址"`
	Port     int    `json:"port" note:"端口, 389或636"`
	Base     string `json:"base" note:"根路径，如: DC=example,DC=com"`
	Account  string `json:"account" note:"账号，全路径，如: CN=Administrator,CN=Users,DC=example,DC=com"`
	Password string `json:"password" note:"账号密码"`
}
