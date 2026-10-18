package top101

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// BM001 反转链表 · 测试
// ---------------------------------------------------------------------------
// buildList 按切片顺序构造单链表，空切片返回 nil（仅测试使用）。
func buildList(vals []int) *ListNode {
	var head *ListNode
	for i := len(vals) - 1; i >= 0; i-- {
		head = &ListNode{Val: vals[i], Next: head}
	}
	return head
}

// listToSlice 把单链表转成切片，便于比较（仅测试使用，遇到环会死循环）。
func listToSlice(head *ListNode) []int {
	var out []int
	for cur := head; cur != nil; cur = cur.Next {
		out = append(out, cur.Val)
	}
	return out
}

func TestBM001(t *testing.T) {
	// 长链表用例：1..1000，反转后应为 1000..1
	long := make([]int, 0, 1000)
	wantLong := make([]int, 0, 1000)
	for i := 1; i <= 1000; i++ {
		long = append(long, i)
		wantLong = append(wantLong, 1001-i)
	}

	tests := []struct {
		name string
		in   []int
		want []int
	}{
		{"官方示例-三结点", []int{1, 2, 3}, []int{3, 2, 1}},
		{"官方示例-空链表", nil, nil},
		{"单结点", []int{7}, []int{7}},
		{"两结点", []int{1, 2}, []int{2, 1}},
		{"常规用例-五结点", []int{1, 2, 3, 4, 5}, []int{5, 4, 3, 2, 1}},
		{"含重复值", []int{1, 1, 2, 2, 3}, []int{3, 2, 2, 1, 1}},
		{"全零值", []int{0, 0, 0}, []int{0, 0, 0}},
		{"长链表-1000结点", long, wantLong},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := listToSlice(SolveBM001(buildList(tt.in)))
			if !slices.Equal(got, tt.want) {
				t.Fatalf("期望 %v，实际 %v", tt.want, got)
			}
		})
	}
	// 反转必须就地改指针，不能新建结点：倒序遍历原结点，逐个比对。
	t.Run("就地反转-复用原结点", func(t *testing.T) {
		nodes := make([]*ListNode, 0, 3)
		for i := 1; i <= 3; i++ {
			nodes = append(nodes, &ListNode{Val: i})
		}
		nodes[0].Next = nodes[1]
		nodes[1].Next = nodes[2]

		got := SolveBM001(nodes[0])
		cur := got
		for i := len(nodes) - 1; i >= 0; i-- {
			if cur != nodes[i] {
				t.Fatalf("第 %d 个位置期望复用结点 %d，实际不是", len(nodes)-1-i, nodes[i].Val)
			}
			cur = cur.Next
		}
		if cur != nil {
			t.Fatalf("期望链表长度为 %d，实际更长", len(nodes))
		}
	})
}

// ---------------------------------------------------------------------------
// BM002 链表内指定区间反转 · 测试
// ---------------------------------------------------------------------------
// bm002testBuildList 按切片顺序构造单链表，空切片返回 nil（仅测试使用）。
func bm002testBuildList(vals []int) *ListNode {
	var head *ListNode
	for i := len(vals) - 1; i >= 0; i-- {
		head = &ListNode{Val: vals[i], Next: head}
	}
	return head
}

// bm002testListToSlice 把单链表转成切片，便于比较（仅测试使用）。
func bm002testListToSlice(head *ListNode) []int {
	var out []int
	for cur := head; cur != nil; cur = cur.Next {
		out = append(out, cur.Val)
	}
	return out
}

func TestBM002(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		m    int
		n    int
		want []int
	}{
		{"官方示例", []int{1, 2, 3, 4, 5}, 2, 4, []int{1, 4, 3, 2, 5}},
		{"官方示例-两结点", []int{1, 5}, 1, 2, []int{5, 1}},
		{"空链表", nil, 1, 1, nil},
		{"单结点-m等于n", []int{9}, 1, 1, []int{9}},
		{"m等于n不翻转", []int{1, 2, 3, 4, 5}, 3, 3, []int{1, 2, 3, 4, 5}},
		{"m等于1-翻转前缀", []int{1, 2, 3, 4, 5}, 1, 3, []int{3, 2, 1, 4, 5}},
		{"n等于长度-翻转后缀", []int{1, 2, 3, 4, 5}, 3, 5, []int{1, 2, 5, 4, 3}},
		{"m等于1且n等于长度-整体翻转", []int{1, 2, 3, 4, 5}, 1, 5, []int{5, 4, 3, 2, 1}},
		{"区间只含两个结点", []int{1, 2, 3, 4}, 2, 3, []int{1, 3, 2, 4}},
		{"区间在链表正中", []int{1, 2, 3, 4, 5, 6}, 3, 4, []int{1, 2, 4, 3, 5, 6}},
		{"含重复值与负值", []int{-1, 0, -1, 2}, 1, 4, []int{2, -1, 0, -1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := bm002testListToSlice(SolveBM002(bm002testBuildList(tt.in), tt.m, tt.n))
			if !slices.Equal(got, tt.want) {
				t.Fatalf("输入 %v, m=%d, n=%d 期望 %v，实际 %v", tt.in, tt.m, tt.n, tt.want, got)
			}
		})
	}
	// 区间反转应就地改指针：翻转后区间内的结点仍是原来的结点。
	t.Run("就地反转-复用原结点", func(t *testing.T) {
		nodes := make([]*ListNode, 0, 5)
		for i := 1; i <= 5; i++ {
			nodes = append(nodes, &ListNode{Val: i})
		}
		for i := 0; i < 4; i++ {
			nodes[i].Next = nodes[i+1]
		}

		got := SolveBM002(nodes[0], 2, 4)
		// 期望 1 -> 4 -> 3 -> 2 -> 5，且各位置结点复用原结点
		want := []*ListNode{nodes[0], nodes[3], nodes[2], nodes[1], nodes[4]}
		cur := got
		for i, w := range want {
			if cur == nil {
				t.Fatalf("第 %d 个位置期望结点值 %d，实际为 nil", i, w.Val)
			}
			if cur != w {
				t.Fatalf("第 %d 个位置期望结点值 %d，实际 %d", i, w.Val, cur.Val)
			}
			cur = cur.Next
		}
		if cur != nil {
			t.Fatalf("期望链表长度为 %d，实际更长", len(want))
		}
	})
}

// ---------------------------------------------------------------------------
// BM003 链表中的节点每k个一组翻转 · 测试
// ---------------------------------------------------------------------------
// bm003testBuildList 按切片顺序构造单链表，空切片返回 nil（仅测试使用）。
func bm003testBuildList(vals []int) *ListNode {
	var head *ListNode
	for i := len(vals) - 1; i >= 0; i-- {
		head = &ListNode{Val: vals[i], Next: head}
	}
	return head
}

// bm003testListToSlice 把单链表转成切片，便于比较（仅测试使用）。
func bm003testListToSlice(head *ListNode) []int {
	var out []int
	for cur := head; cur != nil; cur = cur.Next {
		out = append(out, cur.Val)
	}
	return out
}

func TestBM003(t *testing.T) {
	// 长链表用例：0..1999，每 100 个一组翻转
	long := make([]int, 0, 2000)
	wantLong := make([]int, 0, 2000)
	for i := 0; i < 2000; i++ {
		long = append(long, i)
	}
	for i := 0; i+100 <= 2000; i += 100 {
		for j := i + 99; j >= i; j-- {
			wantLong = append(wantLong, j)
		}
	}

	tests := []struct {
		name string
		in   []int
		k    int
		want []int
	}{
		{"官方示例-k等于2", []int{1, 2, 3, 4, 5}, 2, []int{2, 1, 4, 3, 5}},
		{"官方示例-k等于3", []int{1, 2, 3, 4, 5}, 3, []int{3, 2, 1, 4, 5}},
		{"空链表", nil, 3, nil},
		{"单结点-k等于1", []int{1}, 1, []int{1}},
		{"单结点-k大于长度", []int{1}, 5, []int{1}},
		{"k等于1不翻转", []int{1, 2, 3, 4, 5}, 1, []int{1, 2, 3, 4, 5}},
		{"k等于链表长度-整体翻转", []int{1, 2, 3, 4}, 4, []int{4, 3, 2, 1}},
		{"k大于链表长度-保持原样", []int{1, 2, 3}, 4, []int{1, 2, 3}},
		{"k大于链表长度-单结点", []int{1, 2}, 3, []int{1, 2}},
		{"恰好整除-两组", []int{1, 2, 3, 4}, 2, []int{2, 1, 4, 3}},
		{"末尾不足k个", []int{1, 2, 3, 4, 5, 6, 7}, 3, []int{3, 2, 1, 6, 5, 4, 7}},
		{"含重复值与零", []int{0, 0, 1, 1}, 2, []int{0, 0, 1, 1}},
		{"长链表-2000结点-k等于100", long, 100, wantLong},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := bm003testListToSlice(SolveBM003(bm003testBuildList(tt.in), tt.k))
			if !slices.Equal(got, tt.want) {
				t.Fatalf("输入 %v, k=%d 期望 %v，实际 %v", tt.in, tt.k, tt.want, got)
			}
		})
	}
	// 只能改指针，不能新建结点：翻转后的结点集合必须与原链表一致。
	t.Run("复用原结点", func(t *testing.T) {
		nodes := make([]*ListNode, 0, 6)
		for i := 1; i <= 6; i++ {
			nodes = append(nodes, &ListNode{Val: i})
		}
		for i := 0; i < 5; i++ {
			nodes[i].Next = nodes[i+1]
		}

		got := SolveBM003(nodes[0], 2)
		// 期望 2 -> 1 -> 4 -> 3 -> 6 -> 5
		want := []*ListNode{nodes[1], nodes[0], nodes[3], nodes[2], nodes[5], nodes[4]}
		cur := got
		for i, w := range want {
			if cur != w {
				t.Fatalf("第 %d 个位置期望结点值 %d，实际不是同一个结点", i, w.Val)
			}
			cur = cur.Next
		}
		if cur != nil {
			t.Fatalf("期望链表长度为 %d，实际更长", len(want))
		}
	})
}

// ---------------------------------------------------------------------------
// BM004 合并两个排序的链表 · 测试
// ---------------------------------------------------------------------------
// bm004testBuildList 按切片顺序构造单链表，空切片返回 nil（仅测试使用）。
func bm004testBuildList(vals []int) *ListNode {
	var head *ListNode
	for i := len(vals) - 1; i >= 0; i-- {
		head = &ListNode{Val: vals[i], Next: head}
	}
	return head
}

// bm004testListToSlice 把单链表转成切片，便于比较（仅测试使用）。
func bm004testListToSlice(head *ListNode) []int {
	var out []int
	for cur := head; cur != nil; cur = cur.Next {
		out = append(out, cur.Val)
	}
	return out
}

// isSorted 校验链表是否非递减（仅测试使用）。
func isSorted(head *ListNode) bool {
	for cur := head; cur != nil && cur.Next != nil; cur = cur.Next {
		if cur.Val > cur.Next.Val {
			return false
		}
	}
	return true
}

func TestBM004(t *testing.T) {
	tests := []struct {
		name  string
		list1 []int
		list2 []int
		want  []int
	}{
		{"官方示例", []int{1, 3, 5}, []int{2, 4, 6}, []int{1, 2, 3, 4, 5, 6}},
		{"两个空链表", nil, nil, nil},
		{"第一条为空", nil, []int{1, 2, 3}, []int{1, 2, 3}},
		{"第二条为空", []int{1, 2, 3}, nil, []int{1, 2, 3}},
		{"单结点合并", []int{2}, []int{1}, []int{1, 2}},
		{"长度不等-前短后长", []int{1, 2}, []int{3, 4, 5, 6}, []int{1, 2, 3, 4, 5, 6}},
		{"长度不等-前长后短", []int{1, 3, 5, 7}, []int{2, 4}, []int{1, 2, 3, 4, 5, 7}},
		{"含相等值", []int{1, 2, 2}, []int{2, 2, 3}, []int{1, 2, 2, 2, 2, 3}},
		{"含负值", []int{-5, -1, 0}, []int{-3, 4}, []int{-5, -3, -1, 0, 4}},
		{"交错合并", []int{1, 3, 5, 7, 9}, []int{2, 4, 6, 8, 10}, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := bm004testListToSlice(SolveBM004(bm004testBuildList(tt.list1), bm004testBuildList(tt.list2)))
			if !slices.Equal(got, tt.want) {
				t.Fatalf("合并 %v 与 %v 期望 %v，实际 %v", tt.list1, tt.list2, tt.want, got)
			}
		})
	}
	// 合并只能拼接原结点，不能新建结点。
	t.Run("拼接原结点", func(t *testing.T) {
		a := []*ListNode{{Val: 1}, {Val: 3}}
		b := []*ListNode{{Val: 2}, {Val: 4}}
		a[0].Next = a[1]
		b[0].Next = b[1]

		got := SolveBM004(a[0], b[0])
		want := []*ListNode{a[0], b[0], a[1], b[1]}
		cur := got
		for i, w := range want {
			if cur != w {
				t.Fatalf("第 %d 个位置期望结点值 %d，实际不是同一个结点", i, w.Val)
			}
			cur = cur.Next
		}
		if cur != nil {
			t.Fatalf("期望链表长度为 %d，实际更长", len(want))
		}
		if !isSorted(got) {
			t.Fatalf("合并结果 %v 不是递增有序的", bm004testListToSlice(got))
		}
	})
}

// ---------------------------------------------------------------------------
// BM005 合并k个已排序的链表 · 测试
// ---------------------------------------------------------------------------
// bm005testBuildList 按切片顺序构造单链表，空切片返回 nil（仅测试使用）。
func bm005testBuildList(vals []int) *ListNode {
	var head *ListNode
	for i := len(vals) - 1; i >= 0; i-- {
		head = &ListNode{Val: vals[i], Next: head}
	}
	return head
}

// bm005testListToSlice 把单链表转成切片，便于比较（仅测试使用）。
func bm005testListToSlice(head *ListNode) []int {
	var out []int
	for cur := head; cur != nil; cur = cur.Next {
		out = append(out, cur.Val)
	}
	return out
}

// buildLists 把二维切片批量转成链表数组（仅测试使用）。
func buildLists(matrix [][]int) []*ListNode {
	lists := make([]*ListNode, 0, len(matrix))
	for _, vals := range matrix {
		lists = append(lists, bm005testBuildList(vals))
	}
	return lists
}

func TestBM005(t *testing.T) {
	tests := []struct {
		name  string
		input [][]int
		want  []int
	}{
		{"官方示例-两条链表", [][]int{{1, 2, 3}, {4, 5, 6, 7}}, []int{1, 2, 3, 4, 5, 6, 7}},
		{"空输入", [][]int{}, nil},
		{"n等于0-全为nil", [][]int{nil, nil}, nil},
		{"只有一条链表", [][]int{{1, 2, 3}}, []int{1, 2, 3}},
		{"只有一条且为空", [][]int{nil}, nil},
		{"三条链表", [][]int{{1, 4, 7}, {2, 5, 8}, {3, 6, 9}}, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{"夹带空链表", [][]int{{1, 5}, nil, {2, 3}, nil}, []int{1, 2, 3, 5}},
		{"含重复值", [][]int{{1, 1}, {1, 2}, {2, 2}}, []int{1, 1, 1, 2, 2, 2}},
		{"含负值", [][]int{{-5, -1}, {-3, 0}, {-2, 4}}, []int{-5, -3, -2, -1, 0, 4}},
		{"长度悬殊", [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, {5}}, []int{1, 2, 3, 4, 5, 5, 6, 7, 8, 9, 10}},
		{"单结点链表若干", [][]int{{3}, {1}, {2}}, []int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := bm005testListToSlice(SolveBM005(buildLists(tt.input)))
			if !slices.Equal(got, tt.want) {
				t.Fatalf("输入 %v 期望 %v，实际 %v", tt.input, tt.want, got)
			}
		})
	}
	// 不应改写调用方传入的切片。
	t.Run("不改写入参切片", func(t *testing.T) {
		lists := buildLists([][]int{{1, 5}, {2, 6}, {3, 7}})
		orig := make([]*ListNode, len(lists))
		copy(orig, lists)

		SolveBM005(lists)

		for i := range lists {
			if lists[i] != orig[i] {
				t.Fatalf("第 %d 个位置的头指针被改写了", i)
			}
		}
	})

	// 固定随机种子的大规模用例：k 条链表、总结点数 1e5，
	// 校验结果是递增的、且与所有元素排序后的序列完全一致。
	t.Run("随机大规模-1000条链表", func(t *testing.T) {
		const (
			k       = 1000
			totalIn = 100000
		)
		rng := rand.New(rand.NewSource(20240601))

		matrix := make([][]int, 0, k)
		all := make([]int, 0, totalIn)
		perList := totalIn / k
		for i := 0; i < k; i++ {
			vals := make([]int, 0, perList)
			for j := 0; j < perList; j++ {
				v := rng.Intn(2001) - 1000 // 与题面一致：|val| <= 1000
				vals = append(vals, v)
				all = append(all, v)
			}
			slices.Sort(vals) // 每条链表自身递增
			matrix = append(matrix, vals)
		}
		slices.Sort(all)

		got := bm005testListToSlice(SolveBM005(buildLists(matrix)))
		if len(got) != len(all) {
			t.Fatalf("期望 %d 个结点，实际 %d 个", len(all), len(got))
		}
		if !slices.Equal(got, all) {
			for i := range got {
				if got[i] != all[i] {
					t.Fatalf("第 %d 个元素期望 %d，实际 %d", i, all[i], got[i])
				}
			}
		}
	})
}

// ---------------------------------------------------------------------------
// BM006 判断链表中是否有环 · 测试
// ---------------------------------------------------------------------------
// bm006testBuildList 按切片顺序构造单链表，空切片返回 nil（仅测试使用）。
func bm006testBuildList(vals []int) *ListNode {
	var head *ListNode
	for i := len(vals) - 1; i >= 0; i-- {
		head = &ListNode{Val: vals[i], Next: head}
	}
	return head
}

// buildCycleList 构造链表，并让尾结点指向下标为 pos 的结点（pos < 0 表示无环）。
// 仅测试使用。
func buildCycleList(vals []int, pos int) *ListNode {
	head := bm006testBuildList(vals)
	if head == nil || pos < 0 {
		return head
	}
	var entry, tail *ListNode
	cur := head
	for i := 0; cur != nil; i++ {
		if i == pos {
			entry = cur
		}
		tail = cur
		cur = cur.Next
	}
	tail.Next = entry // 收尾成环
	return head
}

func TestBM006(t *testing.T) {
	// 大规模无环链表：10000 个结点
	longNoCycle := make([]int, 0, 10000)
	for i := 0; i < 10000; i++ {
		longNoCycle = append(longNoCycle, i)
	}
	// 大规模有环链表：10000 个结点，尾结点连回下标 0
	longCycle := make([]int, 0, 10000)
	for i := 0; i < 10000; i++ {
		longCycle = append(longCycle, i)
	}

	tests := []struct {
		name string
		in   []int
		pos  int
		want bool
	}{
		{"官方示例-有环", []int{3, 2, 0, -4}, 1, true},
		{"单结点无环", []int{1}, -1, false},
		{"空链表", nil, -1, false},
		{"单结点自环", []int{1}, 0, true},
		{"两结点无环", []int{1, 2}, -1, false},
		{"两结点成环", []int{1, 2}, 0, true},
		{"常规用例-无环", []int{1, 2, 3, 4, 5}, -1, false},
		{"常规用例-尾连头", []int{1, 2, 3, 4, 5}, 0, true},
		{"环在中间", []int{1, 2, 3, 4, 5}, 2, true},
		{"环入口是尾结点", []int{1, 2, 3, 4, 5}, 4, true},
		{"长链表-无环", longNoCycle, -1, false},
		{"长链表-有环", longCycle, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SolveBM006(buildCycleList(tt.in, tt.pos))
			if got != tt.want {
				t.Fatalf("输入 %v, pos=%d 期望 %v，实际 %v", tt.in, tt.pos, tt.want, got)
			}
		})
	}
	// 固定种子的随机链表：随机决定是否有环、随机环入口。
	t.Run("随机链表-200组", func(t *testing.T) {
		rng := rand.New(rand.NewSource(20240602))

		for i := 0; i < 200; i++ {
			n := rng.Intn(50) // 0..49 个结点
			vals := make([]int, 0, n)
			for j := 0; j < n; j++ {
				vals = append(vals, rng.Intn(100000))
			}

			pos := -1
			want := false
			if n > 0 && rng.Intn(2) == 0 {
				pos = rng.Intn(n)
				want = true
			}

			got := SolveBM006(buildCycleList(vals, pos))
			if got != want {
				t.Fatalf("第 %d 组：输入 %v, pos=%d 期望 %v，实际 %v", i, vals, pos, want, got)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// BM007 链表中环的入口结点 · 测试
// ---------------------------------------------------------------------------
// bm007testBuildCycleList 由 vals 按顺序构造单链表。
// 若 pos >= 0，则把表尾结点的 Next 指向第 pos 个结点（下标从 0 开始），使链表成环；
// 返回表头以及「环的入口结点」（无环时为 nil）。pos < 0 表示无环。
func bm007testBuildCycleList(vals []int, pos int) (head, entry *ListNode) {
	if len(vals) == 0 {
		return nil, nil
	}
	nodes := make([]*ListNode, len(vals))
	for i, v := range vals {
		nodes[i] = &ListNode{Val: v}
	}
	for i := 0; i+1 < len(nodes); i++ {
		nodes[i].Next = nodes[i+1]
	}
	if pos >= 0 {
		nodes[len(nodes)-1].Next = nodes[pos]
		entry = nodes[pos]
	}
	return nodes[0], entry
}

func TestBM007(t *testing.T) {
	cases := []struct {
		name string
		vals []int
		pos  int
	}{
		{"官方示例1_尾结点指回值为3的结点", []int{1, 2, 3, 4, 5}, 2},
		{"官方示例2_无环单结点", []int{1}, -1},
		{"官方示例3_单结点自环", []int{2}, 0},
		{"空链表", nil, -1},
		{"无环多结点", []int{1, 2, 3, 4, 5, 6}, -1},
		{"环入口即头结点", []int{1, 2, 3, 4}, 0},
		{"两结点尾首成环", []int{1, 2}, 0},
		{"两结点尾部自环", []int{1, 2}, 1},
		{"单结点尾自环", []int{7}, 0},
		{"环入口在倒数第二个", []int{1, 2, 3, 4}, 2},
		{"长链表环在尾部", []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}, 9},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			head, want := bm007testBuildCycleList(c.vals, c.pos)
			got := SolveBM007(head)
			if got != want {
				t.Fatalf("期望入口结点 %v，实际 %v", want, got)
			}
			if want != nil && got.Val != want.Val {
				t.Fatalf("期望入口结点值 %d，实际 %d", want.Val, got.Val)
			}
		})
	}

	// 随机用例：固定种子，覆盖「有环」「无环」两类形状
	t.Run("随机链表", func(t *testing.T) {
		r := rand.New(rand.NewSource(20240507))
		for i := 0; i < 30; i++ {
			n := r.Intn(60) + 1
			vals := make([]int, n)
			for j := range vals {
				vals[j] = j + 1
			}
			hasCycle := i%2 == 0
			pos := -1
			if hasCycle {
				pos = r.Intn(n)
			}
			head, want := bm007testBuildCycleList(vals, pos)
			got := SolveBM007(head)
			if got != want {
				t.Fatalf("第 %d 组（n=%d, pos=%d）：期望 %v，实际 %v", i, n, pos, want, got)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// BM008 链表中倒数最后k个结点 · 测试
// ---------------------------------------------------------------------------
// bm008testBuildList 由 vals 按顺序构造单链表，返回表头。
func bm008testBuildList(vals []int) *ListNode {
	var head, tail *ListNode
	for _, v := range vals {
		n := &ListNode{Val: v}
		if head == nil {
			head = n
		} else {
			tail.Next = n
		}
		tail = n
	}
	return head
}

// bm008testListToSlice 把链表转成切片，用于断言（题目输入无环，不会死循环）。
func bm008testListToSlice(head *ListNode) []int {
	out := []int{}
	for p := head; p != nil; p = p.Next {
		out = append(out, p.Val)
	}
	return out
}

// equalInts 比较两个 int 切片是否完全相等。
func equalInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestBM008(t *testing.T) {
	cases := []struct {
		name string
		vals []int
		k    int
		want []int // 期望返回值从该结点到表尾的内容；nil 表示期望返回 nil
	}{
		{"官方示例1_倒数第2个", []int{1, 2, 3, 4, 5}, 2, []int{4, 5}},
		{"官方示例2_单结点k为1", []int{2}, 1, []int{2}},
		{"官方示例3_k大于长度", []int{1, 2}, 3, nil},
		{"空链表", nil, 1, nil},
		{"空链表k为0", nil, 0, nil},
		{"k等于长度返回头结点", []int{1, 2, 3, 4, 5}, 5, []int{1, 2, 3, 4, 5}},
		{"k为1返回尾结点", []int{1, 2, 3, 4, 5}, 1, []int{5}},
		{"k为0", []int{1, 2, 3}, 0, nil},
		{"k超过长度很多", []int{7, 8}, 100, nil},
		{"含0与极值", []int{0, 1000000000, 0}, 2, []int{1000000000, 0}},
		{"两结点倒数第二", []int{1, 2}, 2, []int{1, 2}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			head := bm008testBuildList(c.vals)
			got := bm008testListToSlice(SolveBM008(head, c.k))
			if c.want == nil {
				if len(got) != 0 {
					t.Fatalf("期望返回空链表，实际 %v", got)
				}
				return
			}
			if !equalInts(got, c.want) {
				t.Fatalf("期望 %v，实际 %v", c.want, got)
			}
		})
	}

	// 随机用例：用下标作为结点值，便于核对结果
	t.Run("随机链表", func(t *testing.T) {
		r := rand.New(rand.NewSource(20240508))
		for i := 0; i < 30; i++ {
			n := r.Intn(80) + 1
			vals := make([]int, n)
			for j := range vals {
				vals[j] = j
			}
			k := r.Intn(n+3) + 1 // 可能大于长度
			got := bm008testListToSlice(SolveBM008(bm008testBuildList(vals), k))
			if k > n {
				if len(got) != 0 {
					t.Fatalf("第 %d 组（n=%d, k=%d）：期望空链表，实际 %v", i, n, k, got)
				}
				continue
			}
			want := vals[n-k:]
			if !equalInts(got, want) {
				t.Fatalf("第 %d 组（n=%d, k=%d）：期望 %v，实际 %v", i, n, k, want, got)
			}
		}
	})

	// 大链表：验证 k 接近长度时一次遍历的正确性
	t.Run("大规模链表", func(t *testing.T) {
		const n = 100000
		vals := make([]int, n)
		for i := range vals {
			vals[i] = i
		}
		got := bm008testListToSlice(SolveBM008(bm008testBuildList(vals), 99999))
		if len(got) != 99999 || got[0] != 1 || got[len(got)-1] != n-1 {
			t.Fatalf("大规模链表结果不正确：长度 %d", len(got))
		}
	})
}

// ---------------------------------------------------------------------------
// BM009 删除链表的倒数第n个节点 · 测试
// ---------------------------------------------------------------------------
// bm009testBuildList 由 vals 按顺序构造单链表，返回表头。
func bm009testBuildList(vals []int) *ListNode {
	var head, tail *ListNode
	for _, v := range vals {
		n := &ListNode{Val: v}
		if head == nil {
			head = n
		} else {
			tail.Next = n
		}
		tail = n
	}
	return head
}

// bm009testListToSlice 把链表转成切片，用于断言。
func bm009testListToSlice(head *ListNode) []int {
	out := []int{}
	for p := head; p != nil; p = p.Next {
		out = append(out, p.Val)
	}
	return out
}

// bm009testEqualInts 比较两个 int 切片是否完全相等。
func bm009testEqualInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestBM009(t *testing.T) {
	cases := []struct {
		name string
		vals []int
		n    int
		want []int
	}{
		{"官方示例1_删除倒数第2个", []int{1, 2, 3, 4, 5}, 2, []int{1, 2, 3, 5}},
		{"官方示例2_删除唯一结点", []int{1}, 1, []int{}},
		{"官方示例3_删除倒数第1个", []int{1, 2}, 1, []int{1}},
		{"空链表", nil, 1, []int{}},
		{"删除头结点", []int{1, 2, 3}, 3, []int{2, 3}},
		{"删除尾结点", []int{1, 2, 3}, 1, []int{1, 2}},
		{"n等于长度_两个结点", []int{1, 2}, 2, []int{2}},
		{"删除中间的结点", []int{5, 4, 3, 2, 1}, 3, []int{5, 4, 2, 1}},
		{"n非法大于长度", []int{1, 2}, 5, []int{1, 2}},
		{"n为0不删除", []int{1, 2, 3}, 0, []int{1, 2, 3}},
		{"值全相同", []int{7, 7, 7, 7}, 2, []int{7, 7, 7}},
		{"两结点删倒数第二", []int{9, 8}, 2, []int{8}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := bm009testListToSlice(SolveBM009(bm009testBuildList(c.vals), c.n))
			if !bm009testEqualInts(got, c.want) {
				t.Fatalf("期望 %v，实际 %v", c.want, got)
			}
		})
	}

	// 随机用例：删除第 n 个（倒数），与切片模拟结果对比
	t.Run("随机链表", func(t *testing.T) {
		r := rand.New(rand.NewSource(20240509))
		for i := 0; i < 30; i++ {
			n := r.Intn(60) + 1
			vals := make([]int, n)
			for j := range vals {
				vals[j] = j
			}
			k := r.Intn(n) + 1
			want := make([]int, 0, n-1)
			want = append(want, vals[:n-k]...)
			want = append(want, vals[n-k+1:]...)
			got := bm009testListToSlice(SolveBM009(bm009testBuildList(vals), k))
			if !bm009testEqualInts(got, want) {
				t.Fatalf("第 %d 组（n=%d, k=%d）：期望 %v，实际 %v", i, n, k, want, got)
			}
		}
	})

	// 大链表：删除倒数第 1 个，验证 O(n) 一次遍历
	t.Run("大规模链表", func(t *testing.T) {
		const size = 100000
		vals := make([]int, size)
		for i := range vals {
			vals[i] = i
		}
		got := bm009testListToSlice(SolveBM009(bm009testBuildList(vals), 1))
		if len(got) != size-1 || got[len(got)-1] != size-2 {
			t.Fatalf("大规模链表结果不正确：长度 %d", len(got))
		}
	})
}

// ---------------------------------------------------------------------------
// BM010 两个链表的第一个公共结点 · 测试
// ---------------------------------------------------------------------------
// chain 由 vals 按顺序构造单链表，返回表头与表尾（vals 为空时二者都是 nil）。
func chain(vals []int) (head, tail *ListNode) {
	for _, v := range vals {
		n := &ListNode{Val: v}
		if head == nil {
			head = n
		} else {
			tail.Next = n
		}
		tail = n
	}
	return head, tail
}

// buildSharedLists 构造两条共尾的单链表：
// head1 = prefixA + common，head2 = prefixB + common，两者共享同一段 common 结点。
// 返回两个表头与公共部分的首结点（common 为空时公共首结点为 nil）。
func buildSharedLists(prefixA, prefixB, common []int) (head1, head2, commonHead *ListNode) {
	commonHead, _ = chain(common)
	aHead, aTail := chain(prefixA)
	bHead, bTail := chain(prefixB)
	if aHead == nil {
		aHead = commonHead
	} else {
		aTail.Next = commonHead
	}
	if bHead == nil {
		bHead = commonHead
	} else {
		bTail.Next = commonHead
	}
	return aHead, bHead, commonHead
}

func TestBM010(t *testing.T) {
	cases := []struct {
		name          string
		prefixA       []int
		prefixB       []int
		common        []int
		wantCommonNil bool
	}{
		{"官方示例1_A短B长共享尾段", []int{1, 2, 3}, []int{4, 5}, []int{6, 7}, false},
		{"官方示例2_无公共结点", []int{1}, []int{2, 3}, nil, true},
		{"A比B长", []int{1, 2, 3, 4, 5}, []int{6}, []int{7, 8}, false},
		{"两链表完全相同", nil, nil, []int{1, 2, 3}, false},
		{"公共部分只有一个结点", []int{1, 2}, []int{3, 4, 5}, []int{9}, false},
		{"空链表A", nil, []int{1, 2}, nil, true},
		{"空链表B", []int{1, 2}, nil, nil, true},
		{"两条都为空", nil, nil, nil, true},
		{"长度相同且公共", []int{1}, []int{2}, []int{3, 4, 5}, false},
		{"值重复但无公共结点", []int{1, 1, 1}, []int{1, 1, 1}, nil, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			h1, h2, want := buildSharedLists(c.prefixA, c.prefixB, c.common)
			got := SolveBM010(h1, h2)
			if c.wantCommonNil {
				if got != nil {
					t.Fatalf("期望 nil，实际 %v", got)
				}
				return
			}
			if got != want {
				t.Fatalf("期望结点 %v，实际 %v（必须比较指针而不是值）", want, got)
			}
			if got.Val != c.common[0] {
				t.Fatalf("期望公共结点值 %d，实际 %d", c.common[0], got.Val)
			}
		})
	}

	// 随机用例：随机前缀长度，公共部分随机，并额外验证「公共结点之后再无分叉」
	t.Run("随机链表", func(t *testing.T) {
		r := rand.New(rand.NewSource(20240510))
		for i := 0; i < 30; i++ {
			la := r.Intn(30)
			lb := r.Intn(30)
			lc := r.Intn(30)
			prefixA := make([]int, la)
			prefixB := make([]int, lb)
			common := make([]int, lc)
			for j := range prefixA {
				prefixA[j] = 1000 + j
			}
			for j := range prefixB {
				prefixB[j] = 2000 + j
			}
			for j := range common {
				common[j] = 3000 + j
			}
			h1, h2, want := buildSharedLists(prefixA, prefixB, common)
			got := SolveBM010(h1, h2)
			if got != want {
				t.Fatalf("第 %d 组（la=%d, lb=%d, lc=%d）：期望 %v，实际 %v", i, la, lb, lc, want, got)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// BM011 链表相加(二) · 测试
// ---------------------------------------------------------------------------
// bm011testBuildList 由 digits 构造单链表（digits 高位在前），返回表头。
func bm011testBuildList(digits []int) *ListNode {
	var head, tail *ListNode
	for _, d := range digits {
		n := &ListNode{Val: d}
		if head == nil {
			head = n
		} else {
			tail.Next = n
		}
		tail = n
	}
	return head
}

// bm011testListToSlice 把链表转成切片（高位在前），用于断言。
func bm011testListToSlice(head *ListNode) []int {
	out := []int{}
	for p := head; p != nil; p = p.Next {
		out = append(out, p.Val)
	}
	return out
}

// bm011testEqualInts 比较两个 int 切片是否完全相等。
func bm011testEqualInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// addDigits 是测试内的参考实现：按竖式加法求和（高位在前的切片），用来交叉验证。
func addDigits(a, b []int) []int {
	i, j := len(a)-1, len(b)-1
	carry := 0
	rev := []int{}
	for i >= 0 || j >= 0 || carry != 0 {
		s := carry
		if i >= 0 {
			s += a[i]
			i--
		}
		if j >= 0 {
			s += b[j]
			j--
		}
		carry = s / 10
		rev = append(rev, s%10)
	}
	out := make([]int, 0, len(rev))
	for k := len(rev) - 1; k >= 0; k-- {
		out = append(out, rev[k])
	}
	for len(out) > 1 && out[0] == 0 {
		out = out[1:]
	}
	return out
}

func TestBM011(t *testing.T) {
	cases := []struct {
		name string
		a    []int
		b    []int
		want []int
	}{
		{"官方示例1_937加63", []int{9, 3, 7}, []int{6, 3}, []int{1, 0, 0, 0}},
		{"官方示例2_0加63", []int{0}, []int{6, 3}, []int{6, 3}},
		{"进位溢出_99加1", []int{9, 9}, []int{1}, []int{1, 0, 0}},
		{"进位连续_999加1", []int{9, 9, 9}, []int{1}, []int{1, 0, 0, 0}},
		{"零加零", []int{0}, []int{0}, []int{0}},
		{"长度相差很多", []int{1}, []int{9, 9, 9, 9, 9}, []int{1, 0, 0, 0, 0, 0}},
		{"输入带前导零", []int{0, 0, 1}, []int{2}, []int{3}},
		{"输入全是前导零", []int{0, 0, 0}, []int{0}, []int{0}},
		{"两链表都为空", nil, nil, []int{}},
		{"只有一条为空", nil, []int{1, 2, 3}, []int{1, 2, 3}},
		{"另一条为空", []int{4, 5}, nil, []int{4, 5}},
		{"无进位等长", []int{1, 2, 3}, []int{4, 5, 6}, []int{5, 7, 9}},
		{"结果比两条都长", []int{5}, []int{5}, []int{1, 0}},
		{"中间进位", []int{1, 9, 9}, []int{1}, []int{2, 0, 0}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			head1, head2 := bm011testBuildList(c.a), bm011testBuildList(c.b)
			got := bm011testListToSlice(SolveBM011(head1, head2))
			if !bm011testEqualInts(got, c.want) {
				t.Fatalf("期望 %v，实际 %v", c.want, got)
			}
			// 入参必须被还原成原顺序（实现只做临时反转）
			if !bm011testEqualInts(bm011testListToSlice(head1), c.a) {
				t.Fatalf("链表 head1 被修改：期望 %v，实际 %v", c.a, bm011testListToSlice(head1))
			}
			if !bm011testEqualInts(bm011testListToSlice(head2), c.b) {
				t.Fatalf("链表 head2 被修改：期望 %v，实际 %v", c.b, bm011testListToSlice(head2))
			}
		})
	}

	// 随机用例：与测试内参考实现交叉验证
	t.Run("随机链表", func(t *testing.T) {
		r := rand.New(rand.NewSource(20240511))
		for i := 0; i < 50; i++ {
			la := r.Intn(40) + 1
			lb := r.Intn(40) + 1
			a := make([]int, la)
			b := make([]int, lb)
			for j := range a {
				a[j] = r.Intn(10)
			}
			for j := range b {
				b[j] = r.Intn(10)
			}
			want := addDigits(a, b)
			got := bm011testListToSlice(SolveBM011(bm011testBuildList(a), bm011testBuildList(b)))
			if !bm011testEqualInts(got, want) {
				t.Fatalf("第 %d 组：%v + %v 期望 %v，实际 %v", i, a, b, want, got)
			}
		}
	})

	// 压力用例：位数 10^6（题目上限），验证迭代写法不会栈溢出
	t.Run("百万位相加不爆栈", func(t *testing.T) {
		const n = 1000000
		nines := make([]int, n)
		for i := range nines {
			nines[i] = 9
		}
		got := bm011testListToSlice(SolveBM011(bm011testBuildList(nines), bm011testBuildList([]int{1})))
		if len(got) != n+1 {
			t.Fatalf("期望结果长度 %d，实际 %d", n+1, len(got))
		}
		if got[0] != 1 {
			t.Fatalf("最高位期望 1，实际 %d", got[0])
		}
		for i := 1; i < len(got); i++ {
			if got[i] != 0 {
				t.Fatalf("第 %d 位期望 0，实际 %d", i, got[i])
			}
		}
	})
}

// ---------------------------------------------------------------------------
// BM012 单链表的排序 · 测试
// ---------------------------------------------------------------------------
// bm012testBuildList 由 vals 按顺序构造单链表，返回表头。
func bm012testBuildList(vals []int) *ListNode {
	var head, tail *ListNode
	for _, v := range vals {
		n := &ListNode{Val: v}
		if head == nil {
			head = n
		} else {
			tail.Next = n
		}
		tail = n
	}
	return head
}

// bm012testListToSlice 把链表转成切片，用于断言。
func bm012testListToSlice(head *ListNode) []int {
	out := []int{}
	for p := head; p != nil; p = p.Next {
		out = append(out, p.Val)
	}
	return out
}

// bm012testEqualInts 比较两个 int 切片是否完全相等。
func bm012testEqualInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// isAscending 校验切片是否非严格升序。
func isAscending(a []int) bool {
	for i := 1; i < len(a); i++ {
		if a[i-1] > a[i] {
			return false
		}
	}
	return true
}

func TestBM012(t *testing.T) {
	cases := []struct {
		name string
		vals []int
		want []int
	}{
		{"官方示例1", []int{1, 3, 2, 4, 5}, []int{1, 2, 3, 4, 5}},
		{"官方示例2_含负数", []int{-1, 0, 3, 2, 5}, []int{-1, 0, 2, 3, 5}},
		{"官方示例3_空链表", nil, []int{}},
		{"单结点", []int{42}, []int{42}},
		{"两结点逆序", []int{2, 1}, []int{1, 2}},
		{"两结点有序", []int{1, 2}, []int{1, 2}},
		{"三结点逆序", []int{3, 2, 1}, []int{1, 2, 3}},
		{"已经有序", []int{1, 2, 3, 4, 5, 6, 7, 8}, []int{1, 2, 3, 4, 5, 6, 7, 8}},
		{"完全逆序", []int{9, 8, 7, 6, 5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{"含大量重复", []int{3, 1, 3, 1, 2, 2, 3}, []int{1, 1, 2, 2, 3, 3, 3}},
		{"全部相同", []int{5, 5, 5, 5, 5}, []int{5, 5, 5, 5, 5}},
		{"极值", []int{1000000000, -1000000000, 0}, []int{-1000000000, 0, 1000000000}},
		{"奇数个元素", []int{7, 2, 9, 4, 1}, []int{1, 2, 4, 7, 9}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := bm012testListToSlice(SolveBM012(bm012testBuildList(c.vals)))
			if !bm012testEqualInts(got, c.want) {
				t.Fatalf("期望 %v，实际 %v", c.want, got)
			}
		})
	}

	// 随机用例：与标准库排序结果对比，同时校验升序性质
	t.Run("随机链表", func(t *testing.T) {
		r := rand.New(rand.NewSource(20240512))
		for i := 0; i < 50; i++ {
			n := r.Intn(200)
			vals := make([]int, n)
			for j := range vals {
				vals[j] = r.Intn(201) - 100 // 取值有重复、有正负
			}
			want := make([]int, n)
			copy(want, vals)
			sort.Ints(want)

			got := bm012testListToSlice(SolveBM012(bm012testBuildList(vals)))
			if !bm012testEqualInts(got, want) {
				t.Fatalf("第 %d 组（n=%d）期望 %v，实际 %v", i, n, want, got)
			}
		}
	})

	// 压力用例：10^5 个结点、逆序输入（最坏情况之一），验证 O(n log n) 不退化
	t.Run("十万结点逆序", func(t *testing.T) {
		const n = 100000
		vals := make([]int, n)
		for i := range vals {
			vals[i] = n - i
		}
		got := bm012testListToSlice(SolveBM012(bm012testBuildList(vals)))
		if len(got) != n {
			t.Fatalf("期望 %d 个结点，实际 %d", n, len(got))
		}
		if !isAscending(got) {
			t.Fatal("结果不是升序")
		}
		if got[0] != 1 || got[n-1] != n {
			t.Fatalf("结果首尾期望 1 和 %d，实际 %d 和 %d", n, got[0], got[n-1])
		}
	})
}

// ---------------------------------------------------------------------------
// BM013 判断一个链表是否为回文结构 · 测试
// ---------------------------------------------------------------------------
// bm013testBuildList 用值切片构造单链表，空切片返回 nil（仅测试使用）。
func bm013testBuildList(vals []int) *ListNode {
	var head *ListNode
	for i := len(vals) - 1; i >= 0; i-- {
		head = &ListNode{Val: vals[i], Next: head}
	}
	return head
}

// bm013testListToSlice 把链表转换成值切片，空链表返回 nil（仅测试使用）。
func bm013testListToSlice(head *ListNode) []int {
	var out []int
	for p := head; p != nil; p = p.Next {
		out = append(out, p.Val)
	}
	return out
}

// bm013testEqualInts 比较两个切片是否逐元素相等（nil 与空切片视为相等）。
func bm013testEqualInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// isPalindromeBySlice 用 O(n) 空间的朴素办法求正确答案，用于随机用例交叉验证。
func isPalindromeBySlice(vals []int) bool {
	for i, j := 0, len(vals)-1; i < j; i, j = i+1, j-1 {
		if vals[i] != vals[j] {
			return false
		}
	}
	return true
}

func TestBM013(t *testing.T) {
	cases := []struct {
		name string
		in   []int
		want bool
	}{
		{"示例1_单结点", []int{1}, true},
		{"示例2_两结点非回文", []int{2, 1}, false},
		{"示例3_偶数长度回文", []int{1, 2, 2, 1}, true},
		{"空链表", nil, true},
		{"两结点回文", []int{3, 3}, true},
		{"奇数长度回文", []int{1, 2, 3, 2, 1}, true},
		{"奇数长度非回文", []int{1, 2, 3, 1, 2}, false},
		{"奇数长度首尾不等", []int{1, 2, 3, 2, 2}, false},
		{"全部相同", []int{7, 7, 7, 7, 7, 7}, true},
		{"奇数长度仅一处不同", []int{1, 2, 3, 9, 4, 2, 1}, false},
		{"奇数长度回文含重复值", []int{1, 2, 3, 9, 3, 2, 1}, true},
		{"负数回文", []int{-5, -3, -5}, true},
		{"奇数长度只差一个", []int{1, 1, 1}, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			head := bm013testBuildList(c.in)
			got := SolveBM013(head)
			if got != c.want {
				t.Fatalf("输入 %v，期望 %v，实际 %v", c.in, c.want, got)
			}
			// 解法内部会反转后半段后再恢复，这里顺带校验入参链表结构未被破坏。
			if tail := bm013testListToSlice(head); !bm013testEqualInts(tail, c.in) {
				t.Fatalf("输入 %v 被修改，实际 %v", c.in, tail)
			}
		})
	}

	t.Run("随机用例", func(t *testing.T) {
		rnd := rand.New(rand.NewSource(20240613))
		for iter := 0; iter < 500; iter++ {
			vals := make([]int, rnd.Intn(13))
			for i := range vals {
				vals[i] = rnd.Intn(4) // 取值域很小，便于随机生成回文串
			}
			want := isPalindromeBySlice(vals)
			if got := SolveBM013(bm013testBuildList(vals)); got != want {
				t.Fatalf("第 %d 组输入 %v，期望 %v，实际 %v", iter, vals, want, got)
			}
		}
	})

	t.Run("大规模回文", func(t *testing.T) {
		const n = 20000
		vals := make([]int, 0, 2*n)
		for i := 0; i < n; i++ {
			vals = append(vals, i)
		}
		for i := n - 1; i >= 0; i-- {
			vals = append(vals, i)
		}
		if !SolveBM013(bm013testBuildList(vals)) {
			t.Fatalf("长度 %d 的回文链表被判为非回文", len(vals))
		}
		vals[0] = vals[len(vals)-1] + 1 // 只破坏首尾一处，应判为非回文
		if SolveBM013(bm013testBuildList(vals)) {
			t.Fatalf("长度 %d 的非回文链表被判为回文", len(vals))
		}
	})
}

// ---------------------------------------------------------------------------
// BM014 链表的奇偶重排 · 测试
// ---------------------------------------------------------------------------
// bm014testBuildList 用值切片构造单链表，空切片返回 nil（仅测试使用）。
func bm014testBuildList(vals []int) *ListNode {
	var head *ListNode
	for i := len(vals) - 1; i >= 0; i-- {
		head = &ListNode{Val: vals[i], Next: head}
	}
	return head
}

// bm014testListToSlice 把链表转换成值切片，空链表返回 nil（仅测试使用）。
func bm014testListToSlice(head *ListNode) []int {
	var out []int
	for p := head; p != nil; p = p.Next {
		out = append(out, p.Val)
	}
	return out
}

// bm014testEqualInts 比较两个切片是否逐元素相等（nil 与空切片视为相等）。
func bm014testEqualInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestBM014(t *testing.T) {
	cases := []struct {
		name string
		in   []int
		want []int
	}{
		{"示例1_偶数个结点", []int{1, 2, 3, 4, 5, 6}, []int{1, 3, 5, 2, 4, 6}},
		{"示例2_奇数个结点", []int{1, 4, 6, 3, 7}, []int{1, 6, 7, 4, 3}},
		{"示例3_空链表", nil, nil},
		{"单结点", []int{42}, []int{42}},
		{"两个结点", []int{1, 2}, []int{1, 2}},
		{"三个结点", []int{1, 2, 3}, []int{1, 3, 2}},
		{"全部相同", []int{5, 5, 5, 5, 5}, []int{5, 5, 5, 5, 5}},
		{"升序四个结点", []int{1, 2, 3, 4}, []int{1, 3, 2, 4}},
		{"含重复值", []int{2, 1, 3, 5, 6, 4, 7}, []int{2, 3, 6, 7, 1, 5, 4}},
		{"取值到上界", []int{1000000000, 0, 1000000000, 1}, []int{1000000000, 1000000000, 0, 1}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := bm014testListToSlice(SolveBM014(bm014testBuildList(c.in)))
			if !bm014testEqualInts(got, c.want) {
				t.Fatalf("输入 %v，期望 %v，实际 %v", c.in, c.want, got)
			}
		})
	}

	t.Run("结点对象未被重建", func(t *testing.T) {
		// 重排只能改指针、不能新建结点，这里校验返回链表复用的仍是原来的结点指针。
		in := []int{1, 2, 3, 4, 5}
		head := bm014testBuildList(in)
		origin := make(map[*ListNode]bool)
		for p := head; p != nil; p = p.Next {
			origin[p] = true
		}
		p := SolveBM014(head)
		count := 0
		for ; p != nil; p = p.Next {
			if !origin[p] {
				t.Fatalf("结果链表出现了新结点（Val=%d），重排必须复用原结点", p.Val)
			}
			count++
		}
		if count != len(in) {
			t.Fatalf("结果结点数期望 %d，实际 %d", len(in), count)
		}
	})
}

// ---------------------------------------------------------------------------
// BM015 删除有序链表中重复的元素-I · 测试
// ---------------------------------------------------------------------------
// bm015testBuildList 用值切片构造单链表，空切片返回 nil（仅测试使用）。
func bm015testBuildList(vals []int) *ListNode {
	var head *ListNode
	for i := len(vals) - 1; i >= 0; i-- {
		head = &ListNode{Val: vals[i], Next: head}
	}
	return head
}

// bm015testListToSlice 把链表转换成值切片，空链表返回 nil（仅测试使用）。
func bm015testListToSlice(head *ListNode) []int {
	var out []int
	for p := head; p != nil; p = p.Next {
		out = append(out, p.Val)
	}
	return out
}

// bm015testEqualInts 比较两个切片是否逐元素相等（nil 与空切片视为相等）。
func bm015testEqualInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestBM015(t *testing.T) {
	cases := []struct {
		name string
		in   []int
		want []int
	}{
		{"示例1_一对重复", []int{1, 1, 2}, []int{1, 2}},
		{"示例2_无重复", []int{1, 2, 3}, []int{1, 2, 3}},
		{"示例3_空链表", nil, nil},
		{"单结点", []int{1}, []int{1}},
		{"两结点相同", []int{1, 1}, []int{1}},
		{"多段重复", []int{1, 1, 2, 3, 3}, []int{1, 2, 3}},
		{"全部相同", []int{7, 7, 7, 7, 7}, []int{7}},
		{"重复在头尾", []int{0, 0, 1, 2, 2, 2, 9, 9}, []int{0, 1, 2, 9}},
		{"长重复段后接单值", []int{5, 5, 5, 5, 5, 6}, []int{5, 6}},
		{"取值到上界", []int{100, 100}, []int{100}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := bm015testListToSlice(SolveBM015(bm015testBuildList(c.in)))
			if !bm015testEqualInts(got, c.want) {
				t.Fatalf("输入 %v，期望 %v，实际 %v", c.in, c.want, got)
			}
		})
	}

	t.Run("就地删除复用原结点", func(t *testing.T) {
		// 要求就地删除：结果链表的每个结点都必须来自原链表。
		in := []int{1, 1, 2, 2, 2, 3}
		head := bm015testBuildList(in)
		origin := make(map[*ListNode]bool)
		for p := head; p != nil; p = p.Next {
			origin[p] = true
		}
		for p := SolveBM015(head); p != nil; p = p.Next {
			if !origin[p] {
				t.Fatalf("结果链表出现了新结点（Val=%d），去重必须就地删除", p.Val)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// BM016 删除有序链表中重复的元素-II · 测试
// ---------------------------------------------------------------------------
// bm016testBuildList 用值切片构造单链表，空切片返回 nil（仅测试使用）。
func bm016testBuildList(vals []int) *ListNode {
	var head *ListNode
	for i := len(vals) - 1; i >= 0; i-- {
		head = &ListNode{Val: vals[i], Next: head}
	}
	return head
}

// bm016testListToSlice 把链表转换成值切片，空链表返回 nil（仅测试使用）。
func bm016testListToSlice(head *ListNode) []int {
	var out []int
	for p := head; p != nil; p = p.Next {
		out = append(out, p.Val)
	}
	return out
}

// bm016testEqualInts 比较两个切片是否逐元素相等（nil 与空切片视为相等）。
func bm016testEqualInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestBM016(t *testing.T) {
	cases := []struct {
		name string
		in   []int
		want []int
	}{
		{"示例1_尾部重复", []int{1, 2, 2}, []int{1}},
		{"示例2_空链表", nil, nil},
		{"示例3_头部重复", []int{1, 1, 2}, []int{2}},
		{"单结点", []int{5}, []int{5}},
		{"全部相同", []int{3, 3, 3}, nil},
		{"两结点相同", []int{1, 1}, nil},
		{"两结点不同", []int{1, 2}, []int{1, 2}},
		{"示例_中间两段重复", []int{1, 2, 3, 3, 4, 4, 5}, []int{1, 2, 5}},
		{"示例_头部连续重复", []int{1, 1, 1, 2, 3}, []int{2, 3}},
		{"删掉头结点", []int{1, 1, 2, 3}, []int{2, 3}},
		{"删掉尾结点", []int{1, 2, 3, 3}, []int{1, 2}},
		{"重复段在中间", []int{1, 2, 2, 2, 3}, []int{1, 3}},
		{"交替重复", []int{1, 1, 2, 3, 3, 4, 4, 5, 5, 5}, []int{2}},
		{"取值到上界", []int{100, 100, 100}, nil},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := bm016testListToSlice(SolveBM016(bm016testBuildList(c.in)))
			if !bm016testEqualInts(got, c.want) {
				t.Fatalf("输入 %v，期望 %v，实际 %v", c.in, c.want, got)
			}
		})
	}

	t.Run("哨兵处理头结点被删除", func(t *testing.T) {
		// head 的值重复必须整段删除，此时返回的应是第一个不重复的结点。
		head := bm016testBuildList([]int{9, 9, 9, 1, 2})
		got := SolveBM016(head)
		if got == nil || got.Val != 1 {
			t.Fatalf("期望头结点值为 1，实际 %v", bm016testListToSlice(got))
		}
		if got.Next == nil || got.Next.Val != 2 || got.Next.Next != nil {
			t.Fatalf("期望链表为 {1,2}，实际 %v", bm016testListToSlice(got))
		}
	})
}

// ---------------------------------------------------------------------------
// BM017 二分查找-I · 测试
// ---------------------------------------------------------------------------
// linearSearch 线性扫描求正确答案，用于随机用例交叉验证（仅测试使用）。
func linearSearch(nums []int, target int) int {
	for i, v := range nums {
		if v == target {
			return i
		}
	}
	return -1
}

func TestBM017(t *testing.T) {
	cases := []struct {
		name   string
		nums   []int
		target int
		want   int
	}{
		{"示例1_命中中间", []int{-1, 0, 3, 4, 6, 10, 13, 14}, 6, 4},
		{"示例2_不存在", []int{-1, 0, 3, 4, 6, 10, 13, 14}, 2, -1},
		{"空数组", nil, 1, -1},
		{"单元素命中", []int{5}, 5, 0},
		{"单元素不命中", []int{5}, 4, -1},
		{"命中首元素", []int{1, 3, 5, 7}, 1, 0},
		{"命中尾元素", []int{1, 3, 5, 7}, 7, 3},
		{"小于所有元素", []int{1, 3, 5, 7}, 0, -1},
		{"大于所有元素", []int{1, 3, 5, 7}, 8, -1},
		{"落在两元素之间", []int{1, 3, 5, 7}, 4, -1},
		{"偶数长度左中", []int{1, 2, 3, 4}, 2, 1},
		{"负数数组", []int{-100, -50, -1, 0, 7}, -50, 1},
		{"极值", []int{-1000000000, 1000000000}, 1000000000, 1},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := SolveBM017(c.nums, c.target)
			if got != c.want {
				t.Fatalf("nums=%v target=%d，期望下标 %d，实际 %d", c.nums, c.target, c.want, got)
			}
		})
	}

	t.Run("规模上限2e5", func(t *testing.T) {
		const n = 200000
		nums := make([]int, n)
		for i := range nums {
			nums[i] = i*2 - 100000 // 递增且无重复，含负值
		}
		if got := SolveBM017(nums, nums[n-1]); got != n-1 {
			t.Fatalf("查找尾元素期望下标 %d，实际 %d", n-1, got)
		}
		if got := SolveBM017(nums, nums[0]); got != 0 {
			t.Fatalf("查找首元素期望下标 0，实际 %d", got)
		}
		if got := SolveBM017(nums, 1); got != -1 { // 1 是奇数，不在数组内
			t.Fatalf("查找不存在的元素期望 -1，实际 %d", got)
		}
	})

	t.Run("随机用例对照线性扫描", func(t *testing.T) {
		rnd := rand.New(rand.NewSource(20240617))
		for iter := 0; iter < 300; iter++ {
			n := rnd.Intn(50)
			nums := make([]int, 0, n)
			v := rnd.Intn(20) - 10
			for i := 0; i < n; i++ {
				nums = append(nums, v)
				v += rnd.Intn(3) + 1 // 严格递增
			}
			sort.Ints(nums) // 保险：保证有序且无重复
			target := rnd.Intn(30) - 15
			want := linearSearch(nums, target)
			if got := SolveBM017(nums, target); got != want {
				t.Fatalf("第 %d 组 nums=%v target=%d，期望 %d，实际 %d", iter, nums, target, want, got)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// BM018 二维数组中的查找 · 测试
// ---------------------------------------------------------------------------
// bruteForce 双重循环求正确答案，用于随机用例交叉验证（仅测试使用）。
func bruteForce(matrix [][]int, target int) bool {
	for _, row := range matrix {
		for _, v := range row {
			if v == target {
				return true
			}
		}
	}
	return false
}

func TestBM018(t *testing.T) {
	// 牛客官方示例使用的 4×4 矩阵：每行从左到右递增、每列从上到下递增。
	official := [][]int{
		{1, 2, 8, 9},
		{2, 4, 9, 12},
		{4, 7, 10, 13},
		{6, 8, 11, 15},
	}

	cases := []struct {
		name   string
		matrix [][]int
		target int
		want   bool
	}{
		{"示例1_命中7", official, 7, true},
		{"示例2_不含3", official, 3, false},
		{"示例3_空行矩阵", [][]int{{}}, 1, false},
		{"nil矩阵", nil, 1, false},
		{"空切片矩阵", [][]int{}, 1, false},
		{"1x1矩阵命中", [][]int{{5}}, 5, true},
		{"1x1矩阵不命中", [][]int{{5}}, 3, false},
		{"单行命中", [][]int{{1, 3, 5, 7}}, 7, true},
		{"单行不命中", [][]int{{1, 3, 5, 7}}, 4, false},
		{"单列命中", [][]int{{1}, {3}, {5}}, 3, true},
		{"单列不命中", [][]int{{1}, {3}, {5}}, 4, false},
		{"命中左上角", official, 1, true},
		{"命中右下角", official, 15, true},
		{"大于所有元素", official, 16, false},
		{"小于所有元素", official, 0, false},
		{"命中右上角", official, 9, true},
		{"命中左下角", official, 6, true},
		{"负数矩阵", [][]int{{-9, -5, -1}, {-8, -4, 0}, {-7, -3, 2}}, -4, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := SolveBM018(c.matrix, c.target)
			if got != c.want {
				t.Fatalf("matrix=%v target=%d，期望 %v，实际 %v", c.matrix, c.target, c.want, got)
			}
		})
	}

	t.Run("行列均含大量重复值", func(t *testing.T) {
		// 非严格递增（元素可重复）同样满足题目条件，检查单调性判断不会漏解。
		matrix := [][]int{
			{1, 1, 2, 2},
			{1, 2, 2, 3},
			{2, 2, 3, 3},
			{2, 3, 3, 4},
		}
		for target := 0; target <= 5; target++ {
			want := bruteForce(matrix, target)
			if got := SolveBM018(matrix, target); got != want {
				t.Fatalf("target=%d，期望 %v，实际 %v", target, want, got)
			}
		}
	})

	t.Run("随机用例对照双重循环", func(t *testing.T) {
		rnd := rand.New(rand.NewSource(20240618))
		for iter := 0; iter < 300; iter++ {
			n, m := rnd.Intn(6), rnd.Intn(6)
			matrix := make([][]int, n)
			for i := 0; i < n; i++ {
				row := make([]int, m)
				for j := 0; j < m; j++ {
					// (i*m + j) / 2 在同一行、同一列上都非递减，构造出合法的矩阵。
					row[j] = (i*m + j) / 2
				}
				matrix[i] = row
			}
			target := rnd.Intn(2*n*m+2) - 2
			want := bruteForce(matrix, target)
			if got := SolveBM018(matrix, target); got != want {
				t.Fatalf("第 %d 组 matrix=%v target=%d，期望 %v，实际 %v", iter, matrix, target, want, got)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// BM019 寻找峰值 · 测试
// ---------------------------------------------------------------------------
// checkPeak 校验 idx 是否为 nums 的合法峰值下标。
// 峰值位置不唯一，所以测试只校验「返回的下标确实是一个峰值」，不写死具体下标。
func checkPeak(t *testing.T, nums []int, idx int) {
	t.Helper()
	if idx < 0 || idx >= len(nums) {
		t.Fatalf("峰值下标越界：数组长度 %d，返回下标 %d", len(nums), idx)
	}
	if idx > 0 && nums[idx] <= nums[idx-1] {
		t.Fatalf("下标 %d 不是峰值：nums[%d]=%d 不大于左邻 nums[%d]=%d",
			idx, idx, nums[idx], idx-1, nums[idx-1])
	}
	if idx < len(nums)-1 && nums[idx] <= nums[idx+1] {
		t.Fatalf("下标 %d 不是峰值：nums[%d]=%d 不大于右邻 nums[%d]=%d",
			idx, idx, nums[idx], idx+1, nums[idx+1])
	}
}

func TestBM019(t *testing.T) {
	// 注意：题设要求「相邻元素不相等」，因此不构造全相同元素这类非法输入。
	cases := []struct {
		name string
		nums []int
	}{
		{name: "官方示例1_2_4_1_2_7_8_4", nums: []int{2, 4, 1, 2, 7, 8, 4}},
		{name: "官方示例2_1_2_3_1", nums: []int{1, 2, 3, 1}},
		{name: "峰值在中间_2_4_1_2_5_3", nums: []int{2, 4, 1, 2, 5, 3}},
		{name: "单元素", nums: []int{1}},
		{name: "峰值在左端_降序", nums: []int{3, 2, 1}},
		{name: "峰值在右端_升序", nums: []int{1, 2, 3, 4, 5}},
		{name: "两个元素_上升", nums: []int{1, 2}},
		{name: "两个元素_下降", nums: []int{2, 1}},
		{name: "中间是谷底_两端都是峰值", nums: []int{1, 0, 1}},
		{name: "包含负数", nums: []int{-5, -3, -4, -1, -2}},
		{name: "多个峰值", nums: []int{1, 3, 2, 4, 1}},
		{name: "极值_int32边界", nums: []int{-2147483648, 2147483647}},
	}

	for _, tc := range cases {
		tc := tc // Go 1.21 及更早版本的循环变量需要注意共享问题
		t.Run(tc.name, func(t *testing.T) {
			checkPeak(t, tc.nums, SolveBM019(tc.nums))
		})
	}

	t.Run("空切片返回-1", func(t *testing.T) {
		if got := SolveBM019(nil); got != -1 {
			t.Fatalf("期望 -1，实际 %d", got)
		}
	})

	t.Run("随机数组", func(t *testing.T) {
		r := rand.New(rand.NewSource(20240101))
		for round := 0; round < 200; round++ {
			n := r.Intn(50) + 1
			// 用两两不同的随机值构造，天然满足「相邻元素不相等」的题设。
			perm := r.Perm(n)
			nums := make([]int, n)
			for i, v := range perm {
				nums[i] = v*7 - 100
			}
			checkPeak(t, nums, SolveBM019(nums))
		}
	})

	t.Run("大规模_2e5个元素", func(t *testing.T) {
		r := rand.New(rand.NewSource(7))
		n := 200000
		nums := make([]int, n)
		for i, v := range r.Perm(n) {
			nums[i] = v
		}
		checkPeak(t, nums, SolveBM019(nums))
	})
}

// ---------------------------------------------------------------------------
// BM020 数组中的逆序对 · 测试
// ---------------------------------------------------------------------------
// bruteInversePairs 暴力统计逆序对，作为小规模用例的参照实现。
func bruteInversePairs(nums []int) int {
	var cnt int64
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i] > nums[j] {
				cnt++
			}
		}
	}
	return int(cnt % int64(inversePairsMod))
}

func TestBM020(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		want int
	}{
		{name: "官方示例_1", nums: []int{1, 2, 3, 4, 5, 6, 7, 0}, want: 7},
		{name: "官方示例_2", nums: []int{1, 2, 3, 0}, want: 3},
		{name: "空数组", nums: nil, want: 0},
		{name: "单元素", nums: []int{9}, want: 0},
		{name: "全相同元素", nums: []int{4, 4, 4, 4, 4}, want: 0},
		{name: "升序_无逆序对", nums: []int{1, 2, 3, 4, 5}, want: 0},
		{name: "降序_全部逆序", nums: []int{5, 4, 3, 2, 1}, want: 10},
		{name: "两个元素", nums: []int{2, 1}, want: 1},
		{name: "两段各自有序", nums: []int{4, 5, 6, 1, 2, 3}, want: 9},
		{name: "含零值且有重复", nums: []int{0, 1, 2, 0}, want: 2},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if got := SolveBM020(tc.nums); got != tc.want {
				t.Fatalf("期望 %d，实际 %d", tc.want, got)
			}
		})
	}

	t.Run("不修改入参", func(t *testing.T) {
		nums := []int{3, 1, 2, 5, 4, 0}
		snapshot := make([]int, len(nums))
		copy(snapshot, nums)
		SolveBM020(nums)
		for i := range snapshot {
			if nums[i] != snapshot[i] {
				t.Fatalf("下标 %d 的元素被修改：期望 %d，实际 %d", i, snapshot[i], nums[i])
			}
		}
	})

	t.Run("随机数组_与暴力解对照", func(t *testing.T) {
		r := rand.New(rand.NewSource(20240202))
		for round := 0; round < 300; round++ {
			n := r.Intn(60)
			nums := make([]int, n)
			for i := range nums {
				nums[i] = r.Intn(20) // 允许重复，验证等号不计入逆序对
			}
			want := bruteInversePairs(nums)
			if got := SolveBM020(nums); got != want {
				t.Fatalf("输入 %v：期望 %d，实际 %d", nums, want, got)
			}
		}
	})

	t.Run("大规模_1e5个元素_逆序对超过模数", func(t *testing.T) {
		const n = 100000
		nums := make([]int, n)
		for i := range nums {
			nums[i] = n - i // 完全降序，逆序对总数为 n*(n-1)/2 = 4999950000
		}
		// 期望值是 n*(n-1)/2 再对 1000000007 取模。
		want := int(int64(n) * int64(n-1) / 2 % int64(inversePairsMod))
		if got := SolveBM020(nums); got != want {
			t.Fatalf("期望 %d，实际 %d", want, got)
		}
	})
}

// ---------------------------------------------------------------------------
// BM021 旋转数组的最小数字 · 测试
// ---------------------------------------------------------------------------
func TestBM021(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		want int
	}{
		{name: "官方示例_3_4_5_1_2", nums: []int{3, 4, 5, 1, 2}, want: 1},
		{name: "官方示例_有重复_3_100_200_3", nums: []int{3, 100, 200, 3}, want: 3},
		{name: "未旋转_整体非降序", nums: []int{1, 2, 3, 4, 5}, want: 1},
		{name: "旋转1位", nums: []int{5, 1, 2, 3, 4}, want: 1},
		{name: "旋转n-1位", nums: []int{2, 3, 4, 5, 1}, want: 1},
		{name: "单元素", nums: []int{7}, want: 7},
		{name: "两个元素_已旋转", nums: []int{2, 1}, want: 1},
		{name: "两个元素_未旋转", nums: []int{1, 2}, want: 1},
		{name: "全相同元素", nums: []int{2, 2, 2, 2, 2}, want: 2},
		{name: "重复元素_1_1_1_0_1", nums: []int{1, 1, 1, 0, 1}, want: 0},
		{name: "重复元素_1_0_1_1_1", nums: []int{1, 0, 1, 1, 1}, want: 0},
		{name: "重复元素_2_2_0_2", nums: []int{2, 2, 0, 2}, want: 0},
		{name: "重复元素_3_3_1_3", nums: []int{3, 3, 1, 3}, want: 1},
		{name: "重复元素_未旋转", nums: []int{0, 1, 1, 1, 1}, want: 0},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if got := SolveBM021(tc.nums); got != tc.want {
				t.Fatalf("输入 %v：期望 %d，实际 %d", tc.nums, tc.want, got)
			}
		})
	}

	t.Run("空切片返回0", func(t *testing.T) {
		if got := SolveBM021(nil); got != 0 {
			t.Fatalf("期望 0，实际 %d", got)
		}
	})

	// 生成随机非降序数组并随机旋转（允许大量重复值），与线性扫描得到的最小值对照。
	t.Run("随机旋转数组_与线性扫描对照", func(t *testing.T) {
		r := rand.New(rand.NewSource(20240303))
		for round := 0; round < 500; round++ {
			n := r.Intn(40) + 1
			sorted := make([]int, n)
			v := r.Intn(10)
			for i := 0; i < n; i++ {
				v += r.Intn(3) // 递增但允许相等，构造非降序数组
				sorted[i] = v
			}
			rot := r.Intn(n) // 旋转位数，0 表示不旋转
			rotated := make([]int, 0, n)
			rotated = append(rotated, sorted[rot:]...)
			rotated = append(rotated, sorted[:rot]...)

			want := sorted[0] // 非降序数组的最小值必是首个元素，与是否旋转无关
			if got := SolveBM021(rotated); got != want {
				t.Fatalf("旋转 %d 位后输入 %v：期望 %d，实际 %d", rot, rotated, want, got)
			}
		}
	})

	t.Run("大规模_1e4个元素", func(t *testing.T) {
		const n = 10000
		nums := make([]int, 0, n)
		for i := 0; i < n/2; i++ {
			nums = append(nums, i+5000) // 前一段：5000 ~ 9999
		}
		for i := 0; i < n-n/2; i++ {
			nums = append(nums, i) // 后一段：0 ~ 4999
		}
		if got := SolveBM021(nums); got != 0 {
			t.Fatalf("期望 0，实际 %d", got)
		}
	})
}

// ---------------------------------------------------------------------------
// BM022 比较版本号 · 测试
// ---------------------------------------------------------------------------
// referenceCompareVersion 用「切分 + Atoi」实现的参照版本，仅用于随机用例对照。
// 随机用例中每段都是很小的数字，不会溢出 int，因此这里可以放心用 Atoi。
func referenceCompareVersion(v1, v2 string) int {
	toInts := func(s string) []int {
		parts := strings.Split(s, ".")
		out := make([]int, len(parts))
		for i, p := range parts {
			n, err := strconv.Atoi(p)
			if err != nil {
				return nil
			}
			out[i] = n
		}
		return out
	}
	a, b := toInts(v1), toInts(v2)
	if a == nil || b == nil || len(a) == 0 || len(b) == 0 {
		return 1 << 30 // 参照实现无法处理的输入，返回哨兵值让调用方跳过
	}
	n := len(a)
	if len(b) > n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		x, y := 0, 0
		if i < len(a) {
			x = a[i]
		}
		if i < len(b) {
			y = b[i]
		}
		if x < y {
			return -1
		}
		if x > y {
			return 1
		}
	}
	return 0
}

func TestBM022(t *testing.T) {
	cases := []struct {
		name   string
		v1, v2 string
		want   int
	}{
		{name: "官方示例_前导零", v1: "1.1", v2: "1.01", want: 0},
		{name: "官方示例_缺失段补0", v1: "1.0.1", v2: "1", want: 1},
		{name: "尾部补0相等", v1: "1.0", v2: "1.0.0", want: 0},
		{name: "第一段更小", v1: "0.1", v2: "1.1", want: -1},
		{name: "相同版本号", v1: "1.2.3", v2: "1.2.3", want: 0},
		{name: "单段整数", v1: "2", v2: "10", want: -1},
		{name: "段数不同且前者大", v1: "7.5.3", v2: "7.5.2.4", want: 1},
		{name: "段数不同且前者小", v1: "7.5.2.4", v2: "7.5.3", want: -1},
		{name: "多前导零相等", v1: "001.0", v2: "1", want: 0},
		{name: "前导零后仍更大", v1: "1.002", v2: "1.1", want: 1},
		{name: "最长500位数字_相等", v1: "1." + strings.Repeat("0", 497) + "2", v2: "1.2", want: 0},
		{name: "最长500位数字_更小", v1: "0." + strings.Repeat("9", 498), v2: "1", want: -1},
		{name: "超长数字_位数多者更大", v1: "1." + strings.Repeat("1", 300), v2: "1." + strings.Repeat("9", 299), want: 1},
		{name: "全部为零", v1: "0.0.0", v2: "0", want: 0},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if got := SolveBM022(tc.v1, tc.v2); got != tc.want {
				t.Fatalf("compare(%q, %q)：期望 %d，实际 %d", tc.v1, tc.v2, tc.want, got)
			}
			// 交换参数后结果应当取反，保证比较是反对称的。
			if got := SolveBM022(tc.v2, tc.v1); got != -tc.want {
				t.Fatalf("compare(%q, %q)：期望 %d，实际 %d", tc.v2, tc.v1, -tc.want, got)
			}
		})
	}

	t.Run("随机版本号_与参照实现对照", func(t *testing.T) {
		r := rand.New(rand.NewSource(20240404))
		randomVersion := func() string {
			segments := r.Intn(4) + 1
			parts := make([]string, segments)
			for i := range parts {
				val := r.Intn(1000)
				// 随机加前导零，覆盖 "001" 这类输入。
				parts[i] = strings.Repeat("0", r.Intn(3)) + strconv.Itoa(val)
			}
			return strings.Join(parts, ".")
		}
		for round := 0; round < 500; round++ {
			v1, v2 := randomVersion(), randomVersion()
			want := referenceCompareVersion(v1, v2)
			if want == 1<<30 {
				continue // 参照实现无法处理（理论上不会发生）
			}
			if got := SolveBM022(v1, v2); got != want {
				t.Fatalf("compare(%q, %q)：期望 %d，实际 %d", v1, v2, want, got)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// BM023 二叉树的前序遍历 · 测试
// ---------------------------------------------------------------------------
// buildChain 构造一条链状树：根的值是 1，第 k 个结点的值是 k，每个结点只有一个孩子。
// leftChild 为 true 时全部是左孩子，否则全部是右孩子。
func buildChain(n int, leftChild bool) *TreeNode {
	if n <= 0 {
		return nil
	}
	root := &TreeNode{Val: 1}
	cur := root
	for v := 2; v <= n; v++ {
		node := &TreeNode{Val: v}
		if leftChild {
			cur.Left = node
		} else {
			cur.Right = node
		}
		cur = node
	}
	return root
}

// buildFullTree 构造 {4,2,6,1,3,5,7} 这棵满二叉树（也恰好是一棵二叉搜索树）。
func buildFullTree() *TreeNode {
	return &TreeNode{
		Val: 4,
		Left: &TreeNode{
			Val:   2,
			Left:  &TreeNode{Val: 1},
			Right: &TreeNode{Val: 3},
		},
		Right: &TreeNode{
			Val:   6,
			Left:  &TreeNode{Val: 5},
			Right: &TreeNode{Val: 7},
		},
	}
}

// buildSameValueTree 构造所有结点值都相同的小树，用于确认遍历不会因值相同而出错。
func buildSameValueTree(val, n int) *TreeNode {
	if n <= 0 {
		return nil
	}
	root := &TreeNode{Val: val}
	cur := root
	for i := 1; i < n; i++ {
		node := &TreeNode{Val: val}
		if i%2 == 1 {
			cur.Left = node
		} else {
			cur.Right = node
		}
		cur = node
	}
	return root
}

// buildRandomTree 用固定种子随机挂载父子关系，构造一棵有 n 个结点的二叉树。
func buildRandomTree(r *rand.Rand, n int) *TreeNode {
	if n <= 0 {
		return nil
	}
	nodes := make([]*TreeNode, n)
	for i := range nodes {
		nodes[i] = &TreeNode{Val: i + 1}
	}
	for i := 1; i < n; i++ {
		for {
			parent := nodes[r.Intn(i)]
			if r.Intn(2) == 0 && parent.Left == nil {
				parent.Left = nodes[i]
				break
			}
			if parent.Right == nil {
				parent.Right = nodes[i]
				break
			}
		}
	}
	return nodes[0]
}

// assertInts 依次比较两个切片，失败时只打印前若干个元素，避免大用例刷屏。
func assertInts(t *testing.T, want, got []int) {
	t.Helper()
	if len(want) != len(got) {
		t.Fatalf("长度不符：期望 %d 个元素，实际 %d 个", len(want), len(got))
	}
	for i := range want {
		if want[i] != got[i] {
			t.Fatalf("第 %d 个元素不符：期望 %d，实际 %d（期望前 10 个 %v，实际前 10 个 %v）",
				i, want[i], got[i], head10(want), head10(got))
		}
	}
}

// head10 返回切片前 10 个元素，用于断言失败时的可读输出。
func head10(s []int) []int {
	if len(s) > 10 {
		return s[:10]
	}
	return s
}

// seq 返回 1,2,...,n。
func seq(n int) []int {
	out := make([]int, n)
	for i := range out {
		out[i] = i + 1
	}
	return out
}

func TestBM023(t *testing.T) {
	cases := []struct {
		name string
		root *TreeNode
		want []int
	}{
		{
			name: "官方示例_1_hash_2_3",
			root: &TreeNode{Val: 1, Right: &TreeNode{Val: 2, Left: &TreeNode{Val: 3}}},
			want: []int{1, 2, 3},
		},
		{name: "官方示例_空树", root: nil, want: []int{}},
		{name: "官方示例_单结点", root: &TreeNode{Val: 1}, want: []int{1}},
		{
			name: "只有左孩子的链状树",
			root: buildChain(5, true),
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "只有右孩子的链状树",
			root: buildChain(5, false),
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "满二叉树",
			root: buildFullTree(),
			want: []int{4, 2, 1, 3, 6, 5, 7},
		},
		{
			name: "全相同结点值",
			root: buildSameValueTree(5, 5),
			want: []int{5, 5, 5, 5, 5},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			assertInts(t, tc.want, SolveBM023(tc.root))
			// 递归版只作为对照：两种实现的结果必须完全一致。
			assertInts(t, tc.want, preorderRecursive(tc.root))
		})
	}

	t.Run("随机树_迭代与递归一致", func(t *testing.T) {
		r := rand.New(rand.NewSource(20240505))
		for round := 0; round < 100; round++ {
			root := buildRandomTree(r, r.Intn(50)+1)
			want := preorderRecursive(root)
			assertInts(t, want, SolveBM023(root))
			// 遍历结果必须是所有结点值的一个排列（结点值互不相同）。
			seen := make(map[int]bool, len(want))
			for _, v := range want {
				if seen[v] {
					t.Fatalf("结点值 %d 重复出现", v)
				}
				seen[v] = true
			}
		}
	})

	t.Run("大规模左链_1e5个结点", func(t *testing.T) {
		const n = 100000
		root := buildChain(n, true)
		assertInts(t, seq(n), SolveBM023(root))
		assertInts(t, seq(n), preorderRecursive(root))
	})

	t.Run("大规模右链_1e5个结点", func(t *testing.T) {
		const n = 100000
		root := buildChain(n, false)
		assertInts(t, seq(n), SolveBM023(root))
	})
}

// ---------------------------------------------------------------------------
// BM024 二叉树的中序遍历 · 测试
// ---------------------------------------------------------------------------
// bm024testBuildChain 构造一条链状树：根的值是 1，第 k 个结点的值是 k，每个结点只有一个孩子。
// leftChild 为 true 时全部是左孩子，否则全部是右孩子。
func bm024testBuildChain(n int, leftChild bool) *TreeNode {
	if n <= 0 {
		return nil
	}
	root := &TreeNode{Val: 1}
	cur := root
	for v := 2; v <= n; v++ {
		node := &TreeNode{Val: v}
		if leftChild {
			cur.Left = node
		} else {
			cur.Right = node
		}
		cur = node
	}
	return root
}

// bm024testBuildFullTree 构造 {4,2,6,1,3,5,7} 这棵满二叉树（同时是一棵二叉搜索树，
// 因此中序遍历结果必然是 1..7）。
func bm024testBuildFullTree() *TreeNode {
	return &TreeNode{
		Val: 4,
		Left: &TreeNode{
			Val:   2,
			Left:  &TreeNode{Val: 1},
			Right: &TreeNode{Val: 3},
		},
		Right: &TreeNode{
			Val:   6,
			Left:  &TreeNode{Val: 5},
			Right: &TreeNode{Val: 7},
		},
	}
}

// bm024testBuildSameValueTree 构造所有结点值都相同的小树，用于确认遍历不会因值相同而出错。
func bm024testBuildSameValueTree(val, n int) *TreeNode {
	if n <= 0 {
		return nil
	}
	root := &TreeNode{Val: val}
	cur := root
	for i := 1; i < n; i++ {
		node := &TreeNode{Val: val}
		if i%2 == 1 {
			cur.Left = node
		} else {
			cur.Right = node
		}
		cur = node
	}
	return root
}

// bm024testBuildRandomTree 用固定种子随机挂载父子关系，构造一棵有 n 个结点的二叉树。
func bm024testBuildRandomTree(r *rand.Rand, n int) *TreeNode {
	if n <= 0 {
		return nil
	}
	nodes := make([]*TreeNode, n)
	for i := range nodes {
		nodes[i] = &TreeNode{Val: i + 1}
	}
	for i := 1; i < n; i++ {
		for {
			parent := nodes[r.Intn(i)]
			if r.Intn(2) == 0 && parent.Left == nil {
				parent.Left = nodes[i]
				break
			}
			if parent.Right == nil {
				parent.Right = nodes[i]
				break
			}
		}
	}
	return nodes[0]
}

// bm024testAssertInts 依次比较两个切片，失败时只打印前若干个元素，避免大用例刷屏。
func bm024testAssertInts(t *testing.T, want, got []int) {
	t.Helper()
	if len(want) != len(got) {
		t.Fatalf("长度不符：期望 %d 个元素，实际 %d 个", len(want), len(got))
	}
	for i := range want {
		if want[i] != got[i] {
			t.Fatalf("第 %d 个元素不符：期望 %d，实际 %d（期望前 10 个 %v，实际前 10 个 %v）",
				i, want[i], got[i], bm024testHead10(want), bm024testHead10(got))
		}
	}
}

// bm024testHead10 返回切片前 10 个元素，用于断言失败时的可读输出。
func bm024testHead10(s []int) []int {
	if len(s) > 10 {
		return s[:10]
	}
	return s
}

// bm024testSeq 返回 1,2,...,n。
func bm024testSeq(n int) []int {
	out := make([]int, n)
	for i := range out {
		out[i] = i + 1
	}
	return out
}

// revSeq 返回 n,n-1,...,1。
func revSeq(n int) []int {
	out := make([]int, n)
	for i := range out {
		out[i] = n - i
	}
	return out
}

func TestBM024(t *testing.T) {
	cases := []struct {
		name string
		root *TreeNode
		want []int
	}{
		{
			name: "官方示例_1_hash_2_3",
			root: &TreeNode{Val: 1, Right: &TreeNode{Val: 2, Left: &TreeNode{Val: 3}}},
			want: []int{1, 3, 2},
		},
		{name: "官方示例_空树", root: nil, want: []int{}},
		{name: "官方示例_单结点", root: &TreeNode{Val: 1}, want: []int{1}},
		{
			name: "只有左孩子的链状树",
			root: bm024testBuildChain(5, true),
			want: []int{5, 4, 3, 2, 1},
		},
		{
			name: "只有右孩子的链状树",
			root: bm024testBuildChain(5, false),
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "满二叉搜索树",
			root: bm024testBuildFullTree(),
			want: []int{1, 2, 3, 4, 5, 6, 7},
		},
		{
			name: "全相同结点值",
			root: bm024testBuildSameValueTree(5, 5),
			want: []int{5, 5, 5, 5, 5},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			bm024testAssertInts(t, tc.want, SolveBM024(tc.root))
			// 递归版只作为对照：两种实现的结果必须完全一致。
			bm024testAssertInts(t, tc.want, inorderRecursive(tc.root))
		})
	}

	t.Run("空树返回空切片", func(t *testing.T) {
		if got := SolveBM024(nil); len(got) != 0 {
			t.Fatalf("期望空结果，实际 %v", got)
		}
	})

	t.Run("随机树_迭代与递归一致", func(t *testing.T) {
		r := rand.New(rand.NewSource(20240606))
		for round := 0; round < 100; round++ {
			root := bm024testBuildRandomTree(r, r.Intn(50)+1)
			want := inorderRecursive(root)
			bm024testAssertInts(t, want, SolveBM024(root))
			// 遍历结果必须是所有结点值的一个排列（结点值互不相同）。
			seen := make(map[int]bool, len(want))
			for _, v := range want {
				if seen[v] {
					t.Fatalf("结点值 %d 重复出现", v)
				}
				seen[v] = true
			}
		}
	})

	t.Run("大规模左链_1e5个结点", func(t *testing.T) {
		const n = 100000
		root := bm024testBuildChain(n, true)
		bm024testAssertInts(t, revSeq(n), SolveBM024(root))
		bm024testAssertInts(t, revSeq(n), inorderRecursive(root))
	})

	t.Run("大规模右链_1e5个结点", func(t *testing.T) {
		const n = 100000
		root := bm024testBuildChain(n, false)
		bm024testAssertInts(t, bm024testSeq(n), SolveBM024(root))
	})
}

// ---------------------------------------------------------------------------
// BM025 二叉树的后序遍历 · 测试
// ---------------------------------------------------------------------------
// buildTree 按层序（含 nil 占位）构造二叉树，仅用于测试。
func buildTree(vals []any) *TreeNode {
	if len(vals) == 0 || vals[0] == nil {
		return nil
	}
	root := &TreeNode{Val: vals[0].(int)}
	queue := []*TreeNode{root}
	i := 1
	for len(queue) > 0 && i < len(vals) {
		cur := queue[0]
		queue = queue[1:]
		if i < len(vals) {
			if vals[i] != nil {
				cur.Left = &TreeNode{Val: vals[i].(int)}
				queue = append(queue, cur.Left)
			}
			i++
		}
		if i < len(vals) {
			if vals[i] != nil {
				cur.Right = &TreeNode{Val: vals[i].(int)}
				queue = append(queue, cur.Right)
			}
			i++
		}
	}
	return root
}

// bm025testBuildChain 构造一条单链：order 为从根向下的结点值，leftFirst 为 true 时全部挂在左孩子上。
func bm025testBuildChain(order []int, leftFirst bool) *TreeNode {
	if len(order) == 0 {
		return nil
	}
	root := &TreeNode{Val: order[0]}
	cur := root
	for _, v := range order[1:] {
		next := &TreeNode{Val: v}
		if leftFirst {
			cur.Left = next
		} else {
			cur.Right = next
		}
		cur = next
	}
	return root
}

// postorderRecursive 递归版后序遍历，作为随机用例的参考实现（仅测试使用）。
func postorderRecursive(root *TreeNode) []int {
	res := make([]int, 0)
	var walk func(node *TreeNode)
	walk = func(node *TreeNode) {
		if node == nil {
			return
		}
		walk(node.Left)
		walk(node.Right)
		res = append(res, node.Val)
	}
	walk(root)
	return res
}

// randomTree 以固定种子构造 n 个结点的随机形状二叉树（仅测试使用）。
func randomTree(rnd *rand.Rand, n int) *TreeNode {
	if n == 0 {
		return nil
	}
	root := &TreeNode{Val: rnd.Intn(100)}
	nodes := []*TreeNode{root}
	for i := 1; i < n; i++ {
		node := &TreeNode{Val: rnd.Intn(100)}
		for {
			parent := nodes[rnd.Intn(len(nodes))]
			if rnd.Intn(2) == 0 {
				if parent.Left == nil {
					parent.Left = node
					break
				}
			} else if parent.Right == nil {
				parent.Right = node
				break
			}
		}
		nodes = append(nodes, node)
	}
	return root
}

func TestBM025(t *testing.T) {
	tests := []struct {
		name string
		tree []any
		want []int
	}{
		{name: "官方示例1_右链带左孩子", tree: []any{1, nil, 2, 3}, want: []int{3, 2, 1}},
		{name: "空树", tree: []any{}, want: []int{}},
		{name: "单结点", tree: []any{7}, want: []int{7}},
		{name: "只有左子树的链", tree: []any{1, 2, nil, 3}, want: []int{3, 2, 1}},
		{name: "只有右子树的链", tree: []any{1, nil, 2, nil, 3}, want: []int{3, 2, 1}},
		{name: "满二叉树", tree: []any{1, 2, 3, 4, 5, 6, 7}, want: []int{4, 5, 2, 6, 7, 3, 1}},
		{name: "常规用例", tree: []any{3, 9, 20, nil, nil, 15, 7}, want: []int{9, 15, 7, 20, 3}},
		{name: "全部值相同", tree: []any{2, 2, 2}, want: []int{2, 2, 2}},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := SolveBM025(buildTree(tc.tree))
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("期望 %v，实际 %v", tc.want, got)
			}
		})
	}

	// 大链状树：迭代写法不受递归深度限制，同时也校验顺序。
	t.Run("左链一万结点", func(t *testing.T) {
		const n = 10000
		order := make([]int, n)
		for i := 0; i < n; i++ {
			order[i] = i
		}
		got := SolveBM025(bm025testBuildChain(order, true))
		if len(got) != n {
			t.Fatalf("期望长度 %d，实际 %d", n, len(got))
		}
		for i := 0; i < n; i++ {
			if want := n - 1 - i; got[i] != want {
				t.Fatalf("下标 %d 期望 %d，实际 %d", i, want, got[i])
			}
		}
	})

	// 随机形状二叉树：与递归参考实现比对。
	t.Run("随机形状与递归实现一致", func(t *testing.T) {
		rnd := rand.New(rand.NewSource(20240525))
		for round := 0; round < 50; round++ {
			tree := randomTree(rnd, rnd.Intn(60))
			want := postorderRecursive(tree)
			got := SolveBM025(tree)
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("第 %d 轮随机用例，期望 %v，实际 %v", round, want, got)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// BM026 求二叉树的层序遍历 · 测试
// ---------------------------------------------------------------------------
// bm026testBuildTree 按层序（含 nil 占位）构造二叉树，仅用于测试。
func bm026testBuildTree(vals []any) *TreeNode {
	if len(vals) == 0 || vals[0] == nil {
		return nil
	}
	root := &TreeNode{Val: vals[0].(int)}
	queue := []*TreeNode{root}
	i := 1
	for len(queue) > 0 && i < len(vals) {
		cur := queue[0]
		queue = queue[1:]
		if i < len(vals) {
			if vals[i] != nil {
				cur.Left = &TreeNode{Val: vals[i].(int)}
				queue = append(queue, cur.Left)
			}
			i++
		}
		if i < len(vals) {
			if vals[i] != nil {
				cur.Right = &TreeNode{Val: vals[i].(int)}
				queue = append(queue, cur.Right)
			}
			i++
		}
	}
	return root
}

// bm026testBuildChain 构造一条单链：order 为从根向下的结点值，leftFirst 为 true 时全部挂在左孩子上。
func bm026testBuildChain(order []int, leftFirst bool) *TreeNode {
	if len(order) == 0 {
		return nil
	}
	root := &TreeNode{Val: order[0]}
	cur := root
	for _, v := range order[1:] {
		next := &TreeNode{Val: v}
		if leftFirst {
			cur.Left = next
		} else {
			cur.Right = next
		}
		cur = next
	}
	return root
}

func TestBM026(t *testing.T) {
	tests := []struct {
		name string
		tree []any
		want [][]int
	}{
		{
			name: "官方示例1",
			tree: []any{1, 2, 3, nil, nil, 4, 5},
			want: [][]int{{1}, {2, 3}, {4, 5}},
		},
		{
			name: "空树",
			tree: []any{},
			want: [][]int{},
		},
		{
			name: "单结点",
			tree: []any{9},
			want: [][]int{{9}},
		},
		{
			name: "只有左子树的链",
			tree: []any{1, 2, nil, 3},
			want: [][]int{{1}, {2}, {3}},
		},
		{
			name: "只有右子树的链",
			tree: []any{1, nil, 2, nil, 3},
			want: [][]int{{1}, {2}, {3}},
		},
		{
			name: "满二叉树",
			tree: []any{1, 2, 3, 4, 5, 6, 7},
			want: [][]int{{1}, {2, 3}, {4, 5, 6, 7}},
		},
		{
			name: "常规用例",
			tree: []any{3, 9, 20, nil, nil, 15, 7},
			want: [][]int{{3}, {9, 20}, {15, 7}},
		},
		{
			name: "末层结点稀疏",
			tree: []any{1, 2, 3, 4, nil, nil, 7},
			want: [][]int{{1}, {2, 3}, {4, 7}},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := SolveBM026(bm026testBuildTree(tc.tree))
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("期望 %v，实际 %v", tc.want, got)
			}
		})
	}

	// 大链状树：每层只有一个结点，层数等于结点数。
	t.Run("左链一千层", func(t *testing.T) {
		const n = 1000
		order := make([]int, n)
		for i := 0; i < n; i++ {
			order[i] = i + 1
		}
		got := SolveBM026(bm026testBuildChain(order, true))
		if len(got) != n {
			t.Fatalf("期望层数 %d，实际 %d", n, len(got))
		}
		for i := 0; i < n; i++ {
			want := []int{order[i]}
			if !reflect.DeepEqual(got[i], want) {
				t.Fatalf("第 %d 层期望 %v，实际 %v", i, want, got[i])
			}
		}
	})
}

// ---------------------------------------------------------------------------
// BM027 按之字形顺序打印二叉树 · 测试
// ---------------------------------------------------------------------------
// bm027testBuildTree 按层序（含 nil 占位）构造二叉树，仅用于测试。
func bm027testBuildTree(vals []any) *TreeNode {
	if len(vals) == 0 || vals[0] == nil {
		return nil
	}
	root := &TreeNode{Val: vals[0].(int)}
	queue := []*TreeNode{root}
	i := 1
	for len(queue) > 0 && i < len(vals) {
		cur := queue[0]
		queue = queue[1:]
		if i < len(vals) {
			if vals[i] != nil {
				cur.Left = &TreeNode{Val: vals[i].(int)}
				queue = append(queue, cur.Left)
			}
			i++
		}
		if i < len(vals) {
			if vals[i] != nil {
				cur.Right = &TreeNode{Val: vals[i].(int)}
				queue = append(queue, cur.Right)
			}
			i++
		}
	}
	return root
}

// bm027testBuildChain 构造一条单链：order 为从根向下的结点值，leftFirst 为 true 时全部挂在左孩子上。
func bm027testBuildChain(order []int, leftFirst bool) *TreeNode {
	if len(order) == 0 {
		return nil
	}
	root := &TreeNode{Val: order[0]}
	cur := root
	for _, v := range order[1:] {
		next := &TreeNode{Val: v}
		if leftFirst {
			cur.Left = next
		} else {
			cur.Right = next
		}
		cur = next
	}
	return root
}

func TestBM027(t *testing.T) {
	tests := []struct {
		name string
		tree []any
		want [][]int
	}{
		{
			name: "官方示例1",
			tree: []any{1, 2, 3, nil, nil, 4, 5},
			want: [][]int{{1}, {3, 2}, {4, 5}},
		},
		{
			name: "空树",
			tree: []any{},
			want: [][]int{},
		},
		{
			name: "单结点",
			tree: []any{5},
			want: [][]int{{5}},
		},
		{
			name: "只有左子树的链",
			tree: []any{1, 2, nil, 3},
			want: [][]int{{1}, {2}, {3}},
		},
		{
			name: "只有右子树的链",
			tree: []any{1, nil, 2, nil, 3},
			want: [][]int{{1}, {2}, {3}},
		},
		{
			name: "三层满二叉树",
			tree: []any{1, 2, 3, 4, 5, 6, 7},
			want: [][]int{{1}, {3, 2}, {4, 5, 6, 7}},
		},
		{
			name: "四层满二叉树",
			tree: []any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			want: [][]int{
				{1},
				{3, 2},
				{4, 5, 6, 7},
				{15, 14, 13, 12, 11, 10, 9, 8},
			},
		},
		{
			name: "常规二叉搜索树",
			tree: []any{8, 6, 10, 5, 7, 9, 11},
			want: [][]int{{8}, {10, 6}, {5, 7, 9, 11}},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := SolveBM027(bm027testBuildTree(tc.tree))
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("期望 %v，实际 %v", tc.want, got)
			}
		})
	}

	// 性质校验：每一层排序后应与普通层序遍历的对应层一致，只是方向交替。
	t.Run("每层元素与层序遍历一致", func(t *testing.T) {
		tree := bm027testBuildTree([]any{1, 2, 3, 4, nil, 6, 7, 8, 9, nil, nil, 12})
		got := SolveBM027(tree)
		levelOrder := levelOrderRef(tree)
		if len(got) != len(levelOrder) {
			t.Fatalf("期望层数 %d，实际 %d", len(levelOrder), len(got))
		}
		for i := range got {
			if len(got[i]) != len(levelOrder[i]) {
				t.Fatalf("第 %d 层期望 %d 个元素，实际 %d 个", i+1, len(levelOrder[i]), len(got[i]))
			}
			a := append([]int(nil), got[i]...)
			b := append([]int(nil), levelOrder[i]...)
			sortInts(a)
			sortInts(b)
			if !reflect.DeepEqual(a, b) {
				t.Fatalf("第 %d 层元素集合期望 %v，实际 %v", i+1, b, a)
			}
		}
	})

	// 大链状树：每层一个元素，方向交替不影响结果。
	t.Run("左链一千层", func(t *testing.T) {
		const n = 1000
		order := make([]int, n)
		for i := 0; i < n; i++ {
			order[i] = i + 1
		}
		got := SolveBM027(bm027testBuildChain(order, true))
		if len(got) != n {
			t.Fatalf("期望层数 %d，实际 %d", n, len(got))
		}
		for i := 0; i < n; i++ {
			if want := []int{order[i]}; !reflect.DeepEqual(got[i], want) {
				t.Fatalf("第 %d 层期望 %v，实际 %v", i, want, got[i])
			}
		}
	})
}

// levelOrderRef 普通层序遍历，仅用于测试中做性质比对。
func levelOrderRef(root *TreeNode) [][]int {
	res := make([][]int, 0)
	if root == nil {
		return res
	}
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		size := len(queue)
		level := make([]int, 0, size)
		for i := 0; i < size; i++ {
			node := queue[0]
			queue = queue[1:]
			level = append(level, node.Val)
			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
		res = append(res, level)
	}
	return res
}

// sortInts 对整数切片做插入排序（层内元素很少，避免为测试引入额外依赖）。
func sortInts(nums []int) {
	for i := 1; i < len(nums); i++ {
		for j := i; j > 0 && nums[j] < nums[j-1]; j-- {
			nums[j], nums[j-1] = nums[j-1], nums[j]
		}
	}
}

// ---------------------------------------------------------------------------
// BM028 二叉树的最大深度 · 测试
// ---------------------------------------------------------------------------
// bm028testBuildTree 按层序（含 nil 占位）构造二叉树，仅用于测试。
func bm028testBuildTree(vals []any) *TreeNode {
	if len(vals) == 0 || vals[0] == nil {
		return nil
	}
	root := &TreeNode{Val: vals[0].(int)}
	queue := []*TreeNode{root}
	i := 1
	for len(queue) > 0 && i < len(vals) {
		cur := queue[0]
		queue = queue[1:]
		if i < len(vals) {
			if vals[i] != nil {
				cur.Left = &TreeNode{Val: vals[i].(int)}
				queue = append(queue, cur.Left)
			}
			i++
		}
		if i < len(vals) {
			if vals[i] != nil {
				cur.Right = &TreeNode{Val: vals[i].(int)}
				queue = append(queue, cur.Right)
			}
			i++
		}
	}
	return root
}

// bm028testBuildChain 构造一条单链：order 为从根向下的结点值，leftFirst 为 true 时全部挂在左孩子上。
func bm028testBuildChain(order []int, leftFirst bool) *TreeNode {
	if len(order) == 0 {
		return nil
	}
	root := &TreeNode{Val: order[0]}
	cur := root
	for _, v := range order[1:] {
		next := &TreeNode{Val: v}
		if leftFirst {
			cur.Left = next
		} else {
			cur.Right = next
		}
		cur = next
	}
	return root
}

// levelDepth 用 BFS 逐层计数求深度，作为随机用例的参考实现（仅测试使用）。
func levelDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}
	depth := 0
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		size := len(queue)
		for i := 0; i < size; i++ {
			node := queue[0]
			queue = queue[1:]
			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
		depth++
	}
	return depth
}

// bm028testRandomTree 以固定种子构造 n 个结点的随机形状二叉树（仅测试使用）。
func bm028testRandomTree(rnd *rand.Rand, n int) *TreeNode {
	if n == 0 {
		return nil
	}
	root := &TreeNode{Val: rnd.Intn(100)}
	nodes := []*TreeNode{root}
	for i := 1; i < n; i++ {
		node := &TreeNode{Val: rnd.Intn(100)}
		for {
			parent := nodes[rnd.Intn(len(nodes))]
			if rnd.Intn(2) == 0 {
				if parent.Left == nil {
					parent.Left = node
					break
				}
			} else if parent.Right == nil {
				parent.Right = node
				break
			}
		}
		nodes = append(nodes, node)
	}
	return root
}

func TestBM028(t *testing.T) {
	tests := []struct {
		name string
		tree []any
		want int
	}{
		{name: "官方示例1", tree: []any{1, 2, 3, nil, nil, 4, 5}, want: 3},
		{name: "空树", tree: []any{}, want: 0},
		{name: "单结点", tree: []any{1}, want: 1},
		{name: "只有左子树的链", tree: []any{1, 2, nil, 3, nil, 4}, want: 4},
		{name: "只有右子树的链", tree: []any{1, nil, 2, nil, 3, nil, 4}, want: 4},
		{name: "满二叉树三层", tree: []any{1, 2, 3, 4, 5, 6, 7}, want: 3},
		{
			name: "满二叉树四层",
			tree: []any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			want: 4,
		},
		{name: "左子树更深", tree: []any{3, 9, 20, 7, nil, nil, nil, 5}, want: 4},
		{name: "常规用例", tree: []any{3, 9, 20, nil, nil, 15, 7}, want: 3},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := SolveBM028(bm028testBuildTree(tc.tree))
			if got != tc.want {
				t.Fatalf("期望 %d，实际 %d", tc.want, got)
			}
		})
	}

	// 长链：深度等于结点数，验证长链下不会退化。
	t.Run("左链五万结点", func(t *testing.T) {
		const n = 50000
		order := make([]int, n)
		for i := 0; i < n; i++ {
			order[i] = i
		}
		if got := SolveBM028(bm028testBuildChain(order, true)); got != n {
			t.Fatalf("期望 %d，实际 %d", n, got)
		}
		if got := SolveBM028(bm028testBuildChain(order, false)); got != n {
			t.Fatalf("右链期望 %d，实际 %d", n, got)
		}
	})

	// 随机形状二叉树：与 BFS 逐层计数的结果比对。
	t.Run("随机形状与BFS结果一致", func(t *testing.T) {
		rnd := rand.New(rand.NewSource(20240528))
		for round := 0; round < 50; round++ {
			tree := bm028testRandomTree(rnd, rnd.Intn(80))
			want := levelDepth(tree)
			if got := SolveBM028(tree); got != want {
				t.Fatalf("第 %d 轮随机用例，期望 %d，实际 %d", round, want, got)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// BM029 二叉树中和为某一值的路径(一) · 测试
// ---------------------------------------------------------------------------
// bm029testBuildTree 按层序（含 nil 占位）构造二叉树，仅用于测试。
func bm029testBuildTree(vals []any) *TreeNode {
	if len(vals) == 0 || vals[0] == nil {
		return nil
	}
	root := &TreeNode{Val: vals[0].(int)}
	queue := []*TreeNode{root}
	i := 1
	for len(queue) > 0 && i < len(vals) {
		cur := queue[0]
		queue = queue[1:]
		if i < len(vals) {
			if vals[i] != nil {
				cur.Left = &TreeNode{Val: vals[i].(int)}
				queue = append(queue, cur.Left)
			}
			i++
		}
		if i < len(vals) {
			if vals[i] != nil {
				cur.Right = &TreeNode{Val: vals[i].(int)}
				queue = append(queue, cur.Right)
			}
			i++
		}
	}
	return root
}

// bm029testBuildChain 构造一条单链：order 为从根向下的结点值，leftFirst 为 true 时全部挂在左孩子上。
func bm029testBuildChain(order []int, leftFirst bool) *TreeNode {
	if len(order) == 0 {
		return nil
	}
	root := &TreeNode{Val: order[0]}
	cur := root
	for _, v := range order[1:] {
		next := &TreeNode{Val: v}
		if leftFirst {
			cur.Left = next
		} else {
			cur.Right = next
		}
		cur = next
	}
	return root
}

func TestBM029(t *testing.T) {
	tests := []struct {
		name string
		tree []any
		sum  int
		want bool
	}{
		{
			name: "官方示例1_存在和为22的路径",
			tree: []any{5, 4, 8, 1, 11, nil, 9, nil, nil, 2, 7},
			sum:  22,
			want: true,
		},
		{
			name: "官方示例1_不存在的和",
			tree: []any{5, 4, 8, 1, 11, nil, 9, nil, nil, 2, 7},
			sum:  23,
			want: false,
		},
		{name: "空树", tree: []any{}, sum: 0, want: false},
		{name: "空树非零sum", tree: []any{}, sum: 5, want: false},
		{name: "单结点命中", tree: []any{5}, sum: 5, want: true},
		{name: "单结点未命中", tree: []any{5}, sum: 0, want: false},
		{name: "单结点负数叶子", tree: []any{-3}, sum: -3, want: true},
		{name: "左链命中", tree: []any{1, 2, nil, 3}, sum: 6, want: true},
		{name: "左链中途不截断", tree: []any{1, 2, nil, 3}, sum: 3, want: false},
		{name: "右链命中", tree: []any{1, nil, 2, nil, 3}, sum: 6, want: true},
		{name: "右链中途不截断", tree: []any{1, nil, 2, nil, 3}, sum: 3, want: false},
		{
			name: "满二叉树命中左叶路径",
			tree: []any{1, 2, 3, 4, 5, 6, 7},
			sum:  7,
			want: true,
		},
		{
			name: "满二叉树命中右叶路径",
			tree: []any{1, 2, 3, 4, 5, 6, 7},
			sum:  11,
			want: true,
		},
		{
			name: "满二叉树未命中",
			tree: []any{1, 2, 3, 4, 5, 6, 7},
			sum:  100,
			want: false,
		},
		{name: "负数和_命中左叶", tree: []any{-2, -3, -1}, sum: -5, want: true},
		{name: "负数和_命中右叶", tree: []any{-2, -3, -1}, sum: -3, want: true},
		{name: "负数和_未命中", tree: []any{-2, -3, -1}, sum: -6, want: false},
		{name: "正负抵消命中", tree: []any{10, -5, 3, -2}, sum: 3, want: true},
		{name: "正负抵消未命中", tree: []any{10, -5, 3, -2}, sum: 5, want: false},
		{name: "非叶子结点不能作为终点", tree: []any{1, 2, nil, 4}, sum: 3, want: false},
		{name: "左斜树全零", tree: []any{0, 0, nil, 0}, sum: 0, want: true},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := SolveBM029(bm029testBuildTree(tc.tree), tc.sum)
			if got != tc.want {
				t.Fatalf("期望 %v，实际 %v", tc.want, got)
			}
		})
	}

	// 长链：整条链只有一条根到叶子的路径。
	t.Run("左链一千结点", func(t *testing.T) {
		const n = 1000
		order := make([]int, n)
		total := 0
		for i := 0; i < n; i++ {
			order[i] = i + 1
			total += i + 1
		}
		if got := SolveBM029(bm029testBuildChain(order, true), total); !got {
			t.Fatalf("期望 true，实际 %v", got)
		}
		if got := SolveBM029(bm029testBuildChain(order, false), total-1); got {
			t.Fatalf("期望 false，实际 %v", got)
		}
	})
}

// ---------------------------------------------------------------------------
// BM030 二叉搜索树与双向链表 · 测试
// ---------------------------------------------------------------------------
// bm030testBuildTree 按层序（含 nil 占位）构造二叉搜索树，仅用于测试。
func bm030testBuildTree(vals []any) *TreeNode {
	if len(vals) == 0 || vals[0] == nil {
		return nil
	}
	root := &TreeNode{Val: vals[0].(int)}
	queue := []*TreeNode{root}
	i := 1
	for len(queue) > 0 && i < len(vals) {
		cur := queue[0]
		queue = queue[1:]
		if i < len(vals) {
			if vals[i] != nil {
				cur.Left = &TreeNode{Val: vals[i].(int)}
				queue = append(queue, cur.Left)
			}
			i++
		}
		if i < len(vals) {
			if vals[i] != nil {
				cur.Right = &TreeNode{Val: vals[i].(int)}
				queue = append(queue, cur.Right)
			}
			i++
		}
	}
	return root
}

// bm030testBuildChain 构造一条单链：order 为从根向下的结点值，leftFirst 为 true 时全部挂在左孩子上。
func bm030testBuildChain(order []int, leftFirst bool) *TreeNode {
	if len(order) == 0 {
		return nil
	}
	root := &TreeNode{Val: order[0]}
	cur := root
	for _, v := range order[1:] {
		next := &TreeNode{Val: v}
		if leftFirst {
			cur.Left = next
		} else {
			cur.Right = next
		}
		cur = next
	}
	return root
}

// collectNodes 先序收集树中全部结点指针，用于校验转换后没有新建或丢失结点。
func collectNodes(root *TreeNode) []*TreeNode {
	nodes := make([]*TreeNode, 0)
	var walk func(node *TreeNode)
	walk = func(node *TreeNode) {
		if node == nil {
			return
		}
		nodes = append(nodes, node)
		walk(node.Left)
		walk(node.Right)
	}
	walk(root)
	return nodes
}

// checkDoublyList 校验 head 出发的双向链表：值递增、指针互为前后驱、首尾为 nil、
// 且链表中每个结点都来自 origin 中的原始树结点（未新建结点）。
func checkDoublyList(t *testing.T, head *TreeNode, want []int, origin []*TreeNode) {
	t.Helper()

	originSet := make(map[*TreeNode]bool, len(origin))
	for _, node := range origin {
		originSet[node] = true
	}

	if len(want) == 0 {
		if head != nil {
			t.Fatalf("期望空链表，实际头结点值为 %d", head.Val)
		}
		return
	}
	if head == nil {
		t.Fatalf("期望链表 %v，实际为 nil", want)
	}
	if head.Left != nil {
		t.Fatalf("头结点 Left 期望 nil，实际指向值 %d 的结点", head.Left.Val)
	}

	nodes := make([]*TreeNode, 0, len(want)+1)
	vals := make([]int, 0, len(want))
	for cur := head; cur != nil; cur = cur.Right {
		if len(nodes) > len(origin) {
			t.Fatalf("链表出现环或结点数超过原树结点数 %d", len(origin))
		}
		if !originSet[cur] {
			t.Fatalf("链表中出现非原树的结点（值 %d），说明新建了结点", cur.Val)
		}
		nodes = append(nodes, cur)
		vals = append(vals, cur.Val)
	}

	if len(vals) != len(want) {
		t.Fatalf("期望链表长度 %d（%v），实际 %d（%v）", len(want), want, len(vals), vals)
	}
	for i := range want {
		if vals[i] != want[i] {
			t.Fatalf("期望链表 %v，实际 %v", want, vals)
		}
	}
	// 中序结果非递减；BST 中若有重复值也允许（牛客用例一般无重复，这里按更宽松的性质校验）。
	for i := 1; i < len(nodes); i++ {
		if vals[i] < vals[i-1] {
			t.Fatalf("链表未保持有序：下标 %d 值为 %d，前一结点值为 %d", i, vals[i], vals[i-1])
		}
		if nodes[i].Left != nodes[i-1] {
			t.Fatalf("下标 %d 的 Left 未指向直接前驱", i)
		}
		if nodes[i-1].Right != nodes[i] {
			t.Fatalf("下标 %d 的 Right 未指向直接后继", i-1)
		}
	}
	tail := nodes[len(nodes)-1]
	if tail.Right != nil {
		t.Fatalf("尾结点 Right 期望 nil，实际指向值 %d 的结点", tail.Right.Val)
	}
}

func TestBM030(t *testing.T) {
	tests := []struct {
		name string
		tree []any
		want []int
	}{
		{
			name: "官方示例1",
			tree: []any{10, 6, 14, 4, 8, 12, 16},
			want: []int{4, 6, 8, 10, 12, 14, 16},
		},
		{name: "空树", tree: []any{}, want: []int{}},
		{name: "单结点", tree: []any{5}, want: []int{5}},
		{
			name: "只有左链的BST",
			tree: []any{5, 4, nil, 3, nil, 2},
			want: []int{2, 3, 4, 5},
		},
		{
			name: "只有右链的BST",
			tree: []any{1, nil, 2, nil, 3, nil, 4},
			want: []int{1, 2, 3, 4},
		},
		{
			name: "满二叉树BST",
			tree: []any{8, 4, 12, 2, 6, 10, 14},
			want: []int{2, 4, 6, 8, 10, 12, 14},
		},
		{
			name: "不规则BST",
			tree: []any{20, 10, 30, 5, 15, 25, 35, 3},
			want: []int{3, 5, 10, 15, 20, 25, 30, 35},
		},
		{
			name: "结点值可重复的BST形状",
			tree: []any{2, 1, 3, 1, nil, nil, 3},
			want: []int{1, 1, 2, 3, 3},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tree := bm030testBuildTree(tc.tree)
			origin := collectNodes(tree)
			head := SolveBM030(tree)
			checkDoublyList(t, head, tc.want, origin)
		})
	}

	// 单结点树：首尾都必须是 nil，且左右指针为空。
	t.Run("单结点首尾均为nil", func(t *testing.T) {
		tree := bm030testBuildTree([]any{42})
		head := SolveBM030(tree)
		if head == nil || head != tree {
			t.Fatalf("期望返回原结点作为头结点")
		}
		if head.Left != nil || head.Right != nil {
			t.Fatalf("单结点链表的 Left/Right 都应为 nil")
		}
	})

	// 长链：BST 退化为左链与右链，转换后应为连续整数序列。
	t.Run("左链一千结点", func(t *testing.T) {
		const n = 1000
		order := make([]int, n)
		for i := 0; i < n; i++ {
			order[i] = n - i // 从大到小挂在左孩子上，中序即 1..n
		}
		tree := bm030testBuildChain(order, true)
		origin := collectNodes(tree)
		head := SolveBM030(tree)
		want := make([]int, 0, n)
		for i := 1; i <= n; i++ {
			want = append(want, i)
		}
		checkDoublyList(t, head, want, origin)
	})

	t.Run("右链一千结点", func(t *testing.T) {
		const n = 1000
		order := make([]int, n)
		for i := 0; i < n; i++ {
			order[i] = i + 1
		}
		tree := bm030testBuildChain(order, false)
		origin := collectNodes(tree)
		head := SolveBM030(tree)
		want := make([]int, 0, n)
		for i := 1; i <= n; i++ {
			want = append(want, i)
		}
		checkDoublyList(t, head, want, origin)
	})
}

// ---------------------------------------------------------------------------
// BM031 对称的二叉树 · 测试
// ---------------------------------------------------------------------------
// bm031testBuildTree 按牛客的层序数组构造二叉树：nil 表示该位置为空结点，空结点不再展开子结点。
func bm031testBuildTree(vals []any) *TreeNode {
	if len(vals) == 0 || vals[0] == nil {
		return nil
	}
	root := &TreeNode{Val: vals[0].(int)}
	queue := []*TreeNode{root}
	i := 1
	for len(queue) > 0 && i < len(vals) {
		cur := queue[0]
		queue = queue[1:]
		if i < len(vals) {
			if v, ok := vals[i].(int); ok {
				cur.Left = &TreeNode{Val: v}
				queue = append(queue, cur.Left)
			}
			i++
		}
		if i < len(vals) {
			if v, ok := vals[i].(int); ok {
				cur.Right = &TreeNode{Val: v}
				queue = append(queue, cur.Right)
			}
			i++
		}
	}
	return root
}

// buildPerfectTree 构造高度为 depth、所有结点值均为 val 的满二叉树。
func buildPerfectTree(depth, val int) *TreeNode {
	if depth <= 0 {
		return nil
	}
	return &TreeNode{
		Val:   val,
		Left:  buildPerfectTree(depth-1, val),
		Right: buildPerfectTree(depth-1, val),
	}
}

// leftmostLeaf 返回最左侧的叶子结点，用于构造「结构对称但某个值不同」的用例。
func leftmostLeaf(root *TreeNode) *TreeNode {
	for root != nil && (root.Left != nil || root.Right != nil) {
		if root.Left != nil {
			root = root.Left
		} else {
			root = root.Right
		}
	}
	return root
}

// treeString 把二叉树序列化成牛客风格的层序字符串（空结点记为 #，末尾的 # 省略），便于失败时报错。
func treeString(root *TreeNode) string {
	if root == nil {
		return "{}"
	}
	var parts []string
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur == nil {
			parts = append(parts, "#")
			continue
		}
		parts = append(parts, strconv.Itoa(cur.Val))
		queue = append(queue, cur.Left, cur.Right)
	}
	for len(parts) > 0 && parts[len(parts)-1] == "#" {
		parts = parts[:len(parts)-1]
	}
	return "{" + strings.Join(parts, ",") + "}"
}

func TestBM031(t *testing.T) {
	// 高度 8 的满二叉树、值全相同，必然对称；再把最左侧叶子改掉，破坏对称性。
	perfect := buildPerfectTree(8, 7)
	broken := buildPerfectTree(8, 7)
	leftmostLeaf(broken).Val = 8

	tests := []struct {
		name string
		root *TreeNode
		want bool
	}{
		{"空树视为对称", nil, true},
		{"单结点", bm031testBuildTree([]any{1}), true},
		{"官方示例: 结构镜像且值相等", bm031testBuildTree([]any{1, 2, 2, 3, 4, 4, 3}), true},
		{"镜像结构但对称位置值不等", bm031testBuildTree([]any{1, 2, 2, 3, 4, 4, 5}), false},
		{"一侧缺孩子", bm031testBuildTree([]any{1, 2, 2, nil, 3, nil, 3}), false},
		{"值全相同但结构不对称", bm031testBuildTree([]any{1, 1, 1, 1, nil, 1}), false},
		{"三层对称含负数", bm031testBuildTree([]any{-1, -2, -2, -3, 4, 4, -3}), true},
		{"叶子值不同", bm031testBuildTree([]any{5, 3, 3, 1, 4, 4, 2}), false},
		{"单侧链", bm031testBuildTree([]any{1, 2, nil, 3}), false},
		{"两结点左右各一", bm031testBuildTree([]any{1, 2, 2}), true},
		{"满二叉树深度8", perfect, true},
		{"满二叉树改动最左叶子", broken, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := SolveBM031(tc.root)
			if got != tc.want {
				t.Fatalf("期望 %v，实际 %v（输入树 %s）", tc.want, got, treeString(tc.root))
			}
		})
	}
}

// ---------------------------------------------------------------------------
// BM032 合并二叉树 · 测试
// ---------------------------------------------------------------------------
// bm032testBuildTree 按牛客的层序数组构造二叉树：nil 表示该位置为空结点，空结点不再展开子结点。
func bm032testBuildTree(vals []any) *TreeNode {
	if len(vals) == 0 || vals[0] == nil {
		return nil
	}
	root := &TreeNode{Val: vals[0].(int)}
	queue := []*TreeNode{root}
	i := 1
	for len(queue) > 0 && i < len(vals) {
		cur := queue[0]
		queue = queue[1:]
		if i < len(vals) {
			if v, ok := vals[i].(int); ok {
				cur.Left = &TreeNode{Val: v}
				queue = append(queue, cur.Left)
			}
			i++
		}
		if i < len(vals) {
			if v, ok := vals[i].(int); ok {
				cur.Right = &TreeNode{Val: v}
				queue = append(queue, cur.Right)
			}
			i++
		}
	}
	return root
}

// bm032testTreeString 把二叉树序列化成牛客风格的层序字符串（空结点记为 #，末尾的 # 省略）。
func bm032testTreeString(root *TreeNode) string {
	if root == nil {
		return "{}"
	}
	var parts []string
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur == nil {
			parts = append(parts, "#")
			continue
		}
		parts = append(parts, strconv.Itoa(cur.Val))
		queue = append(queue, cur.Left, cur.Right)
	}
	for len(parts) > 0 && parts[len(parts)-1] == "#" {
		parts = parts[:len(parts)-1]
	}
	return "{" + strings.Join(parts, ",") + "}"
}

func TestBM032(t *testing.T) {
	tests := []struct {
		name string
		a    []any
		b    []any
		want string
	}{
		{
			name: "官方示例",
			a:    []any{1, 3, 2, 5},
			b:    []any{2, 1, 3, nil, 4, nil, 7},
			want: "{3,4,5,5,4,#,7}",
		},
		{
			name: "两棵都是空树",
			a:    nil,
			b:    nil,
			want: "{}",
		},
		{
			name: "只有第二棵非空",
			a:    nil,
			b:    []any{1, 2, 3},
			want: "{1,2,3}",
		},
		{
			name: "只有第一棵非空",
			a:    []any{1, 2, 3},
			b:    nil,
			want: "{1,2,3}",
		},
		{
			name: "单结点相加",
			a:    []any{1},
			b:    []any{2},
			want: "{3}",
		},
		{
			name: "相加得零（含负数）",
			a:    []any{-5},
			b:    []any{5},
			want: "{0}",
		},
		{
			name: "一边缺子树时整段复制",
			a:    []any{1, 2, 3},
			b:    []any{1, nil, 4, 5},
			want: "{2,2,7,#,#,5}",
		},
		{
			// a 为 1 -> 左 2 -> 左 3，b 只有根结点 1：
			// 合并后根为 2，左子树整体复制 a 的左子树（2 带着左孩子 3）。
			name: "两边深度不同",
			a:    []any{1, 2, nil, 3},
			b:    []any{1},
			want: "{2,2,#,3}",
		},
		{
			name: "负值常规用例",
			a:    []any{-1, -2, -3},
			b:    []any{-1, -2, -3},
			want: "{-2,-4,-6}",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			root1, root2 := bm032testBuildTree(tc.a), bm032testBuildTree(tc.b)
			got := bm032testTreeString(SolveBM032(root1, root2))
			if got != tc.want {
				t.Fatalf("期望 %s，实际 %s", tc.want, got)
			}
			// 合并过程不能破坏入参：两棵原始树的序列化结果必须保持不变。
			if s := bm032testTreeString(root1); s != bm032testTreeString(bm032testBuildTree(tc.a)) {
				t.Fatalf("第一棵入参树被修改，期望 %s，实际 %s", bm032testTreeString(bm032testBuildTree(tc.a)), s)
			}
			if s := bm032testTreeString(root2); s != bm032testTreeString(bm032testBuildTree(tc.b)) {
				t.Fatalf("第二棵入参树被修改，期望 %s，实际 %s", bm032testTreeString(bm032testBuildTree(tc.b)), s)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// BM033 二叉树的镜像 · 测试
// ---------------------------------------------------------------------------
// bm033testBuildTree 按牛客的层序数组构造二叉树：nil 表示该位置为空结点，空结点不再展开子结点。
func bm033testBuildTree(vals []any) *TreeNode {
	if len(vals) == 0 || vals[0] == nil {
		return nil
	}
	root := &TreeNode{Val: vals[0].(int)}
	queue := []*TreeNode{root}
	i := 1
	for len(queue) > 0 && i < len(vals) {
		cur := queue[0]
		queue = queue[1:]
		if i < len(vals) {
			if v, ok := vals[i].(int); ok {
				cur.Left = &TreeNode{Val: v}
				queue = append(queue, cur.Left)
			}
			i++
		}
		if i < len(vals) {
			if v, ok := vals[i].(int); ok {
				cur.Right = &TreeNode{Val: v}
				queue = append(queue, cur.Right)
			}
			i++
		}
	}
	return root
}

// bm033testTreeString 把二叉树序列化成牛客风格的层序字符串（空结点记为 #，末尾的 # 省略）。
func bm033testTreeString(root *TreeNode) string {
	if root == nil {
		return "{}"
	}
	var parts []string
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur == nil {
			parts = append(parts, "#")
			continue
		}
		parts = append(parts, strconv.Itoa(cur.Val))
		queue = append(queue, cur.Left, cur.Right)
	}
	for len(parts) > 0 && parts[len(parts)-1] == "#" {
		parts = parts[:len(parts)-1]
	}
	return "{" + strings.Join(parts, ",") + "}"
}

func TestBM033(t *testing.T) {
	tests := []struct {
		name string
		vals []any
		want string
	}{
		{"官方示例", []any{8, 6, 10, 5, 7, 9, 11}, "{8,10,6,11,9,7,5}"},
		{"空树", nil, "{}"},
		{"单结点", []any{1}, "{1}"},
		{"只有左孩子", []any{1, 2}, "{1,#,2}"},
		{"只有右孩子", []any{1, nil, 2}, "{1,2}"},
		{"四结点左偏", []any{1, 2, 3, 4}, "{1,3,2,#,#,#,4}"},
		{"值全相同", []any{1, 1, 1}, "{1,1,1}"},
		{"单侧链", []any{1, 2, nil, 3}, "{1,#,2,#,3}"},
		{"含负数", []any{-1, -2, -3, -4, -5}, "{-1,-3,-2,#,#,-5,-4}"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			root := bm033testBuildTree(tc.vals)
			got := bm033testTreeString(SolveBM033(root))
			if got != tc.want {
				t.Fatalf("期望 %s，实际 %s", tc.want, got)
			}
		})
	}

	t.Run("两次镜像还原原树", func(t *testing.T) {
		vals := []any{1, 2, 3, 4, nil, 5, 6, 7}
		want := bm033testTreeString(bm033testBuildTree(vals))
		got := bm033testTreeString(SolveBM033(SolveBM033(bm033testBuildTree(vals))))
		if got != want {
			t.Fatalf("期望 %s，实际 %s", want, got)
		}
	})

	t.Run("镜像后与原树不同", func(t *testing.T) {
		vals := []any{1, 2, 3, 4, nil, 5, 6, 7}
		mirrored := bm033testTreeString(SolveBM033(bm033testBuildTree(vals)))
		original := bm033testTreeString(bm033testBuildTree(vals))
		if mirrored == original {
			t.Fatalf("镜像结果不应与原树相同：%s", mirrored)
		}
	})
}

// ---------------------------------------------------------------------------
// BM034 判断是不是二叉搜索树 · 测试
// ---------------------------------------------------------------------------
// bm034testBuildTree 按牛客的层序数组构造二叉树：nil 表示该位置为空结点，空结点不再展开子结点。
func bm034testBuildTree(vals []any) *TreeNode {
	if len(vals) == 0 || vals[0] == nil {
		return nil
	}
	root := &TreeNode{Val: vals[0].(int)}
	queue := []*TreeNode{root}
	i := 1
	for len(queue) > 0 && i < len(vals) {
		cur := queue[0]
		queue = queue[1:]
		if i < len(vals) {
			if v, ok := vals[i].(int); ok {
				cur.Left = &TreeNode{Val: v}
				queue = append(queue, cur.Left)
			}
			i++
		}
		if i < len(vals) {
			if v, ok := vals[i].(int); ok {
				cur.Right = &TreeNode{Val: v}
				queue = append(queue, cur.Right)
			}
			i++
		}
	}
	return root
}

// bm034testTreeString 把二叉树序列化成牛客风格的层序字符串（空结点记为 #，末尾的 # 省略）。
func bm034testTreeString(root *TreeNode) string {
	if root == nil {
		return "{}"
	}
	var parts []string
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur == nil {
			parts = append(parts, "#")
			continue
		}
		parts = append(parts, strconv.Itoa(cur.Val))
		queue = append(queue, cur.Left, cur.Right)
	}
	for len(parts) > 0 && parts[len(parts)-1] == "#" {
		parts = parts[:len(parts)-1]
	}
	return "{" + strings.Join(parts, ",") + "}"
}

// buildLeftChain 构造一条向左延伸、根结点值为 n、依次递减到 1 的链（合法 BST）。
func buildLeftChain(n int) *TreeNode {
	var root *TreeNode
	for v := 1; v <= n; v++ {
		root = &TreeNode{Val: v, Left: root}
	}
	return root
}

// buildBadLeftChain 构造根为 1、向下递增的左链（非法 BST）。
func buildBadLeftChain(n int) *TreeNode {
	var root *TreeNode
	for v := n; v >= 1; v-- {
		root = &TreeNode{Val: v, Left: root}
	}
	return root
}

func TestBM034(t *testing.T) {
	tests := []struct {
		name string
		root *TreeNode
		want bool
	}{
		{"空树视为 BST", nil, true},
		{"单结点", bm034testBuildTree([]any{1}), true},
		{"官方示例: 左孩子不小于根", bm034testBuildTree([]any{1, 2, 3}), false},
		{"官方示例: 标准 BST", bm034testBuildTree([]any{2, 1, 3}), true},
		{"官方示例: 右子树的左孩子小于根", bm034testBuildTree([]any{5, 4, 6, nil, nil, 3, 7}), false},
		{"局部父子正确但全局越界", bm034testBuildTree([]any{10, 5, 15, nil, nil, 6, 20}), false},
		{"只有左孩子", bm034testBuildTree([]any{2, 1}), true},
		{"只有右孩子", bm034testBuildTree([]any{2, nil, 3}), true},
		{"只有右孩子但值更小", bm034testBuildTree([]any{2, nil, 1}), false},
		{"左孩子等于根", bm034testBuildTree([]any{2, 2}), false},
		{"右孩子等于根", bm034testBuildTree([]any{5, 3, 5}), false},
		{"深层重复值", bm034testBuildTree([]any{5, 3, 7, nil, nil, 5, 8}), false},
		{"标准 BST 七结点", bm034testBuildTree([]any{5, 3, 8, 1, 4, 7, 9}), true},
		{"左子树里的谷底", bm034testBuildTree([]any{8, 4, 12, 2, 6, 10, 14, 1, 3, 5, 7, 9, 11, 13, 15}), true},
		{"右子树中出现极小值", bm034testBuildTree([]any{8, 4, 12, 2, 6, 10, 14, 1, 3, 5, 7, 9, 11, 13, 0}), false},
		// 上下界用 nil 表示无界，所以结点取到 int 的极值也不受影响。
		{"结点值取 MinInt", bm034testBuildTree([]any{math.MinInt + 1, math.MinInt}), true},
		{"结点值取 MaxInt", bm034testBuildTree([]any{math.MaxInt - 1, math.MaxInt - 2, math.MaxInt}), true},
		{"结点值取 MaxInt 且越界", bm034testBuildTree([]any{math.MaxInt, 0, math.MaxInt - 1}), false},
		{"左偏链 1000 结点合法", buildLeftChain(1000), true},
		{"左偏链 1000 结点非法", buildBadLeftChain(1000), false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := SolveBM034(tc.root)
			if got != tc.want {
				t.Fatalf("期望 %v，实际 %v（输入树 %s）", tc.want, got, bm034testTreeString(tc.root))
			}
		})
	}
}

// ---------------------------------------------------------------------------
// BM035 判断是不是完全二叉树 · 测试
// ---------------------------------------------------------------------------
// bm035testBuildTree 按牛客的层序数组构造二叉树：nil 表示该位置为空结点，空结点不再展开子结点。
func bm035testBuildTree(vals []any) *TreeNode {
	if len(vals) == 0 || vals[0] == nil {
		return nil
	}
	root := &TreeNode{Val: vals[0].(int)}
	queue := []*TreeNode{root}
	i := 1
	for len(queue) > 0 && i < len(vals) {
		cur := queue[0]
		queue = queue[1:]
		if i < len(vals) {
			if v, ok := vals[i].(int); ok {
				cur.Left = &TreeNode{Val: v}
				queue = append(queue, cur.Left)
			}
			i++
		}
		if i < len(vals) {
			if v, ok := vals[i].(int); ok {
				cur.Right = &TreeNode{Val: v}
				queue = append(queue, cur.Right)
			}
			i++
		}
	}
	return root
}

// buildCompleteTree 构造一棵结点数恰好为 n 的完全二叉树（结点值为 1..n，按层序依次填入）。
func buildCompleteTree(n int) *TreeNode {
	if n <= 0 {
		return nil
	}
	nodes := make([]*TreeNode, n)
	for i := range nodes {
		nodes[i] = &TreeNode{Val: i + 1}
	}
	for i := 0; i < n; i++ {
		if 2*i+1 < n {
			nodes[i].Left = nodes[2*i+1]
		}
		if 2*i+2 < n {
			nodes[i].Right = nodes[2*i+2]
		}
	}
	return nodes[0]
}

// bm035testBuildLeftChain 构造向下一直走左孩子的链，结点值从 n 递减到 1。
func bm035testBuildLeftChain(n int) *TreeNode {
	var root *TreeNode
	for v := 1; v <= n; v++ {
		root = &TreeNode{Val: v, Left: root}
	}
	return root
}

func TestBM035(t *testing.T) {
	tests := []struct {
		name string
		root *TreeNode
		want bool
	}{
		{"空树视为完全二叉树", nil, true},
		{"单结点", bm035testBuildTree([]any{1}), true},
		{"官方示例: 最后一层连续靠左", bm035testBuildTree([]any{1, 2, 3, 4, 5, 6}), true},
		{"官方示例: 空结点后又出现非空结点", bm035testBuildTree([]any{1, 2, 3, nil, nil, 4}), false},
		{"左空右不空", bm035testBuildTree([]any{1, nil, 2}), false},
		{"右空左不空", bm035testBuildTree([]any{1, 2, nil}), true},
		{"满二叉树三层", bm035testBuildTree([]any{1, 2, 3, 4, 5, 6, 7}), true},
		{"左侧缺位但右侧有孩子", bm035testBuildTree([]any{1, 2, 3, 4, 5, nil, 7}), false},
		{"左侧缺位且右侧无孩子", bm035testBuildTree([]any{1, 2, 3, 4, 5, nil, nil}), true},
		{"左链三结点", bm035testBuildTree([]any{1, 2, nil, 3}), false},
		{"右链三结点", bm035testBuildTree([]any{1, nil, 2, nil, 3}), false},
		{"最后一层差最右一个", bm035testBuildTree([]any{1, 2, 3, 4, 5, 6, nil}), true},
		{"新增结点挂到最深处", bm035testBuildTree([]any{1, 2, 3, 4, 5, 6, 7, 8}), true},
		{"第八个结点挂在右子树上", bm035testBuildTree([]any{1, 2, 3, 4, 5, 6, 7, nil, 8}), false},
		{"一千结点完全二叉树", buildCompleteTree(1000), true},
		{"一千结点左链", bm035testBuildLeftChain(1000), false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := SolveBM035(tc.root)
			if got != tc.want {
				t.Fatalf("期望 %v，实际 %v", tc.want, got)
			}
		})
	}

	t.Run("完全二叉树末尾多挂一个左孩子", func(t *testing.T) {
		// 1000 个结点的完全二叉树再往最后一个结点（编号 1000）下面挂左孩子，
		// 新结点编号是 2000，与 1000 之间出现断档，因此不再是完全二叉树。
		root := buildCompleteTree(1000)
		nodes := make([]*TreeNode, 0, 1000)
		queue := []*TreeNode{root}
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			nodes = append(nodes, cur)
			if cur.Left != nil {
				queue = append(queue, cur.Left)
			}
			if cur.Right != nil {
				queue = append(queue, cur.Right)
			}
		}
		last := nodes[len(nodes)-1]
		last.Left = &TreeNode{Val: -1}
		if SolveBM035(root) {
			t.Fatalf("期望 false，实际 true")
		}
	})
}

// ---------------------------------------------------------------------------
// BM036 判断是不是平衡二叉树 · 测试
// ---------------------------------------------------------------------------
// bm036testBuildTree 按牛客的层序数组构造二叉树：nil 表示该位置为空结点，空结点不再展开子结点。
func bm036testBuildTree(vals []any) *TreeNode {
	if len(vals) == 0 || vals[0] == nil {
		return nil
	}
	root := &TreeNode{Val: vals[0].(int)}
	queue := []*TreeNode{root}
	i := 1
	for len(queue) > 0 && i < len(vals) {
		cur := queue[0]
		queue = queue[1:]
		if i < len(vals) {
			if v, ok := vals[i].(int); ok {
				cur.Left = &TreeNode{Val: v}
				queue = append(queue, cur.Left)
			}
			i++
		}
		if i < len(vals) {
			if v, ok := vals[i].(int); ok {
				cur.Right = &TreeNode{Val: v}
				queue = append(queue, cur.Right)
			}
			i++
		}
	}
	return root
}

// bm036testBuildPerfectTree 构造高度为 depth、结点值逐层递增的满二叉树。
func bm036testBuildPerfectTree(depth int) *TreeNode {
	if depth <= 0 {
		return nil
	}
	node := &TreeNode{Val: depth}
	node.Left = bm036testBuildPerfectTree(depth - 1)
	node.Right = bm036testBuildPerfectTree(depth - 1)
	return node
}

// bm036testBuildLeftChain 构造一条向下一直走左孩子的链，结点值从 n 递增到 1 的逆序。
func bm036testBuildLeftChain(n int) *TreeNode {
	var root *TreeNode
	for v := 1; v <= n; v++ {
		root = &TreeNode{Val: v, Left: root}
	}
	return root
}

func TestBM036(t *testing.T) {
	// 满二叉树高度 10 一定是平衡的；往最左侧叶子下面挂一条两个结点的链，
	// 该层以下多出 2 高度，最左边那条路径与兄弟子树的高度差就会变成 2。
	perfect := bm036testBuildPerfectTree(10)
	extra := bm036testBuildPerfectTree(10)
	leaf := extra
	for leaf.Left != nil {
		leaf = leaf.Left
	}
	leaf.Left = &TreeNode{Val: -1}
	leaf.Left.Left = &TreeNode{Val: -2}

	tests := []struct {
		name string
		root *TreeNode
		want bool
	}{
		{"空树视为平衡", nil, true},
		{"单结点", bm036testBuildTree([]any{1}), true},
		{"官方示例", bm036testBuildTree([]any{1, 2, 3, nil, nil, 4, 5}), true},
		{"高度差恰为 2", bm036testBuildTree([]any{1, 2, nil, 3}), false},
		{"高度差恰为 1", bm036testBuildTree([]any{1, 2, 3, 4}), true},
		{"左右子树高度差 2 但子结点各自平衡", bm036testBuildTree([]any{1, 2, 2, 3, 3, nil, nil, 4, 4}), false},
		{"根平衡且子树也平衡", bm036testBuildTree([]any{3, 9, 20, nil, nil, 15, 7}), true},
		{"两侧对称深度 2", bm036testBuildTree([]any{1, 2, 2, 3, nil, nil, 3}), true},
		{"两侧对称但一侧多一层", bm036testBuildTree([]any{1, 2, 2, 3, nil, nil, 3, 4}), false},
		{"深层出现高度差 2", bm036testBuildTree([]any{1, 2, 3, 4, nil, nil, nil, 5}), false},
		{"单结点链三层", bm036testBuildTree([]any{1, nil, 2, nil, 3}), false},
		{"满二叉树深度 10", perfect, true},
		{"满二叉树最左叶子下多一个孩子", extra, false},
		{"一千结点左链", bm036testBuildLeftChain(1000), false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := SolveBM036(tc.root)
			if got != tc.want {
				t.Fatalf("期望 %v，实际 %v", tc.want, got)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// BM037 二叉搜索树的最近公共祖先 · 测试
// ---------------------------------------------------------------------------
// nullVal 是测试用的哨兵值，表示层序数组中的空结点（即题目里的 '#'）。
const nullVal = -1 << 62

// bm037testBuildTree 由层序数组构造二叉树，nullVal 表示空结点；空数组或首元素为 nullVal 时返回空树。
func bm037testBuildTree(vals []int) *TreeNode {
	if len(vals) == 0 || vals[0] == nullVal {
		return nil
	}
	root := &TreeNode{Val: vals[0]}
	queue := []*TreeNode{root}
	idx := 1
	for len(queue) > 0 && idx < len(vals) {
		node := queue[0]
		queue = queue[1:]
		// 先补左孩子、再补右孩子，与牛客 {...,#,...} 的层序书写顺序一致。
		if idx < len(vals) {
			if vals[idx] != nullVal {
				node.Left = &TreeNode{Val: vals[idx]}
				queue = append(queue, node.Left)
			}
			idx++
		}
		if idx < len(vals) {
			if vals[idx] != nullVal {
				node.Right = &TreeNode{Val: vals[idx]}
				queue = append(queue, node.Right)
			}
			idx++
		}
	}
	return root
}

// containsNode 判断 target 是否是 root 这棵树上的结点（按指针身份比较）。
func containsNode(root, target *TreeNode) bool {
	if root == nil {
		return false
	}
	if root == target {
		return true
	}
	return containsNode(root.Left, target) || containsNode(root.Right, target)
}

func TestBM037(t *testing.T) {
	// 牛客官方示例使用的二叉搜索树：{7,1,12,0,4,11,14,#,#,3,5}
	official := []int{7, 1, 12, 0, 4, 11, 14, nullVal, nullVal, 3, 5}
	// 只有右孩子的退化成链的 BST：5 -> 6 -> 7
	rightChain := []int{5, nullVal, 6, nullVal, 7}

	cases := []struct {
		name    string
		vals    []int
		p, q    int
		want    int
		wantNil bool
	}{
		{name: "官方示例1_分居两侧", vals: official, p: 3, q: 12, want: 7},
		{name: "官方示例2_一个是另一个的父结点", vals: official, p: 3, q: 4, want: 4},
		{name: "交换参数顺序结果不变", vals: official, p: 12, q: 3, want: 7},
		{name: "其中一个是根结点", vals: official, p: 7, q: 5, want: 7},
		{name: "同在左子树", vals: official, p: 0, q: 5, want: 1},
		{name: "同在右子树", vals: official, p: 11, q: 14, want: 12},
		{name: "单结点_查询自己", vals: []int{1}, p: 1, q: 1, want: 1},
		{name: "右斜链_根与最深结点", vals: rightChain, p: 5, q: 7, want: 5},
		{name: "右斜链_中间结点", vals: rightChain, p: 6, q: 7, want: 6},
		{name: "空树", vals: nil, wantNil: true},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			root := bm037testBuildTree(tt.vals)
			got := SolveBM037(root, tt.p, tt.q)
			if tt.wantNil {
				if got != nil {
					t.Fatalf("期望 nil，实际 %d", got.Val)
				}
				return
			}
			if got == nil {
				t.Fatalf("期望 %d，实际 nil", tt.want)
			}
			if got.Val != tt.want {
				t.Fatalf("期望 %d，实际 %d", tt.want, got.Val)
			}
			// 返回的必须是原树上的结点本身，而不是新建的同值结点。
			if !containsNode(root, got) {
				t.Fatalf("返回的结点 %d 不是原树上的结点", got.Val)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// BM038 在二叉树中找到两个节点的最近公共祖先 · 测试
// ---------------------------------------------------------------------------
// bm038testNullVal 是测试用的哨兵值，表示层序数组中的空结点（即题目里的 '#'）。
const bm038testNullVal = -1 << 62

// bm038testBuildTree 由层序数组构造二叉树，bm038testNullVal 表示空结点；空数组或首元素为 bm038testNullVal 时返回空树。
func bm038testBuildTree(vals []int) *TreeNode {
	if len(vals) == 0 || vals[0] == bm038testNullVal {
		return nil
	}
	root := &TreeNode{Val: vals[0]}
	queue := []*TreeNode{root}
	idx := 1
	for len(queue) > 0 && idx < len(vals) {
		node := queue[0]
		queue = queue[1:]
		// 先补左孩子、再补右孩子，与牛客 {...,#,...} 的层序书写顺序一致。
		if idx < len(vals) {
			if vals[idx] != bm038testNullVal {
				node.Left = &TreeNode{Val: vals[idx]}
				queue = append(queue, node.Left)
			}
			idx++
		}
		if idx < len(vals) {
			if vals[idx] != bm038testNullVal {
				node.Right = &TreeNode{Val: vals[idx]}
				queue = append(queue, node.Right)
			}
			idx++
		}
	}
	return root
}

func TestBM038(t *testing.T) {
	// 牛客官方示例使用的二叉树：{3,5,1,6,2,0,8,#,#,7,4}
	official := []int{3, 5, 1, 6, 2, 0, 8, bm038testNullVal, bm038testNullVal, 7, 4}
	// 只有右孩子的链：1 -> 2 -> 3
	rightChain := []int{1, bm038testNullVal, 2, bm038testNullVal, 3}
	// 只有左孩子的链：1 -> 2 -> 3
	leftChain := []int{1, 2, bm038testNullVal, 3}

	cases := []struct {
		name string
		vals []int
		o1   int
		o2   int
		want int // -1 表示期望「找不到」，即空树
	}{
		{name: "官方示例1_分居根的两侧", vals: official, o1: 5, o2: 1, want: 3},
		{name: "官方示例2_一个是另一个的父结点", vals: official, o1: 2, o2: 7, want: 2},
		{name: "同在左子树_跨两层", vals: official, o1: 6, o2: 4, want: 5},
		{name: "其中一个是根结点", vals: official, o1: 3, o2: 4, want: 3},
		{name: "分居左子树两侧", vals: official, o1: 0, o2: 8, want: 1},
		{name: "两个兄弟结点", vals: official, o1: 7, o2: 4, want: 2},
		{name: "交换参数顺序结果不变", vals: official, o1: 1, o2: 5, want: 3},
		{name: "单结点_查询自己", vals: []int{1}, o1: 1, o2: 1, want: 1},
		{name: "右斜链_根与最深结点", vals: rightChain, o1: 1, o2: 3, want: 1},
		{name: "右斜链_中间结点", vals: rightChain, o1: 2, o2: 3, want: 2},
		{name: "左斜链_中间结点", vals: leftChain, o1: 2, o2: 3, want: 2},
		{name: "空树", vals: nil, want: -1},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			root := bm038testBuildTree(tt.vals)
			got := SolveBM038(root, tt.o1, tt.o2)
			if got != tt.want {
				t.Fatalf("期望 %d，实际 %d", tt.want, got)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// BM039 序列化二叉树 · 测试
// ---------------------------------------------------------------------------
// bm039testNullVal 是测试用的哨兵值，表示层序数组中的空结点（即序列化结果里的 '#'）。
const bm039testNullVal = -1 << 62

// bm039testBuildTree 由层序数组构造二叉树，bm039testNullVal 表示空结点；空数组或首元素为 bm039testNullVal 时返回空树。
func bm039testBuildTree(vals []int) *TreeNode {
	if len(vals) == 0 || vals[0] == bm039testNullVal {
		return nil
	}
	root := &TreeNode{Val: vals[0]}
	queue := []*TreeNode{root}
	idx := 1
	for len(queue) > 0 && idx < len(vals) {
		node := queue[0]
		queue = queue[1:]
		// 先补左孩子、再补右孩子，与牛客 {...,#,...} 的层序书写顺序一致。
		if idx < len(vals) {
			if vals[idx] != bm039testNullVal {
				node.Left = &TreeNode{Val: vals[idx]}
				queue = append(queue, node.Left)
			}
			idx++
		}
		if idx < len(vals) {
			if vals[idx] != bm039testNullVal {
				node.Right = &TreeNode{Val: vals[idx]}
				queue = append(queue, node.Right)
			}
			idx++
		}
	}
	return root
}

// levelOrder 按层序输出结点值（空结点记为 bm039testNullVal，末尾连续的空结点去掉），用于整树比较。
func levelOrder(root *TreeNode) []int {
	res := make([]int, 0)
	if root == nil {
		return res
	}
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		if node == nil {
			res = append(res, bm039testNullVal)
			continue // 空结点没有孩子，不再入队
		}
		res = append(res, node.Val)
		queue = append(queue, node.Left, node.Right)
	}
	for len(res) > 0 && res[len(res)-1] == bm039testNullVal {
		res = res[:len(res)-1]
	}
	return res
}

// renderVals 把层序结果渲染成牛客风格（空结点显示为 '#'），便于失败信息阅读。
func renderVals(vals []int) []string {
	out := make([]string, len(vals))
	for i, v := range vals {
		if v == bm039testNullVal {
			out[i] = "#"
			continue
		}
		out[i] = strconv.Itoa(v)
	}
	return out
}

// diffLevelOrder 比较两棵树的整体形状，一致时返回空串。
func diffLevelOrder(want, got *TreeNode) string {
	a, b := levelOrder(want), levelOrder(got)
	if slices.Equal(a, b) {
		return ""
	}
	return fmt.Sprintf("期望 %v，实际 %v", renderVals(a), renderVals(b))
}

// bm039testRandomTree 把 vals 依次挂到随机空位置上，生成任意形状的二叉树（用于随机往返测试）。
func bm039testRandomTree(rnd *rand.Rand, vals []int) *TreeNode {
	if len(vals) == 0 {
		return nil
	}
	root := &TreeNode{Val: vals[0]}
	// slots 保存「还能挂孩子的空位置」，即孩子指针的地址，取随机位置即可 O(1) 插入。
	slots := []**TreeNode{&root.Left, &root.Right}
	for _, v := range vals[1:] {
		i := rnd.Intn(len(slots))
		slot := slots[i]
		*slot = &TreeNode{Val: v}
		// 被占用的坑位换成新结点自己的两个孩子位置，其余候选保持不变。
		slots[i] = &(*slot).Left
		slots = append(slots, &(*slot).Right)
	}
	return root
}

func TestBM039(t *testing.T) {
	cases := []struct {
		name string
		vals []int
		want string // 期望的序列化字符串
	}{
		{name: "空树", vals: nil, want: ""},
		{name: "单结点", vals: []int{1}, want: "1"},
		{name: "官方示例1_带空结点", vals: []int{1, 2, 3, bm039testNullVal, bm039testNullVal, 6, 7}, want: "1,2,3,#,#,6,7"},
		{name: "官方示例2_满二叉树", vals: []int{8, 6, 10, 5, 7, 9, 11}, want: "8,6,10,5,7,9,11"},
		{name: "只有右孩子的链", vals: []int{1, bm039testNullVal, 2, bm039testNullVal, 3}, want: "1,#,2,#,3"},
		{name: "只有左孩子的链", vals: []int{1, 2, bm039testNullVal, 3}, want: "1,2,#,3"},
		{name: "根只有右孩子", vals: []int{1, bm039testNullVal, 2}, want: "1,#,2"},
		{name: "含负数与空结点", vals: []int{-1, -2, -3, bm039testNullVal, bm039testNullVal, -4}, want: "-1,-2,-3,#,#,-4"},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			root := bm039testBuildTree(tt.vals)

			got := SerializeBM039(root)
			if got != tt.want {
				t.Fatalf("序列化期望 %q，实际 %q", tt.want, got)
			}

			back := DeserializeBM039(got)
			if diff := diffLevelOrder(root, back); diff != "" {
				t.Fatalf("反序列化后树形不一致：%s（序列化结果 %q）", diff, got)
			}

			// 幂等性：还原出的树再序列化一次，必须得到完全相同的字符串。
			if again := SerializeBM039(back); again != tt.want {
				t.Fatalf("二次序列化期望 %q，实际 %q", tt.want, again)
			}
		})
	}

	// 手写的输入串：反序列化需要容忍尾部占位符、多余空格、以及 '#' 表示的空树。
	t.Run("手工输入_尾部占位符被忽略", func(t *testing.T) {
		got := DeserializeBM039("1,#,2,#,#")
		want := bm039testBuildTree([]int{1, bm039testNullVal, 2})
		if diff := diffLevelOrder(want, got); diff != "" {
			t.Fatalf("树形不一致：%s", diff)
		}
	})

	t.Run("手工输入_含空格", func(t *testing.T) {
		got := DeserializeBM039("1, #, 2")
		want := bm039testBuildTree([]int{1, bm039testNullVal, 2})
		if diff := diffLevelOrder(want, got); diff != "" {
			t.Fatalf("树形不一致：%s", diff)
		}
	})

	t.Run("手工输入_仅井号表示空树", func(t *testing.T) {
		if got := DeserializeBM039("#"); got != nil {
			t.Fatalf("期望 nil，实际根结点为 %d", got.Val)
		}
	})

	t.Run("随机形状树往返", func(t *testing.T) {
		rnd := rand.New(rand.NewSource(39))
		for i := 0; i < 30; i++ {
			n := 1 + rnd.Intn(30)
			vals := rand.Perm(n)
			for j := range vals {
				vals[j] -= n / 2 // 让结点值包含负数
			}
			root := bm039testRandomTree(rnd, vals)
			ser := SerializeBM039(root)
			back := DeserializeBM039(ser)
			if diff := diffLevelOrder(root, back); diff != "" {
				t.Fatalf("第 %d 棵随机树往返失败：%s（序列化结果 %q）", i, diff, ser)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// BM040 重建二叉树 · 测试
// ---------------------------------------------------------------------------
// bm040testNullVal 是测试用的哨兵值，表示层序数组中的空结点（即题目里的 '#'）。
const bm040testNullVal = -1 << 62

// bm040testLevelOrder 按层序输出结点值（空结点记为 bm040testNullVal，末尾连续的空结点去掉），用于整树比较。
func bm040testLevelOrder(root *TreeNode) []int {
	res := make([]int, 0)
	if root == nil {
		return res
	}
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		if node == nil {
			res = append(res, bm040testNullVal)
			continue // 空结点没有孩子，不再入队
		}
		res = append(res, node.Val)
		queue = append(queue, node.Left, node.Right)
	}
	for len(res) > 0 && res[len(res)-1] == bm040testNullVal {
		res = res[:len(res)-1]
	}
	return res
}

// bm040testRenderVals 把层序结果渲染成牛客风格（空结点显示为 '#'），便于失败信息阅读。
func bm040testRenderVals(vals []int) []string {
	out := make([]string, len(vals))
	for i, v := range vals {
		if v == bm040testNullVal {
			out[i] = "#"
			continue
		}
		out[i] = strconv.Itoa(v)
	}
	return out
}

// bm040testDiffLevelOrder 比较两棵树的整体形状，一致时返回空串。
func bm040testDiffLevelOrder(want, got *TreeNode) string {
	a, b := bm040testLevelOrder(want), bm040testLevelOrder(got)
	if slices.Equal(a, b) {
		return ""
	}
	return fmt.Sprintf("期望 %v，实际 %v", bm040testRenderVals(a), bm040testRenderVals(b))
}

// preorder 返回前序遍历结果（测试内部独立实现，用于随机用例构造输入）。
func preorder(root *TreeNode, out []int) []int {
	if root == nil {
		return out
	}
	out = append(out, root.Val)
	out = preorder(root.Left, out)
	out = preorder(root.Right, out)
	return out
}

// inorder 返回中序遍历结果（测试内部独立实现，用于随机用例构造输入）。
func inorder(root *TreeNode, out []int) []int {
	if root == nil {
		return out
	}
	out = inorder(root.Left, out)
	out = append(out, root.Val)
	out = inorder(root.Right, out)
	return out
}

// bm040testRandomTree 把 vals 依次挂到随机空位置上，生成任意形状的二叉树（用于随机用例）。
func bm040testRandomTree(rnd *rand.Rand, vals []int) *TreeNode {
	if len(vals) == 0 {
		return nil
	}
	root := &TreeNode{Val: vals[0]}
	// slots 保存「还能挂孩子的空位置」，即孩子指针的地址，取随机位置即可 O(1) 插入。
	slots := []**TreeNode{&root.Left, &root.Right}
	for _, v := range vals[1:] {
		i := rnd.Intn(len(slots))
		slot := slots[i]
		*slot = &TreeNode{Val: v}
		// 被占用的坑位换成新结点自己的两个孩子位置，其余候选保持不变。
		slots[i] = &(*slot).Left
		slots = append(slots, &(*slot).Right)
	}
	return root
}

func TestBM040(t *testing.T) {
	cases := []struct {
		name string
		pre  []int
		vin  []int
		want []int // 期望重建结果的层序表示，bm040testNullVal 表示空结点
	}{
		{
			name: "官方示例1",
			pre:  []int{1, 2, 4, 5, 3, 6, 7},
			vin:  []int{4, 2, 5, 1, 6, 3, 7},
			want: []int{1, 2, 3, 4, 5, 6, 7},
		},
		{
			name: "常规用例_不平衡",
			pre:  []int{3, 9, 20, 15, 7},
			vin:  []int{9, 3, 15, 20, 7},
			want: []int{3, 9, 20, bm040testNullVal, bm040testNullVal, 15, 7},
		},
		{
			name: "官方示例2_单结点",
			pre:  []int{1},
			vin:  []int{1},
			want: []int{1},
		},
		{
			name: "空树",
			pre:  nil,
			vin:  nil,
			want: nil,
		},
		{
			name: "只有左孩子的链",
			pre:  []int{1, 2, 3},
			vin:  []int{3, 2, 1},
			want: []int{1, 2, bm040testNullVal, 3},
		},
		{
			name: "只有右孩子的链",
			pre:  []int{1, 2, 3},
			vin:  []int{1, 2, 3},
			want: []int{1, bm040testNullVal, 2, bm040testNullVal, 3},
		},
		{
			name: "单个左孩子",
			pre:  []int{1, 2},
			vin:  []int{2, 1},
			want: []int{1, 2},
		},
		{
			name: "单个右孩子",
			pre:  []int{1, 2},
			vin:  []int{1, 2},
			want: []int{1, bm040testNullVal, 2},
		},
		{
			name: "长度不一致_返回空",
			pre:  []int{1, 2},
			vin:  []int{2},
			want: nil,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := SolveBM040(tt.pre, tt.vin)
			want := buildTreeFromLevel(tt.want)
			if diff := bm040testDiffLevelOrder(want, got); diff != "" {
				t.Fatalf("重建结果不一致：%s", diff)
			}
		})
	}

	t.Run("随机树_前序中序往返", func(t *testing.T) {
		rnd := rand.New(rand.NewSource(40))
		for i := 0; i < 30; i++ {
			n := 1 + rnd.Intn(20)
			vals := rand.Perm(n)
			for j := range vals {
				vals[j]++ // 题目约定 val >= 1
			}
			root := bm040testRandomTree(rnd, vals)
			pre := preorder(root, nil)
			vin := inorder(root, nil)

			got := SolveBM040(pre, vin)
			if diff := bm040testDiffLevelOrder(root, got); diff != "" {
				t.Fatalf("第 %d 棵随机树重建失败：%s（pre=%v, vin=%v）", i, diff, pre, vin)
			}
			// 重建结果的前序、中序也必须与原输入完全一致。
			if !slices.Equal(preorder(got, nil), pre) {
				t.Fatalf("第 %d 棵随机树前序不一致：期望 %v，实际 %v", i, pre, preorder(got, nil))
			}
			if !slices.Equal(inorder(got, nil), vin) {
				t.Fatalf("第 %d 棵随机树中序不一致：期望 %v，实际 %v", i, vin, inorder(got, nil))
			}
		}
	})
}

// buildTreeFromLevel 由层序数组构造二叉树，bm040testNullVal 表示空结点；空数组或首元素为 bm040testNullVal 时返回空树。
func buildTreeFromLevel(vals []int) *TreeNode {
	if len(vals) == 0 || vals[0] == bm040testNullVal {
		return nil
	}
	root := &TreeNode{Val: vals[0]}
	queue := []*TreeNode{root}
	idx := 1
	for len(queue) > 0 && idx < len(vals) {
		node := queue[0]
		queue = queue[1:]
		// 先补左孩子、再补右孩子，与牛客 {...,#,...} 的层序书写顺序一致。
		if idx < len(vals) {
			if vals[idx] != bm040testNullVal {
				node.Left = &TreeNode{Val: vals[idx]}
				queue = append(queue, node.Left)
			}
			idx++
		}
		if idx < len(vals) {
			if vals[idx] != bm040testNullVal {
				node.Right = &TreeNode{Val: vals[idx]}
				queue = append(queue, node.Right)
			}
			idx++
		}
	}
	return root
}

// ---------------------------------------------------------------------------
// BM041 输出二叉树的右视图 · 测试
// ---------------------------------------------------------------------------
// bm041testNullVal 是测试用的哨兵值，表示层序数组中的空结点（即题目里的 '#'）。
const bm041testNullVal = -1 << 62

// bm041testBuildTree 由层序数组构造二叉树，bm041testNullVal 表示空结点；空数组或首元素为 bm041testNullVal 时返回空树。
func bm041testBuildTree(vals []int) *TreeNode {
	if len(vals) == 0 || vals[0] == bm041testNullVal {
		return nil
	}
	root := &TreeNode{Val: vals[0]}
	queue := []*TreeNode{root}
	idx := 1
	for len(queue) > 0 && idx < len(vals) {
		node := queue[0]
		queue = queue[1:]
		// 先补左孩子、再补右孩子，与牛客 {...,#,...} 的层序书写顺序一致。
		if idx < len(vals) {
			if vals[idx] != bm041testNullVal {
				node.Left = &TreeNode{Val: vals[idx]}
				queue = append(queue, node.Left)
			}
			idx++
		}
		if idx < len(vals) {
			if vals[idx] != bm041testNullVal {
				node.Right = &TreeNode{Val: vals[idx]}
				queue = append(queue, node.Right)
			}
			idx++
		}
	}
	return root
}

// bm041testPreorder 返回前序遍历结果（测试内部独立实现，用于随机用例构造输入）。
func bm041testPreorder(root *TreeNode, out []int) []int {
	if root == nil {
		return out
	}
	out = append(out, root.Val)
	out = bm041testPreorder(root.Left, out)
	out = bm041testPreorder(root.Right, out)
	return out
}

// bm041testInorder 返回中序遍历结果（测试内部独立实现，用于随机用例构造输入）。
func bm041testInorder(root *TreeNode, out []int) []int {
	if root == nil {
		return out
	}
	out = bm041testInorder(root.Left, out)
	out = append(out, root.Val)
	out = bm041testInorder(root.Right, out)
	return out
}

// refRightView 用与解答不同的算法（DFS 先右后左，记录每个深度第一次遇到的结点）求右视图，
// 用于和层序写法交叉验证。
func refRightView(root *TreeNode) []int {
	res := make([]int, 0)
	var dfs func(node *TreeNode, depth int)
	dfs = func(node *TreeNode, depth int) {
		if node == nil {
			return
		}
		if depth == len(res) {
			// 每个深度第一次被访问到的结点，因为先走右子树，它就是该层最右结点。
			res = append(res, node.Val)
		}
		dfs(node.Right, depth+1)
		dfs(node.Left, depth+1)
	}
	dfs(root, 0)
	return res
}

// bm041testRandomTree 把 vals 依次挂到随机空位置上，生成任意形状的二叉树（用于随机用例）。
func bm041testRandomTree(rnd *rand.Rand, vals []int) *TreeNode {
	if len(vals) == 0 {
		return nil
	}
	root := &TreeNode{Val: vals[0]}
	// slots 保存「还能挂孩子的空位置」，即孩子指针的地址，取随机位置即可 O(1) 插入。
	slots := []**TreeNode{&root.Left, &root.Right}
	for _, v := range vals[1:] {
		i := rnd.Intn(len(slots))
		slot := slots[i]
		*slot = &TreeNode{Val: v}
		// 被占用的坑位换成新结点自己的两个孩子位置，其余候选保持不变。
		slots[i] = &(*slot).Left
		slots = append(slots, &(*slot).Right)
	}
	return root
}

func TestBM041(t *testing.T) {
	// 牛客原题入口：输入前序 + 中序，返回右视图。
	cases := []struct {
		name string
		pre  []int
		vin  []int
		want []int
	}{
		{name: "官方示例1", pre: []int{1, 2, 4, 5, 3}, vin: []int{4, 2, 5, 1, 3}, want: []int{1, 3, 5}},
		{name: "常规用例", pre: []int{1, 2, 4, 5, 3, 6, 7}, vin: []int{4, 2, 5, 1, 6, 3, 7}, want: []int{1, 3, 7}},
		{name: "空树", pre: nil, vin: nil, want: []int{}},
		{name: "单结点", pre: []int{1}, vin: []int{1}, want: []int{1}},
		{name: "只有右孩子的链", pre: []int{1, 2, 3}, vin: []int{1, 2, 3}, want: []int{1, 2, 3}},
		{name: "只有左孩子的链", pre: []int{1, 2, 3}, vin: []int{3, 2, 1}, want: []int{1, 2, 3}},
		{name: "左子树比右子树深", pre: []int{1, 2, 4, 3}, vin: []int{4, 2, 1, 3}, want: []int{1, 3, 4}},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := SolveBM041(tt.pre, tt.vin)
			if got == nil {
				t.Fatalf("期望非 nil 切片（可能为空），实际 nil")
			}
			if !slices.Equal(got, tt.want) {
				t.Fatalf("期望 %v，实际 %v", tt.want, got)
			}
		})
	}

	// 直接用建好的树取右视图的入口。
	treeCases := []struct {
		name string
		vals []int
		want []int
	}{
		{name: "空树", vals: nil, want: []int{}},
		{name: "单结点", vals: []int{7}, want: []int{7}},
		{name: "左子树更深", vals: []int{1, 2, 3, 4}, want: []int{1, 3, 4}},
		{name: "右子树更深", vals: []int{1, 2, 3, bm041testNullVal, bm041testNullVal, bm041testNullVal, 4}, want: []int{1, 3, 4}},
		{name: "只有右孩子的链", vals: []int{1, bm041testNullVal, 2, bm041testNullVal, 3}, want: []int{1, 2, 3}},
		{name: "只有左孩子的链", vals: []int{1, 2, bm041testNullVal, 3}, want: []int{1, 2, 3}},
	}

	for _, tt := range treeCases {
		tt := tt
		t.Run("直接取右视图_"+tt.name, func(t *testing.T) {
			root := bm041testBuildTree(tt.vals)
			got := RightViewBM041(root)
			if got == nil {
				t.Fatalf("期望非 nil 切片（可能为空），实际 nil")
			}
			if !slices.Equal(got, tt.want) {
				t.Fatalf("期望 %v，实际 %v", tt.want, got)
			}
		})
	}

	t.Run("随机树_两种算法交叉验证", func(t *testing.T) {
		rnd := rand.New(rand.NewSource(41))
		for i := 0; i < 30; i++ {
			n := 1 + rnd.Intn(25)
			vals := rand.Perm(n)
			for j := range vals {
				vals[j] *= 3 // 取值范围无所谓的任意正整数
			}
			root := bm041testRandomTree(rnd, vals)

			// 层序（解答）与 DFS（测试内部参考实现）必须给出同一个右视图。
			got := RightViewBM041(root)
			want := refRightView(root)
			if !slices.Equal(got, want) {
				t.Fatalf("第 %d 棵随机树右视图不一致：期望 %v，实际 %v", i, want, got)
			}

			// 由随机树的前序 + 中序走牛客原题入口，结果必须与直接取右视图一致。
			fromRebuild := SolveBM041(bm041testPreorder(root, nil), bm041testInorder(root, nil))
			if !slices.Equal(fromRebuild, want) {
				t.Fatalf("第 %d 棵随机树重建后右视图不一致：期望 %v，实际 %v", i, want, fromRebuild)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// BM042 用两个栈实现队列 · 测试
// ---------------------------------------------------------------------------
func TestBM042(t *testing.T) {
	// 牛客原题形式：用 "PSHx" / "POP" 操作序列驱动队列，返回所有出队元素。
	opCases := []struct {
		name string
		ops  []string
		want []int
	}{
		{name: "官方示例1", ops: []string{"PSH1", "PSH2", "POP", "POP"}, want: []int{1, 2}},
		{name: "官方示例2_中途入队不插队", ops: []string{"PSH1", "PSH2", "POP", "PSH3", "POP", "POP"}, want: []int{1, 2, 3}},
		{name: "push负数", ops: []string{"PSH-5", "PSH3", "POP", "POP"}, want: []int{-5, 3}},
		{name: "没有操作", ops: nil, want: []int{}},
		{name: "只push不pop", ops: []string{"PSH1", "PSH2"}, want: []int{}},
	}

	for _, tt := range opCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := SolveBM042(tt.ops)
			if !slices.Equal(got, tt.want) {
				t.Fatalf("期望 %v，实际 %v", tt.want, got)
			}
		})
	}

	t.Run("空队列", func(t *testing.T) {
		q := NewMyQueue()
		if !q.Empty() {
			t.Fatal("新建的队列应当为空")
		}
		if got := q.Pop(); got != -1 {
			t.Fatalf("空队列 Pop 期望 -1，实际 %d", got)
		}
		if got := q.Peek(); got != -1 {
			t.Fatalf("空队列 Peek 期望 -1，实际 %d", got)
		}
		if !q.Empty() {
			t.Fatal("空队列上的 Pop / Peek 不应改变队列状态")
		}
	})

	t.Run("单元素_peek不出队", func(t *testing.T) {
		q := NewMyQueue()
		q.Push(5)
		if q.Empty() {
			t.Fatal("push 之后队列不应为空")
		}
		// 连续 Peek 两次都不能把元素取走。
		for i := 0; i < 2; i++ {
			if got := q.Peek(); got != 5 {
				t.Fatalf("第 %d 次 Peek 期望 5，实际 %d", i+1, got)
			}
		}
		if got := q.Pop(); got != 5 {
			t.Fatalf("Pop 期望 5，实际 %d", got)
		}
		if !q.Empty() {
			t.Fatal("弹出唯一元素后队列应为空")
		}
	})

	// push / pop / peek / empty 交叉操作：重点覆盖「出队栈非空时不能倒腾」的时机。
	t.Run("交叉操作", func(t *testing.T) {
		q := NewMyQueue()
		steps := []struct {
			op       string
			val      int  // push 的参数
			wantInt  int  // pop / peek 的期望值
			wantBool bool // empty 的期望值
		}{
			{op: "push", val: 1},
			{op: "push", val: 2},
			{op: "peek", wantInt: 1},
			{op: "pop", wantInt: 1},
			{op: "empty", wantBool: false},
			{op: "push", val: 3},
			{op: "peek", wantInt: 2},
			{op: "pop", wantInt: 2},
			{op: "empty", wantBool: false},
			{op: "pop", wantInt: 3},
			{op: "empty", wantBool: true},
		}
		for i, st := range steps {
			switch st.op {
			case "push":
				q.Push(st.val)
			case "pop":
				if got := q.Pop(); got != st.wantInt {
					t.Fatalf("第 %d 步 pop 期望 %d，实际 %d", i+1, st.wantInt, got)
				}
			case "peek":
				if got := q.Peek(); got != st.wantInt {
					t.Fatalf("第 %d 步 peek 期望 %d，实际 %d", i+1, st.wantInt, got)
				}
			case "empty":
				if got := q.Empty(); got != st.wantBool {
					t.Fatalf("第 %d 步 empty 期望 %v，实际 %v", i+1, st.wantBool, got)
				}
			default:
				t.Fatalf("第 %d 步出现了未定义的操作 %q", i+1, st.op)
			}
		}
	})

	// 大批量：先全 push 再全 pop，验证倒腾一次后顺序依然正确。
	t.Run("连续push再连续pop", func(t *testing.T) {
		const n = 5000
		q := NewMyQueue()
		for i := 1; i <= n; i++ {
			q.Push(i)
		}
		for i := 1; i <= n; i++ {
			if got := q.Pop(); got != i {
				t.Fatalf("第 %d 次 Pop 期望 %d，实际 %d", i, i, got)
			}
		}
		if !q.Empty() {
			t.Fatal("全部元素出队后队列应为空")
		}
	})

	// 与基于切片的朴素队列（作为基准模型）逐步对照。
	t.Run("随机交叉操作对照切片模型", func(t *testing.T) {
		rnd := rand.New(rand.NewSource(42))
		q := NewMyQueue()
		model := make([]int, 0)
		for i := 0; i < 2000; i++ {
			if len(model) == 0 || rnd.Intn(2) == 0 {
				v := rnd.Intn(200) - 100 // 取值包含负数
				q.Push(v)
				model = append(model, v)
			} else if rnd.Intn(4) == 0 {
				if got, want := q.Peek(), model[0]; got != want {
					t.Fatalf("第 %d 步 Peek 期望 %d，实际 %d", i+1, want, got)
				}
			} else {
				got, want := q.Pop(), model[0]
				model = model[1:]
				if got != want {
					t.Fatalf("第 %d 步 Pop 期望 %d，实际 %d", i+1, want, got)
				}
			}
			if q.Empty() != (len(model) == 0) {
				t.Fatalf("第 %d 步 Empty 与实际元素个数 %d 不一致", i+1, len(model))
			}
		}
		// 收尾：剩余元素必须按入队顺序出队。
		for len(model) > 0 {
			want := model[0]
			model = model[1:]
			if got := q.Pop(); got != want {
				t.Fatalf("收尾 Pop 期望 %d，实际 %d", want, got)
			}
		}
		if !q.Empty() {
			t.Fatal("全部元素出队后队列应为空")
		}
	})
}

// ---------------------------------------------------------------------------
// BM043 包含min函数的栈 · 测试
// ---------------------------------------------------------------------------
func TestBM043(t *testing.T) {
	cases := []struct {
		name   string
		ops    []string
		params [][]int
		want   []int
	}{
		{
			name:   "官方示例",
			ops:    []string{"push", "push", "min", "top", "pop", "push", "top", "min"},
			params: [][]int{{-1}, {2}, nil, nil, nil, {1}, nil, nil},
			want:   []int{-1, 2, 1, -1},
		},
		{
			name:   "空操作序列",
			ops:    []string{},
			params: nil,
			want:   []int{},
		},
		{
			name:   "单元素",
			ops:    []string{"push", "min", "top", "pop", "min", "top"},
			params: [][]int{{7}, nil, nil, nil, nil, nil},
			want:   []int{7, 7, 0, 0}, // 弹出后栈空，min / top 返回 0 兜底
		},
		{
			name:   "单调递减入栈后逐一出栈",
			ops:    []string{"push", "push", "push", "min", "pop", "min", "top", "pop", "min"},
			params: [][]int{{5}, {4}, {3}, nil, nil, nil, nil, nil, nil},
			want:   []int{3, 4, 4, 5},
		},
		{
			name:   "单调递增入栈最小值不变",
			ops:    []string{"push", "push", "push", "min", "pop", "top", "min"},
			params: [][]int{{1}, {2}, {3}, nil, nil, nil, nil},
			want:   []int{1, 2, 1},
		},
		{
			name:   "重复元素",
			ops:    []string{"push", "push", "push", "min", "top", "pop", "min"},
			params: [][]int{{2}, {2}, {2}, nil, nil, nil, nil},
			want:   []int{2, 2, 2},
		},
		{
			name:   "压入更小的元素后 min 立即更新",
			ops:    []string{"push", "min", "push", "min", "push", "min", "pop", "min"},
			params: [][]int{{3}, nil, {1}, nil, {2}, nil, nil, nil},
			want:   []int{3, 1, 1, 1},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := SolveBM043(c.ops, c.params)
			if !bm043testEqualInts(got, c.want) {
				t.Fatalf("期望 %v，实际 %v", c.want, got)
			}
		})
	}

	// 空栈上的 pop / top / min 必须安全返回零值，且不能让栈进入异常状态。
	t.Run("空栈操作", func(t *testing.T) {
		st := NewMinStack()
		if !st.Empty() {
			t.Fatalf("新建栈应为空")
		}
		if got := st.Pop(); got != 0 {
			t.Fatalf("空栈 Pop 期望 0，实际 %d", got)
		}
		if got := st.Top(); got != 0 {
			t.Fatalf("空栈 Top 期望 0，实际 %d", got)
		}
		if got := st.Min(); got != 0 {
			t.Fatalf("空栈 Min 期望 0，实际 %d", got)
		}
		st.Push(9)
		if got := st.Pop(); got != 9 {
			t.Fatalf("空栈操作后 Push(9) 再 Pop 期望 9，实际 %d", got)
		}
		if !st.Empty() {
			t.Fatalf("Pop 后栈应为空")
		}
	})

	// Pop 必须返回真正被弹出的元素，且两个栈同步出栈。
	t.Run("Pop返回值与同步出栈", func(t *testing.T) {
		st := NewMinStack()
		for _, v := range []int{4, 1, 3} {
			st.Push(v)
		}
		if got := st.Min(); got != 1 {
			t.Fatalf("期望 min=1，实际 %d", got)
		}
		if got := st.Pop(); got != 3 {
			t.Fatalf("期望弹出 3，实际 %d", got)
		}
		if got := st.Min(); got != 1 {
			t.Fatalf("期望 min=1，实际 %d", got)
		}
		if got := st.Pop(); got != 1 {
			t.Fatalf("期望弹出 1，实际 %d", got)
		}
		if got := st.Min(); got != 4 {
			t.Fatalf("弹出最小值后期望 min=4，实际 %d", got)
		}
	})

	// 随机压力用例：与暴力求最小值的朴素栈对拍（固定种子，可复现）。
	t.Run("随机对拍", func(t *testing.T) {
		r := rand.New(rand.NewSource(43))
		st := NewMinStack()
		var plain []int // 朴素栈，用于暴力求最小值
		for step := 0; step < 2000; step++ {
			if len(plain) == 0 || r.Intn(3) != 0 {
				v := r.Intn(2001) - 1000
				st.Push(v)
				plain = append(plain, v)
				continue
			}
			want := plain[len(plain)-1]
			plain = plain[:len(plain)-1]
			if got := st.Pop(); got != want {
				t.Fatalf("第 %d 步 Pop 期望 %d，实际 %d", step, want, got)
			}
			if len(plain) == 0 {
				continue
			}
			wantMin := plain[0]
			for _, v := range plain {
				if v < wantMin {
					wantMin = v
				}
			}
			if got := st.Min(); got != wantMin {
				t.Fatalf("第 %d 步 Min 期望 %d，实际 %d", step, wantMin, got)
			}
			if got := st.Top(); got != plain[len(plain)-1] {
				t.Fatalf("第 %d 步 Top 期望 %d，实际 %d", step, plain[len(plain)-1], got)
			}
		}
	})
}

// bm043testEqualInts 比较两个 int 切片是否完全一致，nil 与空切片视为相等。
func bm043testEqualInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// ---------------------------------------------------------------------------
// BM044 有效括号序列 · 测试
// ---------------------------------------------------------------------------
func TestBM044(t *testing.T) {
	cases := []struct {
		name string
		s    string
		want bool
	}{
		{"官方示例1-只有左括号", "[", false},
		{"官方示例2-一对圆括号", "[]", true},
		{"官方示例3-交错嵌套", "([)]", false},
		{"官方示例4-正确嵌套", "([])", true},
		{"空字符串", "", true},
		{"三种括号并列", "()[]{}", true},
		{"三种括号嵌套", "{[()]}", true},
		{"多余右括号", "]", false},
		{"右括号在前", "}{", false},
		{"多余左括号", "(((", false},
		{"只有一个左括号", "(", false},
		{"只有左括号且配对", "((()))", true},
		{"缺少闭合", "()[](", false},
		{"长度奇数", "(){", false},
		{"类型不匹配", "(]", false},
		{"浅嵌套但错位", "[({})](", false},
		{"括号在大括号内", "{[()]}[{}]", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := SolveBM044(c.s); got != c.want {
				t.Fatalf("输入 %q，期望 %v，实际 %v", c.s, c.want, got)
			}
		})
	}

	// 大规模合法序列：验证不会越界，且能把栈撑到 O(n)。
	t.Run("大规模合法序列", func(t *testing.T) {
		s := strings.Repeat("(){}[]", 2000)
		if !SolveBM044(s) {
			t.Fatalf("长度 %d 的合法序列被判为非法", len(s))
		}
	})
	// 大规模非法序列：最后一个字符改为左括号，末尾必然留下未闭合的括号。
	t.Run("大规模非法序列", func(t *testing.T) {
		s := strings.Repeat("[]", 5000) + "["
		if SolveBM044(s) {
			t.Fatalf("长度 %d 的非法序列被判为合法", len(s))
		}
	})
}

// ---------------------------------------------------------------------------
// BM045 滑动窗口的最大值 · 测试
// ---------------------------------------------------------------------------
func TestBM045(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		size int
		want []int
	}{
		{"官方示例1", []int{2, 3, 4, 2, 6, 2, 5, 1}, 3, []int{4, 4, 6, 6, 6, 5}},
		{"官方示例2", []int{9, 10, 9, -7, -3, 8, 2, -6}, 5, []int{10, 10, 9, 8}},
		{"空数组", []int{}, 1, []int{}},
		{"size大于数组长度", []int{1, 2, 3}, 5, []int{}},
		{"size为0", []int{1, 2, 3}, 0, []int{}},
		{"size为负数", []int{1, 2, 3}, -2, []int{}},
		{"size等于数组长度", []int{3, 1, 4, 1, 5}, 5, []int{5}},
		{"size为1即原数组", []int{3, 1, 4, 1, 5}, 1, []int{3, 1, 4, 1, 5}},
		{"单元素数组", []int{42}, 1, []int{42}},
		{"全部重复", []int{5, 5, 5, 5}, 2, []int{5, 5, 5}},
		{"严格递减", []int{5, 4, 3, 2, 1}, 3, []int{5, 4, 3}},
		{"严格递增", []int{1, 2, 3, 4, 5}, 3, []int{3, 4, 5}},
		{"最大值在窗口中间", []int{1, 3, 1, 2, 0, 5}, 3, []int{3, 3, 2, 5}},
		{"含负数和0", []int{0, -1, -2, -3}, 2, []int{0, -1, -2}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := SolveBM045(c.nums, c.size)
			if !bm045testEqualInts(got, c.want) {
				t.Fatalf("nums=%v size=%d，期望 %v，实际 %v", c.nums, c.size, c.want, got)
			}
		})
	}

	// 与暴力解法对拍：固定种子生成随机数组和窗口，验证 O(n) 单调队列的正确性。
	t.Run("随机对拍暴力解", func(t *testing.T) {
		r := rand.New(rand.NewSource(45))
		for round := 0; round < 300; round++ {
			n := r.Intn(40) + 1
			nums := make([]int, n)
			for i := range nums {
				nums[i] = r.Intn(21) - 10
			}
			size := r.Intn(n+2) + 1 // 允许 size 恰好超过 n
			want := bruteMaxInWindows(nums, size)
			if got := SolveBM045(nums, size); !bm045testEqualInts(got, want) {
				t.Fatalf("第 %d 轮 nums=%v size=%d，期望 %v，实际 %v", round, nums, size, want, got)
			}
		}
	})

	// 大规模用例：n = 10000 且严格递增，验证不会退化成 O(n*size)。
	t.Run("大规模递增数组", func(t *testing.T) {
		const n = 10000
		nums := make([]int, n)
		for i := range nums {
			nums[i] = i
		}
		got := SolveBM045(nums, 100)
		if len(got) != n-100+1 {
			t.Fatalf("期望 %d 个窗口，实际 %d", n-100+1, len(got))
		}
		if got[0] != 99 || got[len(got)-1] != n-1 {
			t.Fatalf("首尾结果错误：got[0]=%d got[last]=%d", got[0], got[len(got)-1])
		}
	})
}

// bruteMaxInWindows 是 O(n*size) 的暴力解法，仅用于测试对拍。
func bruteMaxInWindows(nums []int, size int) []int {
	if size <= 0 || size > len(nums) {
		return []int{}
	}
	res := make([]int, 0, len(nums)-size+1)
	for i := 0; i+size <= len(nums); i++ {
		best := nums[i]
		for j := i + 1; j < i+size; j++ {
			if nums[j] > best {
				best = nums[j]
			}
		}
		res = append(res, best)
	}
	return res
}

// bm045testEqualInts 比较两个 int 切片是否完全一致，nil 与空切片视为相等。
func bm045testEqualInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// ---------------------------------------------------------------------------
// BM046 最小的K个数 · 测试
// ---------------------------------------------------------------------------
func TestBM046(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		k    int
		want []int
	}{
		{"官方示例1", []int{4, 5, 1, 6, 2, 7, 3, 8}, 4, []int{1, 2, 3, 4}},
		{"官方示例2-k为0", []int{1}, 0, []int{}},
		{"k为1取最小", []int{4, 5, 1, 6, 2, 7, 3, 8}, 1, []int{1}},
		{"k等于长度即整体升序", []int{4, 5, 1, 6, 2, 7, 3, 8}, 8, []int{1, 2, 3, 4, 5, 6, 7, 8}},
		{"k大于长度返回空", []int{1, 2, 3}, 5, []int{}},
		{"k为负数", []int{1, 2, 3}, -1, []int{}},
		{"空数组", []int{}, 0, []int{}},
		{"空数组且k大于0", []int{}, 2, []int{}},
		{"重复元素", []int{2, 2, 1, 1, 3}, 3, []int{1, 1, 2}},
		{"全部相同", []int{7, 7, 7, 7}, 2, []int{7, 7}},
		{"含负数与0", []int{0, -3, 5, -1, 2}, 3, []int{-3, -1, 0}},
		{"已升序", []int{1, 2, 3, 4, 5}, 3, []int{1, 2, 3}},
		{"已降序", []int{5, 4, 3, 2, 1}, 2, []int{1, 2}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			origin := append([]int(nil), c.nums...)
			got := SolveBM046(c.nums, c.k)
			if !bm046testEqualInts(got, c.want) {
				t.Fatalf("nums=%v k=%d，期望 %v，实际 %v", c.nums, c.k, c.want, got)
			}
			// 题目要求不能修改调用方数组，这里顺带校验。
			if !bm046testEqualInts(c.nums, origin) {
				t.Fatalf("SolveBM046 不应修改入参：原 %v，现在 %v", origin, c.nums)
			}
		})
	}

	// 随机对拍：与「整体排序后取前 k 个」的参考实现比较（固定种子，可复现）。
	t.Run("随机对拍排序解", func(t *testing.T) {
		r := rand.New(rand.NewSource(46))
		for round := 0; round < 300; round++ {
			n := r.Intn(60) + 1
			nums := make([]int, n)
			for i := range nums {
				nums[i] = r.Intn(11) // 故意取小范围，制造大量重复值
			}
			k := r.Intn(n + 1)
			want := bruteLeastK(nums, k)
			if got := SolveBM046(nums, k); !bm046testEqualInts(got, want) {
				t.Fatalf("第 %d 轮 nums=%v k=%d，期望 %v，实际 %v", round, nums, k, want, got)
			}
		}
	})

	// 特殊输入：全相同的大数组（三路划分应保持线性）与已排序的大数组（随机基准避免退化）。
	t.Run("大规模全相同元素", func(t *testing.T) {
		const n = 10000
		nums := make([]int, n)
		for i := range nums {
			nums[i] = 9
		}
		got := SolveBM046(nums, 3000)
		if len(got) != 3000 {
			t.Fatalf("期望 3000 个元素，实际 %d", len(got))
		}
		for i, v := range got {
			if v != 9 {
				t.Fatalf("第 %d 个元素期望 9，实际 %d", i, v)
			}
		}
	})
	t.Run("大规模已升序", func(t *testing.T) {
		const n = 10000
		nums := make([]int, n)
		for i := range nums {
			nums[i] = i
		}
		got := SolveBM046(nums, 10)
		if len(got) != 10 {
			t.Fatalf("期望 10 个元素，实际 %d", len(got))
		}
		for i, v := range got {
			if v != i {
				t.Fatalf("第 %d 个元素期望 %d，实际 %d", i, i, v)
			}
		}
	})
}

// bruteLeastK 是「整体排序后取前 k 个」的参考实现，仅用于测试对拍。
func bruteLeastK(nums []int, k int) []int {
	if k <= 0 || k > len(nums) {
		return []int{}
	}
	cp := append([]int(nil), nums...)
	sort.Ints(cp)
	return cp[:k]
}

// bm046testEqualInts 比较两个 int 切片是否完全一致，nil 与空切片视为相等。
func bm046testEqualInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// ---------------------------------------------------------------------------
// BM047 寻找第K大 · 测试
// ---------------------------------------------------------------------------
func TestBM047(t *testing.T) {
	cases := []struct {
		name string
		a    []int
		k    int
		want int
	}{
		{"官方示例1", []int{1, 3, 5, 2, 2}, 3, 2},
		{"官方示例2-重复元素不去重", []int{10, 10, 9, 9, 8, 7, 5, 6, 4, 3, 4, 2}, 3, 9},
		{"k为1即最大值", []int{1, 3, 5, 2, 2}, 1, 5},
		{"k等于n即最小值", []int{1, 3, 5, 2, 2}, 5, 1},
		{"单元素", []int{42}, 1, 42},
		{"全部相同", []int{2, 2, 2}, 2, 2},
		{"已升序", []int{1, 2, 3, 4, 5}, 2, 4},
		{"已降序", []int{5, 4, 3, 2, 1}, 2, 4},
		{"含负数", []int{-3, -1, 0, 2, 5}, 2, 2},
		{"只有两个不同值", []int{1, 1, 1, 2, 2, 2}, 3, 2},
		{"重复值取到第一个", []int{3, 3, 3, 1}, 2, 3},
		{"k为0越界兜底", []int{1, 2, 3}, 0, 0},
		{"k大于n越界兜底", []int{1, 2, 3}, 4, 0},
		{"空数组越界兜底", []int{}, 1, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			origin := append([]int(nil), c.a...)
			got := SolveBM047(c.a, c.k)
			if got != c.want {
				t.Fatalf("a=%v k=%d，期望 %d，实际 %d", c.a, c.k, c.want, got)
			}
			// 不应修改调用方数组。
			if len(origin) != len(c.a) {
				t.Fatalf("SolveBM047 不应修改入参长度")
			}
			for i := range origin {
				if origin[i] != c.a[i] {
					t.Fatalf("SolveBM047 不应修改入参：原 %v，现在 %v", origin, c.a)
				}
			}
		})
	}

	// 随机对拍：与「整体排序后取下标 n-k」的参考实现比较（固定种子，可复现）。
	t.Run("随机对拍排序解", func(t *testing.T) {
		r := rand.New(rand.NewSource(47))
		for round := 0; round < 500; round++ {
			n := r.Intn(50) + 1
			a := make([]int, n)
			for i := range a {
				a[i] = r.Intn(9) // 小范围取值，制造大量重复元素
			}
			k := r.Intn(n) + 1
			want := bruteKthLargest(a, k)
			if got := SolveBM047(a, k); got != want {
				t.Fatalf("第 %d 轮 a=%v k=%d，期望 %d，实际 %d", round, a, k, want, got)
			}
		}
	})

	// 大规模有序 / 逆序数组：验证伪随机基准不会让划分退化到 O(n^2)。
	t.Run("大规模升序数组", func(t *testing.T) {
		const n = 10000
		a := make([]int, n)
		for i := range a {
			a[i] = i
		}
		if got := SolveBM047(a, 1); got != n-1 {
			t.Fatalf("第 1 大期望 %d，实际 %d", n-1, got)
		}
		if got := SolveBM047(a, n); got != 0 {
			t.Fatalf("第 %d 大期望 0，实际 %d", n, got)
		}
		if got := SolveBM047(a, 5000); got != n-5000 {
			t.Fatalf("第 5000 大期望 %d，实际 %d", n-5000, got)
		}
	})
	t.Run("大规模全相同元素", func(t *testing.T) {
		const n = 10000
		a := make([]int, n)
		for i := range a {
			a[i] = 6
		}
		for _, k := range []int{1, 5000, n} {
			if got := SolveBM047(a, k); got != 6 {
				t.Fatalf("第 %d 大期望 6，实际 %d", k, got)
			}
		}
	})
}

// bruteKthLargest 是「整体排序后取下标 n-k」的参考实现，仅用于测试对拍。
func bruteKthLargest(a []int, k int) int {
	if k < 1 || k > len(a) {
		return 0
	}
	cp := append([]int(nil), a...)
	sort.Ints(cp)
	return cp[len(cp)-k]
}

// ---------------------------------------------------------------------------
// BM048 数据流中的中位数 · 测试
// ---------------------------------------------------------------------------
func TestBM048(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		want []float64
	}{
		{"官方示例", []int{5, 2, 3, 4, 1, 6, 7, 0, 8}, []float64{5, 3.5, 3, 3.5, 3, 3.5, 4, 3.5, 4}},
		{"空数据流", []int{}, []float64{}},
		{"仅一个元素", []int{7}, []float64{7}},
		{"两个元素取平均", []int{2, 4}, []float64{2, 3}},
		{"全部相同", []int{5, 5, 5, 5}, []float64{5, 5, 5, 5}},
		{"严格递减", []int{5, 4, 3, 2, 1}, []float64{5, 4.5, 4, 3.5, 3}},
		{"严格递增", []int{1, 2, 3, 4, 5}, []float64{1, 1.5, 2, 2.5, 3}},
		{"交替大小", []int{1, 100, 2, 99, 3}, []float64{1, 50.5, 2, 50.5, 3}},
		{"偶数个相同中位数", []int{1, 2, 2, 3}, []float64{1, 1.5, 2, 2}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := SolveBM048(c.nums)
			if !equalFloats(got, c.want) {
				t.Fatalf("nums=%v，期望 %v，实际 %v", c.nums, c.want, got)
			}
		})
	}

	// 空数据流单独验证：GetMedian 必须安全返回 0，且之后仍可正常 Insert。
	t.Run("空流后再插入", func(t *testing.T) {
		f := NewMedianFinder()
		if got := f.GetMedian(); got != 0 {
			t.Fatalf("空数据流期望 0，实际 %v", got)
		}
		f.Insert(10)
		if got := f.GetMedian(); got != 10 {
			t.Fatalf("插入 10 后期望 10，实际 %v", got)
		}
		f.Insert(20)
		if got := f.GetMedian(); got != 15 {
			t.Fatalf("再插入 20 后期望 15，实际 %v", got)
		}
	})

	// 随机对拍：与「每次排序后取中位数」的参考实现比较（固定种子，可复现）。
	t.Run("随机对拍排序解", func(t *testing.T) {
		r := rand.New(rand.NewSource(48))
		for round := 0; round < 200; round++ {
			n := r.Intn(40) + 1
			nums := make([]int, n)
			for i := range nums {
				nums[i] = r.Intn(21) // 小范围取值，覆盖大量重复
			}
			want := bruteMedians(nums)
			if got := SolveBM048(nums); !equalFloats(got, want) {
				t.Fatalf("第 %d 轮 nums=%v，期望 %v，实际 %v", round, nums, want, got)
			}
		}
	})

	// 大规模用例：验证两个堆的平衡约定在长数据流上始终成立（不会出现某个堆为空却取堆顶）。
	t.Run("大规模递减数据流", func(t *testing.T) {
		const n = 5000
		nums := make([]int, n)
		for i := range nums {
			nums[i] = n - i
		}
		got := SolveBM048(nums)
		if len(got) != n {
			t.Fatalf("期望 %d 个中位数，实际 %d", n, len(got))
		}
		for i := 0; i < n; i++ {
			cnt := i + 1
			var want float64
			if cnt%2 == 1 {
				// 递减序列前 cnt 个数的中位数就是 nums[cnt/2]
				want = float64(nums[cnt/2])
			} else {
				want = float64(nums[cnt/2-1]+nums[cnt/2]) / 2
			}
			if got[i] != want {
				t.Fatalf("第 %d 个数之后期望中位数 %v，实际 %v", cnt, want, got[i])
			}
		}
	})
}

// bruteMedians 是「每次插入后重新排序取中位数」的参考实现，仅用于测试对拍。
func bruteMedians(nums []int) []float64 {
	res := make([]float64, 0, len(nums))
	sorted := make([]int, 0, len(nums))
	for _, x := range nums {
		sorted = append(sorted, x)
		sort.Ints(sorted)
		n := len(sorted)
		if n%2 == 1 {
			res = append(res, float64(sorted[n/2]))
			continue
		}
		res = append(res, float64(sorted[n/2-1]+sorted[n/2])/2)
	}
	return res
}

// equalFloats 比较两个 float64 切片是否完全一致，nil 与空切片视为相等。
// 本题中位数由整数相加后除以 2 得到，float64 可以精确表示，因此直接比较即可。
func equalFloats(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// ---------------------------------------------------------------------------
// BM049 表达式求值 · 测试
// ---------------------------------------------------------------------------
// TestBM049 表驱动测试表达式求值。
// 输入约束：只含 + - * 与括号，不含除号，数字为非负整数。
func TestBM049(t *testing.T) {
	cases := []struct {
		name string
		expr string
		want int
	}{
		{"官方示例1", "(2*(3-4))*5", -10},
		{"官方示例2", "3+2*3*4-1", 26},
		{"官方示例3", "100+2*30-40*2", 80},
		{"单个数字", "42", 42},
		{"多位数解析", "123+456*10", 4683},
		{"乘法优先级高于加减", "1+10*2-3*4", 9},
		{"同级从左到右", "10-3-2", 5},
		{"括号嵌套", "((1+2)*(3+(4-1)))", 18},
		{"括号内加减", "(1+2)*(3-4)", -3},
		{"连续括号", "2*(3*(4-1))", 18},
		{"零与空格", "  1 + 2 * 3 ", 7},
		{"结果非零且为负", "5-5*(1+1)", -5},
		{"全为乘法", "1*2*3*4", 24},
		{"空表达式", "", 0},
		{"括号包裹整体", "(100)", 100},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := SolveBM049(tc.expr)
			if got != tc.want {
				t.Fatalf("表达式 %q：期望 %v，实际 %v", tc.expr, tc.want, got)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// BM050 两数之和 · 测试
// ---------------------------------------------------------------------------
// TestBM050 表驱动测试两数之和，返回的下标从 1 开始且升序。
func TestBM050(t *testing.T) {
	cases := []struct {
		name    string
		numbers []int
		target  int
		want    []int
	}{
		{"官方示例1", []int{3, 2, 4}, 6, []int{2, 3}},
		{"官方示例2", []int{20, 70, 110, 150}, 90, []int{1, 2}},
		{"官方示例3", []int{-3, 4, 3, 90}, 0, []int{1, 3}},
		{"最小规模两元素", []int{1, 2}, 3, []int{1, 2}},
		{"重复元素", []int{2, 2, 3}, 4, []int{1, 2}},
		{"答案在尾部", []int{1, 5, 8, 13, 21}, 34, []int{4, 5}},
		{"全负数", []int{-5, -1, -12, -2, -3}, -8, []int{1, 5}},
		{"目标为零且含零", []int{5, 0, 0, 3}, 0, []int{2, 3}},
		{"间距较大", []int{1, 3, 5, 7, 9, 11}, 18, []int{4, 6}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := SolveBM050(tc.numbers, tc.target)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("numbers=%v target=%d：期望 %v，实际 %v", tc.numbers, tc.target, tc.want, got)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// BM051 数组中出现次数超过一半的数字 · 测试
// ---------------------------------------------------------------------------
// TestBM051 表驱动测试摩尔投票；所有用例都保证「众数存在且出现次数超过一半」。
func TestBM051(t *testing.T) {
	cases := []struct {
		name    string
		numbers []int
		want    int
	}{
		{"官方示例1", []int{1, 2, 3, 2, 2, 2, 5, 4, 2}, 2},
		{"官方示例2", []int{3, 2, 3}, 3},
		{"单元素", []int{7}, 7},
		{"两元素相同", []int{5, 5}, 5},
		{"全相同", []int{9, 9, 9, 9, 9}, 9},
		{"偶数长度众数过半", []int{4, 4, 4, 1, 2, 3, 4, 4}, 4},
		{"零也合法", []int{0, 0, 0, 1, 2}, 0},
		{"大值元素", []int{10000, 1, 10000, 2, 10000}, 10000},
		{"众数在末尾", []int{1, 2, 1, 3, 1, 4, 1}, 1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := SolveBM051(tc.numbers)
			if got != tc.want {
				t.Fatalf("numbers=%v：期望 %v，实际 %v", tc.numbers, tc.want, got)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// BM052 数组中只出现一次的两个数字 · 测试
// ---------------------------------------------------------------------------
// TestBM052 表驱动测试「只出现一次的两个数字」，返回值需按非降序排列。
func TestBM052(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		want []int
	}{
		{"官方示例1", []int{1, 4, 1, 6}, []int{4, 6}},
		{"官方示例2", []int{1, 2, 3, 3, 2, 9}, []int{1, 9}},
		{"最小规模", []int{2, 3}, []int{2, 3}},
		{"输入顺序颠倒", []int{6, 1, 4, 1}, []int{4, 6}},
		{"无成对元素", []int{5, 3}, []int{3, 5}},
		{"范围外健壮性_含零", []int{0, 5, 0, 7}, []int{5, 7}},
		{"范围外健壮性_含负数", []int{-1, 2, -1, 4}, []int{2, 4}},
		{"多个成对元素", []int{9, 9, 1, 2, 3, 3, 8, 8, 5, 5}, []int{1, 2}},
		{"大值元素", []int{1000000, 999999, 1000000, 1}, []int{1, 999999}},
		{"空数组", []int{}, nil},
		{"单元素", []int{7}, nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := SolveBM052(tc.nums)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("nums=%v：期望 %v，实际 %v", tc.nums, tc.want, got)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// BM053 缺失的第一个正整数 · 测试
// ---------------------------------------------------------------------------
// TestBM053 表驱动测试缺失的第一个正整数。
// 解题函数会就地修改切片，因此每个用例都用函数重新构造一份输入。
func TestBM053(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		want int
	}{
		{"官方示例1", []int{1, 0, 2}, 3},
		{"官方示例2", []int{}, 1},
		{"官方示例3", []int{-2, 3, 4, 1}, 2},
		{"单元素一", []int{1}, 2},
		{"单元素非一", []int{2}, 1},
		{"全负数", []int{-1, -2, -3}, 1},
		{"含零", []int{0, 1, 2}, 3},
		{"完整排列", []int{3, 1, 2}, 4},
		{"乱序缺中间", []int{7, 8, 9, 11, 12, 1, 2, 3, 4, 5}, 6},
		{"只有大数", []int{100000, 200000}, 1},
		{"逆序排列", []int{5, 4, 3, 2, 1}, 6},
		{"答案在一", []int{3, 4, 5}, 1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// 复制一份输入，避免用例之间互相影响（解题函数是原地算法）。
			input := append([]int(nil), tc.nums...)
			got := SolveBM053(input)
			if got != tc.want {
				t.Fatalf("nums=%v：期望 %v，实际 %v", tc.nums, tc.want, got)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// BM054 三数之和 · 测试
// ---------------------------------------------------------------------------
// normalizeTriples 把每个三元组内部升序、三元组之间按字典序排序，
// 这样比较时就不依赖解的枚举顺序（题目只要求结果不重复、组内非递减）。
func normalizeTriples(in [][]int) [][]int {
	out := make([][]int, len(in))
	for i, triple := range in {
		cp := append([]int(nil), triple...)
		sort.Ints(cp)
		out[i] = cp
	}
	sort.Slice(out, func(i, j int) bool {
		for k := 0; k < len(out[i]) && k < len(out[j]); k++ {
			if out[i][k] != out[j][k] {
				return out[i][k] < out[j][k]
			}
		}
		return len(out[i]) < len(out[j])
	})
	return out
}

// TestBM054 表驱动测试三数之和，结果排序后再比较。
func TestBM054(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		want [][]int
	}{
		{"官方示例1", []int{-10, 0, 10, 20, -10, -40}, [][]int{{-10, -10, 20}, {-10, 0, 10}}},
		{"官方示例2", []int{1, 2, 3}, [][]int{}},
		{"官方示例3", []int{0, 0, 0}, [][]int{{0, 0, 0}}},
		{"空数组", []int{}, [][]int{}},
		{"两元素", []int{0, 0}, [][]int{}},
		{"多个零只出一组", []int{0, 0, 0, 0}, [][]int{{0, 0, 0}}},
		{"全正数", []int{1, 2, 3, 4}, [][]int{}},
		{"全负数", []int{-5, -4, -3}, [][]int{}},
		{"一正一负需第三者", []int{-1, 1}, [][]int{}},
		{"含重复但解唯一", []int{-2, 0, 1, 1, 2}, [][]int{{-2, 0, 2}, {-2, 1, 1}}},
		{"经典去重用例", []int{-1, 0, 1, 2, -1, -4}, [][]int{{-1, -1, 2}, {-1, 0, 1}}},
		{"多组解", []int{-4, -2, -1, 0, 1, 2, 3, 4}, [][]int{{-4, 0, 4}, {-4, 1, 3}, {-2, -1, 3}, {-2, 0, 2}, {-1, 0, 1}}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := SolveBM054(tc.nums)
			gotNorm := normalizeTriples(got)
			wantNorm := normalizeTriples(tc.want)
			if !reflect.DeepEqual(gotNorm, wantNorm) {
				t.Fatalf("nums=%v：期望 %v，实际 %v", tc.nums, wantNorm, gotNorm)
			}
			// 结果内部必须非递减，且不能有重复三元组。
			for _, triple := range got {
				if len(triple) != 3 {
					t.Fatalf("nums=%v：三元组长度应为 3，实际 %v", tc.nums, triple)
				}
				if triple[0] > triple[1] || triple[1] > triple[2] {
					t.Fatalf("nums=%v：三元组未按非递减排列 %v", tc.nums, triple)
				}
				if triple[0]+triple[1]+triple[2] != 0 {
					t.Fatalf("nums=%v：三元组之和不为 0：%v", tc.nums, triple)
				}
			}
		})
	}

	// 解题函数不应该修改调用方传入的切片。
	t.Run("不修改输入切片", func(t *testing.T) {
		nums := []int{3, -1, -2, 0, 1, 2}
		original := append([]int(nil), nums...)
		SolveBM054(nums)
		if !reflect.DeepEqual(nums, original) {
			t.Fatalf("输入切片被修改：期望 %v，实际 %v", original, nums)
		}
	})
}

// ---------------------------------------------------------------------------
// BM055 没有重复项数字的全排列 · 测试
// ---------------------------------------------------------------------------
// factorial 返回 n!，用于核对全排列的数量。
func factorial(n int) int {
	res := 1
	for i := 2; i <= n; i++ {
		res *= i
	}
	return res
}

// equalMatrix 按行比较两个二维整型切片是否完全相同。
func equalMatrix(a, b [][]int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

// isPermutation 判断 perm 与 src 是否由同一组数字（含重复次数）组成。
func isPermutation(perm, src []int) bool {
	if len(perm) != len(src) {
		return false
	}
	a := append([]int(nil), perm...)
	b := append([]int(nil), src...)
	sort.Ints(a)
	sort.Ints(b)
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// lessInts 判断两个切片的字典序大小。
func lessInts(a, b []int) bool {
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
	}
	return len(a) < len(b)
}

// isLexicographic 校验结果整体按字典序严格升序（题目要求的输出顺序）。
func isLexicographic(res [][]int) bool {
	for i := 1; i < len(res); i++ {
		if !lessInts(res[i-1], res[i]) {
			return false
		}
	}
	return true
}

func TestBM055(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		want [][]int
	}{
		{
			name: "官方示例三个数字",
			nums: []int{1, 2, 3},
			want: [][]int{{1, 2, 3}, {1, 3, 2}, {2, 1, 3}, {2, 3, 1}, {3, 1, 2}, {3, 2, 1}},
		},
		{
			name: "单元素",
			nums: []int{7},
			want: [][]int{{7}},
		},
		{
			name: "零单元素",
			nums: []int{0},
			want: [][]int{{0}},
		},
		{
			name: "空数组",
			nums: []int{},
			want: [][]int{},
		},
		{
			name: "输入逆序的两个元素",
			nums: []int{2, 1},
			want: [][]int{{1, 2}, {2, 1}},
		},
		{
			name: "两个元素含零",
			nums: []int{1, 0},
			want: [][]int{{0, 1}, {1, 0}},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := SolveBM055(tc.nums)
			if !equalMatrix(got, tc.want) {
				t.Fatalf("期望 %v，实际 %v", tc.want, got)
			}
		})
	}

	// 规模用例：数量必须等于 n!，且每个结果都是合法排列、互不重复、整体字典序升序。
	t.Run("四元素数量与排列合法性", func(t *testing.T) {
		nums := []int{3, 1, 4, 2}
		got := SolveBM055(nums)
		if want := factorial(len(nums)); len(got) != want {
			t.Fatalf("期望 %d 个排列，实际 %d 个", want, len(got))
		}
		seen := make(map[string]bool, len(got))
		for _, perm := range got {
			if !isPermutation(perm, nums) {
				t.Fatalf("排列 %v 不是 %v 的合法排列", perm, nums)
			}
			key := fmt.Sprint(perm)
			if seen[key] {
				t.Fatalf("排列 %v 重复出现", perm)
			}
			seen[key] = true
		}
		if !isLexicographic(got) {
			t.Fatalf("结果未按字典序输出：%v", got)
		}
	})

	t.Run("不修改入参", func(t *testing.T) {
		nums := []int{3, 1, 2}
		SolveBM055(nums)
		want := []int{3, 1, 2}
		for i := range want {
			if nums[i] != want[i] {
				t.Fatalf("入参被修改，期望 %v，实际 %v", want, nums)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// BM056 有重复项数字的全排列 · 测试
// ---------------------------------------------------------------------------
// bm056testFactorial 返回 n!，用于核对全排列的数量。
func bm056testFactorial(n int) int {
	res := 1
	for i := 2; i <= n; i++ {
		res *= i
	}
	return res
}

// countUniquePermutations 用公式 n! / (k1! * k2! * ...) 计算含重复元素的全排列数量，
// 其中 ki 是第 i 种不同数字的出现次数。
func countUniquePermutations(nums []int) int {
	n := len(nums)
	if n == 0 {
		return 0 // 与 SolveBM056 的约定一致：空数组返回空结果
	}
	res := bm056testFactorial(n)
	counts := make(map[int]int, n)
	for _, v := range nums {
		counts[v]++
	}
	for _, c := range counts {
		res /= bm056testFactorial(c)
	}
	return res
}

// bm056testEqualMatrix 按行比较两个二维整型切片是否完全相同。
func bm056testEqualMatrix(a, b [][]int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

// bm056testIsPermutation 判断 perm 与 src 是否由同一组数字（含重复次数）组成。
func bm056testIsPermutation(perm, src []int) bool {
	if len(perm) != len(src) {
		return false
	}
	a := append([]int(nil), perm...)
	b := append([]int(nil), src...)
	sort.Ints(a)
	sort.Ints(b)
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// bm056testLessInts 判断两个切片的字典序大小。
func bm056testLessInts(a, b []int) bool {
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
	}
	return len(a) < len(b)
}

// bm056testIsLexicographic 校验结果整体按字典序严格升序（题目要求的输出顺序）。
func bm056testIsLexicographic(res [][]int) bool {
	for i := 1; i < len(res); i++ {
		if !bm056testLessInts(res[i-1], res[i]) {
			return false
		}
	}
	return true
}

func TestBM056(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		want [][]int
	}{
		{
			name: "官方示例含重复元素",
			nums: []int{1, 1, 2},
			want: [][]int{{1, 1, 2}, {1, 2, 1}, {2, 1, 1}},
		},
		{
			name: "无重复元素",
			nums: []int{3, 2, 1},
			want: [][]int{{1, 2, 3}, {1, 3, 2}, {2, 1, 3}, {2, 3, 1}, {3, 1, 2}, {3, 2, 1}},
		},
		{
			name: "全相同元素",
			nums: []int{2, 2, 2},
			want: [][]int{{2, 2, 2}},
		},
		{
			name: "两个相同元素",
			nums: []int{1, 1},
			want: [][]int{{1, 1}},
		},
		{
			name: "单元素",
			nums: []int{9},
			want: [][]int{{9}},
		},
		{
			name: "空数组",
			nums: []int{},
			want: [][]int{},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := SolveBM056(tc.nums)
			if !bm056testEqualMatrix(got, tc.want) {
				t.Fatalf("期望 %v，实际 %v", tc.want, got)
			}
		})
	}

	// 用阶乘公式核对数量，并逐个校验是合法排列、互不重复、整体字典序升序。
	t.Run("阶乘公式核对数量", func(t *testing.T) {
		list := [][]int{
			{1, 1, 2, 3},    // 4!/2! = 12
			{1, 1, 2, 2},    // 4!/(2!*2!) = 6
			{2, 2, 3, 3},    // 4!/(2!*2!) = 6
			{1, 1, 1, 2},    // 4!/3! = 4
			{4, 5, 6, 7, 8}, // 5! = 120
		}
		for _, nums := range list {
			nums := nums
			t.Run(fmt.Sprint(nums), func(t *testing.T) {
				got := SolveBM056(nums)
				if want := countUniquePermutations(nums); len(got) != want {
					t.Fatalf("期望 %d 个排列，实际 %d 个", want, len(got))
				}
				seen := make(map[string]bool, len(got))
				for _, perm := range got {
					if !bm056testIsPermutation(perm, nums) {
						t.Fatalf("排列 %v 不是 %v 的合法排列", perm, nums)
					}
					key := fmt.Sprint(perm)
					if seen[key] {
						t.Fatalf("排列 %v 重复出现", perm)
					}
					seen[key] = true
				}
				if !bm056testIsLexicographic(got) {
					t.Fatalf("结果未按字典序输出：%v", got)
				}
			})
		}
	})

	t.Run("不修改入参", func(t *testing.T) {
		nums := []int{2, 1, 1}
		SolveBM056(nums)
		want := []int{2, 1, 1}
		for i := range want {
			if nums[i] != want[i] {
				t.Fatalf("入参被修改，期望 %v，实际 %v", want, nums)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// BM057 岛屿数量 · 测试
// ---------------------------------------------------------------------------
// bm057testEqualMatrix 按行比较两个二维整型切片是否完全相同。
func bm057testEqualMatrix(a, b [][]int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

func TestBM057(t *testing.T) {
	cases := []struct {
		name string
		grid [][]int
		want int
	}{
		{
			name: "官方示例5x5矩阵",
			grid: [][]int{
				{1, 1, 0, 0, 0},
				{0, 1, 0, 1, 1},
				{0, 0, 0, 1, 1},
				{0, 0, 0, 0, 0},
				{0, 0, 1, 1, 1},
			},
			want: 3,
		},
		{
			name: "官方示例单个0",
			grid: [][]int{{0}},
			want: 0,
		},
		{
			name: "1x1全是陆地",
			grid: [][]int{{1}},
			want: 1,
		},
		{
			name: "整块陆地3x3",
			grid: [][]int{
				{1, 1, 1},
				{1, 1, 1},
				{1, 1, 1},
			},
			want: 1,
		},
		{
			name: "全海洋3x3",
			grid: [][]int{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			want: 0,
		},
		{
			name: "对角相邻不算同一个岛",
			grid: [][]int{
				{1, 0},
				{0, 1},
			},
			want: 2,
		},
		{
			name: "一维单行交替陆地",
			grid: [][]int{{1, 0, 1, 0, 1}},
			want: 3,
		},
		{
			name: "空矩阵nil",
			grid: nil,
			want: 0,
		},
		{
			name: "空矩阵空切片",
			grid: [][]int{},
			want: 0,
		},
		{
			name: "隔行分割的岛",
			grid: [][]int{
				{1, 1, 1},
				{0, 0, 0},
				{1, 0, 1},
			},
			want: 3,
		},
		{
			name: "环形陆地包住海洋",
			grid: [][]int{
				{1, 1, 1},
				{1, 0, 1},
				{1, 1, 1},
			},
			want: 1,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := SolveBM057(tc.grid)
			if got != tc.want {
				t.Fatalf("期望 %d，实际 %d", tc.want, got)
			}
		})
	}

	// 洪水填充会修改矩阵，这里确认函数内部做了复制，没有污染调用方数据。
	t.Run("不修改入参矩阵", func(t *testing.T) {
		grid := [][]int{
			{1, 1, 0},
			{0, 1, 0},
			{0, 0, 1},
		}
		want := [][]int{
			{1, 1, 0},
			{0, 1, 0},
			{0, 0, 1},
		}
		if got := SolveBM057(grid); got != 2 {
			t.Fatalf("期望岛屿数 2，实际 %d", got)
		}
		if !bm057testEqualMatrix(grid, want) {
			t.Fatalf("入参矩阵被修改，期望 %v，实际 %v", want, grid)
		}
	})

	// 边界：只有一行但没有列信息。
	t.Run("行内无元素", func(t *testing.T) {
		if got := SolveBM057([][]int{{}, {}}); got != 0 {
			t.Fatalf("期望 0，实际 %d", got)
		}
	})
}

// ---------------------------------------------------------------------------
// BM058 字符串的排列 · 测试
// ---------------------------------------------------------------------------
// equalStrings 比较两个字符串切片是否完全相同。
func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// sortedRunes 返回字符串按码点升序排序后的字符切片。
func sortedRunes(s string) []rune {
	rs := []rune(s)
	sort.Slice(rs, func(i, j int) bool { return rs[i] < rs[j] })
	return rs
}

// isRearrangement 判断候选串与原串是否由同一组字符（含重复次数）组成。
func isRearrangement(candidate, src string) bool {
	a, b := sortedRunes(candidate), sortedRunes(src)
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestBM058(t *testing.T) {
	cases := []struct {
		name string
		s    string
		want []string
	}{
		{
			name: "官方示例abc",
			s:    "abc",
			want: []string{"abc", "acb", "bac", "bca", "cab", "cba"},
		},
		{
			name: "单字符",
			s:    "a",
			want: []string{"a"},
		},
		{
			name: "空字符串",
			s:    "",
			want: []string{},
		},
		{
			name: "两个相同字符",
			s:    "aa",
			want: []string{"aa"},
		},
		{
			name: "含重复字符",
			s:    "aba",
			want: []string{"aab", "aba", "baa"},
		},
		{
			name: "输入逆序的cba",
			s:    "cba",
			want: []string{"abc", "acb", "bac", "bca", "cab", "cba"},
		},
		{
			name: "大小写混合",
			s:    "aA",
			want: []string{"Aa", "aA"},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := SolveBM058(tc.s)
			if !equalStrings(got, tc.want) {
				t.Fatalf("期望 %v，实际 %v", tc.want, got)
			}
		})
	}

	// 含重复字符时用公式 n! / (k1! * k2! * ...) 核对数量，并校验合法性、去重与字典序。
	t.Run("重复字符数量与合法性", func(t *testing.T) {
		s := "aabb"
		got := SolveBM058(s)
		if want := 6; len(got) != want { // 4! / (2! * 2!) = 6
			t.Fatalf("期望 %d 个排列，实际 %d 个", want, len(got))
		}
		seen := make(map[string]bool, len(got))
		for i, candidate := range got {
			if !isRearrangement(candidate, s) {
				t.Fatalf("字符串 %q 不是 %q 的合法排列", candidate, s)
			}
			if seen[candidate] {
				t.Fatalf("字符串 %q 重复出现", candidate)
			}
			seen[candidate] = true
			if i > 0 && got[i-1] >= candidate {
				t.Fatalf("结果未按字典序升序：%q 出现在 %q 之后", candidate, got[i-1])
			}
		}
	})

	// 规模用例：8 个互不相同字符共 8! = 40320 个排列，全部合法且不重复。
	t.Run("八字符规模", func(t *testing.T) {
		s := "abcdefgh"
		got := SolveBM058(s)
		if want := 40320; len(got) != want {
			t.Fatalf("期望 %d 个排列，实际 %d 个", want, len(got))
		}
		seen := make(map[string]bool, len(got))
		for _, candidate := range got {
			if !isRearrangement(candidate, s) {
				t.Fatalf("字符串 %q 不是 %q 的合法排列", candidate, s)
			}
			if seen[candidate] {
				t.Fatalf("字符串 %q 重复出现", candidate)
			}
			seen[candidate] = true
		}
	})
}

// ---------------------------------------------------------------------------
// BM059 N皇后问题 · 测试
// ---------------------------------------------------------------------------
func TestBM059(t *testing.T) {
	cases := []struct {
		name string
		n    int
		want int
	}{
		{name: "n为0", n: 0, want: 0},
		{name: "n为1", n: 1, want: 1},
		{name: "n为2", n: 2, want: 0},
		{name: "n为3", n: 3, want: 0},
		{name: "n为4", n: 4, want: 2},
		{name: "n为5", n: 5, want: 10},
		{name: "n为6", n: 6, want: 4},
		{name: "n为7", n: 7, want: 40},
		{name: "n为8官方示例", n: 8, want: 92},
		{name: "n为9上界", n: 9, want: 352},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := SolveBM059(tc.n)
			if got != tc.want {
				t.Fatalf("n=%d 时期望 %d 种摆法，实际 %d 种", tc.n, tc.want, got)
			}
		})
	}

	// 负数是非法输入，约定返回 0，且不应崩溃。
	t.Run("n为负数", func(t *testing.T) {
		if got := SolveBM059(-1); got != 0 {
			t.Fatalf("期望 0，实际 %d", got)
		}
	})
}

// ---------------------------------------------------------------------------
// BM060 括号生成 · 测试
// ---------------------------------------------------------------------------
// isLegalParentheses 判断由 n 对括号组成的字符串是否合法：
// 长度必须为 2n，且任意前缀中左括号数不少于右括号数，最终左右括号数相等。
func isLegalParentheses(s string, n int) bool {
	if len(s) != 2*n {
		return false
	}
	balance := 0
	for _, ch := range s {
		switch ch {
		case '(':
			balance++
		case ')':
			balance--
			if balance < 0 {
				return false
			}
		default:
			return false
		}
	}
	return balance == 0
}

// sortedCopy 返回字符串切片的升序副本（用于「顺序不作要求」的情况下比较集合）。
func sortedCopy(src []string) []string {
	res := append([]string(nil), src...)
	sort.Strings(res)
	return res
}

// bm060testEqualStrings 比较两个字符串切片是否完全相同。
func bm060testEqualStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestBM060(t *testing.T) {
	cases := []struct {
		name string
		n    int
		want []string
	}{
		{name: "n为0", n: 0, want: []string{}},
		{name: "n为1", n: 1, want: []string{"()"}},
		{name: "n为2", n: 2, want: []string{"(())", "()()"}},
		{
			name: "官方示例n为3",
			n:    3,
			want: []string{"((()))", "(()())", "(())()", "()(())", "()()()"},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := SolveBM060(tc.n)
			// 题目不要求输出顺序，排序后再比较集合内容。
			if !bm060testEqualStrings(sortedCopy(got), sortedCopy(tc.want)) {
				t.Fatalf("n=%d 时期望 %v，实际 %v", tc.n, tc.want, got)
			}
			for _, s := range got {
				if !isLegalParentheses(s, tc.n) {
					t.Fatalf("组合 %q 不是合法的括号序列", s)
				}
			}
		})
	}

	// n 对括号的合法组合数就是第 n 个卡特兰数，用它核对数量（C4=14，C8=1430）。
	t.Run("卡特兰数核对数量", func(t *testing.T) {
		cases := []struct {
			n    int
			want int
		}{
			{n: 4, want: 14},
			{n: 5, want: 42},
			{n: 8, want: 1430},
		}
		for _, tc := range cases {
			tc := tc
			t.Run("n为"+itoa(tc.n), func(t *testing.T) {
				got := SolveBM060(tc.n)
				if len(got) != tc.want {
					t.Fatalf("n=%d 时期望 %d 个组合，实际 %d 个", tc.n, tc.want, len(got))
				}
				seen := make(map[string]bool, len(got))
				for _, s := range got {
					if !isLegalParentheses(s, tc.n) {
						t.Fatalf("组合 %q 不是合法的括号序列", s)
					}
					if seen[s] {
						t.Fatalf("组合 %q 重复出现", s)
					}
					seen[s] = true
				}
			})
		}
	})

	// n 为负数属于非法输入，约定返回空结果且不崩溃。
	t.Run("n为负数", func(t *testing.T) {
		if got := SolveBM060(-3); len(got) != 0 {
			t.Fatalf("期望空结果，实际 %v", got)
		}
	})
}

// itoa 把非负整数转成十进制字符串，用于构造子测试名称。
func itoa(v int) string {
	if v == 0 {
		return "0"
	}
	buf := make([]byte, 0, 4)
	for v > 0 {
		buf = append([]byte{byte('0' + v%10)}, buf...)
		v /= 10
	}
	return string(buf)
}

// ---------------------------------------------------------------------------
// BM061 矩阵最长递增路径 · 测试
// ---------------------------------------------------------------------------
func TestBM061(t *testing.T) {
	cases := []struct {
		name   string
		matrix [][]int
		want   int
	}{
		{"官方示例1_三阶递增矩阵", [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}, 5},
		{"官方示例2_四方向绕行", [][]int{{1, 2}, {4, 3}}, 4},
		{"空矩阵", [][]int{}, 0},
		{"空行的矩阵", [][]int{{}}, 0},
		{"单元素", [][]int{{0}}, 1},
		{"单行_递增", [][]int{{1, 2, 3, 4}}, 4},
		{"单列_递减", [][]int{{3}, {2}, {1}}, 3},
		{"全相同元素", [][]int{{7, 7}, {7, 7}}, 1},
		{"螺旋绕行_覆盖四个方向", [][]int{{1, 2, 3}, {8, 9, 4}, {7, 6, 5}}, 9},
		{"常规用例", [][]int{{9, 9, 4}, {6, 6, 8}, {2, 1, 1}}, 4},
		{"大值矩阵", [][]int{{1000, 0, 999}}, 2},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := SolveBM061(tc.matrix); got != tc.want {
				t.Fatalf("矩阵 %v：期望 %d，实际 %d", tc.matrix, tc.want, got)
			}
		})
	}

	t.Run("不修改输入矩阵", func(t *testing.T) {
		matrix := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
		SolveBM061(matrix)
		want := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
		for i := range want {
			for j := range want[i] {
				if matrix[i][j] != want[i][j] {
					t.Fatalf("下标 (%d,%d)：期望 %d，实际 %d", i, j, want[i][j], matrix[i][j])
				}
			}
		}
	})

	t.Run("随机小矩阵与暴力DFS一致", func(t *testing.T) {
		rnd := rand.New(rand.NewSource(20240801))
		for round := 0; round < 100; round++ {
			rows := rnd.Intn(4) + 1
			cols := rnd.Intn(4) + 1
			grid := make([][]int, rows)
			for i := range grid {
				grid[i] = make([]int, cols)
				for j := range grid[i] {
					grid[i][j] = rnd.Intn(6) // 值域小，更容易出现绕行与分叉
				}
			}
			want := bruteBest061(grid)
			if got := SolveBM061(grid); got != want {
				t.Fatalf("矩阵 %v：期望 %d，实际 %d", grid, want, got)
			}
		}
	})
}

// bruteBest061 不带记忆化的暴力 DFS，只用于小规模矩阵的对照。
func bruteBest061(grid [][]int) int {
	dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	var walk func(i, j int) int
	walk = func(i, j int) int {
		best := 1
		for _, d := range dirs {
			ni, nj := i+d[0], j+d[1]
			if ni < 0 || ni >= len(grid) || nj < 0 || nj >= len(grid[ni]) {
				continue
			}
			if grid[ni][nj] <= grid[i][j] {
				continue
			}
			if v := walk(ni, nj) + 1; v > best {
				best = v
			}
		}
		return best
	}
	best := 0
	for i := range grid {
		for j := range grid[i] {
			if v := walk(i, j); v > best {
				best = v
			}
		}
	}
	return best
}

// ---------------------------------------------------------------------------
// BM062 斐波那契数列 · 测试
// ---------------------------------------------------------------------------
func TestBM062(t *testing.T) {
	cases := []struct {
		name string
		n    int
		want int
	}{
		{"官方示例1_n等于4", 4, 3},
		{"官方示例2_n等于1", 1, 1},
		{"边界_n等于2", 2, 1},
		{"常规_n等于3", 3, 2},
		{"常规_n等于5", 5, 5},
		{"常规_n等于10", 10, 55},
		{"数据范围上界_n等于40", 40, 102334155},
		{"约定_n等于0返回0", 0, 0},
		{"取模生效_n等于100", 100, 687995182},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := SolveBM062(tc.n); got != tc.want {
				t.Fatalf("n=%d：期望 %d，实际 %d", tc.n, tc.want, got)
			}
		})
	}

	t.Run("与递推关系自洽且结果落在取模范围内", func(t *testing.T) {
		f0, f1 := 0, 1 // 参考值：f(0) = 0、f(1) = 1
		for n := 1; n <= 100; n++ {
			if n >= 2 {
				f0, f1 = f1, (f0+f1)%bm062Mod // 每步取模，避免大 n 时溢出
			}
			got := SolveBM062(n)
			if got != f1 {
				t.Fatalf("n=%d：期望 %d，实际 %d", n, f1, got)
			}
			if got < 0 || got >= bm062Mod {
				t.Fatalf("n=%d：结果 %d 不在 [0, %d) 区间内", n, got, bm062Mod)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// BM063 跳台阶 · 测试
// ---------------------------------------------------------------------------
func TestBM063(t *testing.T) {
	cases := []struct {
		name string
		n    int
		want int
	}{
		{"官方示例1_n等于2", 2, 2},
		{"官方示例2_n等于7", 7, 21},
		{"边界_n等于0", 0, 1},
		{"边界_n等于1", 1, 1},
		{"常规_n等于3", 3, 3},
		{"常规_n等于4", 4, 5},
		{"常规_n等于5", 5, 8},
		{"常规_n等于10", 10, 89},
		{"数据范围上界_n等于40", 40, 165580141},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := SolveBM063(tc.n); got != tc.want {
				t.Fatalf("n=%d：期望 %d，实际 %d", tc.n, tc.want, got)
			}
		})
	}

	t.Run("跳法数等于斐波那契数列第n+1项", func(t *testing.T) {
		fib := make([]int, 42) // fib[0..41]，f(1) = 1、f(2) = 1
		fib[1] = 1
		for i := 2; i <= 41; i++ {
			fib[i] = fib[i-1] + fib[i-2]
		}
		for n := 0; n <= 40; n++ {
			if got := SolveBM063(n); got != fib[n+1] {
				t.Fatalf("n=%d：期望 %d，实际 %d", n, fib[n+1], got)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// BM064 最小花费爬楼梯 · 测试
// ---------------------------------------------------------------------------
func TestBM064(t *testing.T) {
	cases := []struct {
		name string
		cost []int
		want int
	}{
		{"官方示例1", []int{2, 5, 20}, 5},
		{"官方示例2", []int{1, 100, 1, 1, 1, 90, 1, 1, 80, 1}, 6},
		{"常规用例_十级", []int{1, 100, 1, 1, 1, 100, 1, 1, 100, 1}, 6},
		{"空数组", []int{}, 0},
		{"单级台阶_值10", []int{10}, 10},
		{"单级台阶_值1", []int{1}, 1},
		{"两级台阶", []int{10, 15}, 10},
		{"常规_三级", []int{1, 2, 3}, 2},
		{"全相同花费", []int{5, 5, 5, 5}, 10},
		{"全部为1_一百级", makeOnes(100), 50},
		{"大花费", []int{10000, 1, 10000, 1, 10000}, 2},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := SolveBM064(tc.cost); got != tc.want {
				t.Fatalf("cost=%v：期望 %d，实际 %d", tc.cost, tc.want, got)
			}
		})
	}

	t.Run("随机用例与完整DP表一致", func(t *testing.T) {
		rnd := rand.New(rand.NewSource(20240802))
		for round := 0; round < 200; round++ {
			n := rnd.Intn(12) + 1
			cost := make([]int, n)
			for i := range cost {
				cost[i] = rnd.Intn(20)
			}
			if got, want := SolveBM064(cost), fullTableBest064(cost); got != want {
				t.Fatalf("cost=%v：期望 %d，实际 %d", cost, want, got)
			}
		}
	})
}

// makeOnes 生成 n 个 1 的花费数组。
func makeOnes(n int) []int {
	cost := make([]int, n)
	for i := range cost {
		cost[i] = 1
	}
	return cost
}

// fullTableBest064 用朴素 dp 表实现同一递推，仅用于对照（逐级保留全部状态）。
func fullTableBest064(cost []int) int {
	n := len(cost)
	if n == 0 {
		return 0
	}
	if n == 1 {
		return cost[0]
	}
	dp := make([]int, n+1) // dp[i]：到达位置 i（i == n 表示楼顶）的最小花费
	for i := 2; i <= n; i++ {
		dp[i] = dp[i-1] + cost[i-1]
		if cand := dp[i-2] + cost[i-2]; cand < dp[i] {
			dp[i] = cand
		}
	}
	return dp[n]
}

// ---------------------------------------------------------------------------
// BM065 最长公共子序列(二) · 测试
// ---------------------------------------------------------------------------
func TestBM065(t *testing.T) {
	cases := []struct {
		name   string
		s1, s2 string
		want   string
	}{
		{"官方示例1", "1A2C3D4B56", "B1D23A456A", "123456"},
		{"官方示例2_无公共子序列", "abc", "def", "-1"},
		{"官方示例3_两串相同", "abc", "abc", "abc"},
		{"空串_左", "", "abc", "-1"},
		{"空串_右", "abc", "", "-1"},
		{"两串都为空", "", "", "-1"},
		{"无公共字符", "xyz", "abc", "-1"},
		{"单字符相同", "a", "a", "a"},
		{"单字符不同", "a", "b", "-1"},
		{"一字符是另一串的子序列", "ace", "abcde", "ace"},
		{"多解取字典序最小_ab与ba", "ab", "ba", "a"},
		{"多解取字典序最小_ba与ab", "ba", "ab", "a"},
		{"多解取字典序最小_abc与cba", "abc", "cba", "a"},
		{"多解取字典序最小_acb与abc", "acb", "abc", "ab"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := SolveBM065(tc.s1, tc.s2)
			if got != tc.want {
				t.Fatalf("s1=%q s2=%q：期望 %q，实际 %q", tc.s1, tc.s2, tc.want, got)
			}
			if tc.want == "-1" {
				if n := refLCSLen(tc.s1, tc.s2); n != 0 {
					t.Fatalf("s1=%q s2=%q 存在长度 %d 的公共子序列，不该返回 -1", tc.s1, tc.s2, n)
				}
				return
			}
			// 结果必须是两串的公共子序列，且长度等于最优值
			if !isSubseq065(got, tc.s1) || !isSubseq065(got, tc.s2) {
				t.Fatalf("%q 不是 s1=%q s2=%q 的公共子序列", got, tc.s1, tc.s2)
			}
			if wantLen := refLCSLen(tc.s1, tc.s2); len(got) != wantLen {
				t.Fatalf("s1=%q s2=%q：期望长度 %d，实际 %d（%q）", tc.s1, tc.s2, wantLen, len(got), got)
			}
		})
	}

	t.Run("结果始终是合法公共子序列且长度最优", func(t *testing.T) {
		rnd := rand.New(rand.NewSource(20240803))
		for round := 0; round < 100; round++ {
			s1 := randomString065(rnd, "abcdef", rnd.Intn(10)+1)
			s2 := randomString065(rnd, "abcdef", rnd.Intn(10)+1)
			got := SolveBM065(s1, s2)
			wantLen := refLCSLen(s1, s2)
			if wantLen == 0 {
				if got != "-1" {
					t.Fatalf("s1=%q s2=%q：无公共子序列时应返回 -1，实际 %q", s1, s2, got)
				}
				continue
			}
			if !isSubseq065(got, s1) || !isSubseq065(got, s2) {
				t.Fatalf("s1=%q s2=%q：%q 不是公共子序列", s1, s2, got)
			}
			if len(got) != wantLen {
				t.Fatalf("s1=%q s2=%q：期望长度 %d，实际 %d（%q）", s1, s2, wantLen, len(got), got)
			}
		}
	})

	t.Run("与小规模暴力枚举的字典序最小解一致", func(t *testing.T) {
		rnd := rand.New(rand.NewSource(20240804))
		for round := 0; round < 150; round++ {
			s1 := randomString065(rnd, "abc", rnd.Intn(8)+1)
			s2 := randomString065(rnd, "abc", rnd.Intn(8)+1)
			want := bruteLCS065(s1, s2)
			if want == "" {
				want = "-1"
			}
			if got := SolveBM065(s1, s2); got != want {
				t.Fatalf("s1=%q s2=%q：期望 %q，实际 %q", s1, s2, want, got)
			}
		}
	})

	t.Run("较大规模仍为合法且长度最优", func(t *testing.T) {
		s1 := "1A2C3D4B56"
		s2 := "B1D23A456A"
		for i := 0; i < 60; i++ {
			s1 += "abcdefghij"[i%10 : i%10+1]
			s2 += "jihgfedcba"[i%10 : i%10+1]
		}
		got := SolveBM065(s1, s2)
		wantLen := refLCSLen(s1, s2)
		if len(got) != wantLen {
			t.Fatalf("期望长度 %d，实际 %d", wantLen, len(got))
		}
		if !isSubseq065(got, s1) || !isSubseq065(got, s2) {
			t.Fatalf("%q 不是公共子序列", got)
		}
	})
}

// isSubseq065 判断 sub 是否为 s 的子序列。
func isSubseq065(sub, s string) bool {
	i := 0
	for j := 0; j < len(s) && i < len(sub); j++ {
		if s[j] == sub[i] {
			i++
		}
	}
	return i == len(sub)
}

// refLCSLen 用朴素 dp 表计算最长公共子序列长度，仅用于校验。
func refLCSLen(s1, s2 string) int {
	n, m := len(s1), len(s2)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if s1[i-1] == s2[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else if dp[i-1][j] >= dp[i][j-1] {
				dp[i][j] = dp[i-1][j]
			} else {
				dp[i][j] = dp[i][j-1]
			}
		}
	}
	return dp[n][m]
}

// bruteLCS065 枚举 s1 的全部子序列，返回「是 s2 子序列」中最长且字典序最小的一个。
// 只用于长度 <= 8 的小规模随机对照。
func bruteLCS065(s1, s2 string) string {
	best := ""
	n := len(s1)
	buf := make([]byte, 0, n)
	for mask := 0; mask < 1<<n; mask++ {
		buf = buf[:0]
		for i := 0; i < n; i++ {
			if mask>>i&1 == 1 {
				buf = append(buf, s1[i])
			}
		}
		if len(buf) < len(best) {
			continue
		}
		if !isSubseq065(string(buf), s2) {
			continue
		}
		if len(buf) > len(best) || string(buf) < best {
			best = string(buf)
		}
	}
	return best
}

// randomString065 生成由 alphabet 中字符组成的随机串。
func randomString065(rnd *rand.Rand, alphabet string, n int) string {
	buf := make([]byte, n)
	for i := range buf {
		buf[i] = alphabet[rnd.Intn(len(alphabet))]
	}
	return string(buf)
}

// ---------------------------------------------------------------------------
// BM066 最长公共子串 · 测试
// ---------------------------------------------------------------------------
func TestBM066(t *testing.T) {
	cases := []struct {
		name   string
		s1, s2 string
		want   string
	}{
		{"官方示例1", "1AB2345CD", "12345EF", "2345"},
		{"示例2_与最长公共子序列区分", "abcde", "bce", "bc"},
		{"两串相同", "abcdef", "abcdef", "abcdef"},
		{"空串_左", "", "abc", "-1"},
		{"空串_右", "abc", "", "-1"},
		{"两串都为空", "", "", "-1"},
		{"无公共字符", "abc", "xyz", "-1"},
		{"单字符相同", "a", "a", "a"},
		{"单字符不同", "a", "b", "-1"},
		{"全相同字符_短串更长", "aaaa", "aaa", "aaa"},
		{"连续与否区分_LCS长但公共子串短", "abcd", "abdc", "ab"},
		{"唯一长公共子串", "xxabcyy", "zzabczz", "abc"},
		{"等长多解取行方向最早", "abab", "baba", "aba"},
		{"交错字符", "abcdefg", "xxdefgyy", "defg"},
		{"重复模式", "abcabcabc", "zzabcabczz", "abcabc"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := SolveBM066(tc.s1, tc.s2)
			if got != tc.want {
				t.Fatalf("s1=%q s2=%q：期望 %q，实际 %q", tc.s1, tc.s2, tc.want, got)
			}
			if got == "-1" {
				return
			}
			// 结果必须是两串的公共连续子串
			if !strings.Contains(tc.s1, got) || !strings.Contains(tc.s2, got) {
				t.Fatalf("%q 不是 s1=%q s2=%q 的公共子串", got, tc.s1, tc.s2)
			}
		})
	}

	t.Run("结果是公共子串且长度与朴素DP一致", func(t *testing.T) {
		grid := [][2]string{
			{"1AB2345CD", "12345EF"},
			{"abcdefg", "xxdefgyy"},
			{"abab", "baba"},
			{"abcd", "abdc"},
			{"abcabc", "cbacba"},
			{"", "abc"},
			{"aaa", "aaaaaaa"},
		}
		for _, p := range grid {
			s1, s2 := p[0], p[1]
			got, wantLen := SolveBM066(s1, s2), refLongestCommonSubstringLen(s1, s2)
			if wantLen == 0 {
				if got != "-1" {
					t.Fatalf("s1=%q s2=%q：无公共子串时应返回 -1，实际 %q", s1, s2, got)
				}
				continue
			}
			if len(got) != wantLen {
				t.Fatalf("s1=%q s2=%q：期望长度 %d，实际 %d（%q）", s1, s2, wantLen, len(got), got)
			}
			if !strings.Contains(s1, got) || !strings.Contains(s2, got) {
				t.Fatalf("%q 不是 s1=%q s2=%q 的公共子串", got, s1, s2)
			}
		}
	})

	t.Run("较大规模", func(t *testing.T) {
		s1 := strings.Repeat("xyz", 100) + "end"   // 303 个字符
		s2 := "start" + strings.Repeat("xyz", 100) // 305 个字符
		got := SolveBM066(s1, s2)
		want := strings.Repeat("xyz", 100)
		if got != want {
			t.Fatalf("期望长度 %d 的重复模式，实际长度 %d（%q）", len(want), len(got), got)
		}
		if !strings.Contains(s1, got) || !strings.Contains(s2, got) {
			t.Fatalf("%q 不是公共子串", got)
		}
	})
}

// refLongestCommonSubstringLen 用朴素二维 dp 计算最长公共子串长度，仅用于校验。
func refLongestCommonSubstringLen(s1, s2 string) int {
	n, m := len(s1), len(s2)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
	}
	best := 0
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if s1[i-1] == s2[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
				if dp[i][j] > best {
					best = dp[i][j]
				}
			}
		}
	}
	return best
}

// ---------------------------------------------------------------------------
// BM067 不同路径的数目(一) · 测试
// ---------------------------------------------------------------------------
// bruteForceUniquePaths 朴素实现：从 (0,0) 出发枚举向右/向下，作为小规模数据的对照。
func bruteForceUniquePaths(m, n int) int {
	if m <= 0 || n <= 0 {
		return 0
	}
	var walk func(i, j int) int
	walk = func(i, j int) int {
		if i == m-1 && j == n-1 {
			return 1
		}
		total := 0
		if i+1 < m {
			total += walk(i+1, j)
		}
		if j+1 < n {
			total += walk(i, j+1)
		}
		return total
	}
	return walk(0, 0)
}

func TestBM067(t *testing.T) {
	cases := []struct {
		name string
		m    int
		n    int
		want int
	}{
		{"官方示例1_两行一列", 2, 1, 1},
		{"官方示例2_两行两列", 2, 2, 2},
		{"最小网格1x1", 1, 1, 1},
		{"单行1x7", 1, 7, 1},
		{"单列7x1", 7, 1, 1},
		{"常规3x7", 3, 7, 28},
		{"常规10x10", 10, 10, 48620},
		{"非法输入_m为0", 0, 5, 0},
		{"非法输入_n为0", 5, 0, 0},
		{"非法输入_负数", -1, 3, 0},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			if got := SolveBM067(c.m, c.n); got != c.want {
				t.Fatalf("SolveBM067(%d,%d) 期望 %v，实际 %v", c.m, c.n, c.want, got)
			}
		})
	}

	t.Run("随机数据对照朴素递归", func(t *testing.T) {
		r := rand.New(rand.NewSource(20240521))
		for k := 0; k < 20; k++ {
			m := r.Intn(7) + 1
			n := r.Intn(7) + 1
			want := bruteForceUniquePaths(m, n)
			if got := SolveBM067(m, n); got != want {
				t.Fatalf("SolveBM067(%d,%d) 期望 %v，实际 %v", m, n, want, got)
			}
		}
	})

	t.Run("行列交换结果不变", func(t *testing.T) {
		for k := 1; k <= 20; k++ {
			if a, b := SolveBM067(k, 13), SolveBM067(13, k); a != b {
				t.Fatalf("SolveBM067(%d,13)=%v 与 SolveBM067(13,%d)=%v 应相等", k, a, k, b)
			}
		}
	})

	t.Run("较大规模_30x30_对照组合数公式", func(t *testing.T) {
		// 用 math/big 独立算出组合数 C(m+n-2, m-1)，确认一维滚动实现与组合数公式一致。
		// 取 30x30 是因为 C(58,29) 约 3e16，仍在 int64 范围内（题目保证答案在 32 位 int 范围内）。
		const m, n = 30, 30
		want := new(big.Int).Binomial(int64(m+n-2), int64(m-1))
		got := big.NewInt(int64(SolveBM067(m, n)))
		if got.Cmp(want) != 0 {
			t.Fatalf("SolveBM067(%d,%d) 期望 %v，实际 %v", m, n, want, got)
		}
	})
}

// ---------------------------------------------------------------------------
// BM068 矩阵的最小路径和 · 测试
// ---------------------------------------------------------------------------
// bruteForceMinPathSum 朴素实现：枚举所有只向右/向下的路径，取最小和。
func bruteForceMinPathSum(matrix [][]int) int {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return 0
	}
	n, m := len(matrix), len(matrix[0])
	var walk func(i, j int) int
	walk = func(i, j int) int {
		if i == n-1 && j == m-1 {
			return matrix[i][j]
		}
		best := -1
		if i+1 < n {
			best = walk(i+1, j)
		}
		if j+1 < m {
			if v := walk(i, j+1); best < 0 || v < best {
				best = v
			}
		}
		return matrix[i][j] + best
	}
	return walk(0, 0)
}

func TestBM068(t *testing.T) {
	cases := []struct {
		name   string
		matrix [][]int
		want   int
	}{
		{
			name:   "官方示例1_4x4",
			matrix: [][]int{{1, 3, 5, 9}, {8, 1, 3, 4}, {5, 0, 6, 1}, {8, 8, 4, 0}},
			want:   12,
		},
		{
			name:   "官方示例2_2x3",
			matrix: [][]int{{1, 2, 3}, {1, 2, 3}},
			want:   7,
		},
		{
			name:   "1x1矩阵",
			matrix: [][]int{{7}},
			want:   7,
		},
		{
			name:   "单行1x4",
			matrix: [][]int{{1, 2, 3, 4}},
			want:   10,
		},
		{
			name:   "单列4x1",
			matrix: [][]int{{1}, {2}, {3}, {4}},
			want:   10,
		},
		{
			name:   "全零矩阵",
			matrix: [][]int{{0, 0, 0}, {0, 0, 0}},
			want:   0,
		},
		{
			name:   "必须下绕避开大数",
			matrix: [][]int{{1, 2, 100, 4}, {1, 100, 100, 4}, {1, 1, 1, 4}},
			want:   9,
		},
		{
			name:   "空矩阵",
			matrix: [][]int{},
			want:   0,
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			if got := SolveBM068(c.matrix); got != c.want {
				t.Fatalf("SolveBM068(%v) 期望 %v，实际 %v", c.matrix, c.want, got)
			}
		})
	}

	t.Run("随机数据对照朴素枚举", func(t *testing.T) {
		r := rand.New(rand.NewSource(20240522))
		for k := 0; k < 30; k++ {
			n := r.Intn(5) + 1
			m := r.Intn(5) + 1
			matrix := make([][]int, n)
			for i := 0; i < n; i++ {
				matrix[i] = make([]int, m)
				for j := 0; j < m; j++ {
					matrix[i][j] = r.Intn(101)
				}
			}
			want := bruteForceMinPathSum(matrix)
			if got := SolveBM068(matrix); got != want {
				t.Fatalf("SolveBM068(%v) 期望 %v，实际 %v", matrix, want, got)
			}
		}
	})

	t.Run("较大规模不越界_200x300", func(t *testing.T) {
		matrix := make([][]int, 200)
		for i := range matrix {
			matrix[i] = make([]int, 300)
			for j := range matrix[i] {
				matrix[i][j] = 1
			}
		}
		// 全 1 时任意路径都经过 200+300-1 个格子
		const want = 200 + 300 - 1
		if got := SolveBM068(matrix); got != want {
			t.Fatalf("SolveBM068(200x300 全1) 期望 %v，实际 %v", want, got)
		}
	})
}

// ---------------------------------------------------------------------------
// BM069 把数字翻译成字符串 · 测试
// ---------------------------------------------------------------------------
// bruteForceTranslateNum 朴素实现：暴力枚举「切一个数字」还是「切两个数字」的所有切法。
func bruteForceTranslateNum(s string) int {
	if len(s) == 0 {
		return 1
	}
	total := 0
	if s[0] != '0' {
		total += bruteForceTranslateNum(s[1:])
	}
	if len(s) >= 2 && s[0] != '0' {
		two := int(s[0]-'0')*10 + int(s[1]-'0')
		if two >= 10 && two <= 26 {
			total += bruteForceTranslateNum(s[2:])
		}
	}
	return total % modBM069
}

func TestBM069(t *testing.T) {
	cases := []struct {
		name string
		s    string
		want int
	}{
		{"官方示例1_12", "12", 2},
		{"官方示例2_长数字串", "31717126241541717", 192},
		{"单个数字1", "1", 1},
		{"单个数字0", "0", 0},
		{"10只有j", "10", 1},
		{"100无法译码", "100", 0},
		{"00无法译码", "00", 0},
		{"30无法译码", "30", 0},
		{"26两位都可", "26", 2},
		{"27只能拆开", "27", 1},
		{"120", "120", 1},
		{"1111", "1111", 5},
		{"含0的常规串", "2101", 1},
		{"全0串", "000", 0},
		{"空串", "", 0},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			if got := SolveBM069(c.s); got != c.want {
				t.Fatalf("SolveBM069(%q) 期望 %v，实际 %v", c.s, c.want, got)
			}
		})
	}

	t.Run("斐波那契型_全1串长度20", func(t *testing.T) {
		// 全 '1' 串的方案数满足 f(n) = f(n-1) + f(n-2)
		f := []int{1, 1}
		for i := 2; i <= 20; i++ {
			f = append(f, f[i-1]+f[i-2])
		}
		s := strings.Repeat("1", 20)
		if got, want := SolveBM069(s), f[20]%modBM069; got != want {
			t.Fatalf("SolveBM069(%q) 期望 %v，实际 %v", s, want, got)
		}
	})

	t.Run("随机数据对照暴力枚举", func(t *testing.T) {
		r := rand.New(rand.NewSource(20240523))
		for k := 0; k < 50; k++ {
			n := r.Intn(10) + 1
			var sb strings.Builder
			for i := 0; i < n; i++ {
				// 刻意放大 0 的出现概率，覆盖含 0 的非法情况
				if r.Intn(4) == 0 {
					sb.WriteByte('0')
				} else {
					sb.WriteByte(byte('0' + r.Intn(10)))
				}
			}
			s := sb.String()
			want := bruteForceTranslateNum(s)
			if got := SolveBM069(s); got != want {
				t.Fatalf("SolveBM069(%q) 期望 %v，实际 %v", s, want, got)
			}
		}
	})

	t.Run("长度90取模不溢出", func(t *testing.T) {
		s := strings.Repeat("1", 90)
		got := SolveBM069(s)
		if got < 0 || got >= modBM069 {
			t.Fatalf("SolveBM069(90个1) 结果 %v 未落在 [0, 1e9+7) 区间", got)
		}
	})
}

// ---------------------------------------------------------------------------
// BM070 兑换零钱(一) · 测试
// ---------------------------------------------------------------------------
// bfsMinCoins 朴素实现：把金额看成图上的点，用 BFS 求从 0 到 aim 的最少步数（每次加一张面值）。
func bfsMinCoins(arr []int, aim int) int {
	if aim < 0 {
		return -1
	}
	dist := make([]int, aim+1)
	for i := range dist {
		dist[i] = -1
	}
	dist[0] = 0
	queue := []int{0}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur == aim {
			return dist[cur]
		}
		for _, coin := range arr {
			if coin <= 0 || cur+coin > aim || dist[cur+coin] != -1 {
				continue
			}
			dist[cur+coin] = dist[cur] + 1
			queue = append(queue, cur+coin)
		}
	}
	return -1
}

func TestBM070(t *testing.T) {
	cases := []struct {
		name string
		arr  []int
		aim  int
		want int
	}{
		{"官方示例1_20元", []int{5, 2, 3}, 20, 4},
		{"官方示例2_金额为0", []int{5, 2, 3}, 0, 0},
		{"官方示例3_无解", []int{3, 5}, 2, -1},
		{"单面值恰好整除", []int{7}, 21, 3},
		{"单面值无解", []int{7}, 20, -1},
		{"只有1元", []int{1}, 100, 100},
		{"空数组有目标", []int{}, 10, -1},
		{"空数组金额0", []int{}, 0, 0},
		{"大面值比目标还大", []int{100, 200}, 50, -1},
		{"经典贪心反例", []int{1, 5, 6, 9}, 11, 2},
		{"面值乱序", []int{9, 6, 5, 1}, 11, 2},
		{"负数金额", []int{1, 2}, -3, -1},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			if got := SolveBM070(c.arr, c.aim); got != c.want {
				t.Fatalf("SolveBM070(%v,%d) 期望 %v，实际 %v", c.arr, c.aim, c.want, got)
			}
		})
	}

	t.Run("随机数据对照BFS", func(t *testing.T) {
		r := rand.New(rand.NewSource(20240524))
		for k := 0; k < 50; k++ {
			n := r.Intn(4) + 1
			seen := make(map[int]bool)
			arr := make([]int, 0, n)
			for len(arr) < n {
				v := r.Intn(15) + 1
				if !seen[v] {
					seen[v] = true
					arr = append(arr, v)
				}
			}
			aim := r.Intn(60)
			want := bfsMinCoins(arr, aim)
			if got := SolveBM070(arr, aim); got != want {
				t.Fatalf("SolveBM070(%v,%d) 期望 %v，实际 %v", arr, aim, want, got)
			}
		}
	})

	t.Run("最大规模_面值含1可保证有解", func(t *testing.T) {
		arr := make([]int, 0, 100)
		for v := 1; v <= 100; v++ {
			arr = append(arr, v)
		}
		// 面值里有 1，aim=5000 一定有解，且答案不超过 5000
		got := SolveBM070(arr, 5000)
		if got <= 0 || got > 5000 {
			t.Fatalf("SolveBM070(1..100,5000) 期望 (0,5000] 区间内的解，实际 %v", got)
		}
	})
}

// ---------------------------------------------------------------------------
// BM071 最长上升子序列(一) · 测试
// ---------------------------------------------------------------------------
// dpLisNSquare 朴素实现：O(n^2) 的经典 DP，用于对照。
func dpLisNSquare(arr []int) int {
	if len(arr) == 0 {
		return 0
	}
	dp := make([]int, len(arr))
	best := 1
	for i := range arr {
		dp[i] = 1
		for j := 0; j < i; j++ {
			if arr[j] < arr[i] && dp[j]+1 > dp[i] {
				dp[i] = dp[j] + 1
			}
		}
		if dp[i] > best {
			best = dp[i]
		}
	}
	return best
}

// bruteForceLis 更暴力的对照实现：枚举所有子集，校验是否严格递增。
func bruteForceLis(arr []int) int {
	best := 0
	for mask := 0; mask < 1<<uint(len(arr)); mask++ {
		length := 0
		prev := 0
		ok := true
		for i := range arr {
			if mask&(1<<uint(i)) == 0 {
				continue
			}
			if length > 0 && arr[i] <= prev {
				ok = false
				break
			}
			prev = arr[i]
			length++
		}
		if ok && length > best {
			best = length
		}
	}
	return best
}

func TestBM071(t *testing.T) {
	cases := []struct {
		name string
		arr  []int
		want int
	}{
		{"官方示例", []int{6, 3, 1, 5, 2, 3, 7}, 4},
		{"空数组", []int{}, 0},
		{"单元素", []int{42}, 1},
		{"严格递增", []int{1, 2, 3, 4, 5}, 5},
		{"严格递减", []int{5, 4, 3, 2, 1}, 1},
		{"全相等_严格上升", []int{2, 2, 2, 2}, 1},
		{"含重复值的上升", []int{1, 2, 2, 3}, 3},
		{"全负数组", []int{-5, -4, -3, -2}, 4},
		{"交替升降", []int{7, 8, 9, 1, 2, 3}, 3},
		{"题面中的例子", []int{1, 5, 3, 7, 3}, 3},
		{"零与正数混合", []int{0, 8, 4, 12, 2, 10, 6, 14, 1, 9, 5, 13, 3, 11, 7, 15}, 6},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			if got := SolveBM071(c.arr); got != c.want {
				t.Fatalf("SolveBM071(%v) 期望 %v，实际 %v", c.arr, c.want, got)
			}
		})
	}

	t.Run("随机数据对照O_n2_DP", func(t *testing.T) {
		r := rand.New(rand.NewSource(20240525))
		for k := 0; k < 50; k++ {
			n := r.Intn(30) + 1
			arr := make([]int, n)
			for i := range arr {
				arr[i] = r.Intn(21) - 10 // 值域小，专门制造重复元素
			}
			want := dpLisNSquare(arr)
			if got := SolveBM071(arr); got != want {
				t.Fatalf("SolveBM071(%v) 期望 %v，实际 %v", arr, want, got)
			}
		}
	})

	t.Run("小规模对照暴力枚举子集", func(t *testing.T) {
		r := rand.New(rand.NewSource(20240526))
		for k := 0; k < 30; k++ {
			n := r.Intn(9) + 1
			arr := make([]int, n)
			for i := range arr {
				arr[i] = r.Intn(9)
			}
			want := bruteForceLis(arr)
			if got := SolveBM071(arr); got != want {
				t.Fatalf("SolveBM071(%v) 期望 %v，实际 %v", arr, want, got)
			}
		}
	})

	t.Run("大规模严格递增_50000", func(t *testing.T) {
		const n = 50000
		arr := make([]int, n)
		for i := range arr {
			arr[i] = i
		}
		if got := SolveBM071(arr); got != n {
			t.Fatalf("SolveBM071(0..%d) 期望 %v，实际 %v", n-1, n, got)
		}
	})
}

// ---------------------------------------------------------------------------
// BM072 连续子数组的最大和 · 测试
// ---------------------------------------------------------------------------
// bruteForceMaxSubArray 朴素实现：枚举左右端点，O(n^2) 求所有子数组和的最大值。
func bruteForceMaxSubArray(array []int) int {
	if len(array) == 0 {
		return 0
	}
	best := array[0]
	for i := 0; i < len(array); i++ {
		sum := 0
		for j := i; j < len(array); j++ {
			sum += array[j]
			if sum > best {
				best = sum
			}
		}
	}
	return best
}

func TestBM072(t *testing.T) {
	cases := []struct {
		name  string
		array []int
		want  int
	}{
		{"官方示例1", []int{1, -2, 3, 10, -4, 7, 2, -5}, 18},
		{"官方示例2_单元素", []int{2}, 2},
		{"官方示例3_单个负数", []int{-10}, -10},
		{"全负数组", []int{-3, -1, -7, -2}, -1},
		{"全正数组", []int{1, 2, 3, 4}, 10},
		{"全零数组", []int{0, 0, 0}, 0},
		{"首尾为负但中间最大", []int{-2, 4, -1, 3, -5}, 6},
		{"最大和是单个元素", []int{-5, 100, -200, 1}, 100},
		{"极值组合", []int{100, -100, 100}, 100},
		{"空数组", []int{}, 0},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			if got := SolveBM072(c.array); got != c.want {
				t.Fatalf("SolveBM072(%v) 期望 %v，实际 %v", c.array, c.want, got)
			}
		})
	}

	t.Run("随机数据对照朴素枚举", func(t *testing.T) {
		r := rand.New(rand.NewSource(20240527))
		for k := 0; k < 100; k++ {
			n := r.Intn(40) + 1
			array := make([]int, n)
			for i := range array {
				array[i] = r.Intn(201) - 100 // 覆盖 [-100,100]
			}
			want := bruteForceMaxSubArray(array)
			if got := SolveBM072(array); got != want {
				t.Fatalf("SolveBM072(%v) 期望 %v，实际 %v", array, want, got)
			}
		}
	})

	t.Run("大规模全正数组_200000", func(t *testing.T) {
		const n = 200000
		array := make([]int, n)
		for i := range array {
			array[i] = 1
		}
		if got := SolveBM072(array); got != n {
			t.Fatalf("SolveBM072(%d 个 1) 期望 %v，实际 %v", n, n, got)
		}
	})
}

// ---------------------------------------------------------------------------
// BM073 最长回文子串 · 测试
// ---------------------------------------------------------------------------
// isPalindrome 判断 s 是否为回文串（仅测试使用）。
func isPalindrome(s string) bool {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		if s[i] != s[j] {
			return false
		}
	}
	return true
}

// leftmostLongestPalindrome 暴力枚举：从长到短、从前往后找第一个回文子串，
// 即"最长且最靠左"的答案，作为 SolveBM073 的参照实现（仅测试使用）。
func leftmostLongestPalindrome(s string) string {
	for length := len(s); length >= 1; length-- {
		for i := 0; i+length <= len(s); i++ {
			if isPalindrome(s[i : i+length]) {
				return s[i : i+length]
			}
		}
	}
	return ""
}

func TestBM073(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"官方示例-aba与bab并列取最靠左", "ababc", "aba"},
		{"官方示例-abbba", "abbba", "abbba"},
		{"官方示例-单字符", "b", "b"},
		{"空串", "", ""},
		{"双字符相同", "aa", "aa"},
		{"双字符不同", "ab", "a"},
		{"全相同字符", "aaaa", "aaaa"},
		{"无长度大于1的回文", "abcde", "a"},
		{"常规用例-cbbd", "cbbd", "bb"},
		{"常规用例-babad", "babad", "bab"},
		{"常规用例-长回文横跨中轴", "forgeeksskeegfor", "geeksskeeg"},
		{"常规用例-含数字与大小写", "a1b2b1aZ", "a1b2b1a"},
		{"全部不同字符", "abcdefghij", "a"},
		{"长串-1000个相同字符", strings.Repeat("x", 1000), strings.Repeat("x", 1000)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SolveBM073(tt.in)
			if got != tt.want {
				t.Fatalf("输入 %q 期望 %q，实际 %q", tt.in, tt.want, got)
			}
			if !isPalindrome(got) || !strings.Contains(tt.in, got) {
				t.Fatalf("输入 %q 的结果 %q 不是回文或不是子串", tt.in, got)
			}
		})
	}

	// 随机对拍：与暴力枚举出的"最长且最靠左"的答案逐个比较，同时验证长度 >= 1。
	t.Run("随机对拍-非空串结果长度至少为1", func(t *testing.T) {
		r := rand.New(rand.NewSource(20240708))
		for iter := 0; iter < 300; iter++ {
			n := 1 + r.Intn(10)
			var sb strings.Builder
			for i := 0; i < n; i++ {
				sb.WriteByte("abc"[r.Intn(3)])
			}
			s := sb.String()
			want := leftmostLongestPalindrome(s)
			got := SolveBM073(s)
			if got != want {
				t.Fatalf("随机串 %q 期望 %q，实际 %q", s, want, got)
			}
			if len(got) == 0 {
				t.Fatalf("随机串 %q 的答案长度不应为 0", s)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// BM074 数字字符串转化成IP地址 · 测试
// ---------------------------------------------------------------------------
// checkIPv4 校验 ip 是否由数字串 digits 按顺序切分而成，且 4 段都合法（仅测试使用）。
func checkIPv4(ip, digits string) bool {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return false
	}
	for _, seg := range parts {
		if len(seg) == 0 || len(seg) > 3 {
			return false
		}
		for i := 0; i < len(seg); i++ {
			if seg[i] < '0' || seg[i] > '9' {
				return false
			}
		}
		if len(seg) > 1 && seg[0] == '0' {
			return false // 多余的（非单独 "0" 的）前导零不合法
		}
		if len(seg) == 3 && seg > "255" {
			return false // 三位数超过 255 不合法
		}
	}
	return strings.Join(parts, "") == digits
}

func TestBM074(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []string
	}{
		{"官方示例-25525522135", "25525522135", []string{"255.255.22.135", "255.255.221.35"}},
		{"官方示例-1111", "1111", []string{"1.1.1.1"}},
		{"官方示例-000256", "000256", nil},
		{"空串", "", nil},
		{"长度不足4", "111", nil},
		{"长度超过12", "1111111111111", nil},
		{"前导零-010010", "010010", []string{"0.10.0.10", "0.100.1.0"}},
		{"前导零-101023", "101023", []string{"1.0.10.23", "1.0.102.3", "10.1.0.23", "10.10.2.3", "101.0.2.3"}},
		{"全零-0000", "0000", []string{"0.0.0.0"}},
		{"前导零-0101", "0101", []string{"0.1.0.1"}},
		{"每段取满255", "255255255255", []string{"255.255.255.255"}},
		{"12位只能切成4段三位数", "111111111111", []string{"111.111.111.111"}},
		{"存在超255的段", "255255255256", nil},
		{"每段都超过255", "999999999999", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SolveBM074(tt.in)
			if !slices.Equal(got, tt.want) {
				t.Fatalf("输入 %q 期望 %v，实际 %v", tt.in, tt.want, got)
			}
			// 逐项校验：每个结果都由原串按顺序切分且每段合法；结果按字典序且无重复
			for i, ip := range got {
				if !checkIPv4(ip, tt.in) {
					t.Fatalf("输入 %q 的结果 %q 不是由原串切分出的合法 IP", tt.in, ip)
				}
				if i > 0 && got[i-1] >= ip {
					t.Fatalf("输入 %q 的结果 %v 未按字典序严格递增", tt.in, got)
				}
			}
		})
	}

	// 6 个字符切成 4 段（每段 1~3 位）共 10 种：段长 (3,1,1,1) 的 4 种排列 + (2,2,1,1) 的 6 种排列，
	// 每段都不超过 255，应全部保留，且天然按字典序排列。
	t.Run("全部合法切分-6个1应得10种", func(t *testing.T) {
		const in = "111111"
		got := SolveBM074(in)
		if len(got) != 10 {
			t.Fatalf("输入 %q 期望 10 种方案，实际 %d 种：%v", in, len(got), got)
		}
		if !sort.StringsAreSorted(got) {
			t.Fatalf("输入 %q 的结果没有按字典序排列：%v", in, got)
		}
		seen := make(map[string]bool, len(got))
		for _, ip := range got {
			if !checkIPv4(ip, in) {
				t.Fatalf("结果 %q 不是合法 IP 切分", ip)
			}
			if seen[ip] {
				t.Fatalf("结果 %q 重复出现", ip)
			}
			seen[ip] = true
		}
	})
}

// ---------------------------------------------------------------------------
// BM075 编辑距离(一) · 测试
// ---------------------------------------------------------------------------
// editDistance2D 用完整二维 DP 实现的参照版本（仅测试使用）。
func editDistance2D(a, b string) int {
	n, m := len(a), len(b)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
		dp[i][0] = i // 目标为空串，只能全删
	}
	for j := 0; j <= m; j++ {
		dp[0][j] = j // 源串为空串，只能全插
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1]
				continue
			}
			best := dp[i-1][j]
			if dp[i][j-1] < best {
				best = dp[i][j-1]
			}
			if dp[i-1][j-1] < best {
				best = dp[i-1][j-1]
			}
			dp[i][j] = best + 1
		}
	}
	return dp[n][m]
}

func TestBM075(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want int
	}{
		{"官方示例-nowcoder与new", "nowcoder", "new", 6},
		{"官方示例-intention与execution", "intention", "execution", 5},
		{"官方示例-now与nowcoder", "now", "nowcoder", 5},
		{"完全相同", "abc", "abc", 0},
		{"都是空串", "", "", 0},
		{"源串为空", "", "abc", 3},
		{"目标为空", "abc", "", 3},
		{"单字符不同", "a", "b", 1},
		{"单字符相同", "a", "a", 0},
		{"常规-ab与ba", "ab", "ba", 2},
		{"常规-kitten与sitting", "kitten", "sitting", 3},
		{"常规-只差一个字符", "abc", "abd", 1},
		{"常规-纯插入", "abc", "abcdef", 3},
		{"常规-纯删除", "abcdef", "abc", 3},
		{"常规-无公共字符", "abc", "xyz", 3},
		{"长串-相同", strings.Repeat("a", 1000), strings.Repeat("a", 1000), 0},
		{"长串-长度差1", strings.Repeat("a", 1000), strings.Repeat("a", 999), 1},
		{"长串-完全不同", strings.Repeat("a", 1000), strings.Repeat("b", 1000), 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SolveBM075(tt.a, tt.b)
			if got != tt.want {
				t.Fatalf("%q -> %q 期望 %d，实际 %d", tt.a, tt.b, tt.want, got)
			}
			// 顺带用二维参照实现交叉验证，避免只有一处实现对
			if ref := editDistance2D(tt.a, tt.b); ref != got {
				t.Fatalf("%q -> %q 二维参照实现为 %d，实际 %d", tt.a, tt.b, ref, got)
			}
		})
	}

	// 随机对拍：小写字母 a-c、长度 0~8，与二维参照实现比较。
	t.Run("随机对拍", func(t *testing.T) {
		r := rand.New(rand.NewSource(20240709))
		for iter := 0; iter < 300; iter++ {
			la, lb := r.Intn(9), r.Intn(9)
			a := randomString(r, la)
			b := randomString(r, lb)
			want := editDistance2D(a, b)
			got := SolveBM075(a, b)
			if got != want {
				t.Fatalf("%q -> %q 期望 %d，实际 %d", a, b, want, got)
			}
			// 编辑距离满足对称性
			if back := SolveBM075(b, a); back != got {
				t.Fatalf("%q -> %q 为 %d，反向为 %d，应对称", a, b, got, back)
			}
		}
	})
}

// randomString 生成只含 a-c、长度为 n 的随机串（仅测试使用）。
func randomString(r *rand.Rand, n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteByte("abc"[r.Intn(3)])
	}
	return sb.String()
}

// ---------------------------------------------------------------------------
// BM076 正则表达式匹配 · 测试
// ---------------------------------------------------------------------------
// regexBrute 指数级递归的参照实现，只用于小规模随机对拍（仅测试使用）。
func regexBrute(s, p string) bool {
	if p == "" {
		return s == ""
	}
	first := s != "" && (p[0] == '.' || p[0] == s[0])
	if len(p) >= 2 && p[1] == '*' {
		// '*' 匹配 0 次：丢掉 "x*"；或匹配 1 次后继续用同一个 "x*" 去匹配剩余字符
		return regexBrute(s, p[2:]) || (first && regexBrute(s[1:], p))
	}
	return first && regexBrute(s[1:], p[1:])
}

func TestBM076(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		pattern string
		want    bool
	}{
		{"官方示例-aaa与a*a", "aaa", "a*a", true},
		{"官方示例-aad与c*a*d", "aad", "c*a*d", true},
		{"官方示例-a与.*", "a", ".*", true},
		{"官方示例-aaab与a*a*a*c", "aaab", "a*a*a*c", false},
		{"空串与空模式", "", "", true},
		{"空串与.*", "", ".*", true},
		{"空串与a*", "", "a*", true},
		{"空串与a*b*", "", "a*b*", true},
		{"空串与a", "", "a", false},
		{"经典题面-aaa与a.a", "aaa", "a.a", true},
		{"经典题面-aaa与ab*ac*a", "aaa", "ab*ac*a", true},
		{"经典题面-aaa与aa.a", "aaa", "aa.a", false},
		{"经典题面-aaa与ab*a", "aaa", "ab*a", false},
		{"aaa与a*", "aaa", "a*", true},
		{"aab与c*a*b", "aab", "c*a*b", true},
		{"abc与.*", "abc", ".*", true},
		{"ab与.*c", "ab", ".*c", false},
		{"a与ab*", "a", "ab*", true},
		{"aaaa与a*aa", "aaaa", "a*aa", true},
		{"abc与abc", "abc", "abc", true},
		{"abc与abd", "abc", "abd", false},
		{"abc与ab", "abc", "ab", false},
		{"abc与abc.", "abc", "abc.", false},
		{"mississippi与mis*is*p*.", "mississippi", "mis*is*p*.", false},
		{"mississippi与mis*is*ip*.", "mississippi", "mis*is*ip*.", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SolveBM076(tt.str, tt.pattern)
			if got != tt.want {
				t.Fatalf("str=%q pattern=%q 期望 %v，实际 %v", tt.str, tt.pattern, tt.want, got)
			}
			// 用指数级递归实现交叉验证
			if ref := regexBrute(tt.str, tt.pattern); ref != got {
				t.Fatalf("str=%q pattern=%q 参照实现为 %v，实际 %v", tt.str, tt.pattern, ref, got)
			}
		})
	}

	// 随机对拍：字符串取 a-c、模式由合法 token 拼出（保证 '*' 前面一定有字符）。
	t.Run("随机对拍", func(t *testing.T) {
		r := rand.New(rand.NewSource(20240710))
		tokens := []string{"a", "b", "c", ".", "a*", "b*", "c*", ".*"}
		for iter := 0; iter < 500; iter++ {
			var s strings.Builder
			for i, sn := 0, r.Intn(6); i < sn; i++ {
				s.WriteByte("abc"[r.Intn(3)])
			}
			var p strings.Builder
			for i, pn := 0, 1+r.Intn(4); i < pn; i++ {
				p.WriteString(tokens[r.Intn(len(tokens))])
			}
			str, pattern := s.String(), p.String()
			want := regexBrute(str, pattern)
			got := SolveBM076(str, pattern)
			if got != want {
				t.Fatalf("str=%q pattern=%q 期望 %v，实际 %v", str, pattern, want, got)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// BM077 最长的括号子串 · 测试
// ---------------------------------------------------------------------------
// longestValidBrute 暴力枚举所有子串，O(n^2)（仅测试使用）。
func longestValidBrute(s string) int {
	best := 0
	for i := 0; i < len(s); i++ {
		balance := 0
		for j := i; j < len(s); j++ {
			if s[j] == '(' {
				balance++
			} else {
				balance--
			}
			if balance < 0 {
				break // 前缀已经非法，更长的子串也不可能合法
			}
			if balance == 0 && j-i+1 > best {
				best = j - i + 1
			}
		}
	}
	return best
}

func TestBM077(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"官方示例-(()", "(()", 2},
		{"官方示例-(())", "(())", 4},
		{"官方示例-)()())", ")()())", 4},
		{"空串", "", 0},
		{"单个左括号-全部非法", "(", 0},
		{"单个右括号-全部非法", ")", 0},
		{"全部是左括号", "(((((", 0},
		{"全部是右括号", ")))))", 0},
		{"只有一对", "()", 2},
		{"连续两对", "()()", 4},
		{"嵌套两对", "(())", 4},
		{"括号全部非法-中缀", ")(((", 0},
		{"常规-前半非法", ")()())", 4},
		{"常规-尾部多左括号", "()(()", 2},
		{"常规-前缀拼接", "()(())", 6},
		{"常规-完整嵌套", "(()())", 6},
		{"常规-多余右括号在中间", "())()", 2},
		{"常规-跨越多个段", "()(())()", 8},
		{"长串-五万对括号", strings.Repeat("()", 50000), 100000},
		{"长串-十万个左括号", strings.Repeat("(", 100000), 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SolveBM077(tt.in)
			if got != tt.want {
				t.Fatalf("输入 %q 期望 %d，实际 %d", tt.in, tt.want, got)
			}
		})
	}

	// 随机对拍：括号串长度 0~20，与暴力枚举比较。
	t.Run("随机对拍", func(t *testing.T) {
		r := rand.New(rand.NewSource(20240711))
		for iter := 0; iter < 400; iter++ {
			var sb strings.Builder
			for i, n := 0, r.Intn(21); i < n; i++ {
				sb.WriteByte("()"[r.Intn(2)])
			}
			s := sb.String()
			want := longestValidBrute(s)
			got := SolveBM077(s)
			if got != want {
				t.Fatalf("输入 %q 期望 %d，实际 %d", s, want, got)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// BM078 打家劫舍(一) · 测试
// ---------------------------------------------------------------------------
// robFullDP 用完整一维数组的 DP 参照实现（仅测试使用）。
func robFullDP(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	dp := make([]int, len(nums)+1)
	dp[1] = nums[0]
	for i := 2; i <= len(nums); i++ {
		dp[i] = dp[i-1]
		if dp[i-2]+nums[i-1] > dp[i] {
			dp[i] = dp[i-2] + nums[i-1]
		}
	}
	return dp[len(nums)]
}

// robExhaustive 枚举所有"不含相邻房间"的子集，指数级但绝对正确（仅测试使用，n 需很小）。
func robExhaustive(nums []int) int {
	best := 0
	for mask := 0; mask < 1<<len(nums); mask++ {
		sum, ok := 0, true
		for i := 0; i < len(nums); i++ {
			if mask&(1<<i) == 0 {
				continue
			}
			if i > 0 && mask&(1<<(i-1)) != 0 {
				ok = false // 相邻两间都被偷，方案非法
				break
			}
			sum += nums[i]
		}
		if ok && sum > best {
			best = sum
		}
	}
	return best
}

func TestBM078(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want int
	}{
		{"官方示例-1,2,3,4", []int{1, 2, 3, 4}, 6},
		{"官方示例-1,3,6", []int{1, 3, 6}, 7},
		{"官方示例-2,10,5", []int{2, 10, 5}, 10},
		{"空数组-没有房间", nil, 0},
		{"单间房", []int{5}, 5},
		{"单间房金额为0", []int{0}, 0},
		{"两间房取较大", []int{2, 1}, 2},
		{"两间房取较大-反过来", []int{1, 2}, 2},
		{"全为0", []int{0, 0, 0, 0}, 0},
		{"常规-隔一间取", []int{2, 7, 9, 3, 1}, 12},
		{"常规-三个房间", []int{2, 1, 1}, 3},
		{"常规-等值序列", []int{4, 4, 4, 4}, 8},
		{"常规-前小后大", []int{1, 1, 1, 1, 9}, 11},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SolveBM078(tt.in)
			if got != tt.want {
				t.Fatalf("输入 %v 期望 %d，实际 %d", tt.in, tt.want, got)
			}
			if ref := robFullDP(tt.in); ref != got {
				t.Fatalf("输入 %v 参照 DP 为 %d，实际 %d", tt.in, ref, got)
			}
		})
	}

	// 随机对拍：小规模用指数级枚举校验，稍大规模用参照 DP 校验。
	t.Run("随机对拍", func(t *testing.T) {
		r := rand.New(rand.NewSource(20240712))
		for iter := 0; iter < 300; iter++ {
			n := r.Intn(9)
			nums := make([]int, n)
			for i := range nums {
				nums[i] = r.Intn(401) // 题目的金额范围 0~400
			}
			want := robExhaustive(nums)
			if ref := robFullDP(nums); ref != want {
				t.Fatalf("输入 %v 两个参照实现不一致：%d vs %d", nums, ref, want)
			}
			got := SolveBM078(nums)
			if got != want {
				t.Fatalf("输入 %v 期望 %d，实际 %d", nums, want, got)
			}
		}
	})

	// 长数组：10 万间房，验证 O(n) 实现不会超时且结果与参照 DP 一致。
	t.Run("长数组-10万间房", func(t *testing.T) {
		nums := make([]int, 100000)
		for i := range nums {
			nums[i] = 400
		}
		want := robFullDP(nums)
		if got := SolveBM078(nums); got != want {
			t.Fatalf("10 万间房期望 %d，实际 %d", want, got)
		}
	})
}

// ---------------------------------------------------------------------------
// BM079 打家劫舍(二) · 测试
// ---------------------------------------------------------------------------
func TestBM079(t *testing.T) {
	// 长数组用例：100000 间房全部是 1，环上最多能隔一间偷一间，答案为 n/2。
	const longN = 100000
	long := make([]int, longN)
	for i := range long {
		long[i] = 1
	}

	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"官方示例-四个房间", []int{1, 2, 3, 4}, 6},
		{"官方示例-首尾相邻只能偷一间", []int{1, 3, 6}, 6},
		{"空数组", nil, 0},
		{"单元素", []int{5}, 5},
		{"两个元素-首尾相邻取较大", []int{7, 2}, 7},
		{"两个元素-后者更大", []int{2, 7}, 7},
		{"三个元素-任两间都相邻只能偷一间", []int{2, 9, 3}, 9},
		{"全部相同-偶数间房", []int{3, 3, 3, 3}, 6},
		{"全部为零", []int{0, 0, 0, 0}, 0},
		{"首尾都大但相邻", []int{10, 1, 1, 10}, 11},
		{"常规用例-五间房", []int{2, 7, 9, 3, 1}, 11},
		{"长数组-100000间房全为1", long, longN / 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SolveBM079(tt.nums)
			if got != tt.want {
				t.Fatalf("期望 %v，实际 %v", tt.want, got)
			}
		})
	}

	// 与暴力枚举对拍：小规模随机数据下逐一枚举「偷哪些房间」的方案，验证环形约束与答案。
	t.Run("与暴力枚举对拍", func(t *testing.T) {
		rnd := rand.New(rand.NewSource(79))
		for c := 0; c < 200; c++ {
			n := 2 + rnd.Intn(9) // 房间数 2..10
			nums := make([]int, n)
			for i := range nums {
				nums[i] = rnd.Intn(20)
			}
			want := bruteForceBM079(nums)
			got := SolveBM079(nums)
			if got != want {
				t.Fatalf("输入 %v：期望 %v，实际 %v", nums, want, got)
			}
		}
	})
}

// bruteForceBM079 枚举所有偷窃方案（至多 2^n 个），返回满足「环上相邻不能同时偷」的最大金额。
func bruteForceBM079(nums []int) int {
	n := len(nums)
	if n == 0 {
		return 0
	}
	if n == 1 {
		return nums[0]
	}

	best := 0
	for mask := 0; mask < 1<<n; mask++ {
		sum := 0
		ok := true
		for i := 0; i < n; i++ {
			if mask&(1<<i) == 0 {
				continue
			}
			// 环形相邻：第 i 间与第 (i+1)%n 间不能同时被选中。
			if mask&(1<<((i+1)%n)) != 0 {
				ok = false
				break
			}
			sum += nums[i]
		}
		if ok && sum > best {
			best = sum
		}
	}
	return best
}

// ---------------------------------------------------------------------------
// BM080 买卖股票的最好时机(一) · 测试
// ---------------------------------------------------------------------------
func TestBM080(t *testing.T) {
	// 长数组用例：价格 1..100000 单调递增，第 1 天买、最后一天卖，利润为 99999。
	const longN = 100000
	long := make([]int, longN)
	for i := range long {
		long[i] = i + 1
	}

	tests := []struct {
		name   string
		prices []int
		want   int
	}{
		{"官方示例1-先跌后涨", []int{8, 9, 2, 5, 4, 7, 1}, 5},
		{"官方示例2-小数组", []int{2, 4, 1}, 2},
		{"官方示例3-所有价格递减", []int{3, 2, 1}, 0},
		{"空数组", nil, 0},
		{"单元素", []int{7}, 0},
		{"所有价格递减-长一点", []int{7, 6, 4, 3, 1}, 0},
		{"平盘-价格全部相同", []int{5, 5, 5, 5}, 0},
		{"单调递增-第一天买最后一天卖", []int{1, 2, 3, 4, 5}, 4},
		{"最低价在最后", []int{9, 8, 2, 1}, 0},
		{"最高价在最后", []int{5, 4, 3, 9}, 6},
		{"常规用例-山谷形", []int{7, 1, 5, 3, 6, 4}, 5},
		{"长数组-单调递增", long, longN - 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SolveBM080(tt.prices)
			if got != tt.want {
				t.Fatalf("期望 %v，实际 %v", tt.want, got)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// BM081 买卖股票的最好时机(二) · 测试
// ---------------------------------------------------------------------------
func TestBM081(t *testing.T) {
	// 长数组用例：价格 1..100000 单调递增，每天低买高卖，总收益为 99999。
	const longN = 100000
	long := make([]int, longN)
	for i := range long {
		long[i] = i + 1
	}

	tests := []struct {
		name   string
		prices []int
		want   int
	}{
		{"官方示例1-多段上涨", []int{8, 9, 2, 5, 4, 7, 1}, 7},
		{"官方示例2-所有价格递减", []int{5, 4, 3, 2, 1}, 0},
		{"官方示例3-单调递增", []int{1, 2, 3, 4, 5}, 4},
		{"空数组", nil, 0},
		{"单元素", []int{7}, 0},
		{"平盘-价格全部相同", []int{5, 5, 5, 5}, 0},
		{"锯齿形-每一段上涨都吃掉", []int{1, 3, 2, 4, 3, 5}, 6},
		{"只涨一天", []int{1, 2}, 1},
		{"只跌一天", []int{2, 1}, 0},
		{"常规用例-先跌后涨", []int{7, 1, 5, 3, 6, 4}, 7},
		{"长数组-单调递增", long, longN - 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SolveBM081(tt.prices)
			if got != tt.want {
				t.Fatalf("期望 %v，实际 %v", tt.want, got)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// BM082 买卖股票的最好时机(三) · 测试
// ---------------------------------------------------------------------------
func TestBM082(t *testing.T) {
	// 长数组用例：价格 1..100000 单调递增，最多两笔但实际只需一笔，收益 99999。
	const longN = 100000
	long := make([]int, longN)
	for i := range long {
		long[i] = i + 1
	}

	tests := []struct {
		name   string
		prices []int
		want   int
	}{
		{"官方示例1-两笔各赚2", []int{8, 9, 3, 5, 1, 3}, 4},
		{"官方示例2-所有价格递减", []int{9, 8, 4, 1}, 0},
		{"官方示例3-两笔各赚7和5", []int{1, 2, 8, 3, 8}, 12},
		{"空数组", nil, 0},
		{"单元素", []int{7}, 0},
		{"平盘-价格全部相同", []int{4, 4, 4}, 0},
		{"单调递增-一笔即可", []int{1, 2, 3, 4, 5}, 4},
		{"最优必须用两笔", []int{1, 4, 2, 6}, 7},
		{"最优必须用两笔-长一点", []int{3, 3, 5, 0, 0, 3, 1, 4}, 6},
		{"官方示例1的变体-最优两笔共6", []int{8, 9, 2, 5, 4, 7, 1}, 6},
		{"一大一小两笔", []int{1, 10, 2, 3}, 10},
		{"长数组-单调递增", long, longN - 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SolveBM082(tt.prices)
			if got != tt.want {
				t.Fatalf("期望 %v，实际 %v", tt.want, got)
			}
		})
	}

	// 与暴力枚举对拍：枚举「在某个切分点之前做第一笔、之后做第二笔」的所有方案，
	// 小规模随机数据下验证状态机的正确性（尤其是只有一笔收益为 0 的边界）。
	t.Run("与暴力枚举对拍", func(t *testing.T) {
		rnd := rand.New(rand.NewSource(82))
		for c := 0; c < 300; c++ {
			n := 1 + rnd.Intn(9) // 长度 1..9
			prices := make([]int, n)
			for i := range prices {
				prices[i] = 1 + rnd.Intn(20)
			}
			want := bruteForceBM082(prices)
			got := SolveBM082(prices)
			if got != want {
				t.Fatalf("输入 %v：期望 %v，实际 %v", prices, want, got)
			}
		}
	})
}

// bruteForceBM082 枚举所有「至多两笔、互不重叠、先买后卖」的方案（仅测试使用）。
// 任意两笔合规交易之间都存在一个切分点：第一笔完全落在左侧，第二笔完全落在右侧。
func bruteForceBM082(prices []int) int {
	best := 0
	for split := -1; split < len(prices); split++ {
		// best 至少为 one(-1) + one(0..n-1) = 只用一笔最优，因此「0 笔」也被覆盖。
		if v := bestOneBM082(prices, 0, split) + bestOneBM082(prices, split+1, len(prices)-1); v > best {
			best = v
		}
	}
	return best
}

// bestOneBM082 返回在闭区间 [lo, hi] 内只做一笔交易的最大收益（空区间返回 0）。
func bestOneBM082(prices []int, lo, hi int) int {
	if lo > hi {
		return 0
	}
	best := 0
	minPrice := prices[lo]
	for i := lo; i <= hi; i++ {
		if prices[i] < minPrice {
			minPrice = prices[i]
		}
		if prices[i]-minPrice > best {
			best = prices[i] - minPrice
		}
	}
	return best
}

// ---------------------------------------------------------------------------
// BM083 字符串变形 · 测试
// ---------------------------------------------------------------------------
func TestBM083(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"官方示例1-四个单词", "This is a sample", "SAMPLE A IS tHIS"},
		{"官方示例2-单个单词", "nowcoder", "NOWCODER"},
		{"官方示例3-大小写混合", "iOS", "Ios"},
		{"官方题目描述里的例子", "Hello World", "wORLD hELLO"},
		{"空字符串", "", ""},
		{"只有一个空格", " ", " "},
		{"只有一个字母", "a", "A"},
		{"单词间单空格", "a b", "B A"},
		{"三个单词", "one two three", "THREE TWO ONE"},
		{"全大写", "ABC DEF", "def abc"},
		{"多空格-单词之间", "a  b", "B  A"},
		{"多空格-首尾前后缀", "  a  b  ", "  B  A  "},
		{"首尾各一个空格", " ab ", " AB "},
		{"全是空格", "    ", "    "},
		{"单词内大小写混合", "ABc dEf", "DeF abC"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SolveBM083(tt.in)
			if got != tt.want {
				t.Fatalf("输入 %q：期望 %q，实际 %q", tt.in, tt.want, got)
			}
		})
	}

	// 与「按空格切词、反序、再逐词反转大小写」的参考实现交叉验证，覆盖随机空格分布。
	t.Run("与逐词参考实现一致", func(t *testing.T) {
		rnd := rand.New(rand.NewSource(83))
		letters := []byte("abzABZ")
		for c := 0; c < 200; c++ {
			var sb strings.Builder
			for k := 0; k < rnd.Intn(12); k++ {
				if rnd.Intn(3) == 0 {
					sb.WriteByte(' ')
					continue
				}
				sb.WriteByte(letters[rnd.Intn(len(letters))])
			}
			in := sb.String()
			want := referenceBM083(in)
			got := SolveBM083(in)
			if got != want {
				t.Fatalf("输入 %q：期望 %q，实际 %q", in, want, got)
			}
		}
	})

	// 长字符串：100000 个字符规模下不应出错。
	t.Run("长字符串", func(t *testing.T) {
		in := strings.Repeat("aB cD ", 20000)
		want := referenceBM083(in)
		got := SolveBM083(in)
		if got != want {
			t.Fatalf("长度 %d 的输入：期望 %q，实际 %q", len(in), want, got)
		}
		if len(got) != len(in) {
			t.Fatalf("期望长度 %d，实际 %d", len(in), len(got))
		}
	})
}

// referenceBM083 是逐词实现的参考版本：按空格切分 → 词序反转 → 每个单词翻转大小写 → 拼回（仅测试使用）。
func referenceBM083(s string) string {
	words := strings.Split(s, " ")
	for i, j := 0, len(words)-1; i < j; i, j = i+1, j-1 {
		words[i], words[j] = words[j], words[i]
	}
	for i := range words {
		b := []byte(words[i])
		for k := range b {
			switch {
			case b[k] >= 'a' && b[k] <= 'z':
				b[k] -= 'a' - 'A'
			case b[k] >= 'A' && b[k] <= 'Z':
				b[k] += 'a' - 'A'
			}
		}
		words[i] = string(b)
	}
	return strings.Join(words, " ")
}

// ---------------------------------------------------------------------------
// BM084 最长公共前缀 · 测试
// ---------------------------------------------------------------------------
func TestBM084(t *testing.T) {
	tests := []struct {
		name string
		strs []string
		want string
	}{
		{"官方示例1-五个字符串", []string{"abca", "abc", "abca", "abc", "abcc"}, "abc"},
		{"官方示例2-只有一个字符串", []string{"abc"}, "abc"},
		{"空数组", nil, ""},
		{"空切片", []string{}, ""},
		{"只有一个空字符串", []string{""}, ""},
		{"含空字符串", []string{"abc", "", "abd"}, ""},
		{"两个空字符串", []string{"", ""}, ""},
		{"无公共前缀", []string{"dog", "racecar", "car"}, ""},
		{"常规用例-前缀 fl", []string{"flower", "flow", "flight"}, "fl"},
		{"第一个字符串本身就最短", []string{"abc", "abcdef", "abcd"}, "abc"},
		{"完全相同的字符串", []string{"abc", "abc", "abc"}, "abc"},
		{"单字符数组", []string{"a", "a", "a"}, "a"},
		{"全不相同的首字符", []string{"a", "b"}, ""},
		{"公共前缀等于整个最短串", []string{"aaaa", "aa", "aaa"}, "aa"},
		{"公共前缀包含空格", []string{" ab", " abc"}, " ab"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SolveBM084(tt.strs)
			if got != tt.want {
				t.Fatalf("输入 %v：期望 %q，实际 %q", tt.strs, tt.want, got)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// BM085 验证IP地址 · 测试
// ---------------------------------------------------------------------------
func TestBM085(t *testing.T) {
	tests := []struct {
		name string
		ip   string
		want string
	}{
		// 官方示例
		{name: "官方示例1-IPv4", ip: "172.16.254.1", want: "IPv4"},
		{name: "官方示例2-IPv6", ip: "2001:0db8:85a3:0:0:8A2E:0370:7334", want: "IPv6"},
		{name: "官方示例3-全256", ip: "256.256.256.256", want: "Neither"},

		// IPv4 合法
		{name: "IPv4-常规", ip: "1.0.0.1", want: "IPv4"},
		{name: "IPv4-全零", ip: "0.0.0.0", want: "IPv4"},
		{name: "IPv4-全255", ip: "255.255.255.255", want: "IPv4"},
		{name: "IPv4-三段位混合", ip: "192.168.0.255", want: "IPv4"},

		// IPv4 非法
		{name: "IPv4-超范围256", ip: "256.1.1.1", want: "Neither"},
		{name: "IPv4-前导零01", ip: "01.1.1.1", want: "Neither"},
		{name: "IPv4-前导零后置段", ip: "1.1.1.01", want: "Neither"},
		{name: "IPv4-前导零00", ip: "1.1.1.00", want: "Neither"},
		{name: "IPv4-段数不足", ip: "1.1.1", want: "Neither"},
		{name: "IPv4-段数过多", ip: "1.1.1.1.1", want: "Neither"},
		{name: "IPv4-空段", ip: "1..1.1", want: "Neither"},
		{name: "IPv4-含非数字", ip: "1.1.1.a", want: "Neither"},
		{name: "IPv4-数值越界", ip: "999.1.1.1", want: "Neither"},
		{name: "IPv4-末尾带点", ip: "1.1.1.1.", want: "Neither"},

		// IPv6 合法
		{name: "IPv6-忽略前导零与大小写", ip: "2001:db8:85a3:0:0:8A2E:0370:7334", want: "IPv6"},
		{name: "IPv6-全写", ip: "2001:0db8:85a3:0000:0000:8a2e:0370:7334", want: "IPv6"},
		{name: "IPv6-全零", ip: "0:0:0:0:0:0:0:0", want: "IPv6"},
		{name: "IPv6-全f", ip: "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff", want: "IPv6"},

		// IPv6 非法
		{name: "IPv6-压缩双冒号", ip: "2001:db8:85a3::8A2E:0370:7334", want: "Neither"},
		{name: "IPv6-仅回环缩写", ip: "::1", want: "Neither"},
		{name: "IPv6-组数不足", ip: "2001:db8:85a3:0:0:8A2E:0370", want: "Neither"},
		{name: "IPv6-组数过多", ip: "2001:db8:85a3:0:0:8A2E:0370:7334:1234", want: "Neither"},
		{name: "IPv6-多余的前导零", ip: "02001:0db8:85a3:0000:0000:8a2e:0370:7334", want: "Neither"},
		{name: "IPv6-非十六进制", ip: "2001:db8:85a3:0:0:8A2E:0370:733G", want: "Neither"},
		{name: "IPv6-末尾冒号", ip: "2001:db8:85a3:0:0:8A2E:0370:7334:", want: "Neither"},

		// 其它边界
		{name: "空串", ip: "", want: "Neither"},
		{name: "混合分隔符", ip: "1.1.1.1:80", want: "Neither"},
		{name: "纯数字", ip: "1234", want: "Neither"},
		{name: "纯字母", ip: "invalid", want: "Neither"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := SolveBM085(tc.ip)
			if got != tc.want {
				t.Fatalf("输入 %q：期望 %q，实际 %q", tc.ip, tc.want, got)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// BM086 大数加法 · 测试
// ---------------------------------------------------------------------------
func TestBM086(t *testing.T) {
	tests := []struct {
		name string
		s    string
		t    string
		want string
	}{
		// 官方示例
		{name: "官方示例1", s: "1", t: "99", want: "100"},
		{name: "官方示例2-空串", s: "114514", t: "", want: "114514"},

		// 进位相关
		{name: "进位到新的一位", s: "99", t: "1", want: "100"},
		{name: "连续进位", s: "999", t: "1", want: "1000"},
		{name: "等长无进位", s: "100", t: "200", want: "300"},
		{name: "等长有进位", s: "123456789", t: "987654321", want: "1111111110"},

		// 边界
		{name: "零加零", s: "0", t: "0", want: "0"},
		{name: "零加数", s: "0", t: "9527", want: "9527"},
		{name: "两个空串", s: "", t: "", want: "0"},
		{name: "左侧为空", s: "", t: "114514", want: "114514"},
		{name: "长度相差很大", s: "1", t: "999999999", want: "1000000000"},
		{name: "前导零输入", s: "007", t: "3", want: "10"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := SolveBM086(tc.s, tc.t)
			if got != tc.want {
				t.Fatalf("%q + %q：期望 %q，实际 %q", tc.s, tc.t, tc.want, got)
			}
		})
	}

	// 大数用例：1000 个 9 加 1，结果应为 1 后面跟 1000 个 0（验证不溢出、不丢位）。
	t.Run("一千位全九加一", func(t *testing.T) {
		s := strings.Repeat("9", 1000)
		want := "1" + strings.Repeat("0", 1000)
		got := SolveBM086(s, "1")
		if got != want {
			t.Fatalf("长度 %d 的全九串 + 1：期望 %q，实际 %q", len(s), want, got)
		}
	})

	// 交换律自检：a + b 应与 b + a 一致（10^50 + (10^50-1) = 2*10^50 - 1）。
	t.Run("加法交换律", func(t *testing.T) {
		a := "1" + strings.Repeat("0", 50)
		b := strings.Repeat("9", 50)
		want := "1" + strings.Repeat("9", 50)
		if got := SolveBM086(a, b); got != want {
			t.Fatalf("%q + %q：期望 %q，实际 %q", a, b, want, got)
		}
		if SolveBM086(a, b) != SolveBM086(b, a) {
			t.Fatalf("交换律不满足：%q != %q", SolveBM086(a, b), SolveBM086(b, a))
		}
	})
}

// ---------------------------------------------------------------------------
// BM087 合并两个有序的数组 · 测试
// ---------------------------------------------------------------------------
// buildMergedA 按题设构造 A：前 m 个是有效值，后 n 个用 0 占位，长度恰为 m+n。
func buildMergedA(valid []int, n int) []int {
	a := make([]int, len(valid)+n)
	copy(a, valid)
	return a
}

func TestBM087(t *testing.T) {
	tests := []struct {
		name string
		a    []int // A 的有效部分（前 m 个）
		b    []int
		want []int
	}{
		// 官方示例
		{name: "官方示例1", a: []int{4, 5, 6}, b: []int{1, 2, 3}, want: []int{1, 2, 3, 4, 5, 6}},
		{name: "官方示例2", a: []int{1, 2, 3}, b: []int{2, 5, 6}, want: []int{1, 2, 2, 3, 5, 6}},

		// 边界：一边为空
		{name: "A为空", a: []int{}, b: []int{1, 2, 3}, want: []int{1, 2, 3}},
		{name: "B为空", a: []int{1, 2, 3}, b: []int{}, want: []int{1, 2, 3}},
		{name: "A和B都为空", a: []int{}, b: []int{}, want: []int{}},
		{name: "单元素A为空", a: []int{}, b: []int{7}, want: []int{7}},
		{name: "单元素B为空", a: []int{7}, b: []int{}, want: []int{7}},

		// 边界：整体大小关系
		{name: "A全部小于B", a: []int{1, 2, 3}, b: []int{4, 5, 6}, want: []int{1, 2, 3, 4, 5, 6}},
		{name: "B全部小于A", a: []int{4, 5, 6}, b: []int{1, 2, 3}, want: []int{1, 2, 3, 4, 5, 6}},
		{name: "全为相同值", a: []int{2, 2}, b: []int{2, 2}, want: []int{2, 2, 2, 2}},

		// 常规用例
		{name: "交错合并", a: []int{1, 3, 5, 7}, b: []int{2, 4, 6, 8}, want: []int{1, 2, 3, 4, 5, 6, 7, 8}},
		{name: "长度不等", a: []int{1, 2}, b: []int{3, 4, 5, 6, 7}, want: []int{1, 2, 3, 4, 5, 6, 7}},
		{name: "含极值", a: []int{0, 1000000000}, b: []int{5}, want: []int{0, 5, 1000000000}},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			m, n := len(tc.a), len(tc.b)
			A := buildMergedA(tc.a, n)
			B := make([]int, len(tc.b))
			copy(B, tc.b)

			SolveBM087(A, m, B, n)

			if !reflect.DeepEqual(A, tc.want) {
				t.Fatalf("期望 %v，实际 %v", tc.want, A)
			}
			// 不应修改 B
			if !reflect.DeepEqual(B, tc.b) {
				t.Fatalf("B 不应被修改：期望 %v，实际 %v", tc.b, B)
			}
		})
	}

	// 规模用例：A、B 各 1000 个元素，验证整体有序且元素计数正确。
	t.Run("规模用例一千加一千", func(t *testing.T) {
		const n = 1000
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = 2 * i   // 0,2,4,...
			b[i] = 2*i + 1 // 1,3,5,...
		}
		A := buildMergedA(a, n)
		SolveBM087(A, n, b, n)
		for i := 0; i < 2*n; i++ {
			if A[i] != i {
				t.Fatalf("下标 %d：期望 %d，实际 %d", i, i, A[i])
			}
		}
	})
}

// ---------------------------------------------------------------------------
// BM088 判断是否为回文字符串 · 测试
// ---------------------------------------------------------------------------
func TestBM088(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		// 官方示例
		{name: "官方示例1-absba", s: "absba", want: true},
		{name: "官方示例2-ranko", s: "ranko", want: false},
		{name: "官方示例3-yamatomaya", s: "yamatomaya", want: false},
		{name: "官方示例4-单字符", s: "a", want: true},

		// 边界
		{name: "空串", s: "", want: true},
		{name: "两字符相同", s: "aa", want: true},
		{name: "两字符不同", s: "ab", want: false},
		{name: "全同字符", s: "aaaaa", want: true},

		// 常规
		{name: "偶数长度回文", s: "abccba", want: true},
		{name: "偶数长度非回文", s: "abccbb", want: false},
		{name: "奇数长度回文", s: "abcba", want: true},
		{name: "仅中间不同", s: "abcda", want: false},
		{name: "含数字回文", s: "a1b2b1a", want: true},
		{name: "含数字非回文", s: "a1b2c1a", want: false},
		{name: "大小写敏感", s: "Aa", want: false},
		{name: "全数字回文", s: "12321", want: true},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := SolveBM088(tc.s)
			if got != tc.want {
				t.Fatalf("输入 %q：期望 %v，实际 %v", tc.s, tc.want, got)
			}
		})
	}

	// 长串用例：长度 100000 的回文与非回文，验证不超时、不越界。
	t.Run("长回文串", func(t *testing.T) {
		half := strings.Repeat("ab", 25000) // 长度 50000
		s := half + "z" + reverse(half)     // 长度 100001
		if !SolveBM088(s) {
			t.Fatalf("长度 %d 的回文串被判定为非回文", len(s))
		}
	})

	t.Run("长非回文串", func(t *testing.T) {
		s := strings.Repeat("ab", 50000)
		if SolveBM088(s) {
			t.Fatalf("长度 %d 的非回文串被判定为回文", len(s))
		}
	})
}

// reverse 按字节反转字符串，仅用于测试构造用例。
func reverse(s string) string {
	b := []byte(s)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

// ---------------------------------------------------------------------------
// BM089 合并区间 · 测试
// ---------------------------------------------------------------------------
// normalizeIntervals 复制并按 (Start, End) 升序排序，用于不受输入顺序影响的比较。
func normalizeIntervals(in []Interval) []Interval {
	out := make([]Interval, len(in))
	copy(out, in)
	sort.Slice(out, func(i, j int) bool {
		if out[i].Start != out[j].Start {
			return out[i].Start < out[j].Start
		}
		return out[i].End < out[j].End
	})
	return out
}

func TestBM089(t *testing.T) {
	tests := []struct {
		name string
		in   []Interval
		want []Interval
	}{
		// 官方示例
		{
			name: "官方示例1",
			in:   []Interval{{10, 30}, {20, 60}, {80, 100}, {150, 180}},
			want: []Interval{{10, 60}, {80, 100}, {150, 180}},
		},
		{
			name: "官方示例2-端点相接",
			in:   []Interval{{0, 10}, {10, 20}},
			want: []Interval{{0, 20}},
		},

		// 边界
		{name: "空区间组", in: []Interval{}, want: []Interval{}},
		{name: "单区间", in: []Interval{{5, 7}}, want: []Interval{{5, 7}}},
		{name: "单点区间", in: []Interval{{0, 0}}, want: []Interval{{0, 0}}},
		{name: "重合的单点区间", in: []Interval{{0, 0}, {0, 0}}, want: []Interval{{0, 0}}},

		// 包含与相交
		{name: "完全包含", in: []Interval{{1, 10}, {2, 3}}, want: []Interval{{1, 10}}},
		{name: "被包含", in: []Interval{{2, 3}, {1, 10}}, want: []Interval{{1, 10}}},
		{name: "部分相交", in: []Interval{{1, 3}, {2, 6}, {8, 10}, {15, 18}}, want: []Interval{{1, 6}, {8, 10}, {15, 18}}},
		{name: "链式重叠", in: []Interval{{1, 4}, {2, 3}, {5, 7}}, want: []Interval{{1, 4}, {5, 7}}},
		{name: "全部互相重叠", in: []Interval{{1, 4}, {2, 5}, {3, 8}}, want: []Interval{{1, 8}}},
		{name: "起止完全相同", in: []Interval{{1, 2}, {1, 2}}, want: []Interval{{1, 2}}},
		{name: "起点相同终点不同", in: []Interval{{1, 5}, {1, 3}}, want: []Interval{{1, 5}}},

		// 顺序与常规
		{
			name: "逆序输入",
			in:   []Interval{{150, 180}, {80, 100}, {20, 60}, {10, 30}},
			want: []Interval{{10, 60}, {80, 100}, {150, 180}},
		},
		{name: "互不相交需排序", in: []Interval{{3, 4}, {1, 2}, {5, 6}}, want: []Interval{{1, 2}, {3, 4}, {5, 6}}},
		{name: "相邻但不重叠", in: []Interval{{1, 2}, {3, 4}}, want: []Interval{{1, 2}, {3, 4}}},
		{name: "最大取值", in: []Interval{{190000, 200000}, {0, 200000}}, want: []Interval{{0, 200000}}},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// 深拷贝一份原始输入，用于检查函数没有就地修改入参
			origin := make([]Interval, len(tc.in))
			copy(origin, tc.in)

			got := normalizeIntervals(SolveBM089(tc.in))
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("期望 %v，实际 %v", tc.want, got)
			}
			if !reflect.DeepEqual(tc.in, origin) {
				t.Fatalf("入参不应被修改：原始 %v，现在 %v", origin, tc.in)
			}
		})
	}

	// 规模用例：2500 个首尾相连的区间，应当被合并成一个。
	t.Run("首尾相连全部合并", func(t *testing.T) {
		const n = 2500
		in := make([]Interval, n)
		for i := 0; i < n; i++ {
			in[i] = Interval{Start: i * 2, End: i*2 + 2}
		}
		got := SolveBM089(in)
		want := []Interval{{0, 2 * n}}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("期望 %v，实际 %v（共 %d 段）", want, got, len(got))
		}
	})
}

// ---------------------------------------------------------------------------
// BM090 最小覆盖子串 · 测试
// ---------------------------------------------------------------------------
// covers 校验 window 是否覆盖了 whole 中全部字符（含重复次数），用于性质断言。
func covers(window, whole string) bool {
	if len(window) == 0 {
		return len(whole) == 0
	}
	need := make(map[byte]int, len(whole))
	for i := 0; i < len(whole); i++ {
		need[whole[i]]++
	}
	for c, cnt := range need {
		if strings.Count(window, string(c)) < cnt {
			return false
		}
	}
	return true
}

func TestBM090(t *testing.T) {
	tests := []struct {
		name string
		s    string
		t    string
		want string
	}{
		// 官方示例
		{name: "官方示例1", s: "XDOYEZODEYXNZ", t: "XYZ", want: "YXNZ"},
		{name: "官方示例2", s: "abcAbA", t: "AA", want: "AbA"},
		{name: "示例3-无解", s: "abcAbc", t: "AB", want: ""},

		// 经典用例
		{name: "经典ADOBECODEBANC", s: "ADOBECODEBANC", t: "ABC", want: "BANC"},
		{name: "滑动窗口需回退", s: "cabwefgewcwaefgcf", t: "cae", want: "cwae"},
		{name: "重复字符需求", s: "aaflslflsldkalskaaa", t: "aaa", want: "aaa"},
		{name: "整串即答案", s: "ab", t: "ab", want: "ab"},
		{name: "单字符覆盖", s: "ab", t: "a", want: "a"},
		{name: "取尾部单字符", s: "ab", t: "b", want: "b"},
		{name: "T含重复字符", s: "aab", t: "aa", want: "aa"},
		{name: "顺序无关", s: "bba", t: "ab", want: "ba"},

		// 边界
		{name: "S为空", s: "", t: "abc", want: ""},
		{name: "T为空", s: "abc", t: "", want: ""},
		{name: "S和T都为空", s: "", t: "", want: ""},
		{name: "T比S长", s: "a", t: "aa", want: ""},
		{name: "字符集不匹配", s: "aaaa", t: "b", want: ""},
		{name: "单字符匹配", s: "a", t: "a", want: "a"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := SolveBM090(tc.s, tc.t)
			if got != tc.want {
				t.Fatalf("S=%q T=%q：期望 %q，实际 %q", tc.s, tc.t, tc.want, got)
			}
			// 性质校验：有解时必是 S 的子串，且覆盖 T 的全部字符
			if got != "" {
				if !strings.Contains(tc.s, got) {
					t.Fatalf("结果 %q 不是 S=%q 的子串", got, tc.s)
				}
				if !covers(got, tc.t) {
					t.Fatalf("结果 %q 未覆盖 T=%q 的全部字符", got, tc.t)
				}
			}
		})
	}

	// 规模用例：S 为 10000 个 'a' 与 'b' 交替，T = "ab"，最短答案长度为 2。
	t.Run("规模用例一万长度", func(t *testing.T) {
		s := strings.Repeat("ab", 5000)
		got := SolveBM090(s, "ba")
		if len(got) != 2 || !covers(got, "ba") {
			t.Fatalf("期望长度 2 且覆盖 \"ba\"，实际 %q（长度 %d）", got, len(got))
		}
	})
}

// ---------------------------------------------------------------------------
// BM091 反转字符串 · 测试
// ---------------------------------------------------------------------------
func TestBM091(t *testing.T) {
	// 构造上下界规模的长字符串（1000 个字符），用于覆盖题目给出的最大长度。
	var longIn strings.Builder
	var longWant strings.Builder
	for i := 0; i < 1000; i++ {
		longIn.WriteByte(byte('a' + i%26))
		longWant.WriteByte(byte('a' + (999-i)%26))
	}

	tests := []struct {
		name string
		in   string
		want string
	}{
		{"官方示例-普通字符串", "abcd", "dcba"},
		{"官方示例-空串", "", ""},
		{"单字符", "a", "a"},
		{"全相同字符", "aaaaaa", "aaaaaa"},
		{"回文字符串", "level", "level"},
		{"常规用例-大小写数字混合", "A1b2C3", "3C2b1A"},
		{"含空格字符串", "ab cd", "dc ba"},
		{"中文字符-按字符反转不拆字节", "程序员牛客", "客牛员序程"},
		{"长字符串-1000字符", longIn.String(), longWant.String()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SolveBM091(tt.in)
			if got != tt.want {
				t.Fatalf("期望 %q，实际 %q", tt.want, got)
			}
			// 反转不改变长度（按 rune 计算），顺便校验没有丢字符。
			if len([]rune(got)) != len([]rune(tt.in)) {
				t.Fatalf("长度应保持 %d，实际 %d", len([]rune(tt.in)), len([]rune(got)))
			}
		})
	}
}

// ---------------------------------------------------------------------------
// BM092 最长无重复子数组 · 测试
// ---------------------------------------------------------------------------
// bruteForceMaxLength 暴力枚举所有连续子数组求最长无重复长度，仅用于与滑窗结果交叉验证。
func bruteForceMaxLength(arr []int) int {
	best := 0
	for i := range arr {
		seen := make(map[int]struct{}, len(arr))
		for j := i; j < len(arr); j++ {
			if _, ok := seen[arr[j]]; ok {
				break
			}
			seen[arr[j]] = struct{}{}
			if l := j - i + 1; l > best {
				best = l
			}
		}
	}
	return best
}

func TestBM092(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want int
	}{
		{"官方示例-全部不重复", []int{2, 3, 4, 5}, 4},
		{"官方示例-含重复元素", []int{2, 2, 3, 4, 3}, 3},
		{"边界-空数组", nil, 0},
		{"边界-单元素", []int{1}, 1},
		{"边界-全相同元素", []int{7, 7, 7, 7}, 1},
		{"边界-全递增", []int{1, 2, 3, 4, 5}, 5},
		{"边界-全递减", []int{5, 4, 3, 2, 1}, 5},
		{"常规-重复值跨窗口", []int{1, 2, 3, 1, 2, 3}, 3},
		{"常规-含负数与零", []int{-1, 0, -1, 2}, 3},
		{"常规-稀疏大值", []int{100000, 1, 100000}, 2},
		{"常规-重复值在窗口外", []int{1, 2, 1, 3, 4}, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SolveBM092(tt.in)
			if got != tt.want {
				t.Fatalf("期望 %v，实际 %v", tt.want, got)
			}
		})
	}

	// 固定种子生成随机数组，与 O(n^2) 暴力解交叉验证滑动窗口的正确性。
	t.Run("随机交叉验证", func(t *testing.T) {
		rnd := rand.New(rand.NewSource(20240501))
		for round := 0; round < 300; round++ {
			n := rnd.Intn(30)
			arr := make([]int, n)
			for i := range arr {
				arr[i] = rnd.Intn(8) // 值域刻意取小，制造大量重复以覆盖窗口收缩分支
			}
			want := bruteForceMaxLength(arr)
			got := SolveBM092(arr)
			if got != want {
				t.Fatalf("第 %d 轮：arr=%v，期望 %d，实际 %d", round, arr, want, got)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// BM093 盛水最多的容器 · 测试
// ---------------------------------------------------------------------------
// bruteForceMaxArea 枚举所有柱子对求最大容量，仅用于与双指针结果交叉验证。
func bruteForceMaxArea(height []int) int {
	best := 0
	for i := 0; i < len(height); i++ {
		for j := i + 1; j < len(height); j++ {
			h := height[i]
			if height[j] < h {
				h = height[j]
			}
			if area := (j - i) * h; area > best {
				best = area
			}
		}
	}
	return best
}

func TestBM093(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want int
	}{
		{"官方示例1", []int{1, 7, 3, 2, 4, 5, 8, 2, 7}, 49},
		{"官方示例2", []int{5, 4, 3, 2, 1, 5}, 25},
		{"示例3-两根等高柱", []int{2, 2}, 2},
		{"常规-两根不等高柱", []int{1, 2}, 1},
		{"边界-空数组", nil, 0},
		{"边界-只有一根柱子", []int{9}, 0},
		{"边界-全相同柱形", []int{3, 3, 3, 3}, 9},
		{"边界-全递增柱形", []int{1, 2, 3, 4, 5}, 6},
		{"边界-全递减柱形", []int{5, 4, 3, 2, 1}, 6},
		{"边界-含零高度柱", []int{0, 5}, 0},
		{"常规-两端最高", []int{5, 1, 1, 1, 5}, 20},
		{"常规-中间短板卡住两侧", []int{9, 1, 8}, 16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SolveBM093(tt.in)
			if got != tt.want {
				t.Fatalf("期望 %v，实际 %v", tt.want, got)
			}
		})
	}

	// 固定种子随机柱形，与暴力解交叉验证「移动较短一侧」不会漏掉最优解。
	t.Run("随机交叉验证", func(t *testing.T) {
		rnd := rand.New(rand.NewSource(20240502))
		for round := 0; round < 300; round++ {
			n := rnd.Intn(30)
			height := make([]int, n)
			for i := range height {
				height[i] = rnd.Intn(50)
			}
			want := bruteForceMaxArea(height)
			got := SolveBM093(height)
			if got != want {
				t.Fatalf("第 %d 轮：height=%v，期望 %d，实际 %d", round, height, want, got)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// BM094 接雨水问题 · 测试
// ---------------------------------------------------------------------------
// bruteForceTrap 用前缀最高/后缀最高直接套公式求接水量，仅用于交叉验证双指针写法。
func bruteForceTrap(arr []int) int {
	n := len(arr)
	if n < 3 {
		return 0
	}

	leftMax := make([]int, n)
	rightMax := make([]int, n)
	for i := 0; i < n; i++ {
		leftMax[i] = arr[i]
		if i > 0 && leftMax[i-1] > leftMax[i] {
			leftMax[i] = leftMax[i-1]
		}
	}
	for i := n - 1; i >= 0; i-- {
		rightMax[i] = arr[i]
		if i < n-1 && rightMax[i+1] > rightMax[i] {
			rightMax[i] = rightMax[i+1]
		}
	}

	water := 0
	for i := 0; i < n; i++ {
		bound := leftMax[i]
		if rightMax[i] < bound {
			bound = rightMax[i]
		}
		water += bound - arr[i]
	}
	return water
}

func TestBM094(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want int
	}{
		{"官方示例-普通柱形", []int{3, 1, 2, 5, 2, 4}, 5},
		{"官方示例-经典柱形", []int{0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1}, 6},
		{"官方示例-右侧高墙", []int{4, 2, 0, 3, 2, 5}, 9},
		{"边界-空数组", nil, 0},
		{"边界-只有一根柱子", []int{4}, 0},
		{"边界-两根柱子", []int{1, 2}, 0},
		{"边界-全相同柱形", []int{3, 3, 3, 3}, 0},
		{"边界-全递增柱形", []int{1, 2, 3, 4, 5}, 0},
		{"边界-全递减柱形", []int{5, 4, 3, 2, 1}, 0},
		{"边界-全部为零", []int{0, 0, 0, 0}, 0},
		{"常规-中间凹槽", []int{5, 0, 5}, 5},
		{"常规-多段凹槽", []int{5, 0, 5, 0, 5}, 10},
		{"常规-锯齿柱形", []int{2, 0, 3, 0, 2}, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SolveBM094(tt.in)
			if got != tt.want {
				t.Fatalf("期望 %v，实际 %v", tt.want, got)
			}
		})
	}

	// 固定种子随机柱形，与前后缀最高公式的暴力解交叉验证。
	t.Run("随机交叉验证", func(t *testing.T) {
		rnd := rand.New(rand.NewSource(20240503))
		for round := 0; round < 300; round++ {
			n := rnd.Intn(30)
			arr := make([]int, n)
			for i := range arr {
				arr[i] = rnd.Intn(10)
			}
			want := bruteForceTrap(arr)
			got := SolveBM094(arr)
			if got != want {
				t.Fatalf("第 %d 轮：arr=%v，期望 %d，实际 %d", round, arr, want, got)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// BM095 分糖果问题 · 测试
// ---------------------------------------------------------------------------
// bruteForceCandy 用另一种等价定义求最少糖果数（O(n^2)）：
// 位置 i 的糖果数 = max(左边连续严格递增段长度, 右边连续严格递减段长度)，
// 仅用于与两趟扫描的写法交叉验证。
func bruteForceCandy(arr []int) int {
	total := 0
	for i := range arr {
		// 以 i 结尾的连续严格递增段长度（至少为 1）
		up := 1
		for j := i; j > 0 && arr[j-1] < arr[j]; j-- {
			up++
		}
		// 以 i 开头的连续严格递减段长度（至少为 1）
		down := 1
		for j := i; j < len(arr)-1 && arr[j] > arr[j+1]; j++ {
			down++
		}
		if up > down {
			total += up
		} else {
			total += down
		}
	}
	return total
}

func TestBM095(t *testing.T) {
	// 两组大规模用例：严格递增 / 严格递减，答案都是 1+2+...+n。
	const largeN = 1000
	increasing := make([]int, largeN)
	decreasing := make([]int, largeN)
	same := make([]int, largeN)
	for i := 0; i < largeN; i++ {
		increasing[i] = i
		decreasing[i] = largeN - i
		same[i] = 7
	}
	wantLarge := largeN * (largeN + 1) / 2

	tests := []struct {
		name string
		in   []int
		want int
	}{
		{"官方示例-上升后下降", []int{1, 0, 2}, 5},
		{"官方示例-尾部评分相同", []int{1, 2, 2}, 4},
		{"边界-空数组", nil, 0},
		{"边界-只有一个孩子", []int{5}, 1},
		{"边界-全同评分", []int{1, 1, 1}, 3},
		{"边界-全同评分四个孩子", []int{7, 7, 7, 7}, 4},
		{"评分递增-1到5", []int{1, 2, 3, 4, 5}, 15},
		{"评分递减-5到1", []int{5, 4, 3, 2, 1}, 15},
		{"常规-对称山谷", []int{3, 2, 1, 2, 3}, 11},
		{"常规-对称山峰", []int{1, 2, 3, 2, 1}, 9},
		{"常规-首尾相等中间低", []int{2, 1, 0, 1, 2}, 11},
		{"常规-多段起伏", []int{1, 3, 2, 2, 1}, 7},
		{"大规模-严格递增1000", increasing, wantLarge},
		{"大规模-严格递减1000", decreasing, wantLarge},
		{"大规模-全同评分1000", same, largeN},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SolveBM095(tt.in)
			if got != tt.want {
				t.Fatalf("期望 %v，实际 %v", tt.want, got)
			}
		})
	}

	// 固定种子随机评分，与「取左右单调段较长者」的暴力定义交叉验证。
	t.Run("随机交叉验证", func(t *testing.T) {
		rnd := rand.New(rand.NewSource(20240504))
		for round := 0; round < 300; round++ {
			n := 1 + rnd.Intn(20)
			arr := make([]int, n)
			for i := range arr {
				arr[i] = rnd.Intn(5) // 小值域制造大量相等与相邻重复
			}
			want := bruteForceCandy(arr)
			got := SolveBM095(arr)
			if got != want {
				t.Fatalf("第 %d 轮：arr=%v，期望 %d，实际 %d", round, arr, want, got)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// BM096 主持人调度（二） · 测试
// ---------------------------------------------------------------------------
// bruteForceHosts 暴力求最大重叠活动数：在每个活动的开始时刻统计仍在进行的活动数量，
// 左闭右开区间下「开始时刻仍被占用」的判定是 Start <= s && s < End。仅用于交叉验证。
func bruteForceHosts(activities []Activity) int {
	best := 0
	for _, probe := range activities {
		s := probe.Start
		cnt := 0
		for _, a := range activities {
			if a.Start <= s && s < a.End {
				cnt++
			}
		}
		if cnt > best {
			best = cnt
		}
	}
	return best
}

func TestBM096(t *testing.T) {
	// 大规模串行用例：1000 个首尾相接的区间，理论上一位主持人即可。
	serial := make([]Activity, 0, 1000)
	for i := 0; i < 1000; i++ {
		serial = append(serial, Activity{Start: i, End: i + 1})
	}

	tests := []struct {
		name string
		in   []Activity
		want int
	}{
		{"官方示例-首尾相接可复用", []Activity{{0, 10}, {10, 20}}, 1},
		{"官方示例-存在重叠", []Activity{{0, 10}, {9, 20}}, 2},
		{"边界-空区间列表", nil, 0},
		{"边界-只有一个活动", []Activity{{5, 10}}, 1},
		// 原题数据范围允许 start == end 的退化区间：它不占用任何时间段，
		// 但只要有活动就至少需要 1 位主持人（修复前这些用例会返回 0）。
		{"边界-单个退化区间", []Activity{{5, 5}}, 1},
		{"边界-全部为退化区间", []Activity{{0, 0}, {3, 3}, {7, 7}}, 1},
		{"边界-负时间的退化区间", []Activity{{-3, -3}}, 1},
		{"边界-退化区间接在普通区间尾部", []Activity{{1, 3}, {3, 3}}, 1},
		{"边界-负时间的普通区间", []Activity{{-5, -1}, {-4, 0}}, 2},
		{"串行区间-答案1", []Activity{{1, 2}, {2, 3}, {3, 4}, {4, 5}}, 1},
		{"完全重叠-答案n", []Activity{{1, 5}, {1, 5}, {1, 5}}, 3},
		{"完全相同区间两个", []Activity{{2, 4}, {2, 4}}, 2},
		{"常规-部分重叠", []Activity{{1, 5}, {2, 6}, {6, 8}, {7, 9}}, 2},
		{"常规-嵌套区间", []Activity{{1, 10}, {2, 3}, {4, 5}}, 2},
		{"常规-乱序输入", []Activity{{7, 9}, {1, 5}, {6, 8}, {2, 6}}, 2},
		{"常规-三活动同时重叠", []Activity{{0, 5}, {1, 6}, {2, 7}, {8, 9}}, 3},
		{"大规模-1000个串行区间", serial, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SolveBM096(tt.in)
			if got != tt.want {
				t.Fatalf("期望 %v，实际 %v", tt.want, got)
			}
		})
	}

	// 固定种子随机区间，与暴力最大重叠数交叉验证。
	t.Run("随机交叉验证", func(t *testing.T) {
		rnd := rand.New(rand.NewSource(20240505))
		for round := 0; round < 300; round++ {
			n := rnd.Intn(15)
			activities := make([]Activity, n)
			for i := range activities {
				s := rnd.Intn(20)
				e := s + 1 + rnd.Intn(5) // 保证 Start < End
				activities[i] = Activity{Start: s, End: e}
			}
			want := bruteForceHosts(activities)
			got := SolveBM096(activities)
			if got != want {
				t.Fatalf("第 %d 轮：activities=%v，期望 %d，实际 %d", round, activities, want, got)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// BM097 旋转数组 · 测试
// ---------------------------------------------------------------------------
// TestBM097 覆盖官方示例、m 取模、空/单元素、原地修改性质、随机对照。
func TestBM097(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		m    int
		want []int
	}{
		{"官方示例-右移2位", []int{1, 2, 3, 4, 5, 6}, 2, []int{5, 6, 1, 2, 3, 4}},
		{"m为0", []int{1, 2, 3}, 0, []int{1, 2, 3}},
		{"m等于长度", []int{1, 2, 3}, 3, []int{1, 2, 3}},
		{"m大于长度需取模", []int{1, 2, 3, 4, 5, 6}, 8, []int{5, 6, 1, 2, 3, 4}},
		{"m是长度整数倍", []int{1, 2, 3, 4}, 12, []int{1, 2, 3, 4}},
		{"单元素-取模后为0", []int{9}, 5, []int{9}},
		{"单元素-位移1", []int{9}, 1, []int{9}},
		{"两元素", []int{1, 2}, 1, []int{2, 1}},
		{"空数组", []int{}, 3, []int{}},
		{"nil切片", nil, 1, nil},
		{"常规用例-右移3位", []int{1, 2, 3, 4, 5}, 3, []int{3, 4, 5, 1, 2}},
		{"含重复值", []int{1, 1, 2, 2, 3}, 2, []int{2, 3, 1, 1, 2}},
		{"全部相同", []int{7, 7, 7, 7}, 3, []int{7, 7, 7, 7}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 复制一份用于确认函数是否原地修改（原地修改是本题的约定）。
			in := slices.Clone(tt.nums)
			got := SolveBM097(in, tt.m)
			if !slices.Equal(got, tt.want) {
				t.Fatalf("期望 %v，实际 %v", tt.want, got)
			}
			// 返回的切片必须就是入参本身（同一个底层数组），说明没有额外申请数组。
			if len(got) != len(in) {
				t.Fatalf("期望返回长度 %d，实际 %d", len(in), len(got))
			}
			if len(got) > 0 && &got[0] != &in[0] {
				t.Fatalf("期望原地修改并返回同一个切片，实际返回了新的底层数组")
			}
		})
	}

	// 原地修改语义的显式断言：调用后入参本身必须已经是轮转后的结果。
	t.Run("原地修改入参", func(t *testing.T) {
		in := []int{1, 2, 3, 4, 5, 6}
		SolveBM097(in, 2)
		if want := []int{5, 6, 1, 2, 3, 4}; !slices.Equal(in, want) {
			t.Fatalf("期望入参被原地修改为 %v，实际 %v", want, in)
		}
	})

	// 固定种子的随机对照：用「另开数组」的朴素实现验证三次反转的正确性，可复现。
	t.Run("随机对照", func(t *testing.T) {
		rnd := rand.New(rand.NewSource(20240907))
		for i := 0; i < 500; i++ {
			n := rnd.Intn(20) // 长度 0..19，覆盖空数组
			nums := make([]int, n)
			for j := range nums {
				nums[j] = rnd.Intn(100)
			}
			m := rnd.Intn(50) // 故意让 m 可能远大于 n，覆盖取模分支
			want := rotateNaive(nums, m)
			got := SolveBM097(nums, m)
			if !slices.Equal(got, want) {
				t.Fatalf("第 %d 组（长度 %d、m=%d）：期望 %v，实际 %v", i, n, m, want, got)
			}
		}
	})
}

// rotateNaive 朴素实现：借助辅助数组按 (i+m) % n 放置，仅用于测试对照。
func rotateNaive(nums []int, m int) []int {
	n := len(nums)
	out := make([]int, n)
	if n == 0 {
		return out
	}
	m %= n
	for i := 0; i < n; i++ {
		out[(i+m)%n] = nums[i]
	}
	return out
}

// ---------------------------------------------------------------------------
// BM098 螺旋矩阵 · 测试
// ---------------------------------------------------------------------------
// TestBM098 覆盖官方示例、空矩阵、1x1、单行、单列、非方阵与随机性质校验。
func TestBM098(t *testing.T) {
	tests := []struct {
		name   string
		matrix [][]int
		want   []int
	}{
		{
			"官方示例-3x3",
			[][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			[]int{1, 2, 3, 6, 9, 8, 7, 4, 5},
		},
		{
			"3x4长矩阵",
			[][]int{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}},
			[]int{1, 2, 3, 4, 8, 12, 11, 10, 9, 5, 6, 7},
		},
		{
			"4x3高矩阵",
			[][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10, 11, 12}},
			[]int{1, 2, 3, 6, 9, 12, 11, 10, 7, 4, 5, 8},
		},
		{
			"1x1矩阵",
			[][]int{{7}},
			[]int{7},
		},
		{
			"只有一行",
			[][]int{{1, 2, 3, 4, 5}},
			[]int{1, 2, 3, 4, 5},
		},
		{
			"只有一列",
			[][]int{{1}, {2}, {3}, {4}, {5}},
			[]int{1, 2, 3, 4, 5},
		},
		{
			"2x2矩阵",
			[][]int{{1, 2}, {3, 4}},
			[]int{1, 2, 4, 3},
		},
		{
			"2x3偶数列",
			[][]int{{1, 2, 3}, {4, 5, 6}},
			[]int{1, 2, 3, 6, 5, 4},
		},
		{
			"空矩阵-零行",
			[][]int{},
			[]int{},
		},
		{
			"空矩阵-nil",
			nil,
			[]int{},
		},
		{
			"空矩阵-零列",
			[][]int{{}},
			[]int{},
		},
		{
			"含相同元素",
			[][]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
			[]int{0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SolveBM098(tt.matrix)
			if !slices.Equal(got, tt.want) {
				t.Fatalf("期望 %v，实际 %v", tt.want, got)
			}
		})
	}

	// 随机矩阵：结果必须是全部元素的一个排列（数量与重数都对得上），且首元素为左上角。
	t.Run("随机矩阵-结果覆盖全部元素", func(t *testing.T) {
		rnd := rand.New(rand.NewSource(20240908))
		for round := 0; round < 200; round++ {
			m := 1 + rnd.Intn(6)
			n := 1 + rnd.Intn(6)
			matrix := make([][]int, m)
			var flat []int
			for i := 0; i < m; i++ {
				matrix[i] = make([]int, n)
				for j := 0; j < n; j++ {
					matrix[i][j] = rnd.Intn(50)
					flat = append(flat, matrix[i][j])
				}
			}

			got := SolveBM098(matrix)
			if len(got) != m*n {
				t.Fatalf("%dx%d：期望返回 %d 个元素，实际 %d 个（%v）", m, n, m*n, len(got), got)
			}
			if got[0] != matrix[0][0] {
				t.Fatalf("%dx%d：螺旋顺序的首元素应是左上角 %d，实际 %d", m, n, matrix[0][0], got[0])
			}
			if m > 1 && n > 1 && got[1] != matrix[0][1] {
				t.Fatalf("%dx%d：第二个元素应向右走到 %d，实际 %d", m, n, matrix[0][1], got[1])
			}

			// 元素的多重集合必须一致：排序后逐项比较。
			sortedGot := slices.Clone(got)
			sortedFlat := slices.Clone(flat)
			sort.Ints(sortedGot)
			sort.Ints(sortedFlat)
			if !slices.Equal(sortedGot, sortedFlat) {
				t.Fatalf("%dx%d：结果不是全部元素的一个排列，期望 %v，实际 %v", m, n, sortedFlat, sortedGot)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// BM099 顺时针旋转矩阵 · 测试
// ---------------------------------------------------------------------------
// TestBM099 覆盖官方示例、1x1、空矩阵、原地版本一致性与随机对照。
func TestBM099(t *testing.T) {
	tests := []struct {
		name   string
		matrix [][]int
		want   [][]int
	}{
		{
			"官方示例-3x3",
			[][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			[][]int{{7, 4, 1}, {8, 5, 2}, {9, 6, 3}},
		},
		{
			"官方示例-2x2",
			[][]int{{1, 2}, {3, 4}},
			[][]int{{3, 1}, {4, 2}},
		},
		{
			"1x1矩阵",
			[][]int{{1}},
			[][]int{{1}},
		},
		{
			"4x4矩阵",
			[][]int{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}, {13, 14, 15, 16}},
			[][]int{{13, 9, 5, 1}, {14, 10, 6, 2}, {15, 11, 7, 3}, {16, 12, 8, 4}},
		},
		{
			"空矩阵-零行",
			[][]int{},
			[][]int{},
		},
		{
			"空矩阵-nil",
			nil,
			[][]int{},
		},
		{
			"空矩阵-零列",
			[][]int{{}},
			[][]int{},
		},
		{
			"非方阵-2x3",
			[][]int{{1, 2, 3}, {4, 5, 6}},
			[][]int{{4, 1}, {5, 2}, {6, 3}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 复制一份入参，用于确认 SolveBM099 不修改入参。
			origin := cloneMatrix(tt.matrix)
			got := SolveBM099(tt.matrix)
			if !matrixEqual(got, tt.want) {
				t.Fatalf("期望 %v，实际 %v", tt.want, got)
			}
			if !matrixEqual(tt.matrix, origin) {
				t.Fatalf("SolveBM099 不应修改入参，期望 %v，实际 %v", origin, tt.matrix)
			}
		})
	}

	// 原地版本：只有方阵适用，结果必须与 SolveBM099 一致，并且返回的就是同一个矩阵（无额外数组）。
	t.Run("原地旋转-方阵", func(t *testing.T) {
		squares := [][][]int{
			{{1}},
			{{1, 2}, {3, 4}},
			{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}, {13, 14, 15, 16}},
			{},
		}
		for _, matrix := range squares {
			want := SolveBM099(cloneMatrix(matrix))
			got := SolveBM099InPlace(matrix)
			if !matrixEqual(got, want) {
				t.Fatalf("%dx%d：原地旋转期望 %v，实际 %v", len(matrix), len(matrix), want, got)
			}
			if len(matrix) > 0 && len(got) > 0 && &got[0][0] != &matrix[0][0] {
				t.Fatalf("%dx%d：原地旋转应复用同一块存储，实际换了底层数组", len(matrix), len(matrix))
			}
		}
	})

	// 非方阵无法原地旋转，函数应按约定原样返回，不破坏数据。
	t.Run("原地旋转-非方阵原样返回", func(t *testing.T) {
		matrix := [][]int{{1, 2, 3}, {4, 5, 6}}
		origin := cloneMatrix(matrix)
		got := SolveBM099InPlace(matrix)
		if !matrixEqual(got, origin) {
			t.Fatalf("非方阵应原样返回，期望 %v，实际 %v", origin, got)
		}
	})

	// 随机方阵：与「直接按下标搬迁」的定义式实现对照，并验证旋转 4 次还原。
	t.Run("随机方阵-与定义式实现对照", func(t *testing.T) {
		rnd := rand.New(rand.NewSource(20240909))
		for round := 0; round < 100; round++ {
			n := 1 + rnd.Intn(6)
			matrix := make([][]int, n)
			for i := 0; i < n; i++ {
				matrix[i] = make([]int, n)
				for j := 0; j < n; j++ {
					matrix[i][j] = rnd.Intn(100)
				}
			}

			want := rotateByFormula(matrix)
			if got := SolveBM099(matrix); !matrixEqual(got, want) {
				t.Fatalf("%dx%d：期望 %v，实际 %v", n, n, want, got)
			}

			// 连续旋转 4 次应回到原矩阵。
			four := cloneMatrix(matrix)
			for k := 0; k < 4; k++ {
				four = SolveBM099InPlace(four)
			}
			if !matrixEqual(four, matrix) {
				t.Fatalf("%dx%d：旋转 4 次应还原为 %v，实际 %v", n, n, matrix, four)
			}
		}
	})
}

// rotateByFormula 按定义直接搬迁元素：原第 i 行第 j 列变为新矩阵第 j 行第 n-1-i 列。
func rotateByFormula(matrix [][]int) [][]int {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return [][]int{}
	}
	m, n := len(matrix), len(matrix[0])
	res := make([][]int, n)
	for j := 0; j < n; j++ {
		res[j] = make([]int, m)
		for i := 0; i < m; i++ {
			res[j][m-1-i] = matrix[i][j]
		}
	}
	return res
}

// cloneMatrix 深拷贝矩阵，nil 与空矩阵保持原样（仅测试使用）。
func cloneMatrix(matrix [][]int) [][]int {
	out := make([][]int, len(matrix))
	for i, row := range matrix {
		out[i] = slices.Clone(row)
	}
	return out
}

// matrixEqual 比较两个矩阵的行数、列数与元素；nil 与空矩阵视为相等。
func matrixEqual(a, b [][]int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !slices.Equal(a[i], b[i]) {
			return false
		}
	}
	return true
}

// ---------------------------------------------------------------------------
// BM100 设计LRU缓存结构 · 测试
// ---------------------------------------------------------------------------
// lruStep 一步操作：kind 为 "set" 时忽略 want；为 "get" 时断言返回值等于 want。
type lruStep struct {
	kind  string
	key   int
	value int
	want  int
}

// TestBM100 覆盖官方示例、容量 0/1、更新已有 key、22 步操作序列、随机对照与结构不变量。
func TestBM100(t *testing.T) {
	t.Run("官方示例-容量2", func(t *testing.T) {
		runLRU(t, 2, []lruStep{
			{"set", 1, 1, 0},
			{"set", 2, 2, 0},
			{"get", 1, 0, 1},
			{"set", 3, 3, 0},
			{"get", 2, 0, -1},
			{"set", 4, 4, 0},
			{"get", 1, 0, -1},
			{"get", 3, 0, 3},
			{"get", 4, 0, 4},
		})
	})

	t.Run("容量为1", func(t *testing.T) {
		runLRU(t, 1, []lruStep{
			{"get", 1, 0, -1}, // 空缓存
			{"set", 1, 10, 0},
			{"get", 1, 0, 10},
			{"set", 2, 20, 0}, // 淘汰 1
			{"get", 1, 0, -1},
			{"get", 2, 0, 20},
			{"set", 2, 22, 0}, // 更新已有 key
			{"get", 2, 0, 22},
		})
	})

	t.Run("容量为0", func(t *testing.T) {
		c := NewLRUCache(0)
		c.Set(1, 1)
		c.Set(2, 2)
		if got := c.Get(1); got != -1 {
			t.Fatalf("容量为 0 时 get 期望 -1，实际 %d", got)
		}
		checkLRUInvariants(t, c)
		if len(c.items) != 0 {
			t.Fatalf("容量为 0 时不应存入任何记录，实际 %d 条", len(c.items))
		}
	})

	t.Run("更新已有key不会挤掉其他记录", func(t *testing.T) {
		runLRU(t, 2, []lruStep{
			{"set", 1, 1, 0},
			{"set", 2, 2, 0},
			{"set", 1, 11, 0}, // 更新 1，1 变成最近使用
			{"get", 1, 0, 11},
			{"get", 2, 0, 2}, // 2 仍在
			{"set", 3, 3, 0}, // 淘汰 1（上一步 get(2) 让 2 变新）
			{"get", 1, 0, -1},
			{"get", 2, 0, 2},
			{"get", 3, 0, 3},
		})
	})

	// 22 步连续操作序列（容量 3），逐步断言，覆盖更新、命中前置、淘汰与多次淘汰。
	t.Run("22步操作序列", func(t *testing.T) {
		runLRU(t, 3, []lruStep{
			{"set", 1, 10, 0},
			{"set", 2, 20, 0},
			{"set", 3, 30, 0},
			{"get", 1, 0, 10}, // 1 变最近使用，2 最久未使用
			{"set", 4, 40, 0}, // 淘汰 2
			{"get", 2, 0, -1}, // 已被淘汰
			{"get", 3, 0, 30}, // 3 变最近使用
			{"set", 2, 22, 0}, // 淘汰 1
			{"get", 1, 0, -1}, // 已被淘汰
			{"set", 4, 44, 0}, // 更新 4
			{"get", 4, 0, 44}, // 4 变最近使用
			{"get", 3, 0, 30}, // 3 变最近使用
			{"set", 5, 50, 0}, // 淘汰 2
			{"get", 2, 0, -1}, // 已被淘汰
			{"get", 5, 0, 50}, // 5 命中
			{"set", 3, 33, 0}, // 更新 3
			{"get", 4, 0, 44}, // 4 命中
			{"set", 6, 60, 0}, // 淘汰 5
			{"get", 5, 0, -1}, // 已被淘汰
			{"get", 6, 0, 60}, // 6 命中
			{"get", 3, 0, 33}, // 3 更新后的值
			{"set", 4, 40, 0}, // 更新 4
		})
	})

	// 大量连续操作：随机序列与朴素实现逐步对照，并随时校验内部结构不变量。
	t.Run("随机序列与朴素实现对照", func(t *testing.T) {
		rnd := rand.New(rand.NewSource(20241001))
		for round := 0; round < 200; round++ {
			capacity := rnd.Intn(5) // 0..4，覆盖容量为 0
			c := NewLRUCache(capacity)
			naive := newNaiveLRU(capacity)
			for step := 0; step < 120; step++ {
				key := rnd.Intn(8)
				if rnd.Intn(2) == 0 {
					value := rnd.Intn(1000)
					c.Set(key, value)
					naive.set(key, value)
				} else {
					got, want := c.Get(key), naive.get(key)
					if got != want {
						t.Fatalf("第 %d 轮第 %d 步 get(%d)：期望 %d，实际 %d", round, step, key, want, got)
					}
				}
				checkLRUInvariants(t, c)
			}
		}
	})
}

// runLRU 按步骤驱动 LRU，get 步骤逐步断言返回值，并在每步之后校验内部结构。
func runLRU(t *testing.T, capacity int, steps []lruStep) {
	t.Helper()
	c := NewLRUCache(capacity)
	for i, s := range steps {
		switch s.kind {
		case "set":
			c.Set(s.key, s.value)
		case "get":
			if got := c.Get(s.key); got != s.want {
				t.Fatalf("第 %d 步 get(%d) 期望 %d，实际 %d", i+1, s.key, s.want, got)
			}
		default:
			t.Fatalf("第 %d 步出现未知操作 %q", i+1, s.kind)
		}
		checkLRUInvariants(t, c)
	}
}

// checkLRUInvariants 校验「map 与双向链表」的一致性：长度相同、指针互相指回、key 不重复。
func checkLRUInvariants(t *testing.T, c *LRUCache) {
	t.Helper()
	if len(c.items) > c.capacity {
		t.Fatalf("记录数 %d 超过容量 %d", len(c.items), c.capacity)
	}
	seen := make(map[int]bool, len(c.items))
	count := 0
	for n := c.head.next; n != c.tail; n = n.next {
		if n == nil {
			t.Fatalf("正向遍历时遇到 nil 结点，链表断裂")
		}
		if n.prev.next != n || n.next.prev != n {
			t.Fatalf("结点 key=%d 的前后指针不一致", n.key)
		}
		if got, ok := c.items[n.key]; !ok || got != n {
			t.Fatalf("结点 key=%d 与 map 不一致", n.key)
		}
		if seen[n.key] {
			t.Fatalf("链表里出现重复 key=%d", n.key)
		}
		seen[n.key] = true
		count++
		if count > c.capacity {
			t.Fatalf("链表长度超过容量 %d", c.capacity)
		}
	}
	if count != len(c.items) {
		t.Fatalf("链表长度 %d 与 map 大小 %d 不一致", count, len(c.items))
	}
	back := 0
	for n := c.tail.prev; n != c.head; n = n.prev {
		if n == nil {
			t.Fatalf("反向遍历时遇到 nil 结点，链表断裂")
		}
		back++
		if back > c.capacity {
			t.Fatalf("反向遍历长度超过容量 %d", c.capacity)
		}
	}
	if back != count {
		t.Fatalf("反向遍历长度 %d 与正向 %d 不一致", back, count)
	}
}

// naiveLRUEntry 朴素实现的一条记录：值与最近使用的时间戳。
type naiveLRUEntry struct {
	value  int
	usedAt int
}

// naiveLRU 朴素 LRU：淘汰时线性扫描找最早使用的记录，仅用于测试对照。
type naiveLRU struct {
	capacity int
	tick     int
	entries  map[int]*naiveLRUEntry
}

func newNaiveLRU(capacity int) *naiveLRU {
	return &naiveLRU{capacity: capacity, entries: make(map[int]*naiveLRUEntry)}
}

func (c *naiveLRU) get(key int) int {
	e, ok := c.entries[key]
	if !ok {
		return -1
	}
	c.tick++
	e.usedAt = c.tick
	return e.value
}

func (c *naiveLRU) set(key, value int) {
	if c.capacity <= 0 {
		return
	}
	if e, ok := c.entries[key]; ok {
		e.value = value
		c.tick++
		e.usedAt = c.tick
		return
	}
	if len(c.entries) >= c.capacity {
		victimKey, victim := 0, (*naiveLRUEntry)(nil)
		for k, e := range c.entries {
			if victim == nil || e.usedAt < victim.usedAt {
				victimKey, victim = k, e
			}
		}
		delete(c.entries, victimKey)
	}
	c.tick++
	c.entries[key] = &naiveLRUEntry{value: value, usedAt: c.tick}
}

// ---------------------------------------------------------------------------
// BM101 设计LFU缓存结构 · 测试
// ---------------------------------------------------------------------------
// lfuStep 一步操作：kind 为 "set" 时忽略 want；为 "get" 时断言返回值等于 want。
type lfuStep struct {
	kind  string
	key   int
	value int
	want  int
}

// TestBM101 覆盖牛客官方示例、示例推演、容量 0/1、频次淘汰、同频次淘汰最久未使用、
// 23 步操作序列、随机对照与结构不变量。
func TestBM101(t *testing.T) {
	// 牛客原题的 operators 形式：[1,key,value] 表示 set，[2,key] 表示 get。
	t.Run("官方示例-operators形式", func(t *testing.T) {
		ops := [][]int{
			{1, 1, 1}, {1, 2, 2}, {1, 3, 2}, {1, 2, 4}, {1, 3, 5},
			{2, 2}, {1, 4, 4}, {2, 1},
		}
		want := []int{4, -1} // 官方返回值：第一次 get(2) 得 4，get(1) 时键 1 已被淘汰
		if got := runLFUOps(t, 3, ops); !slices.Equal(got, want) {
			t.Fatalf("期望 %v，实际 %v", want, got)
		}
	})

	// 示例 2：capacity = 2，展示「频次相同则淘汰最久未使用」。
	t.Run("示例2-同频次淘汰最久未使用", func(t *testing.T) {
		runLFU(t, 2, []lfuStep{
			{"set", 1, 1, 0},  // 1 频次 1
			{"set", 2, 2, 0},  // 2 频次 1，晚于 1 被调用
			{"set", 3, 3, 0},  // 两者频次相同，淘汰更久未使用的 1
			{"get", 1, 0, -1}, // 已被淘汰
			{"get", 2, 0, 2},  // 2 频次升为 2
			{"set", 4, 4, 0},  // 3 频次最低，淘汰 3
			{"get", 3, 0, -1}, // 已被淘汰
			{"get", 4, 0, 4},  // 4 频次升为 2
		})
	})

	t.Run("容量为1", func(t *testing.T) {
		runLFU(t, 1, []lfuStep{
			{"get", 1, 0, -1}, // 空缓存
			{"set", 1, 10, 0},
			{"get", 1, 0, 10},
			{"set", 1, 11, 0}, // 更新已有 key，频次 +1
			{"get", 1, 0, 11},
			{"set", 2, 20, 0}, // 唯一的记录 1 被淘汰
			{"get", 1, 0, -1},
			{"get", 2, 0, 20},
		})
	})

	t.Run("容量为0", func(t *testing.T) {
		c := NewLFUCache(0)
		c.Set(1, 1)
		c.Set(2, 2)
		if got := c.Get(1); got != -1 {
			t.Fatalf("容量为 0 时 get 期望 -1，实际 %d", got)
		}
		checkLFUInvariants(t, c)
		if len(c.items) != 0 {
			t.Fatalf("容量为 0 时不应存入任何记录，实际 %d 条", len(c.items))
		}
	})

	// 频次不同时先淘汰频次最低的，即使它更晚被写入。
	t.Run("优先淘汰频次最低", func(t *testing.T) {
		runLFU(t, 3, []lfuStep{
			{"set", 1, 10, 0},
			{"set", 2, 20, 0},
			{"set", 3, 30, 0},
			{"get", 1, 0, 10}, // 1 频次 2
			{"get", 1, 0, 10}, // 1 频次 3
			{"get", 2, 0, 20}, // 2 频次 2
			{"set", 4, 40, 0}, // 3 频次最低（1 次），淘汰 3
			{"get", 3, 0, -1}, // 已被淘汰
			{"get", 1, 0, 10},
			{"get", 2, 0, 20},
			{"get", 4, 0, 40},
		})
	})

	// 23 步连续操作序列（容量 3），逐步断言，覆盖频次提升、更新、同频次淘汰与 minFreq 维护。
	t.Run("23步操作序列", func(t *testing.T) {
		runLFU(t, 3, []lfuStep{
			{"set", 1, 10, 0},
			{"set", 2, 20, 0},
			{"set", 3, 30, 0},
			{"get", 1, 0, 10}, // 1 频次 2
			{"get", 1, 0, 10}, // 1 频次 3
			{"set", 4, 40, 0}, // 2 与 3 频次同为 1 且 2 更久未用，淘汰 2
			{"get", 2, 0, -1}, // 已被淘汰
			{"get", 3, 0, 30}, // 3 频次 2
			{"set", 5, 50, 0}, // 4 频次最低（1 次），淘汰 4
			{"get", 4, 0, -1}, // 已被淘汰
			{"get", 1, 0, 10}, // 1 频次 4
			{"get", 3, 0, 30}, // 3 频次 3
			{"set", 6, 60, 0}, // 5 频次最低（1 次），淘汰 5
			{"get", 5, 0, -1}, // 已被淘汰
			{"get", 6, 0, 60}, // 6 频次 2
			{"set", 1, 11, 0}, // 更新 1，频次 5
			{"get", 1, 0, 11}, // 1 频次 6
			{"get", 3, 0, 30}, // 3 频次 4
			{"set", 7, 70, 0}, // 6 频次最低（2 次），淘汰 6
			{"get", 6, 0, -1}, // 已被淘汰
			{"get", 1, 0, 11},
			{"get", 3, 0, 30},
			{"get", 7, 0, 70},
		})
	})

	// 大量连续操作：随机序列与朴素实现逐步对照，并随时校验桶与 minFreq 的一致性。
	t.Run("随机序列与朴素实现对照", func(t *testing.T) {
		rnd := rand.New(rand.NewSource(20241002))
		for round := 0; round < 200; round++ {
			capacity := rnd.Intn(5) // 0..4，覆盖容量为 0
			c := NewLFUCache(capacity)
			naive := newNaiveLFU(capacity)
			for step := 0; step < 120; step++ {
				key := rnd.Intn(8)
				if rnd.Intn(2) == 0 {
					value := rnd.Intn(1000)
					c.Set(key, value)
					naive.set(key, value)
				} else {
					got, want := c.Get(key), naive.get(key)
					if got != want {
						t.Fatalf("第 %d 轮第 %d 步 get(%d)：期望 %d，实际 %d", round, step, key, want, got)
					}
				}
				checkLFUInvariants(t, c)
			}
		}
	})
}

// runLFUOps 按牛客 operators 格式执行操作：opt=1 表示 set(op[1], op[2])，opt=2 表示 get(op[1])，
// 返回所有 get 的返回值（set 不产生返回值）。
func runLFUOps(t *testing.T, capacity int, ops [][]int) []int {
	t.Helper()
	c := NewLFUCache(capacity)
	var got []int
	for i, op := range ops {
		switch op[0] {
		case 1:
			c.Set(op[1], op[2])
		case 2:
			got = append(got, c.Get(op[1]))
		default:
			t.Fatalf("第 %d 个操作的 opt=%d 非法", i+1, op[0])
		}
		checkLFUInvariants(t, c)
	}
	if got == nil {
		got = []int{}
	}
	return got
}

// runLFU 按步骤驱动 LFU，get 步骤逐步断言返回值，并在每步之后校验内部结构。
func runLFU(t *testing.T, capacity int, steps []lfuStep) {
	t.Helper()
	c := NewLFUCache(capacity)
	for i, s := range steps {
		switch s.kind {
		case "set":
			c.Set(s.key, s.value)
		case "get":
			if got := c.Get(s.key); got != s.want {
				t.Fatalf("第 %d 步 get(%d) 期望 %d，实际 %d", i+1, s.key, s.want, got)
			}
		default:
			t.Fatalf("第 %d 步出现未知操作 %q", i+1, s.kind)
		}
		checkLFUInvariants(t, c)
	}
}

// checkLFUInvariants 校验频次桶与 map 的一致性：结点所在桶的频次等于结点频次、
// 各桶记录的结点总数等于 map 大小、minFreq 等于真实最小频次、不存在空桶。
func checkLFUInvariants(t *testing.T, c *LFUCache) {
	t.Helper()
	if len(c.items) > c.capacity {
		t.Fatalf("记录数 %d 超过容量 %d", len(c.items), c.capacity)
	}
	total := 0
	realMin := math.MaxInt
	for freq, b := range c.buckets {
		if b.size == 0 {
			t.Fatalf("频次 %d 的桶是空桶，应在清空时删除", freq)
		}
		count := 0
		for n := b.head.next; n != b.tail; n = n.next {
			if n == nil {
				t.Fatalf("频次 %d 的桶内链表断裂", freq)
			}
			if n.freq != freq {
				t.Fatalf("结点 key=%d 的频次 %d 与所在桶频次 %d 不一致", n.key, n.freq, freq)
			}
			if n.prev.next != n || n.next.prev != n {
				t.Fatalf("结点 key=%d 的前后指针不一致", n.key)
			}
			if got, ok := c.items[n.key]; !ok || got != n {
				t.Fatalf("结点 key=%d 与 map 不一致", n.key)
			}
			count++
			if count > c.capacity {
				t.Fatalf("频次 %d 的桶长度超过容量 %d", freq, c.capacity)
			}
		}
		if count != b.size {
			t.Fatalf("频次 %d 的桶 size=%d 与实际结点数 %d 不一致", freq, b.size, count)
		}
		total += count
		if freq < realMin {
			realMin = freq
		}
	}
	if total != len(c.items) {
		t.Fatalf("所有桶内的结点总数 %d 与 map 大小 %d 不一致", total, len(c.items))
	}
	if len(c.items) == 0 {
		return
	}
	if c.minFreq != realMin {
		t.Fatalf("minFreq 记录为 %d，实际最小频次为 %d", c.minFreq, realMin)
	}
	if _, ok := c.buckets[c.minFreq]; !ok {
		t.Fatalf("minFreq=%d 对应的频次桶不存在", c.minFreq)
	}
}

// naiveLFUEntry 朴素实现的一条记录：值、调用次数与最近一次调用的时间戳。
type naiveLFUEntry struct {
	value  int
	freq   int
	usedAt int
}

// naiveLFU 朴素 LFU：淘汰时线性扫描找「频次最低、同频次最久未使用」的记录，仅用于测试对照。
type naiveLFU struct {
	capacity int
	tick     int
	entries  map[int]*naiveLFUEntry
}

func newNaiveLFU(capacity int) *naiveLFU {
	return &naiveLFU{capacity: capacity, entries: make(map[int]*naiveLFUEntry)}
}

func (c *naiveLFU) get(key int) int {
	e, ok := c.entries[key]
	if !ok {
		return -1
	}
	c.tick++
	e.freq++
	e.usedAt = c.tick
	return e.value
}

func (c *naiveLFU) set(key, value int) {
	if c.capacity <= 0 {
		return
	}
	if e, ok := c.entries[key]; ok {
		e.value = value
		c.tick++
		e.freq++
		e.usedAt = c.tick
		return
	}
	if len(c.entries) >= c.capacity {
		victimKey, victim := 0, (*naiveLFUEntry)(nil)
		for k, e := range c.entries {
			switch {
			case victim == nil, e.freq < victim.freq,
				e.freq == victim.freq && e.usedAt < victim.usedAt:
				victimKey, victim = k, e
			}
		}
		delete(c.entries, victimKey)
	}
	c.tick++
	c.entries[key] = &naiveLFUEntry{value: value, freq: 1, usedAt: c.tick}
}
