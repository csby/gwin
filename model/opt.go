package model

type Opt struct {
	Dhcp OptEnable `json:"dhcp" note:"DHCP设置"`
	Dns  OptDns    `json:"dns" note:"DNS设置"`
	Svn  OptEnable `json:"svn" note:"SVN设置"`
}
