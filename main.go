package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

func WalkDir(filepath string, wg *sync.WaitGroup) ([]string, error) {
	files, err := ioutil.ReadDir(filepath) // files为当前目录下的所有文件名称【包括文件夹】
	if err != nil {
		return nil, err
	}

	var pathList []string
	for _, v := range files {
		name := v.Name()
		fullPath := filepath + "/"

		if ".git" == name {
			wg.Add(1)
			go fetch(fullPath, wg)
			break
		} else {
			if v.IsDir() {
				// 如果是目录遍历改路径下的所有文件
				a, _ := WalkDir(fullPath+name, wg)
				pathList = append(pathList, a...)
			}
		}
	}
	return pathList, nil
}

func fetch(path string, wg *sync.WaitGroup) {
	defer wg.Done()
	fetchLog, _ := runGitCommand(path, "git", "fetch", "--all")
	fmt.Printf("%s\n%s\n", path, fetchLog)
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func runGitCommand(path string, name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	cmd.Dir = path
	msg, err := cmd.CombinedOutput()
	_ = cmd.Run()
	return string(msg), err
}

func main() {
	dir := ""
	if 1 == len(os.Args) {
		dir = getCurrentDirectory()
	} else if 2 == len(os.Args) {
		dir = os.Args[1]
	} else {
		fmt.Println("please input the path like: /work")
	}

	if dir != "" {
		var wg sync.WaitGroup
		_, _ = WalkDir(dir, &wg)
		wg.Wait()
		fmt.Printf(" %s fetch done.^_^.", dir)
	}
}
