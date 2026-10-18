# 牛客网面试必刷 TOP101 · Go 语言题解（单文件版）

<p>
  <img alt="Go" src="https://img.shields.io/badge/Go-1.21%2B-00ADD8?logo=go&logoColor=white">
  <img alt="Problems" src="https://img.shields.io/badge/题目-101%20题-success">
  <img alt="Tests" src="https://img.shields.io/badge/go%20test-passing-brightgreen">
  <img alt="Dependencies" src="https://img.shields.io/badge/dependencies-仅标准库-informational">
  <img alt="License" src="https://img.shields.io/badge/license-MIT-blue">
</p>

《牛客网 · 面试必刷 TOP101》的 Go 语言实现。

**全部 101 道题解写在根目录的同一个文件 [`solution.go`](solution.go) 里，全部测试写在同一个文件 [`solution_test.go`](solution_test.go) 里**，
不拆包、不分子目录，打开一个文件即可从头刷到尾；也可以按题号直接跳到对应小节。

每道题在文件中是一个独立小节，固定按「题目名 → 题目描述 → 示例 → 数据范围 → 解题思路 → 复杂度分析 → 代码」的顺序书写，
代码中复杂、易错的地方都加了中文注释；对应测试集中在 `solution_test.go` 中，`go test` 直接可跑。

| 规模 | 数量 |
| --- | --- |
| 题目 | 101 道（BM001–BM101，全部收录） |
| 题解文件 | 1 个（`solution.go`，约 5.1k 行，含题目与思路注释） |
| 测试文件 | 1 个（`solution_test.go`，约 10.1k 行） |
| 测试用例 | 1400+ 个（`go test ./...` 中实际运行的测试与子测试） |
| 第三方依赖 | 0（仅 Go 标准库） |

---

## ✨ 项目特点

- **单文件题解**：101 道题全部集中在 `solution.go`，用统一的 `// BMxxx 题目名` 小标题分隔，按题号顺序排列。
- **单文件测试**：101 个 `TestBMxxx` 集中在 `solution_test.go`，与题解小节一一对应，可用 `-run` 精确运行某一道题。
- **一题一入口**：每题对外暴露 `SolveBMxxx`（设计类题目为 `NewXxx` + 方法），签名在 [`docs/problem-list.md`](docs/problem-list.md) 中列明。
- **冲突已被处理**：`ListNode`、`TreeNode` 统一提到文件开头定义；同名的辅助函数按 `bm<题号>xxx` / `bm<题号>testXxx` 前缀区分，互不干扰。
- **零第三方依赖**：只使用 Go 标准库，`go.mod` 里没有任何 `require`。
- **索引自动生成**：`go run ./tools/genindex -write` 会根据题解中的小节标题刷新下面的题目索引表（含跳转行号）。

---

## 📁 目录结构

```
nc-top101/
├── solution.go          # ★ 101 道题解（单文件，按 BM001–BM101 顺序排列）
├── solution_test.go     # ★ 101 道题的测试（单文件，对应 TestBM001–TestBM101）
├── go.mod               # module github.com/yourname/nc-top101（无第三方依赖）
├── README.md            # 本文件
├── CONTRIBUTING.md      # 编写规范：小节格式、命名、测试
├── LICENSE              # MIT
├── docs/
│   └── problem-list.md  # 101 题清单（题号 / 题目 / 核心要求 / 关键约束）
├── tools/
│   └── genindex/        # README 题目索引生成器
└── .github/workflows/
    └── ci.yml           # gofmt / go vet / go test / 索引一致性
```

---

## 🚀 快速开始

环境要求：**Go 1.21+**（仅标准库，无需拉取依赖）。

```bash
git clone https://github.com/yourname/nc-top101.git
cd nc-top101

# 跑全部 101 题
go test ./... -count=1

# 只跑某一题（-v 可以看到每个子测试）
go test -run '^TestBM001$' -v

# 跑某一批题（正则匹配，例如 BM010–BM019）
go test -run '^TestBM01' -v

# 格式检查与静态检查
gofmt -l . && go vet ./...

# 重新生成 README 的题目索引（含行号链接）
go run ./tools/genindex -write
```

> 说明：`go.mod` 里的 module 路径使用了占位符 `github.com/yourname/nc-top101`，
> 推送到你自己的 GitHub 仓库后改成实际路径即可（包内没有跨模块 import，改路径不影响编译）。

---

## 📦 推送到 GitHub

```bash
# 1. 可选：把 module 路径改成自己的仓库地址
#    go.mod: module github.com/<你的用户名>/nc-top101
#    README.md: 同步替换 clone 链接

# 2. 初始化并提交
git init
git add .
git commit -m "feat: 牛客网 TOP101 Go 语言题解（单文件版，101 题）"
git branch -M main

# 3. 关联远程仓库并推送
git remote add origin https://github.com/<你的用户名>/nc-top101.git
git push -u origin main
```

仓库已内置 `.gitignore`（忽略构建缓存与 IDE 目录）和 `.github/workflows/ci.yml`
（push / PR 时自动执行 `gofmt` 检查、`go vet ./...`、`go test ./...`，并校验 README 索引是否过期）。

---

## 📚 题目索引

> 「位置」列可直接跳到 `solution.go` 中对应题目的起始行（GitHub 支持行号锚点）。

<!-- BEGIN:INDEX -->

| 题号 | 题目 | 位置 | 时间复杂度 | 空间复杂度 |
| --- | --- | --- | --- | --- |
| BM001 | 反转链表 | [L30](solution.go#L30) | O(n) | O(1) |
| BM002 | 链表内指定区间反转 | [L72](solution.go#L72) | O(n) | O(1) |
| BM003 | 链表中的节点每k个一组翻转 | [L129](solution.go#L129) | O(n) | O(1) |
| BM004 | 合并两个排序的链表 | [L201](solution.go#L201) | O(n) | O(1) |
| BM005 | 合并k个已排序的链表 | [L256](solution.go#L256) | O(n log k) | O(k) |
| BM006 | 判断链表中是否有环 | [L336](solution.go#L336) | O(n) | O(1) |
| BM007 | 链表中环的入口结点 | [L381](solution.go#L381) | O(n) | O(1) |
| BM008 | 链表中倒数最后k个结点 | [L445](solution.go#L445) | O(n) | O(1) |
| BM009 | 删除链表的倒数第n个节点 | [L502](solution.go#L502) | O(n) | O(1) |
| BM010 | 两个链表的第一个公共结点 | [L569](solution.go#L569) | O(n+m) | O(1) |
| BM011 | 链表相加(二) | [L633](solution.go#L633) | O(n+m) | O(1) |
| BM012 | 单链表的排序 | [L724](solution.go#L724) | O(n log n) | O(1) |
| BM013 | 判断一个链表是否为回文结构 | [L843](solution.go#L843) | O(n) | O(1) |
| BM014 | 链表的奇偶重排 | [L923](solution.go#L923) | O(n) | O(1) |
| BM015 | 删除有序链表中重复的元素-I | [L980](solution.go#L980) | O(n) | O(1) |
| BM016 | 删除有序链表中重复的元素-II | [L1025](solution.go#L1025) | O(n) | O(1) |
| BM017 | 二分查找-I | [L1084](solution.go#L1084) | O(log n) | O(1) |
| BM018 | 二维数组中的查找 | [L1135](solution.go#L1135) | O(m+n) | O(1) |
| BM019 | 寻找峰值 | [L1194](solution.go#L1194) | O(log n) | O(1) |
| BM020 | 数组中的逆序对 | [L1258](solution.go#L1258) | O(n log n) | O(n) |
| BM021 | 旋转数组的最小数字 | [L1353](solution.go#L1353) | O(log n) | O(1) |
| BM022 | 比较版本号 | [L1414](solution.go#L1414) | O(n + m) | O(1) |
| BM023 | 二叉树的前序遍历 | [L1520](solution.go#L1520) | O(n) | O(n) |
| BM024 | 二叉树的中序遍历 | [L1604](solution.go#L1604) | O(n) | O(n) |
| BM025 | 二叉树的后序遍历 | [L1688](solution.go#L1688) | O(n) | O(n) |
| BM026 | 求二叉树的层序遍历 | [L1741](solution.go#L1741) | O(n) | O(n) |
| BM027 | 按之字形顺序打印二叉树 | [L1792](solution.go#L1792) | O(n) | O(n) |
| BM028 | 二叉树的最大深度 | [L1852](solution.go#L1852) | O(n) | O(h) |
| BM029 | 二叉树中和为某一值的路径(一) | [L1895](solution.go#L1895) | O(n) | O(h) |
| BM030 | 二叉搜索树与双向链表 | [L1940](solution.go#L1940) | O(n) | O(h) |
| BM031 | 对称的二叉树 | [L2011](solution.go#L2011) | O(n) | O(h) |
| BM032 | 合并二叉树 | [L2077](solution.go#L2077) | O(n1+n2) | O(h1+h2) |
| BM033 | 二叉树的镜像 | [L2135](solution.go#L2135) | O(n) | O(h) |
| BM034 | 判断是不是二叉搜索树 | [L2189](solution.go#L2189) | O(n) | O(h) |
| BM035 | 判断是不是完全二叉树 | [L2250](solution.go#L2250) | O(n) | O(n) |
| BM036 | 判断是不是平衡二叉树 | [L2315](solution.go#L2315) | O(n) | O(h) |
| BM037 | 二叉搜索树的最近公共祖先 | [L2374](solution.go#L2374) | O(h) | O(1) |
| BM038 | 在二叉树中找到两个节点的最近公共祖先 | [L2423](solution.go#L2423) | O(n) | O(h) |
| BM039 | 序列化二叉树 | [L2484](solution.go#L2484) | O(n) | O(n) |
| BM040 | 重建二叉树 | [L2601](solution.go#L2601) | O(n) | O(n) |
| BM041 | 输出二叉树的右视图 | [L2666](solution.go#L2666) | O(n) | O(n) |
| BM042 | 用两个栈实现队列 | [L2764](solution.go#L2764) | O(1) 均摊 | O(n) |
| BM043 | 包含min函数的栈 | [L2867](solution.go#L2867) | O(1) | O(n) |
| BM044 | 有效括号序列 | [L2973](solution.go#L2973) | O(n) | O(n) |
| BM045 | 滑动窗口的最大值 | [L3033](solution.go#L3033) | O(n) | O(n) |
| BM046 | 最小的K个数 | [L3091](solution.go#L3091) | O(n + k log k) | O(n) |
| BM047 | 寻找第K大 | [L3183](solution.go#L3183) | O(n) | O(n) |
| BM048 | 数据流中的中位数 | [L3271](solution.go#L3271) | Insert 为 O(log n) | O(n) |
| BM049 | 表达式求值 | [L3422](solution.go#L3422) | O(n) | O(n) |
| BM050 | 两数之和 | [L3538](solution.go#L3538) | O(n) | O(n) |
| BM051 | 数组中出现次数超过一半的数字 | [L3594](solution.go#L3594) | O(n) | O(1) |
| BM052 | 数组中只出现一次的两个数字 | [L3651](solution.go#L3651) | O(n) | O(1) |
| BM053 | 缺失的第一个正整数 | [L3720](solution.go#L3720) | O(n) | O(1) |
| BM054 | 三数之和 | [L3776](solution.go#L3776) | O(n^2) | O(n) |
| BM055 | 没有重复项数字的全排列 | [L3869](solution.go#L3869) | O(n * n!) | O(n) |
| BM056 | 有重复项数字的全排列 | [L3943](solution.go#L3943) | O(n * n!) | O(n) |
| BM057 | 岛屿数量 | [L4019](solution.go#L4019) | O(m*n) | O(m*n) |
| BM058 | 字符串的排列 | [L4102](solution.go#L4102) | O(n * n!) | O(n) |
| BM059 | N皇后问题 | [L4175](solution.go#L4175) | O(n!) | O(n) |
| BM060 | 括号生成 | [L4245](solution.go#L4245) | O(4^n / sqrt(n)) | O(n) |
| BM061 | 矩阵最长递增路径 | [L4309](solution.go#L4309) | O(nm) | O(nm) |
| BM062 | 斐波那契数列 | [L4388](solution.go#L4388) | O(n) | O(1) |
| BM063 | 跳台阶 | [L4436](solution.go#L4436) | O(n) | O(1) |
| BM064 | 最小花费爬楼梯 | [L4478](solution.go#L4478) | O(n) | O(1) |
| BM065 | 最长公共子序列(二) | [L4533](solution.go#L4533) | O(nm + L(n + 256)) | O(nm) |
| BM066 | 最长公共子串 | [L4675](solution.go#L4675) | O(nm) | O(min(n, m)) |
| BM067 | 不同路径的数目(一) | [L4743](solution.go#L4743) | O(m*n) | O(n) |
| BM068 | 矩阵的最小路径和 | [L4795](solution.go#L4795) | O(n*m) | O(m) |
| BM069 | 把数字翻译成字符串 | [L4857](solution.go#L4857) | O(n) | O(n) |
| BM070 | 兑换零钱(一) | [L4923](solution.go#L4923) | O(n*aim) | O(aim) |
| BM071 | 最长上升子序列(一) | [L4991](solution.go#L4991) | O(n log n) | O(n) |
| BM072 | 连续子数组的最大和 | [L5051](solution.go#L5051) | O(n) | O(1) |
| BM073 | 最长回文子串 | [L5107](solution.go#L5107) | O(n^2) | O(1) |
| BM074 | 数字字符串转化成IP地址 | [L5174](solution.go#L5174) | O(3^4 * n) | O(n) |
| BM075 | 编辑距离(一) | [L5275](solution.go#L5275) | O(n*m) | O(m) |
| BM076 | 正则表达式匹配 | [L5353](solution.go#L5353) | O(n*m) | O(m) |
| BM077 | 最长的括号子串 | [L5443](solution.go#L5443) | O(n) | O(n) |
| BM078 | 打家劫舍(一) | [L5511](solution.go#L5511) | O(n) | O(1) |
| BM079 | 打家劫舍(二) | [L5570](solution.go#L5570) | O(n) | O(1) |
| BM080 | 买卖股票的最好时机(一) | [L5639](solution.go#L5639) | O(n) | O(1) |
| BM081 | 买卖股票的最好时机(二) | [L5699](solution.go#L5699) | O(n) | O(1) |
| BM082 | 买卖股票的最好时机(三) | [L5761](solution.go#L5761) | O(n) | O(1) |
| BM083 | 字符串变形 | [L5837](solution.go#L5837) | O(n) | O(n) |
| BM084 | 最长公共前缀 | [L5921](solution.go#L5921) | O(n * m) | O(1) |
| BM085 | 验证IP地址 | [L5972](solution.go#L5972) | O(n) | O(n) |
| BM086 | 大数加法 | [L6094](solution.go#L6094) | O(n) | O(n) |
| BM087 | 合并两个有序的数组 | [L6168](solution.go#L6168) | O(m+n) | O(1) |
| BM088 | 判断是否为回文字符串 | [L6229](solution.go#L6229) | O(n) | O(1) |
| BM089 | 合并区间 | [L6279](solution.go#L6279) | O(n log n) | O(n) |
| BM090 | 最小覆盖子串 | [L6348](solution.go#L6348) | O(|S| + |T|) | O(|S| + |T|) 级别的字符表规模 |
| BM091 | 反转字符串 | [L6437](solution.go#L6437) | O(n) | O(n) |
| BM092 | 最长无重复子数组 | [L6477](solution.go#L6477) | O(n) | O(n) |
| BM093 | 盛水最多的容器 | [L6532](solution.go#L6532) | O(n) | O(1) |
| BM094 | 接雨水问题 | [L6592](solution.go#L6592) | O(n) | O(1) |
| BM095 | 分糖果问题 | [L6659](solution.go#L6659) | O(n) | O(n) |
| BM096 | 主持人调度（二） | [L6731](solution.go#L6731) | O(n log n) | O(n) |
| BM097 | 旋转数组 | [L6826](solution.go#L6826) | O(n) | O(1) |
| BM098 | 螺旋矩阵 | [L6892](solution.go#L6892) | O(m*n) | O(m*n) |
| BM099 | 顺时针旋转矩阵 | [L6971](solution.go#L6971) | O(m*n) | O(m*n) |
| BM100 | 设计LRU缓存结构 | [L7058](solution.go#L7058) | O(1) | O(capacity) |
| BM101 | 设计LFU缓存结构 | [L7181](solution.go#L7181) | O(1) | O(capacity) |

<!-- END:INDEX -->

---

## 🗺️ 按专题刷题

| 专题 | 题号 | 说明 |
| --- | --- | --- |
| 链表 | BM001 – BM016 | 反转、区间反转、K 个一组、合并、环、公共结点、相加、排序、回文、去重 |
| 查找与排序 | BM017 – BM022 | 二分查找、二维数组查找、峰值、逆序对、旋转数组最小值、版本号比较 |
| 二叉树 | BM023 – BM041 | 四种遍历、深度、路径和、BST 转链表、镜像、BST/完全/平衡判定、LCA、序列化、重建、右视图 |
| 栈与队列 | BM042 – BM044 | 两个栈实现队列、最小栈、括号匹配 |
| 堆与 Top-K | BM045 – BM048 | 滑动窗口最大值、最小的 K 个数、第 K 大、数据流中位数 |
| 数组与查找技巧 | BM049 – BM054 | 表达式求值、两数之和、多数元素、只出现一次的两个数、缺失正整数、三数之和 |
| 回溯与搜索 | BM055 – BM061 | 全排列、岛屿数量、字符串排列、N 皇后、括号生成、矩阵最长递增路径 |
| 动态规划 | BM062 – BM082 | 斐波那契、跳台阶、爬楼梯、LCS/最长公共子串、路径、翻译、零钱、LIS、子数组和、编辑距离、正则、打家劫舍、股票 |
| 字符串 | BM083 – BM091 | 变形、公共前缀、IP 校验、大数加法、合并数组、回文、区间合并、最小覆盖子串、反转 |
| 双指针与贪心 | BM092 – BM097 | 最长无重复子数组、盛水容器、接雨水、分糖果、主持人调度、旋转数组 |
| 矩阵 | BM098 – BM099 | 螺旋矩阵、顺时针旋转矩阵 |
| 设计题 | BM100 – BM101 | LRU、LFU 缓存 |

---

## 🧩 单个题目长什么样

`solution.go` 中每道题都是这样一个自包含的小节（以 BM001 为例，节选）：

```go
// ---------------------------------------------------------------------------
// BM001 反转链表
//
// 题目描述：
// 给定一个单链表的头结点 pHead（该头结点是有值的），长度为 n，反转该链表后返回新链表的头结点。
//
// 示例：
// 输入：{1,2,3}
// 返回值：{3,2,1}
//
// 解题思路：
//  1. 三指针迭代：pre 指向已反转部分的头，cur 指向当前待处理结点。
//  2. 先用 next 暂存 cur.Next，再把 cur.Next 指向 pre，然后 pre/cur 同步后移。
//  3. 循环结束后 pre 即新链表头结点；空链表直接返回 nil。
//
// 时间复杂度：O(n)
// 空间复杂度：O(1)
func SolveBM001(head *ListNode) *ListNode {
	var pre *ListNode
	cur := head
	for cur != nil {
		next := cur.Next // 先保存后继，否则改指针后就丢了
		cur.Next = pre
		pre = cur
		cur = next
	}
	return pre
}
```

`solution_test.go` 中对应的测试：

```go
func TestBM001(t *testing.T) {
	t.Run("官方示例", func(t *testing.T) { ... })
	t.Run("空链表", func(t *testing.T) { ... })
}
```

文件开头的公共类型（101 道题共用）：

```go
// ListNode 单链表结点定义（全部链表题共用）。
type ListNode struct {
	Val  int
	Next *ListNode
}

// TreeNode 二叉树结点定义（全部二叉树题共用）。
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}
```

---

## 🤝 贡献与约定

新增或修改题目请先阅读 [CONTRIBUTING.md](CONTRIBUTING.md)，要点：

- 题解只写在根目录 `solution.go`，测试只写在 `solution_test.go`，不要新增其它 `.go` 文件；
- 每道题以 `// BMxxx 中文题目名` 开头，按「题目描述 → 示例 → 数据范围 → 解题思路 → 时间复杂度/空间复杂度 → 代码」的顺序书写；
- 对外入口统一 `SolveBMxxx`，私有辅助函数加 `bm<题号>` 前缀，测试辅助函数加 `bm<题号>test` 前缀，避免跨题同名冲突；
- 提交前通过 `gofmt -l .`、`go vet ./...`、`go test ./...`，并执行 `go run ./tools/genindex -write` 同步 README 索引。

CI 配置在 `.github/workflows/ci.yml`，推送后会自动执行上述检查。

---

## 📄 License

本项目代码采用 [MIT License](LICENSE)。

## ⚠️ 声明

题目描述与示例的版权归**牛客网**所有，本仓库仅为个人学习与面试复习整理，
如涉及版权问题请联系删除。
