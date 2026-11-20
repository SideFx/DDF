/////////////////////////////////////////////////////////////////////////////
// Name:        main.go (DDF: Delete Dot Files)
// Purpose:     Recursively delete dot files
// Author:      Jan Buchholz
// Created:     2026-11-19
// Last update: 2026-11-20
/////////////////////////////////////////////////////////////////////////////

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	var option string
	list := flag.Bool("l", false, "List dot files")
	dele := flag.Bool("d", false, "Delete dot files")
	flag.Parse()
	if (*list && *dele) || (!*list && !*dele) {
		fmt.Println("Usage:")
		fmt.Println("  DDF -l directory")
		fmt.Println("  DDF -d directory")
		os.Exit(1)
	}
	if *list {
		option = "l"
	} else {
		option = "d"
	}
	if flag.NArg() != 1 {
		fmt.Println("Error: Missing directory parameter.")
		os.Exit(1)
	}
	path := flag.Arg(0)
	if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		absPath, _ := filepath.Abs(path)
		process(absPath, option)
	} else {
		fmt.Println("Error: Directory not found: ", path)
		os.Exit(1)
	}
}

func process(path string, option string) {
	fmt.Printf("Processing directory: %s\n", path)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		ls, err := os.Lstat(path)
		if err != nil {
			return err
		}
		mode := ls.Mode()
		if mode&os.ModeSymlink != 0 || mode&os.ModeIrregular != 0 {
			fmt.Printf("Found symlink, skipping: %s\n", path)
			if info.IsDir() {
				return filepath.SkipDir
			}
		} else {
			name := info.Name()
			if len(name) > 0 && name[0] == '.' && !info.IsDir() {
				switch option {
				case "l":
					fmt.Printf("Found dot file: %s\n", path)
				case "d":
					err = os.Remove(path)
					if err != nil {
						switch {
						case os.IsPermission(err):
							fmt.Printf("Permission denied for file %s\n", path)
						case os.IsNotExist(err):
							fmt.Printf("File %s does not exist.\n", path)
						default:
							fmt.Printf("Error deleting file: %s\n", err.Error())
						}
					} else {
						fmt.Printf("File %s successfully deleted.\n", path)
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error occurred: %s\n", err.Error())
	}
}
