package main

import (
	"bufio"
	"os"
	"path/filepath"
	"sort"
	"time"
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

func diffPath(a, b []string) ([]string, []string) {

	mapA := make(map[string]struct{})
	for _, p := range a {
		mapA[p] = struct{}{}
	}

	mapB := make(map[string]struct{})
	for _, p := range b {
		mapB[p] = struct{}{}
	}

	var onlyInA, onlyInB []string
	for p := range mapA {
		if _, ok := mapB[p]; !ok {
			onlyInA = append(onlyInA, p)
		}
	}

	for p := range mapB {
		if _, ok := mapA[p]; !ok {
			onlyInB = append(onlyInB, p)
		}
	}

	return onlyInA, onlyInB

}

func displayPath(paths []string) []string {
	sort.Strings(paths)

	return paths
}

func main() {
	filename := time.Now().Format("2006-01-02_15-04-05") + ".txt"

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var ignoreDirs = []string{".git", ".repo"}

	arg1 := os.Args[1]
	arg2 := os.Args[2]

	a, err := getPaths(arg1, ignoreDirs)
	if err != nil {
		panic(err)
	}

	b, err := getPaths(arg2, ignoreDirs)
	if err != nil {
		panic(err)
	}

	onlyInA, onlyInB := diffPath(a, b)

	writer := bufio.NewWriter(file)

	writer.WriteString("Only in " + arg1 + "\n")

	pathsA := displayPath(onlyInA)
	for _, path := range pathsA {
		writer.WriteString(path + "\n")
	}

	writer.WriteString("----\n")

	writer.WriteString("Only in " + arg2 + "\n")

	pathsB := displayPath(onlyInB)
	for _, path := range pathsB {
		writer.WriteString(path + "\n")
	}

	writer.Flush()
}
