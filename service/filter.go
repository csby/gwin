package main

type Filter struct {
	Allow   bool   `json:"allow" note:"ture-运行; false-拒绝"`
	Address string `json:"address" note:"MAC地址"`
	Comment string `json:"comment" note:"描述"`
}

type FilterDelete struct {
	Address string `json:"address" note:"MAC地址"`
}

type FilterModify struct {
	Address string `json:"address" note:"要修改的MAC地址"`

	Filter Filter `json:"filter" note:"新的筛选器"`
}
