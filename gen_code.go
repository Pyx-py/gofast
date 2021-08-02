package gofast

import (
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/pyx-py/gofast/utils"
)

func NewCodeTool(projectPath, gofastPath string) {
	if gofastPath == "" {
		path, err := utils.GetGoFastPath()
		if err != nil {
			fmt.Println("can not find gofast path")
		}
		gofastPath = path
	}
	fmt.Println(gofastPath)
	fmt.Printf("go build -o %s %s/cmd/main.go", projectPath+"/gf\n", gofastPath)
	cmd := exec.Command("go", "build", "-o", projectPath+"/gf", gofastPath+"/cmd/main.go")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
		return
	}
	if err := cmd.Start(); err != nil {
		fmt.Println("Error:The command is err", err)
		return
	}
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Println("ReadAll Stdout:", err.Error())
		return
	}
	if err := cmd.Wait(); err != nil {
		fmt.Println("wait:", err.Error())
		return
	}
	fmt.Printf("stdout:\n\n %s", bytes)
}

// type Argument struct {
// 	ProjectPath        string
// 	ModuleName         string
// 	SqlFilePath        string
// 	ColSearchMapString string
// 	LogPath            string
// }

// func CodeGenarate(args Argument) error {
// 	coder, err := utils.NewAutoCoder(args.ProjectPath, args.ModuleName, args.SqlFilePath, args.ColSearchMapString, args.LogPath)
// 	if err != nil {
// 		return err
// 	}
// 	err = coder.CreateTemp()
// 	return err
// }
