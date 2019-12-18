// Package blog contains blog content and entries
package blog

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/mariomac/goblog/src/fs"
	"github.com/russross/blackfriday/v2"
	nethtml "golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Entry holds the information of a blog entry or page
type Entry struct {
	FileName string
	Title    string
	HTML     template.HTML
	Preview  template.HTML
	Time     *time.Time // may be nil, if it is a page
}

// Content holds the information of all the entries and pages of the blog
type Content struct {
	entries []Entry          // Timestamped entries, sorted from new to old
	all     map[string]Entry // All the pages, mapped by its file name
}

// YYYYMMDDHHMMsome-text_here.md
var entryFormat = regexp.MustCompile(`[0-9]{12}[_\-a-zA-Z0-9]+\.md$`)
var allFormat = regexp.MustCompile(`^[_\-a-zA-Z0-9]+\.md$`)
var allFileFormat = regexp.MustCompile(`[_\-a-zA-Z0-9]+\.md$`)

// GetEntries returns all the entries of the blog.
func (blog *Content) GetEntries() []Entry {
	return blog.entries
}

// Get returned the entry corresponding to the given file name
func (blog *Content) Get(fileName string) (Entry, bool) {
	entry, ok := blog.all[fileName]
	return entry, ok
}

// Load loads all the files in a folder and constructs the entries of the blog.
func (blog *Content) Load(folder string) {
	blog.entries = make([]Entry, 0)
	blog.all = make(map[string]Entry, 0)

	log.Printf("Scanning for entries in folder %s...", folder)
	paths := fs.Search(folder, allFormat)
	for _, path := range paths {

		fileBody, err := ioutil.ReadFile(path)
		if err != nil {
			log.Printf("Error reading file %s: %s", path, err.Error())
			continue
		}

		timestamped := entryFormat.FindString(path)
		title, html, preview := getTitleBodyAndPreview(fileBody)

		var fileName string
		if len(timestamped) > 0 {
			fileName = timestamped
			time := extractTime(fileName)
			log.Printf("Entry found: %s [%s]", path, fileName)
			entry := Entry{
				Time:     &time,
				Title:    title,
				FileName: fileName,
				HTML:     html,
				Preview:  preview,
			}
			blog.entries = append(blog.entries, entry)
			blog.all[fileName] = entry
		} else {
			fileName = allFileFormat.FindString(path)
			blog.all[fileName] = Entry{
				Title:    title,
				FileName: fileName,
				HTML:     html,
				Time:     nil,
				Preview:  preview,
			}
			log.Printf("Page found: %s [%s]", path, fileName)
		}
	}
	// sort entries by time in descending order
	sort.SliceStable(blog.entries[:], func(i, j int) bool {
		return blog.entries[i].Time.After(*blog.entries[j].Time)
	})
}

// TODO: configure by env
var location, _ = time.LoadLocation("Europe/Madrid")

// Extracts a Time from a string beggining with a timestamp in the format YYYYMMDDHHMM...
func extractTime(timestr string) time.Time {
	year, _ := strconv.Atoi(timestr[:4])
	month, _ := strconv.Atoi(timestr[4:6])
	day, _ := strconv.Atoi(timestr[6:8])
	hour, _ := strconv.Atoi(timestr[8:10])
	minute, _ := strconv.Atoi(timestr[10:12])
	parsedTime := time.Date(year, time.Month(month), day, hour, minute, 0, 0, location)
	return parsedTime
}

func getTitleBodyAndPreview(mdBytes []byte) (string, template.HTML, template.HTML) {
	htmlBytes := blackfriday.Run(mdBytes)

	htmlNode, err := nethtml.Parse(bytes.NewReader(htmlBytes))
	if err != nil {
		return err.Error(), "", ""
	}

	firstParagraph := getFirstParagraph(htmlNode)

	h1 := removeFirstH1(htmlNode)
	title, _ := getText(h1)
	log.Printf("Parsed title: %s", title)

	bodyBuf := new(bytes.Buffer)
	nethtml.Render(bodyBuf, htmlNode)
	body := template.HTML(bodyBuf.String())

	var preview template.HTML
	if firstParagraph != nil {
		previewBuf := new(bytes.Buffer)
		nethtml.Render(previewBuf, firstParagraph)
		preview = template.HTML(previewBuf.String())
	} else {
		preview = template.HTML("")
	}

	return title, body, preview
}

// Parameter, parent node. Return type, removed node
func removeFirstH1(parent *nethtml.Node) *nethtml.Node {
	child := parent.FirstChild
	for child != nil {
		if child.DataAtom == atom.H1 {
			parent.RemoveChild(child)
			return child
		}
		removedH1 := removeFirstH1(child)
		if removedH1 != nil {
			return removedH1
		}
		child = child.NextSibling
	}
	return nil
}

func getFirstParagraph(parent *nethtml.Node) *nethtml.Node {
	child := parent.FirstChild
	for child != nil {
		if child.DataAtom == atom.P {
			return child
		}
		paragraph := getFirstParagraph(child)
		if paragraph != nil {
			return paragraph
		}
		child = child.NextSibling
	}
	return nil
}

func getText(parent *nethtml.Node) (string, bool) {
	child := parent.FirstChild
	for child != nil {
		if child.Type == nethtml.TextNode && child.FirstChild == nil {
			return child.Data, true
		}
		text, found := getText(child)
		if found {
			return text, true
		}
		child = child.NextSibling
	}
	return "", false
}
