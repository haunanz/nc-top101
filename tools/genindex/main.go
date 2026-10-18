// genindex 根据根目录 solution.go 中各题的小节注释，自动生成 README 的题目索引表。
//
// 用法（在仓库根目录执行）：
//
//	go run ./tools/genindex          # 打印索引表到标准输出
//	go run ./tools/genindex -write   # 用生成的索引表替换 README.md 中标记区间的内容
//
// 约定：每道题在 solution.go 中以 `// BM001 反转链表` 开头，
// 小节内包含 `// 时间复杂度：O(...)` 与 `// 空间复杂度：O(...)` 两行注释。
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const (
	solutionPath = "solution.go"
	readmePath   = "README.md"
	beginMarker  = "<!-- BEGIN:INDEX -->"
	endMarker    = "<!-- END:INDEX -->"
)

var (
	// 题目小节标题：`// BM001 反转链表`
	titleRe = regexp.MustCompile(`^//\s*(BM(\d{3}))\s+(.+?)\s*$`)
	// 复杂度：兼容全角与半角冒号
	complexityRe = regexp.MustCompile(`^//\s*(时间复杂度|空间复杂度)\s*[:：]\s*(.+?)\s*$`)
)

type entry struct {
	id     string
	title  string
	line   int
	timeC  string
	spaceC string
}

func main() {
	write := flag.Bool("write", false, "把生成的索引表写回 README.md 的标记区间")
	flag.Parse()

	entries, err := collect(solutionPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "解析 %s 失败: %v\n", solutionPath, err)
		os.Exit(1)
	}
	if len(entries) == 0 {
		fmt.Fprintf(os.Stderr, "%s 中没有找到形如 `// BM001 题目名` 的小节标题\n", solutionPath)
		os.Exit(1)
	}

	table := render(entries)
	if !*write {
		fmt.Println(table)
		return
	}
	if err := updateReadme(table); err != nil {
		fmt.Fprintf(os.Stderr, "更新 README 失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("已更新 %s 中的题目索引，共 %d 题\n", readmePath, len(entries))
}

// collect 扫描 solution.go，抽取每题的题号、题目名、起始行号与复杂度。
func collect(path string) ([]entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []entry
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	lineNo := 0
	cur := -1 // 当前题目在 entries 中的下标
	for sc.Scan() {
		lineNo++
		line := strings.TrimRight(sc.Text(), "\r")

		if m := titleRe.FindStringSubmatch(line); m != nil {
			entries = append(entries, entry{id: m[1], title: m[3], line: lineNo})
			cur = len(entries) - 1
			continue
		}
		if cur < 0 {
			continue
		}
		if m := complexityRe.FindStringSubmatch(line); m != nil {
			if m[1] == "时间复杂度" {
				entries[cur].timeC = m[2]
			} else {
				entries[cur].spaceC = m[2]
			}
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	for i := range entries {
		if entries[i].timeC == "" {
			entries[i].timeC = "—"
		}
		if entries[i].spaceC == "" {
			entries[i].spaceC = "—"
		}
	}
	return entries, nil
}

// render 生成 markdown 表格。
func render(entries []entry) string {
	var b strings.Builder
	b.WriteString("| 题号 | 题目 | 位置 | 时间复杂度 | 空间复杂度 |\n")
	b.WriteString("| --- | --- | --- | --- | --- |\n")
	for _, e := range entries {
		fmt.Fprintf(&b, "| %s | %s | [L%d](%s#L%d) | %s | %s |\n",
			e.id, e.title, e.line, solutionPath, e.line,
			shortComplexity(e.timeC), shortComplexity(e.spaceC))
	}
	return strings.TrimRight(b.String(), "\n")
}

// shortComplexity 把「O(n)，其中 n 为…」这类带解释的复杂度收敛成「O(n)」，
// 表格里只保留渐进量级；完整说明在 solution.go 对应小节的注释里。
// 注意 O(capacity) 这种括号属于量级本身，不能截断，所以只在小括号深度为 0 处切分。
func shortComplexity(s string) string {
	s = strings.TrimSpace(s)
	depth := 0
	for i, r := range s {
		switch r {
		case '(', '（':
			if r == '（' && depth == 0 {
				return trimPunct(s[:i])
			}
			depth++
		case ')', '）':
			if depth > 0 {
				depth--
			}
		case '，', '；', ',', ';', '。':
			if depth == 0 {
				return trimPunct(s[:i])
			}
		}
	}
	return trimPunct(s)
}

// trimPunct 去掉尾部多余的空格与标点。
func trimPunct(s string) string {
	return strings.TrimRight(strings.TrimSpace(s), " ，,；;。")
}

// updateReadme 用 table 替换 README 中两个标记之间的内容。
func updateReadme(table string) error {
	data, err := os.ReadFile(readmePath)
	if err != nil {
		return err
	}
	content := string(data)
	start := strings.Index(content, beginMarker)
	end := strings.Index(content, endMarker)
	if start < 0 || end < 0 || end < start {
		return fmt.Errorf("README.md 中缺少 %s / %s 标记", beginMarker, endMarker)
	}
	updated := content[:start+len(beginMarker)] + "\n\n" + table + "\n\n" + content[end:]
	return os.WriteFile(readmePath, []byte(updated), 0o644)
}
