package model

type DnsRecord struct {
	Name string `json:"name" note:"记录名称"`
	Data string `json:"data" note:"记录数据"`
}

type DnsRecordArgument struct {
	ZoneName string `json:"zoneName" required:"true" note:"域名"`
}

type DnsRecordModify struct {
	DnsRecordArgument
	DnsRecord
}
