package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
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

func displayPath(paths []string) {
	sort.Strings(paths)
	for _, path := range paths {
		fmt.Println(path)
	}
}

func main() {
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

	fmt.Println("Only in " + arg1)
	displayPath(onlyInA)

	fmt.Println("----")

	fmt.Println("Only in " + arg2)
	displayPath(onlyInB)
}
