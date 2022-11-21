package model

type OptDns struct {
	OptEnable

	ZoomNames []string `json:"zoomNames" note:"区域名称"`
}
