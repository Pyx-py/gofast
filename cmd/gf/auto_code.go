package main

import (
	"bufio"
	"fmt"
	"go/build"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gobuffalo/packr/v2"
	"github.com/pyx-py/gofast/utils"
	"github.com/xwb1989/sqlparser"
)

const GOFAST = "github.com/pyx-py/gofast"

// var staticBox = packr.New("sBox", "./resource/static")
// var templateBox = packr.New("tBox", "./resource/template")

var staticBox *packr.Box
var templateBox *packr.Box

var staticFileMap = make(map[string]string)
var templateFileMap = make(map[string]string)

var staticFileList = make([]string, 0)
var templateFileList = make([]string, 0)

var (
	SqlGoTypeMap = map[string]string{
		"int":                "int",
		"integer":            "int",
		"tinyint":            "int8",
		"smallint":           "int16",
		"mediumint":          "int32",
		"bigint":             "int64",
		"int unsigned":       "uint",
		"integer unsigned":   "uint",
		"tinyint unsigned":   "uint8",
		"smallint unsigned":  "uint16",
		"mediumint unsigned": "uint32",
		"bigint unsigned":    "uint64",
		"bit":                "byte",
		"bool":               "bool",
		"enum":               "string",
		"set":                "string",
		"varchar":            "string",
		"char":               "string",
		"tinytext":           "string",
		"mediumtext":         "string",
		"text":               "string",
		"longtext":           "string",
		"blob":               "string",
		"tinyblob":           "string",
		"mediumblob":         "string",
		"longblob":           "string",
		"date":               "time.Time",
		"datetime":           "time.Time",
		"timestamp":          "time.Time",
		"time":               "time.Time",
		"float":              "float64",
		"double":             "float64",
		"decimal":            "float64",
		"binary":             "string",
		"varbinary":          "string",
	}
)

// 初始版本自动化代码工具
type AutoCoder struct {
	ModuleName          string        `json:"moduleName"`  // 用户的module名称
	ProjectPath         string        `json:"projectPath"` // 用户的项目路径
	ProjectName         string        `json:"projectName"` // 用户的项目名称
	LogPath             string        `json:"logPath"`     // 日志的存放路径，为空代表不启用默认日志
	TplPath             string        `json:"tplPath"`     // 模板文件dir
	StructName          string        `json:"structName"`
	TableName           string        `json:"tableName"`
	Abbreviation        string        `json:"abbreviation"`
	ImportTime          bool          `json:"importTime"`
	GoStructString      template.HTML `json:"goStructString"`
	Fields              []*Field      `json:"fields"`
	ColSearchTypeString string        `json:"colSearchTypeString"`
	SqlPath             string        `json:"sqlPath"` // 传入的sql文件路径
	// GoFastPath          string        `json:"goFastPath"` //

}

type Field struct {
	FieldName       string `json:"fieldName"`
	FieldType       string `json:"fieldType"`
	ColumnName      string `json:"columnName"`
	FieldSearchType string `json:"fieldSearchType"`
}

type tplData struct {
	template *template.Template
	// locationPath string
	autoCodePath string
	repeat       bool
}

func initData(gfPath string) error {
	if gfPath == "" {
		path, err := getGoFastPath()
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		gfPath = path
	}
	initBox(gfPath)
	if err := getAllResourceFileName(gfPath + "/cmd/gf/resource"); err != nil {
		fmt.Println(err.Error())
		return err
	}
	if err := getAllResourceFile(); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func NewAutoCoder(projectPath, moduleName, sqlFilePath, colSearchTypeString, logPath, gofastPath string) (*AutoCoder, error) {
	if projectPath == "" || moduleName == "" {
		return nil, fmt.Errorf("projectPath or moduleName can not be null")
	}
	projectNames := strings.Split(projectPath, "/")
	projectName := projectNames[len(projectNames)-1]
	// init aurocoder
	autoCoder := &AutoCoder{
		ProjectName: projectName,
		TplPath:     "./resource/template",
		ImportTime:  false,
		Fields:      make([]*Field, 0),
	}
	// autoCoder.GoFastPath = gofastPath
	if err := initData(gofastPath); err != nil {
		return nil, err
	}
	autoCoder.SqlPath = sqlFilePath
	autoCoder.ColSearchTypeString = colSearchTypeString
	colSearchTypeMap := handleSearchMap(colSearchTypeString)
	autoCoder.LogPath = logPath
	autoCoder.ModuleName = moduleName
	if strings.HasSuffix(projectPath, "/") {
		autoCoder.ProjectPath = projectPath[0 : len(projectPath)-1]
	} else {
		autoCoder.ProjectPath = projectPath
	}

	// init sql
	if sqlFilePath != "" {
		sqlData, err := ioutil.ReadFile(sqlFilePath)
		if err != nil {
			fmt.Printf("[error]:%s\n", err.Error())
			return nil, err
		}
		sql := string(sqlData)
		// parse sql
		statement, err := sqlparser.ParseStrictDDL(sql)
		if err != nil {
			return nil, err
		}
		stmt, ok := statement.(*sqlparser.DDL)
		if !ok {
			return nil, fmt.Errorf("sql is not a create statement")
		}
		tableName := stmt.NewName.Name.String()
		autoCoder.TableName = tableName
		primaryIdxMap, uniqueIdxMap, idxMap := buildIdxMaps(stmt.TableSpec.Indexes)
		builder := strings.Builder{}
		structName := snakeToCamel(tableName)
		autoCoder.StructName = structName
		autoCoder.Abbreviation = structName
		structStart := fmt.Sprintf("type %s struct { \n", structName)
		builder.WriteString(structStart)
		for _, col := range stmt.TableSpec.Columns {
			columnType := col.Type.Type
			if col.Type.Unsigned {
				columnType += " unsigned"
			}
			goType := SqlGoTypeMap[columnType]
			if goType == "time.Time" {
				autoCoder.ImportTime = true
			}
			fieldName := snakeToCamel(col.Name.String())
			searchType := colSearchTypeMap[col.Name.String()]
			if colSearchTypeString == "" {
				searchType = "="
			}
			oneField := &Field{
				FieldName:       fieldName,
				FieldType:       goType,
				ColumnName:      col.Name.String(),
				FieldSearchType: searchType,
			}
			autoCoder.Fields = append(autoCoder.Fields, oneField)

			gormStr := buildGormStr(col, primaryIdxMap, uniqueIdxMap, idxMap)
			buildGormStr(col, primaryIdxMap, uniqueIdxMap, idxMap)
			builder.WriteString(fmt.Sprintf("\t%s\t%-25s ", fieldName, goType) + "`json:" + `"` + col.Name.String() + `"` + " gorm:" + `"` + gormStr + `"` + "`\n")
		}
		builder.WriteString("}\n")
		autoCoder.GoStructString = template.HTML(builder.String())
	}

	return autoCoder, nil
}

func initBox(gfPath string) {
	staticBox = packr.New("sBox", gfPath+"/cmd/gf/resource/static")
	templateBox = packr.New("tBox", gfPath+"/cmd/gf/resource/template")
}

func getGoFastPath() (string, error) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
		if gopath == "" {
			return "", fmt.Errorf("can not find gopath, need manual incoming gofast path parameter")
		}
	}
	pathList := strings.Split(GOFAST, "/")
	subPathList := pathList[0 : len(pathList)-1]
	lastPath := pathList[len(pathList)-1]
	prePath := gopath + "/pkg/mod"
	for _, p := range subPathList {
		prePath += "/" + p
	}
	dirInfo, err := ioutil.ReadDir(prePath)
	if err != nil {
		return "", fmt.Errorf("get dir info error:%s", err.Error())
	}
	var gofastPath string
	for _, d := range dirInfo {
		if strings.Contains(d.Name(), lastPath) {
			gofastPath = prePath + "/" + d.Name()
			break
		}
	}
	return gofastPath, nil
}

func handleSearchMap(colSearchMapString string) map[string]string {
	var searchMap = make(map[string]string)
	if colSearchMapString == "" {
		return searchMap
	}

	mapList := strings.Split(colSearchMapString, "#")
	for _, m := range mapList {
		kv := strings.Split(m, ":")
		searchMap[strings.Trim(kv[0], " ")] = strings.Trim(kv[1], " ")
	}
	return searchMap
}

// inner funcs
func buildIdxMaps(indexList []*sqlparser.IndexDefinition) (primaryIdxMap map[string]string, uniqueIdxMap map[string]string, idxMap map[string]string) {
	primaryIdxMap = make(map[string]string)
	uniqueIdxMap = make(map[string]string)
	idxMap = make(map[string]string)
	for idx, _ := range indexList {
		if indexList[idx].Info.Primary {
			for cIdx, _ := range indexList[idx].Columns {
				primaryIdxMap[indexList[idx].Columns[cIdx].Column.String()] = indexList[idx].Info.Name.String()
			}
		} else if indexList[idx].Info.Unique {
			for cIdx, _ := range indexList[idx].Columns {
				uniqueIdxMap[indexList[idx].Columns[cIdx].Column.String()] = indexList[idx].Info.Name.String()
			}
		} else {
			for cIdx, _ := range indexList[idx].Columns {
				idxMap[indexList[idx].Columns[cIdx].Column.String()] = indexList[idx].Info.Name.String()
			}
		}
	}
	return
}

func buildGormStr(col *sqlparser.ColumnDefinition, primaryIdxMap map[string]string, uniqueIdxMap map[string]string, idxMap map[string]string) string {
	builder := strings.Builder{}
	columnStr := fmt.Sprintf("column:%s", col.Name.String())
	builder.WriteString(columnStr)
	switch col.Type.Type {
	case "enum":
		enumBuilder := strings.Builder{}
		for idx, _ := range col.Type.EnumValues {
			if 0 == idx {
				enumBuilder.WriteString(col.Type.EnumValues[idx])
			} else {
				enumBuilder.WriteString("," + col.Type.EnumValues[idx])
			}
		}
		typeStr := fmt.Sprintf(";type:enum(%s)", enumBuilder.String())
		builder.WriteString(typeStr)
	default:
		if nil != col.Type.Length {
			switch int(col.Type.Length.Type) {
			case 1: // int
				typeStr := fmt.Sprintf(";type:%s(%s)", col.Type.Type, col.Type.Length.Val)
				builder.WriteString(typeStr)
			}
		} else {

			typeStr := fmt.Sprintf(";type:%s", col.Type.Type)
			builder.WriteString(typeStr)
		}
	}

	if col.Type.Unsigned {
		builder.WriteString(" unsigned")
	}

	if col.Type.Autoincrement {
		builder.WriteString(" autoIncrement")
	}

	_, ok := primaryIdxMap[col.Name.String()]
	if ok {
		builder.WriteString(";primaryKey")
	}
	_, ok = uniqueIdxMap[col.Name.String()]
	if ok {
		builder.WriteString(";unique")
	}

	if nil != col.Type.Default {
		defaultStr := ""

		if col.Type.Type == "string" {
			defaultStr = fmt.Sprintf(";default:'%s'", col.Type.Default.Val)
		} else {
			defaultStr = fmt.Sprintf(";default:%s", col.Type.Default.Val)
		}
		builder.WriteString(defaultStr)
	}

	if col.Type.NotNull {
		builder.WriteString(";not null")
	}

	idxName, ok := idxMap[col.Name.String()]
	if ok {
		indexStr := fmt.Sprintf(";index:%s", idxName)
		builder.WriteString(indexStr)
	}

	if nil != col.Type.Comment {
		commentStr := fmt.Sprintf(";comment:'%s'", col.Type.Comment.Val)
		builder.WriteString(commentStr)
	}
	builder.WriteString(";")
	return builder.String()
}

func snakeToCamel(str string) string {
	builder := strings.Builder{}
	index := 0
	if str[0] >= 'a' && str[0] <= 'z' {
		builder.WriteByte(str[0] - ('a' - 'A'))
		index = 1
	}
	for i := index; i < len(str); i++ {
		if str[i] == '_' && i+1 < len(str) {
			if str[i+1] >= 'a' && str[i+1] <= 'z' {
				builder.WriteByte(str[i+1] - ('a' - 'A'))
				i++
				continue
			}
		}
		builder.WriteByte(str[i])
	}
	return builder.String()
}

func getAllResourceFile() error {

	for _, s := range staticFileList {
		sc, err := staticBox.FindString(s)
		if err != nil {
			return fmt.Errorf("get static file box err: %s", err.Error())
		}
		staticFileMap[s] = sc
	}
	for _, t := range templateFileList {
		tc, err := templateBox.FindString(t)
		if err != nil {
			return fmt.Errorf("get template file box err: %s", err.Error())
		}
		templateFileMap[t] = tc
	}
	return nil
}

func getAllResourceFileName(pathName string) error {
	files, err := ioutil.ReadDir(pathName)
	if err != nil {
		return fmt.Errorf("get all resource file name err: %s", err.Error())
	}
	for _, fi := range files {
		if fi.IsDir() {
			err = getAllResourceFileName(pathName + "/" + fi.Name())
			if err != nil {
				return err
			}
		} else {
			if strings.HasSuffix(fi.Name(), ".tpl") {
				templateFileList = append(templateFileList, fi.Name())
			}
			if strings.HasSuffix(fi.Name(), ".static") {
				staticFileList = append(staticFileList, fi.Name())
			}
		}
	}
	return nil
}

func getCreateDir(projectPath string) []string {
	var dirList = make([]string, 0)
	dirList = append(dirList, projectPath+"/conf")
	dirList = append(dirList, projectPath+"/centos7")
	dirList = append(dirList, projectPath+"/config")
	dirList = append(dirList, projectPath+"/router")
	dirList = append(dirList, projectPath+"/model")
	dirList = append(dirList, projectPath+"/api")
	dirList = append(dirList, projectPath+"/api/v1")
	dirList = append(dirList, projectPath+"/service")
	dirList = append(dirList, projectPath+"/initialize")
	dirList = append(dirList, projectPath+"/core")
	dirList = append(dirList, projectPath+"/global")
	dirList = append(dirList, projectPath+"/middleware")
	dirList = append(dirList, projectPath+"/utils")
	dirList = append(dirList, projectPath+"/model/request")
	dirList = append(dirList, projectPath+"/model/response")
	dirList = append(dirList, projectPath+"/initialize/internal")
	return dirList
}

func (t *AutoCoder) getTplDataList() ([]tplData, error) {
	var tplDataList = make([]tplData, 0)
	for _, tf := range templateFileList {
		tem, err := template.New(tf).Parse(templateFileMap[tf])
		if err != nil {
			return nil, err
		}
		td := tplData{
			template: tem,
		}
		switch tf {
		case "makefile.go.tpl":
			td.autoCodePath = t.ProjectPath + "/makefile"
		case "api_health.go.tpl":
			td.autoCodePath = t.ProjectPath + "/api/v1/health.go"
		case "global.go.tpl":
			td.autoCodePath = t.ProjectPath + "/global/global.go"
		case "redis.go.tpl":
			td.autoCodePath = t.ProjectPath + "/initialize/redis.go"
		case "system_api.go.tpl":
			td.autoCodePath = t.ProjectPath + "/api/v1/system.go"
		case "system_model.go.tpl":
			td.autoCodePath = t.ProjectPath + "/model/system.go"
		case "system_response.go.tpl":
			td.autoCodePath = t.ProjectPath + "/model/response/system.go"
		case "system_router.go.tpl":
			td.autoCodePath = t.ProjectPath + "/router/system.go"
		case "system_service.go.tpl":
			td.autoCodePath = t.ProjectPath + "/service/system.go"
		case "viper.go.tpl":
			td.autoCodePath = t.ProjectPath + "/core/viper.go"
		case "api.go.tpl":
			td.autoCodePath = t.ProjectPath + "/api/v1/" + t.TableName + ".go"
			td.repeat = true
		case "error.go.tpl":
			td.autoCodePath = t.ProjectPath + "/middleware/error.go"
		case "gorm.go.tpl":
			td.autoCodePath = t.ProjectPath + "/initialize/gorm.go"
		case "health.go.tpl":
			td.autoCodePath = t.ProjectPath + "/router/health.go"
		case "initRouter.go.tpl":
			td.autoCodePath = t.ProjectPath + "/initialize/router.go"
		case "logger.go.tpl":
			td.autoCodePath = t.ProjectPath + "/initialize/internal/logger.go"
		case "main.go.tpl":
			td.autoCodePath = t.ProjectPath + "/main.go"
		case "model.go.tpl":
			td.autoCodePath = t.ProjectPath + "/model/" + t.TableName + ".go"
			td.repeat = true
		case "request.go.tpl":
			td.autoCodePath = t.ProjectPath + "/model/request/" + t.TableName + ".go"
			td.repeat = true
		case "router.go.tpl":
			td.autoCodePath = t.ProjectPath + "/router/" + t.TableName + ".go"
			td.repeat = true
		case "service.go.tpl":
			td.autoCodePath = t.ProjectPath + "/service/" + t.TableName + ".go"
			td.repeat = true
		case "zap.go.tpl":
			td.autoCodePath = t.ProjectPath + "/core/zap.go"
		}
		tplDataList = append(tplDataList, td)
	}
	return tplDataList, nil
}

func (t *AutoCoder) copyAllStaticFile() error {
	var fp string
	for _, s := range staticFileList {
		switch s {
		// case "api_health.static":
		// 	fp = t.ProjectPath + "/api/health.go"
		case "fmt_plus.static":
			fp = t.ProjectPath + "/utils/fmt_plus.go"
		case "config_struct.static":
			fp = t.ProjectPath + "/config/config.go"
		case "config.static":
			fp = t.ProjectPath + "/conf/GF_PROJECT_NAME.conf"
		case "service.static":
			fp = t.ProjectPath + "/centos7/GF_PROJECT_NAME.service"
		case "spec.static":
			fp = t.ProjectPath + "/centos7/GF_PROJECT_NAME.spec"
		case "constant.static":
			fp = t.ProjectPath + "/utils/constant.go"
		case "mysql_struct.static":
			fp = t.ProjectPath + "/config/mysql.go"
		case "redis_struct.static":
			fp = t.ProjectPath + "/config/redis.go"
		case "server.static":
			fp = t.ProjectPath + "/utils/server.go"
		case "system_struct.static":
			fp = t.ProjectPath + "/config/system.go"
		case "zap_struct.static":
			fp = t.ProjectPath + "/config/zap.go"
		case "cors.static":
			fp = t.ProjectPath + "/middleware/cors.go"
		case "global.static":
			fp = t.ProjectPath + "/global/global.go"
		case "loadtls.static":
			fp = t.ProjectPath + "/middleware/loadtls.go"
		case "request.static":
			fp = t.ProjectPath + "/model/request/request.go"
		case "response.static":
			fp = t.ProjectPath + "/model/response/response.go"
		case "file_operation.static":
			fp = t.ProjectPath + "/utils/file_operation.go"
		case "directory.static":
			fp = t.ProjectPath + "/utils/directory.go"
		case "rotatelogs.static":
			fp = t.ProjectPath + "/utils/rotatelogs.go"
		}
		if err := utils.CopyFile(staticFileMap[s], fp); err != nil {
			fmt.Println("copy file error:" + err.Error())
			return err
		}
	}
	return nil
}

// func getAllTplFile(pathName string, fileList []string) ([]string, error) {
// 	files, err := ioutil.ReadDir(pathName)
// 	for _, fi := range files {
// 		if fi.IsDir() {
// 			fileList, err = getAllTplFile(pathName+"/"+fi.Name(), fileList)
// 			if err != nil {
// 				return nil, err
// 			}
// 		} else {
// 			if strings.HasSuffix(fi.Name(), ".tpl") {
// 				fileList = append(fileList, pathName+"/"+fi.Name())
// 			}
// 		}
// 	}
// 	return fileList, err
// }

func (t *AutoCoder) CreateTemp() (err error) {
	needCreateDirs := getCreateDir(t.ProjectPath)

	// dataList, needMKdir, err := t.getNeedList()
	// if err != nil {
	// 	return err
	// }
	tplDataList, err := t.getTplDataList()
	if err != nil {
		return err
	}
	// 写入文件前，先创建文件夹
	fmt.Printf("needMKdirs::%v\n", needCreateDirs)
	if err = utils.CreateDir(needCreateDirs...); err != nil {
		fmt.Println("create dirs error:" + err.Error())
		return err
	}

	// 复制文件
	if err = t.copyAllStaticFile(); err != nil {
		return err
	}

	// 生成文件
	for _, value := range tplDataList {
		// 对于/initialize/router.go文件,需要追加代码而不是清空覆盖
		if strings.Contains(value.autoCodePath, "/initialize/router") {
			fmt.Println("enter init router")
			exist, err := utils.PathExists(value.autoCodePath)
			if err != nil {
				fmt.Println("init router path err")
				return err
			}
			if exist {
				if t.SqlPath != "" {
					// 逐行添加init的router代码
					err := t.writeInitRouterCode(value.autoCodePath)

					if err != nil {
						fmt.Println("write router err")
						return err
					}
				}
			} else {
				if err = executeTemplate(&value, t); err != nil {
					return err
				}
				if t.SqlPath != "" {
					err := t.writeInitRouterCode(value.autoCodePath)
					if err != nil {
						fmt.Printf("first write router err :%s", err.Error())
						return err
					}
				}
			}
		} else if !value.repeat {
			fmt.Println("enter " + value.autoCodePath)
			exist, err := utils.PathExists(value.autoCodePath)
			if err != nil {
				return err
			}
			if !exist {
				if err = executeTemplate(&value, t); err != nil {
					fmt.Println("execute template err: " + value.autoCodePath)
					return err
				}
			}
		} else {
			// 对于需要传入sql文件路径才能进行渲染的部分单独进行判断
			if t.SqlPath != "" {
				fmt.Println("enter else")
				if err = executeTemplate(&value, t); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func executeTemplate(tpl *tplData, coder *AutoCoder) error {
	fmt.Println("autocodepath::" + tpl.autoCodePath)
	f, err := os.OpenFile(tpl.autoCodePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println("open tpl file err::" + tpl.autoCodePath)
		return err
	}
	if err = tpl.template.Execute(f, coder); err != nil {
		fmt.Println("excute err::" + tpl.autoCodePath)
		return err
	}
	_ = f.Close()
	return nil
}

func (t *AutoCoder) writeInitRouterCode(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	var lineList []string
	reader := bufio.NewReader(f)

	for {
		content, _, err := reader.ReadLine()
		lineList = append(lineList, string(content)+"\n")
		if strings.Contains(string(content), "**BEGIN") {
			lineList = append(lineList, "    router.Init"+t.StructName+"Router(group)\n")
		}
		if err == io.EOF {
			break
		}
	}
	f.Close()
	fw, err := os.OpenFile(filePath, os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	fw.Write([]byte(""))
	fw.Close()
	ff, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		return err
	}
	for _, line := range lineList {
		_, err := ff.Write([]byte(line))
		if err != nil {
			return err
		}
	}
	ff.Close()
	return nil
}

// func (t *AutoCoder) getNeedList() (dataList []tplData, needMKDirs []string, err error) {
// 	// 去除所有空格
// 	utils.TrimSpace(t)
// 	for _, field := range t.Fields {
// 		utils.TrimSpace(field)
// 	}
// 	// 获取basePath 文件夹下所有的tpl文件和静态文件
// 	tplFileList, err := getAllTplFile(t.TplPath, nil)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	dataList = make([]tplData, 0)
// 	needMKDirs = make([]string, 0)
// 	// 根据文件路径生成tplData结构体，待填充数据
// 	for _, value := range tplFileList {
// 		dataList = append(dataList, tplData{locationPath: value})
// 	}
// 	// 生成 *template, 填充template字段
// 	for index, value := range dataList {
// 		dataList[index].template, err = template.ParseFiles(value.locationPath)
// 		if err != nil {
// 			return nil, nil, err
// 		}
// 	}
// 	for index, value := range dataList {
// 		if strings.Contains(value.locationPath, "router") {
// 			dataList[index].autoCodePath = t.ProjectPath + "/router/" + t.TableName + ".go"
// 			needMKDirs = append(needMKDirs, t.ProjectPath+"/router")
// 		} else if strings.Contains(value.locationPath, "model") {
// 			dataList[index].autoCodePath = t.ProjectPath + "/model/" + t.TableName + ".go"
// 			needMKDirs = append(needMKDirs, t.ProjectPath+"/model")
// 		} else if strings.Contains(value.locationPath, "api") {
// 			dataList[index].autoCodePath = t.ProjectPath + "/api/" + t.TableName + ".go"
// 			needMKDirs = append(needMKDirs, t.ProjectPath+"/api")
// 		} else if strings.Contains(value.locationPath, "service") {
// 			dataList[index].autoCodePath = t.ProjectPath + "/service/" + t.TableName + ".go"
// 			needMKDirs = append(needMKDirs, t.ProjectPath+"/service")
// 		} else if strings.Contains(value.locationPath, "request") {
// 			dataList[index].autoCodePath = t.ProjectPath + "/model/request/" + t.TableName + ".go"
// 			needMKDirs = append(needMKDirs, t.ProjectPath+"/model/request")
// 		} else if strings.Contains(value.locationPath, "health") {
// 			dataList[index].autoCodePath = t.ProjectPath + "/router/health.go"
// 		} else if strings.Contains(value.locationPath, "main") {
// 			dataList[index].autoCodePath = t.ProjectPath + "/main.go"
// 		} else if strings.Contains(value.locationPath, "initRouter") {
// 			dataList[index].autoCodePath = t.ProjectPath + "/initialize/router.go"
// 		}
// 	}

// // 添加健康检查的router文件渲染
// tpRouter, err := template.ParseFiles(t.GoFastPath + "/resource/template/health.go.tpl")
// if err != nil {
// 	return nil, nil, err
// }
// routerTplData := tplData{
// 	template:     tpRouter,
// 	locationPath: t.GoFastPath + "/resource/template/health.go.tpl",
// 	autoCodePath: t.ProjectPath + "/router/" + "health.go",
// }
// dataList = append(dataList, routerTplData)

// // 添加main文件的渲染
// tplMain, err := template.ParseFiles(t.GoFastPath + "/resource/template/main.go.tpl")
// if err != nil {
// 	return nil, nil, err
// }
// mainTplData := tplData{
// 	template:     tplMain,
// 	locationPath: t.GoFastPath + "/resource/template/main.go.tpl",
// 	autoCodePath: t.ProjectPath + "/main.go",
// }
// dataList = append(dataList, mainTplData)

// // 添加初始化router文件的渲染
// tplInitRouter, err := template.ParseFiles(t.GoFastPath + "/resource/template/initRouter.go.tpl")
// if err != nil {
// 	return nil, nil, err
// }
// initRouterTplData := tplData{
// 	template:     tplInitRouter,
// 	locationPath: t.GoFastPath + "/resource/template/initRouter.go.tpl",
// 	autoCodePath: t.ProjectPath + "/initialize/router.go",
// }
// dataList = append(dataList, initRouterTplData)

// 	needMKDirs = append(needMKDirs, t.ProjectPath+"/initialize")     // 添加initialize文件夹
// 	needMKDirs = append(needMKDirs, t.ProjectPath+"/model/response") // 添加response文件夹
// 	return dataList, needMKDirs, err
// }

// TODO: 添加对main文件的渲染替换
// TODO: 添加在router文件中写入每个api的路由语句
