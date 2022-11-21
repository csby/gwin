package model

type SvnRepositoryItem struct {
	Id         string `json:"id" note:"标识ID"`
	Repository string `json:"repository" note:"存储库名称"`
	Name       string `json:"name" note:"项目名称"`
	Path       string `json:"path" note:"路径"`
	Type       int    `json:"type" note:"类型: 0-存储库; 1-文件夹; 2-文件"`
	Url        string `json:"url" note:"地址"`
	Revisions  int    `json:"revisions" note:"修订次数"`

	Children []*SvnRepositoryItem `json:"children" note:"子项"`
}
