package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"code.cloudfoundry.org/bytefmt"
	. "github.com/logrusorgru/aurora"
)

var maxlvl = flag.Int("level", 3, "")

func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func Process(dir string) error {
	list, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("err reading the directory %s: %s", dir, err)
	}

	pad := strings.Repeat("    ", 0)

	for _, f := range list {
		size := f.Size()
		name := f.Name()

		if f.IsDir() {
			path := path.Join(dir, name)
			size, err = DirSize(path)
			if err != nil {
				return err
			}
			// pretty print
			name = Bold(name).String()
		}

		prettySize := bytefmt.ByteSize(uint64(size))

		fmt.Printf("%s%s %s\n", pad, name, prettySize)
	}

	return nil
}

func main() {
	dir := "." // current directory
	if len(os.Args) == 2 {
		dir = os.Args[1]
	}

	flag.Parse()
	if err := Process(dir); err != nil {
		fmt.Printf("something went wrong: %s", err)
		os.Exit(2)
	}
}
