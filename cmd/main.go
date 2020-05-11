package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"

	"code.cloudfoundry.org/bytefmt"
	"github.com/juju/ansiterm"
	. "github.com/logrusorgru/aurora"
)

var sortBySize = flag.Bool("s", false, "Sort items by their size")

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

type item struct {
	name  string
	size  uint64
	isDir bool
}

func Process(dir string, w io.Writer) error {
	list, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("err reading the directory %s: %s", dir, err)
	}

	var items []*item

	// collect
	for _, f := range list {
		n := f.Name()
		itm := &item{
			name:  n,
			isDir: f.IsDir(),
		}

		if itm.isDir {
			p := path.Join(dir, n)
			s, err := DirSize(p)
			if err != nil {
				return err
			}
			itm.size = uint64(s)
		} else {
			itm.size = uint64(f.Size())
		}

		items = append(items, itm)
	}

	// sort
	if *sortBySize {
		sort.Sort(bySize(items))
	} else {
		sort.Sort(byAlpha(items))
	}

	// print
	for _, i := range items {
		prettySize := bytefmt.ByteSize(i.size)
		n := i.name
		if i.isDir {
			n = Bold(n).String()
		}
		fmt.Fprintf(w, "%s\t%s\t\n", n, prettySize)
	}
	return nil
}

func main() {
	flag.Parse()

	dir := "." // current directory
	args := flag.Args()
	if len(args) == 1 {
		dir = args[0]
	}

	tw := ansiterm.NewTabWriter(os.Stdout, 0, 0, 1, ' ', 0)
	if err := Process(dir, tw); err != nil {
		fmt.Printf("something went wrong: %s", err)
		os.Exit(2)
	}
	tw.Flush()
}
