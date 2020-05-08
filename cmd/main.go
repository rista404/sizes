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

		p := path.Join(dir, n)
		s, err := DirSize(p)
		if err != nil {
			return err
		}
		itm.size = uint64(s)

		items = append(items, itm)
	}

	sort.Sort(byAlpha(items))
	sort.Sort(bySize(items))

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
	dir := "." // current directory
	if len(os.Args) == 2 {
		dir = os.Args[1]
	}

	flag.Parse()
	tw := ansiterm.NewTabWriter(os.Stdout, 0, 0, 1, ' ', 0)
	if err := Process(dir, tw); err != nil {
		fmt.Printf("something went wrong: %s", err)
		os.Exit(2)
	}
	tw.Flush()
}
