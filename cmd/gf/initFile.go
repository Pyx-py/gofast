package gf

import (
	"fmt"

	"github.com/Pyx-py/gofast/utils"
	"github.com/spf13/cobra"
)

var projectPath string
var moduleName string
var sqlFilePath string
var logPath string
var colSearchMapString string
var gofastPath string // gofast包被下载之后的在gopath或者mod下面的路径，用来找渲染文件和静态文件

var fileCmd = &cobra.Command{
	Use:   "init",
	Short: "genarate code file",
	Long:  "genarate crud code file, include api, router, model, service etc.",
	Run: func(cmd *cobra.Command, args []string) {
		// if sqlFilePath == "" {
		// 	fmt.Println("argument sql can't be empty!")
		// 	return
		// }
		if projectPath == "" {
			fmt.Println("argument project path can't be empty!")
			return
		}
		if moduleName == "" {
			fmt.Println("argument module name can't be empty!")
			return
		}
		// sqlData, err := ioutil.ReadFile(sqlFilePath)
		// if err != nil {
		// 	fmt.Printf("[error]:%s\n", err.Error())
		// 	return
		// }
		// sqlString := string(sqlData)
		// searchMap, err := handleSearchMap()
		// if err != nil {
		// 	fmt.Printf("[error]:%s\n", err.Error())
		// 	return
		// }
		coder, err := utils.NewAutoCoder(projectPath, moduleName, sqlFilePath, colSearchMapString, logPath, gofastPath)
		if err != nil {
			fmt.Printf("[error]:%s\n", err.Error())
			return
		}
		err = coder.CreateTemp()
		if err != nil {
			fmt.Printf("[error]:%s\n", err.Error())
			return
		}
	},
}

func init() {
	fileCmd.Flags().StringVarP(&projectPath, "path", "p", "", "project absolute dir path")
	fileCmd.Flags().StringVarP(&moduleName, "module", "m", "", "module name")
	fileCmd.Flags().StringVarP(&sqlFilePath, "sql", "s", "", "sql file absolute path")
	fileCmd.Flags().StringVarP(&logPath, "log", "l", "", "log dir path")
	fileCmd.Flags().StringVarP(&colSearchMapString, "column", "c", "", "column search type")
	fileCmd.Flags().StringVarP(&gofastPath, "gofast", "f", "", "gofast package path")
}
