package config

type Dns struct {
	Enable bool `json:"enable"`

	ZoomNames []string `json:"zoomNames"`
}
