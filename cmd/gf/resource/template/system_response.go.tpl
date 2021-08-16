package response

import "{{.ModuleName}}/config"

type SysConfigResponse struct {
	Config config.Server `json:"config"`
}
