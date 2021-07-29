package utils

import (
	"bufio"
	"fmt"
	"go/build"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/xwb1989/sqlparser"
)

const GOFAST = "github.com/Pyx-py/gofast"

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
	ModuleName          string   `json:"moduleName"`  // 用户的module名称
	ProjectPath         string   `json:"projectPath"` // 用户的项目路径
	LogPath             string   `json:"logPath"`     // 日志的存放路径，为空代表不启用默认日志
	TplPath             string   `json:"tplPath"`     // 模板文件dir
	StructName          string   `json:"structName"`
	TableName           string   `json:"tableName"`
	Abbreviation        string   `json:"abbreviation"`
	ImportTime          bool     `json:"importTime"`
	GoStructString      string   `json:"goStructString"`
	Fields              []*Field `json:"fields"`
	ColSearchTypeString string   `json:"colSearchTypeString"`
	SqlPath             string   `json:"sqlPath"`    // 传入的sql文件路径
	GoFastPath          string   `json:"goFastPath"` //
}

type Field struct {
	FieldName       string `json:"fieldName"`
	FieldType       string `json:"fieldType"`
	ColumnName      string `json:"columnName"`
	FieldSearchType string `json:"fieldSearchType"`
}

type tplData struct {
	template     *template.Template
	locationPath string
	autoCodePath string
	// autoMoveFilePath string
}

func NewAutoCoder(projectPath, moduleName, sqlFilePath, colSearchTypeString, logPath, gofastPath string) (*AutoCoder, error) {
	if projectPath == "" || moduleName == "" {
		return nil, fmt.Errorf("projectPath or moduleName can not be null")
	}
	if gofastPath == "" {
		path, err := GetGoFastPath()
		if err != nil {
			return nil, err
		}
		gofastPath = path
	}
	fmt.Printf("get gofast path: %s", gofastPath)
	// init aurocoder
	autoCoder := &AutoCoder{
		TplPath:    gofastPath + "/resource/template",
		ImportTime: false,
		Fields:     make([]*Field, 0),
	}
	autoCoder.GoFastPath = gofastPath
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
			builder.WriteString(fmt.Sprintf("\t%s\t%-25s `json:\"%s\" gorm:\"%s\"`\n", fieldName, goType, col.Name.String(), gormStr))
		}
		builder.WriteString("}\n")
		autoCoder.GoStructString = builder.String()
	}

	return autoCoder, nil
}

func GetGoFastPath() (string, error) {
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

	mapList := strings.Split(colSearchMapString, "|")
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
		builder.WriteString(" auto_increment")
	}

	_, ok := primaryIdxMap[col.Name.String()]
	if ok {
		builder.WriteString(";primary_key")
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

func getAllTplFile(pathName string, fileList []string) ([]string, error) {
	files, err := ioutil.ReadDir(pathName)
	for _, fi := range files {
		if fi.IsDir() {
			fileList, err = getAllTplFile(pathName+"/"+fi.Name(), fileList)
			if err != nil {
				return nil, err
			}
		} else {
			if strings.HasSuffix(fi.Name(), ".tpl") {
				fileList = append(fileList, pathName+"/"+fi.Name())
			}
		}
	}
	return fileList, err
}

func (t *AutoCoder) CreateTemp() (err error) {
	dataList, needMKdir, err := t.getNeedList()
	if err != nil {
		return err
	}
	// 写入文件前，先创建文件夹
	fmt.Printf("needMKdirs::%v\n", needMKdir)
	if err = CreateDir(needMKdir...); err != nil {
		return err
	}
	// 复制文件
	reqErr := CopyFile(t.GoFastPath+"/resource/static/request.static", t.ProjectPath+"/model/request/request.go")
	if reqErr != nil {
		return reqErr
	}
	resErr := CopyFile(t.GoFastPath+"/resource/static/response.static", t.ProjectPath+"/model/response/response.go")
	if resErr != nil {
		return resErr
	}
	hErr := CopyFile(t.GoFastPath+"/resource/static/api_health.static", t.ProjectPath+"/api/health.go")
	if hErr != nil {
		return hErr
	}
	// 生成文件
	for _, value := range dataList {
		// 对于/initialize/router.go文件,需要追加代码而不是清空覆盖
		if strings.Contains(value.autoCodePath, "/initialize/router") {
			fmt.Println("enter init router")
			exist, err := PathExists(value.autoCodePath)
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
			}
			continue
		}
		// 对于入口main文件，第一次渲染完毕之后，如果main文件存在并且第一行有tag的话，也不会再次渲染
		if strings.Contains(value.autoCodePath, "main") {
			fmt.Println("enter main")
			mExist, err := PathExists(value.autoCodePath)
			if err != nil {
				return err
			}
			if mExist {
				f, err := os.Open(value.autoCodePath)
				if err != nil {
					return err
				}
				reader := bufio.NewReader(f)
				var tag = false
				for {
					content, _, err := reader.ReadLine()
					if strings.Contains(string(content), "**INIT_MAIN**") {
						tag = true
						break
					}
					if err != nil {
						if err != io.EOF {
							return err
						} else {
							break
						}
					}
				}
				if !tag {
					if err = executeTemplate(&value, t); err != nil {
						return err
					}
				}
			}
			continue
		}

		// 对于healthCheck的router，可以单独进行渲染
		if strings.Contains(value.autoCodePath, "health") {
			fmt.Println("enter router health chenk")
			if err = executeTemplate(&value, t); err != nil {
				fmt.Println("execute template err :" + value.autoCodePath)
				return err
			}
		}

		// 对于需要传入sql文件路径才能进行渲染的部分单独进行判断
		if t.SqlPath != "" {
			fmt.Println("enter else")
			if err = executeTemplate(&value, t); err != nil {
				return err
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
		lineList = append(lineList, string(content))
		if strings.Contains(string(content), "**BEGIN") {
			lineList = append(lineList, "    router.Init"+t.StructName+"Router(group)")
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

func (t *AutoCoder) getNeedList() (dataList []tplData, needMKDirs []string, err error) {
	// 去除所有空格
	TrimSpace(t)
	for _, field := range t.Fields {
		TrimSpace(field)
	}
	// 获取basePath 文件夹下所有的tpl文件
	tplFileList, err := getAllTplFile(t.TplPath, nil)
	if err != nil {
		return nil, nil, err
	}
	dataList = make([]tplData, 0)
	needMKDirs = make([]string, 0)
	// 根据文件路径生成tplData结构体，待填充数据
	for _, value := range tplFileList {
		dataList = append(dataList, tplData{locationPath: value})
	}
	// 生成 *template, 填充template字段
	for index, value := range dataList {
		dataList[index].template, err = template.ParseFiles(value.locationPath)
		if err != nil {
			return nil, nil, err
		}
	}
	for index, value := range dataList {
		if strings.Contains(value.locationPath, "router") {
			dataList[index].autoCodePath = t.ProjectPath + "/router/" + t.TableName + ".go"
			needMKDirs = append(needMKDirs, t.ProjectPath+"/router")
		} else if strings.Contains(value.locationPath, "model") {
			dataList[index].autoCodePath = t.ProjectPath + "/model/" + t.TableName + ".go"
			needMKDirs = append(needMKDirs, t.ProjectPath+"/model")
		} else if strings.Contains(value.locationPath, "api") {
			dataList[index].autoCodePath = t.ProjectPath + "/api/" + t.TableName + ".go"
			needMKDirs = append(needMKDirs, t.ProjectPath+"/api")
		} else if strings.Contains(value.locationPath, "service") {
			dataList[index].autoCodePath = t.ProjectPath + "/service/" + t.TableName + ".go"
			needMKDirs = append(needMKDirs, t.ProjectPath+"/service")
		} else if strings.Contains(value.locationPath, "request") {
			dataList[index].autoCodePath = t.ProjectPath + "/model/request/" + t.TableName + ".go"
			needMKDirs = append(needMKDirs, t.ProjectPath+"/model/request")
		} else if strings.Contains(value.locationPath, "health") {
			dataList[index].autoCodePath = t.ProjectPath + "/router/health.go"
		} else if strings.Contains(value.locationPath, "main") {
			dataList[index].autoCodePath = t.ProjectPath + "/main.go"
		} else if strings.Contains(value.locationPath, "initRouter") {
			dataList[index].autoCodePath = t.ProjectPath + "/initialize/router.go"
		}
	}

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

	needMKDirs = append(needMKDirs, t.ProjectPath+"/initialize")     // 添加initialize文件夹
	needMKDirs = append(needMKDirs, t.ProjectPath+"/model/response") // 添加response文件夹
	return dataList, needMKDirs, err
}

// TODO: 添加对main文件的渲染替换
// TODO: 添加在router文件中写入每个api的路由语句
