package filesearch

import (
	"path/filepath"
	"os"
	"log"
	"regexp"
	"time"
	"strconv"
	"html/template"
	"github.com/shurcooL/github_flavored_markdown"
	"io/ioutil"
	"golang.org/x/net/html/atom"
)

// TODO: use pointers & references
type Entry struct {
	Time time.Time
	UrlPath string
	FileLocation string
	Title string
	Html template.HTML
}

var validTemplate = regexp.MustCompile(".*\\.html$")

func GetTemplates(folder string) []string {
	files := make([]string, 0, 8)
	log.Printf("Looking for templates in folder %s\n", folder)
	filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && validTemplate.MatchString(info.Name()) {
			log.Printf("found %s [%s]", path, info.Name())
			files = append(files, path)
		}
		return nil
	})
	return files
}

// YYYYMMDDHHMMsome-text_here.md
var validTimestampedEntry = regexp.MustCompile("^[0-9]{12}[_\\-a-zA-Z0-9]+\\.md$")

// todo: paginate
func ScanTimestampedEntries(folder string) []*Entry {
	entries := make([]*Entry, 0, 16)

	log.Printf("Scanning timestamped entries in folder %s\n", folder)
	filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		fileName := info.Name()
		if !info.IsDir() && validTimestampedEntry.MatchString(fileName) {
			fileBody, _ := ioutil.ReadFile(path)

			html := template.HTML(github_flavored_markdown.Markdown(fileBody))

			entry := &Entry{
				Time:fileNameToTime(fileName),
				UrlPath:fileName[:(len(fileName)-3)], // remove .md
				FileLocation:path,
				Title:title,
				Html:html,
			}
			log.Printf("found %s [%s] %s", entry.FileLocation, entry.UrlPath, entry.Time)
			entries = append(entries, entry)
		}
		return nil
	})

	return entries
}




// TODO: configure by env
var location, _ = time.LoadLocation("Europe/Madrid")

func fileNameToTime(filename string) time.Time {
	year, _ := strconv.Atoi(filename[:4])
	month, _ := strconv.Atoi(filename[4:6])
	day, _ := strconv.Atoi(filename[6:8])
	hour, _ := strconv.Atoi(filename[8:10])
	minute, _ := strconv.Atoi(filename[10:12])
	parsedTime := time.Date(year, time.Month(month), day, hour, minute, 0, 0, location)
	return parsedTime
}


