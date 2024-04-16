package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	flagDir     = flag.Bool("d", false, "Find directories")
	flagFile    = flag.Bool("f", false, "Find files")
	flagSymlink = flag.Bool("sl", false, "Find symbolic links")
	flagExt     = flag.String("ext", "", "Find files with a certain extension")
)

func walkDirFunc(path string, d os.DirEntry, err error) error {
	if err != nil {
		if os.IsPermission(err) {
			return filepath.SkipDir
		}
		// For other types of errors, return them to stop the Walk
		return err
	}

	info, err := d.Info()
	if err != nil {
		fmt.Printf("Error getting info for path %q: %v\n", path, err)
		return err
	}

	if *flagDir && info.IsDir() {
		fmt.Println(path)
	}

	if *flagFile && !info.IsDir() {
		if *flagExt == "" || strings.HasSuffix(info.Name(), "."+*flagExt) {
			fmt.Println(path)
		}
	}

	if *flagSymlink {
		if (info.Mode() & os.ModeSymlink) != 0 {
			realPath, err := filepath.EvalSymlinks(path)
			if err != nil {
				fmt.Printf("%s -> [broken]\n", path)
			} else {
				fmt.Printf("%s -> %s\n", path, realPath)
			}
		}
	}

	return nil
}

func main() {
	flag.Parse()
	root := flag.Arg(0)

	if root == "" {
		fmt.Println("Please provide a path to search.")
		os.Exit(1)
	}

	// Walk the file tree using the separate walkDirFunc
	err := filepath.WalkDir(root, walkDirFunc)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("The path %q does not exist.\n", root)
		} else {
			fmt.Printf("Error walking the path %q: %v\n", root, err)
		}
	}
}