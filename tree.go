/*
treeコマンドは指定したディレクトリ配下のファイル一覧を表示します。
引数を指定しないとカレントディレクトリ以下を表示します。
option:
    -a: .ファイル含めすべてのファイルを表示します。
    -d: ディレクトリのみを表示します。
*/
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const parts1 = "│   "
const parts2 = "    "
const parts3 = "└──"
const parts4 = "├──"

type count struct {
	FileCount int
	DirCount  int
}

var showAll = flag.Bool("a", false, "Show all dir and files")
var dirOnly = flag.Bool("d", false, "Show directory only")

func main() {
	flag.Parse()
	path := getPath(flag.Args())

	c := &count{}
	c.getCount(path)
	countMessage := ""

	fmt.Println(path)

	// ディレクトリのみ
	if *dirOnly {
		showDir(path, "")
		countMessage = fmt.Sprintf("%d directories", c.DirCount)
	} else {
		showTree(path, "")
		countMessage = fmt.Sprintf("%d directories, %d files", c.DirCount, c.FileCount)
	}

	// 最後にヒットした件数を表示
	fmt.Printf("\n%s\n", countMessage)
}

// 通常のツリーを表示します。
func showTree(dir string, connector string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Printf("tree: %s: No such file or directory\n", dir)
		os.Exit(1)
	}

	newConnector := connector
	fileCount := len(files)
	for i, file := range files {

		if string(file.Name()[0]) == "." && !*showAll {
			continue
		}

		if i == fileCount-1 {
			fmt.Println(getOutput(connector, parts3, file.Name()))
			newConnector = connector + parts2
		} else {
			fmt.Println(getOutput(connector, parts4, file.Name()))
			newConnector = connector + parts1
		}

		if file.IsDir() {
			showTree(filepath.Join(dir, file.Name()), newConnector)
		}
	}

}

// ディレクトリのみを表示します。
func showDir(dir string, connector string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Printf("tree: %s: No such file or directory\n", dir)
		os.Exit(1)
	}

	var dirCount int
	// その階層のディレクトリの総数をあらかじめに取得
	for _, file := range files {

		if string(file.Name()[0]) == "." && !*showAll {
			continue
		}
		if file.IsDir() {
			dirCount++
		}
	}

	newConnector := connector
	// ディレクトリ出現回数
	var dirAppearCount int
	for _, file := range files {

		if string(file.Name()[0]) == "." && !*showAll {
			continue
		}

		if file.IsDir() {
			if dirAppearCount == dirCount-1 {
				fmt.Println(getOutput(connector, parts3, file.Name()))
				newConnector = connector + parts2
			} else {
				fmt.Println(getOutput(connector, parts4, file.Name()))
				newConnector = connector + parts1
			}

			dirAppearCount++

			showDir(filepath.Join(dir, file.Name()), newConnector)
		}
	}
}

// ディレクトリとファイルの数を集計します。
func (c *count) getCount(dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Printf("tree: %s: No such file or directory\n", dir)
		os.Exit(1)
	}

	for _, file := range files {

		if string(file.Name()[0]) == "." && !*showAll {
			continue
		}

		if file.IsDir() {
			c.DirCount++
			c.getCount(filepath.Join(dir, file.Name()))
		} else {
			c.FileCount++
		}
	}
}

// 引数からパスを取得します。
// 引数がない場合はカレントディレクトリを指定します。
func getPath(paths []string) string {
	if len(paths) < 1 {
		return "."
	}

	return paths[0]
}

// 表示用の文字列を作成します。
func getOutput(connector string, parts string, fileName string) string {
	return fmt.Sprintf("%s%s %s", connector, parts, fileName)
}
