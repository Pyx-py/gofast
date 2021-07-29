package middleware

import (
	"fmt"
	"go/build"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestSome(*testing.T) {
	f, _ := os.OpenFile("test", os.O_TRUNC, 0755)
	f.Write([]byte(""))
}

type testTem struct {
	Name string
	Age  int
}

func TestTemplate(*testing.T) {
	a := testTem{
		Name: "qq",
		Age:  18,
	}
	f, _ := os.OpenFile("test", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	t, _ := template.ParseFiles("test_template.tpl")
	t.Execute(f, a)
}

func TestMap(*testing.T) {
	a := make(map[string]string)
	// a["11"] = "ww"
	fmt.Println(reflect.TypeOf(a["11"]))
}

func TestFor(*testing.T) {
	var tag = false
	a := []int{1, 2, 3, 4, 5, 6}
	for _, i := range a {
		if i > 10 {
			tag = true
			break
		}
	}
	if tag {
		fmt.Println("okok")
	} else {
		fmt.Println("ererer")
		fmt.Println(fmt.Sprintf("server run on %s:%s", "11", "22"))
	}
}

func TestList(*testing.T) {
	a := []int{1, 2, 3, 4}
	b := a[0 : len(a)-1]
	fmt.Println(b)
	c := a[0 : len(a)-2]
	fmt.Println(c)
	d := a[len(a)-1]
	fmt.Println(d)
}

func getDirList(dirpath string) ([]string, error) {
	var dir_list []string
	dir_err := filepath.Walk(dirpath,
		func(path string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() {
				dir_list = append(dir_list, path)
				return nil
			}

			return nil
		})
	return dir_list, dir_err
}

func TestGetDirList(*testing.T) {
	dirList, err := getDirList("/home/pyx/work/gofast")
	if err != nil {
		panic(err)
	}
	fmt.Println(dirList)
	dirList1, _ := filepath.Glob("/home/pyx/work/gofast/*")
	fmt.Println(dirList1)
	dirList2, _ := ioutil.ReadDir("/home/pyx/work/gofast")
	for _, d := range dirList2 {
		fmt.Println(d.Name())
	}
}

func TestGOPATH(*testing.T) {
	fmt.Println(os.Getenv("GOPATH"))
	fmt.Println(build.Default.GOPATH)
}
