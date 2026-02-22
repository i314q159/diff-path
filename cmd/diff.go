package cmd

import (
	"bufio"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

func getPaths(root string, ignore []string) ([]string, error) {
	var paths []string
	absRoot, _ := filepath.Abs(root)

	err := filepath.Walk(absRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		for _, d := range ignore {
			if info.IsDir() && info.Name() == d {
				return filepath.SkipDir
			}
		}

		relPath, _ := filepath.Rel(absRoot, path)
		if info.IsDir() {
			relPath += string(os.PathSeparator)
		}

		paths = append(paths, relPath)
		return nil
	})

	return paths, err
}

// 双指针法
func diffPathSorted(a, b []string) ([]string, []string) {
	sort.Strings(a)
	sort.Strings(b)

	var onlyInA, onlyInB []string
	i, j := 0, 0

	for i < len(a) && j < len(b) {
		cmp := strings.Compare(a[i], b[j])
		switch {
		case cmp < 0:
			onlyInA = append(onlyInA, a[i])
			i++
		case cmp > 0:
			onlyInB = append(onlyInB, b[j])
			j++
		default:
			i++
			j++
		}
	}

	// 处理剩余元素
	for ; i < len(a); i++ {
		onlyInA = append(onlyInA, a[i])
	}
	for ; j < len(b); j++ {
		onlyInB = append(onlyInB, b[j])
	}

	return onlyInA, onlyInB
}

var DiffCmd = &cobra.Command{
	Use:   "diff [arg1] [arg2]",
	Short: "对比2个AOSP目录差异",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		filename := time.Now().Format("2006-01-02_15-04-05") + ".txt"

		file, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		var ignoreDirs = []string{".git", ".repo"}

		arg1 := args[0]
		arg2 := args[1]

		// 并发获取两个目录的文件列表
		var a, b []string
		var err1, err2 error
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			a, err1 = getPaths(arg1, ignoreDirs)
		}()

		go func() {
			defer wg.Done()
			b, err2 = getPaths(arg2, ignoreDirs)
		}()

		// 等待两个遍历任务完成
		wg.Wait()

		if err1 != nil {
			panic(err1)
		}

		if err2 != nil {
			panic(err2)
		}

		onlyInA, onlyInB := diffPathSorted(a, b)

		writer := bufio.NewWriter(file)

		writer.WriteString("Only in " + arg1 + "\n")

		for _, path := range onlyInA {
			writer.WriteString(path + "\n")
		}

		writer.WriteString(strings.Repeat("-", 100) + "\n")

		writer.WriteString("Only in " + arg2 + "\n")

		for _, path := range onlyInB {
			writer.WriteString(path + "\n")
		}

		writer.Flush()
	},
}
