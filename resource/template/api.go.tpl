package v1

import (
    {{- if ne .LogPath ""}}
    "github.com/pyx-py/gofast/global"
    {{- else}}
    /* "github.com/pyx-py/gofast/global" */
    {{- end}}
    "github.com/gin-gonic/gin"
    "{{.ModuleName}}/model/request"
    "{{.ModuleName}}/model/response"
    "{{.ModuleName}}/model"
    "{{.ModuleName}}/service"
    {{- if .LogPath}}
    "go.uber.org/zap"
    {{- end}}
)

// @Tags {{.StructName}}
// @Summary 创建{{.StructName}}
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.{{.StructName}} true "创建{{.StructName}}"
// @Success 200 {string} string "{"success":true,"data":{}, "msg":"获取成功"}"
// @Router /{{.Abbreviation}}/create{{.StructName}} [post]
func Create{{.StructName}}(c *gin.Context) {
    var {{.Abbreviation}} model.{{.StructName}}
    _ = c.ShouldBindJSON(&{{.Abbreviation}})
    if err := service.Create{{.StructName}}({{.Abbreviation}}); err != nil {
    {{- if ne .LogPath ""}}
        global.GF_LOG.Error("创建失败！", zap.Any("err", err))
    {{- else}}
        /* global.GF_LOG.Error("创建失败！", zap.Any("err", err)) */
    {{- end}}
        response.FailWithMessage("创建失败", c)
    } else {
        response.OkWithMessage("创建成功", c)
    }
}

// Tags {{.StructName}}
// @Summary 删除{{.StructName}}
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.{{.StructName}} true "删除{{.StructName}}"
// @Success 200 {string} string "{"success":true,"data":{}, "msg":"删除成功"}"
// @Router /{{.Abbreviation}}/delete{{.StructName}} [delete]
func Delete{{.StructName}}(c *gin.Context) {
    var {{.Abbreviation}} model.{{.StructName}}
    _ = c.ShouldBindJSON(&{{.Abbreviation}})
    if err := service.Delete{{.StructName}}({{.Abbreviation}}); err != nil {
    {{- if ne .LogPath ""}}
        global.GF_LOG.Error("删除失败！", zap.Any("err", err))
    {{- else}}
        /* global.GF_LOG.Error("删除失败！", zap.Any("err", err)) */
    {{- end}}
        response.FailWithMessage("删除失败", c)
    } else {
        response.OkWithMessage("删除成功", c)
    }
}


// Tags {{.StructName}}
// @Summary 更新{{.StructName}}
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.{{.StructName}} true "更新{{.StructName}}"
// @Success 200 {string} string "{"success":true,"data":{}, "msg":"更新成功"}"
// @Router /{{.Abbreviation}}/update{{.StructName}} [put]
func Update{{.StructName}}(c *gin.Context) {
    var {{.Abbreviation}} model.{{.StructName}}
    _ = c.ShouldBindJSON(&{{.Abbreviation}})
    if err := service.Update{{.StructName}}({{.Abbreviation}}); err != nil {
    {{- if ne .LogPath ""}}
        global.GF_LOG.Error("更新失败", zap.Any("err", err))
    {{- else}}
        /* global.GF_LOG.Error("更新失败", zap.Any("err", err)) */
    {{- end}}
        response.FailWithMessage("更新失败", c)
    } else {
        response.OkWithMessage("更新成功", c)
    }
}


// Tags {{.StructName}}
// @Summary 分页获取{{.StructName}}列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.{{.StructName}}Search true "分页获取{{.StructName}}列表"
// @Success 200 {string} string "{"success":true,"data":{}, "msg":"获取成功"}"
// @Router /{{.Abbreviation}}/get{{.StructName}}List [get]
func Get{{.StructName}}List(c *gin.Context) {
    var pageInfo request.{{.StructName}}Search
    _ = c.ShouldBindJSON(&pageInfo)
    if err, list, total := service.Get{{.StructName}}InfoList(pageInfo); err != nil {
    {{- if ne .LogPath ""}}
        global.GF_LOG.Error("获取失败！", zap.Any("err", err))
    {{- else}}
        /* global.GF_LOG.Error("获取失败！", zap.Any("err", err)) */
    {{- end}}
        response.FailWithMessage("获取失败", c)
    } else {
        response.OKWithDetailed(response.PageResult{
            List:   list,
            Total:  total,
            Page:   pageInfo.Page,
            PageSize: pageInfo.PageSize,
        }, "获取成功", c)
    }
}