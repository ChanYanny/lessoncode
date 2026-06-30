package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	target := "/Users/bytedance/study/lessoncode"
	task := fmt.Sprintf("扫描目录：%s, 列出 .go 文件，并读取第一个 go 文件的前 5 行", target)
	fmt.Println("任务：", task)
	fmt.Println()

	fmt.Println("==== 聊天机器人回答 =======")
	fmt.Println(chatbotReply(task))

	fmt.Println("===== agent 的做法=======")
	fmt.Println(agentRun(target))
}

func chatbotReply(task string) string {
	return "你可以用 `dir *.go` 或 `Get-ChildItem -Filter *.go` 看一下目录里的 Go 文件，\n" +
		"再用编辑器或 `Get-Content` 看第一个文件的前几行。\n" +
		"（我只是在描述做法，没真去执行，也不知道你目录里到底有什么文件。）"
}

func agentRun(target string) string {
	var (
		out     strings.Builder
		goFiles []string
	)
	_ = filepath.WalkDir(target, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if strings.HasSuffix(d.Name(), ".go") {
			goFiles = append(goFiles, path)
		}
		return nil
	})
	out.WriteString("第一步：扫描目录，得到Go文件：" + strings.Join(goFiles, ",") + "\n")
	if len(goFiles) == 0 {
		out.WriteString("当前目录没有go文件")
		return out.String()
	}
	first := goFiles[0]
	out.WriteString("第二步：读取 " + first + " 前5行：\n")
	lines, err := readFirstFileLines(first, 5)
	if err != nil {
		out.WriteString("读取失败：" + err.Error())
		return out.String()
	}
	for _, line := range lines {
		out.WriteString("   " + line + "\n")
	}
	out.WriteString("第三步：整理结果")
	return out.String()
}

func readFirstFileLines(first string, n int) ([]string, error) {
	f, err := os.Open(first)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var lines []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		lines = append(lines, sc.Text())
		if len(lines) >= n {
			break
		}
	}
	return lines, sc.Err()
}
