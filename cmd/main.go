package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"text/tabwriter"

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

func Process(dir string, w io.Writer) error {
	list, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("err reading the directory %s: %s", dir, err)
	}

	var dirs, files []os.FileInfo

	for _, f := range list {
		if f.IsDir() {
			dirs = append(dirs, f)
		} else {
			files = append(files, f)
		}
	}

	// Print dirs
	for _, f := range dirs {
		name := f.Name()
		path := path.Join(dir, name)
		size, err := DirSize(path)
		if err != nil {
			return err
		}
		prettySize := bytefmt.ByteSize(uint64(size))
		fmt.Fprintf(w, "%s\t%s\t\n", Bold(name), prettySize)
	}

	// Print files
	for _, f := range files {
		size := f.Size()
		name := f.Name()

		prettySize := bytefmt.ByteSize(uint64(size))
		fmt.Fprintf(w, "%s\t%s\t\n", name, prettySize)
	}

	return nil
}

func main() {
	dir := "." // current directory
	if len(os.Args) == 2 {
		dir = os.Args[1]
	}

	flag.Parse()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	if err := Process(dir, w); err != nil {
		fmt.Printf("something went wrong: %s", err)
		os.Exit(2)
	}
	w.Flush()
}
