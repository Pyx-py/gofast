// 自动生成模板{{.StructName}}
package model

{{ if .ImportTime }}
import (
    "time"
)
{{ end }}

{{- if ne .GoStructString "" }}
{{.GoStructString}}
{{- end }}

{{ if .TableName}}
func ({{.StructName}}) TableName() string {
    return "{{.TableName}}"
}
{{ end }}

