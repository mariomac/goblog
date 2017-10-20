package files

import (
	//"path/filepath"
	"time"
)

type File struct {
	Path string
}

type Entry struct {
	File
	date time.Time
}

// entries -> the entries per page
// page -> the page number
func ScanEntries(folder string, entries int, page int) []Entry {

}
