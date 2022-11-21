package model

type DhcpLease struct {
	IpV4    string `json:"ipV4" note:"IPv4地址"`
	Address string `json:"address" note:"MAC地址"`
	Comment string `json:"comment" note:"描述"`
}
