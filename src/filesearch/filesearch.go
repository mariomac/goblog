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
	nethtml "golang.org/x/net/html"
	"bytes"
)

const DATE_FORMAT = "Mon Jan _2 15:04:05 2006"

// TODO: use pointers & references
type Entry struct {
	Time string
	UrlPath string
	FileLocation string
	Title string
	Html template.HTML
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
			mdBody := github_flavored_markdown.Markdown(fileBody)

			html := template.HTML(mdBody)

			h1tokenizer := nethtml.NewTokenizerFragment(bytes.NewReader(mdBody), "h1")

			h1tokenizer.Next()
			token := h1tokenizer.Token()
			for token.Type != nethtml.StartTagToken && token.Data != "h1" {
				h1tokenizer.Next()
				token = h1tokenizer.Token()
			}

			h1tokenizer.Next()
			token = h1tokenizer.Token()
			for token.Type != nethtml.TextToken {
				h1tokenizer.Next()
				token = h1tokenizer.Token()
			}

			title := token.Data

			entry := &Entry{
				Time:fileNameToTime(fileName).Format("Jan _2, 2006 at 15:04:05"),
				UrlPath:fileName[:(len(fileName)-3)], // remove .md
				FileLocation:path,
				Title:title,
				Html:html,
			}
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


