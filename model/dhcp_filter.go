package model

type DhcpFilter struct {
	Allow   bool   `json:"allow" note:"ture-运行; false-拒绝"`
	Address string `json:"address" note:"MAC地址"`
	Comment string `json:"comment" note:"描述"`
}

type DhcpFilterDelete struct {
	Address string `json:"address" note:"MAC地址"`
}

type DhcpFilterModify struct {
	Address string `json:"address" note:"要修改的MAC地址"`

	Filter DhcpFilter `json:"filter" note:"新的筛选器"`
}
