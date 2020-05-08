package main

// Sort by name
// dirs first

type byAlpha []*item

func (a byAlpha) Len() int      { return len(a) }
func (a byAlpha) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byAlpha) Less(i, j int) bool {
	if a[i].isDir {
		// both dir
		if a[j].isDir {
			return a[i].name < a[j].name
		}
		// only i dir
		return true
	}
	// only j dir
	if a[j].isDir {
		return false
	}

	return a[i].name < a[j].name
}

// Sort by size

type bySize []*item

func (a bySize) Len() int      { return len(a) }
func (a bySize) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a bySize) Less(i, j int) bool {
	return a[j].size < a[i].size
}
