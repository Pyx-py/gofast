package request

import "{{.ModuleName}}/model"

type {{.StructName}}Search struct {
    model.{{.StructName}}
    PageInfo
}
