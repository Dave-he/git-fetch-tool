package main

import (
	// "bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

type fetch func(string, *sync.WaitGroup)

func fetchLinux(path string, wg *sync.WaitGroup) {
	defer wg.Done()
	runGitCommand(path, "git", "fetch", "--all")
}

func fetchWin(path string, wg *sync.WaitGroup) {
	defer wg.Done()
	runGitCommand(path, "cmd", "/C", "git", "fetch", "--all")
}

func runGitCommand(path string, name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Dir = path
	msg, _ := cmd.CombinedOutput()
	_ = cmd.Run()
	fmt.Printf("%s\n%s\n", path, msg)
}

func WalkDir(filepath string, fetchCMD fetch, wg *sync.WaitGroup) ([]string, error) {
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
			go fetchCMD(fullPath, wg)
			break
		} else {
			if v.IsDir() {
				// 如果是目录遍历改路径下的所有文件
				a, _ := WalkDir(fullPath+name, fetchCMD, wg)
				pathList = append(pathList, a...)
			}
		}
	}
	return pathList, nil
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func process(dir string, fetchCMD fetch) {
	if dir != "" {
		var wg sync.WaitGroup
		_, _ = WalkDir(dir, fetchCMD, &wg)
		wg.Wait()
		fmt.Printf("%s fetch done.^_^.", dir)
	}
}

func main() {
	fetchCMD := fetchLinux
	sysType := runtime.GOOS
	if sysType == "windows" {
		fetchCMD = fetchWin
	}

	switch len(os.Args) {
	case 1:
		process(getCurrentDirectory(), fetchCMD)
	case 2:
		process(os.Args[1], fetchCMD)
	default:
		fmt.Println("please input the path like: C:/work")
	}

	// input := bufio.NewScanner(os.Stdin)
	// input.Scan()
	// fmt.Println(input.Text())
}
