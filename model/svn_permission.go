package model

type SvnPermission struct {
	AccountId   string `json:"accountId" note:"账号ID"`
	AccountName string `json:"accountName" note:"账号名称"`
	AccessLevel int    `json:"accessLevel" note:"访问权限: 0-无; 1-只读; 2-读写"`
	Inherited   bool   `json:"inherited" note:"是否继承"`
}

type SvnPermissionUser struct {
	Repository  string `json:"repository" note:"存储库名称"`
	Path        string `json:"path" note:"路径"`
	AccessLevel int    `json:"accessLevel" note:"访问权限: 0-无; 1-只读; 2-读写"`
}

type SvnPermissionUserArgument struct {
	AccountId string `json:"accountId" required:"true" note:"账号ID"`
}

type SvnPermissionArgument struct {
	Repository string `json:"repository" required:"true" note:"存储库名称"`
	Path       string `json:"path" required:"true" note:"路径"`
	AccountId  string `json:"accountId" required:"true" note:"账号ID"`
}

type SvnPermissionArgumentEdit struct {
	SvnPermissionArgument
	AccessLevel int `json:"accessLevel" note:"访问权限: 0-无; 1-只读; 2-读写"`
}
