package model

const (
	SvnRepositoryItemKindRepository = 0
	SvnRepositoryItemKindFolder     = 1
	SvnRepositoryItemKindFile       = 2
)

const (
	SvnPermissionNoAccess  = 0
	SvnPermissionReadOnly  = 1
	SvnPermissionReadWrite = 2
)

type SvnRepository struct {
	Id   string `json:"id" note:"标识ID"`
	Name string `json:"name" required:"true" note:"名称"`
	Path string `json:"path" required:"true" note:"路径"`
}

type SvnRepositoryNew struct {
	Name string `json:"name" required:"true" note:"名称"`
}
