package bookmark

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"regexp"
	"sync"
)

type Bookmark struct {
	Checksum string
	Version  int
	Roots    struct {
		BookmarkBar rawEntry `json:"bookmark_bar"`
		Other       rawEntry
		Synced      rawEntry
	}

	cachedEntries     []Entry
	lockCachedEntries sync.Mutex
}

type common struct {
	DateAdded    string `json:"date_added"`
	DateModified string `json:"date_modified"`
	GUID         string `json:"guid"`
	ID           string `json:"id"`
	Name         string `json:"name"`
}

type rawEntry struct {
	common
	Type     string
	URL      string
	Children []rawEntry
}

type Entry struct {
	common
	URL  string `json:"url"`
	path []string
}

type EntryFilter struct {
	Name   string
	URL    string
	Folder string
}

func (ef *EntryFilter) filter(entries []Entry) []Entry {
	newEntries := []Entry{}

	// generate matchers
	matchers := make([]func(e Entry) bool, 0, 3)
	addMatcher := func(regexpStr string, attrExtractor func(e Entry) string) {
		if regexpStr == "" {
			return
		}
		re := regexp.MustCompile(regexpStr)
		matchers = append(matchers, func(e Entry) bool {
			return re.MatchString(attrExtractor(e))
		})
	}
	addMatcher(ef.Name, func(e Entry) string {
		return e.Name
	})
	addMatcher(ef.URL, func(e Entry) string {
		return e.URL
	})
	addMatcher(ef.Folder, func(e Entry) string {
		return e.Folder()
	})

	if len(matchers) == 0 {
		return append(newEntries, entries...)
	}

	// perform matching
	for _, e := range entries {
		matched := true
		for _, match := range matchers {
			if !match(e) {
				matched = false
				break
			}
		}
		if matched {
			newEntries = append(newEntries, e)
		}
	}

	return newEntries
}

func (entry *Entry) ToJson() string {
	bytes, err := json.Marshal(entry)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func (e *Entry) Host() string {
	u, err := url.Parse(e.URL)
	if err != nil {
		panic(err)
	}
	return u.Host
}

func (e *Entry) Folder() string {
	return filepath.Join(e.path...)
}

func Load(path string) *Bookmark {
	js, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	b := &Bookmark{}
	json.Unmarshal(js, b)
	return b
}

func (b *Bookmark) roots() map[string]rawEntry {
	return map[string]rawEntry{
		"/":      b.Roots.BookmarkBar,
		"other":  b.Roots.Other,
		"synced": b.Roots.Synced,
	}
}

func (b *Bookmark) Entries() []Entry {
	if b.cachedEntries == nil {
		b.lockCachedEntries.Lock()
		defer b.lockCachedEntries.Unlock()

		if b.cachedEntries == nil {
			es := []Entry{}
			for name, root := range b.roots() {
				es = collectEntries([]string{name}, root.Children, es)
			}
			b.cachedEntries = es
		}
	}

	return b.cachedEntries
}

func collectEntries(path []string, res []rawEntry, es []Entry) []Entry {
	for _, re := range res {
		if re.Type == "folder" {
			newPath := append(path, re.Name)
			es = collectEntries(newPath, re.Children, es)
			continue
		}

		e := Entry{re.common, re.URL, path}
		es = append(es, e)
	}

	return es
}

func (b *Bookmark) Stats() ([]Stats, error) {
	result := make([]Stats, 0, len(statMethods))
	for methodName, method := range statMethods {
		stats, err := b.Stat(methodName)
		if err != nil {
			return nil, err
		}
		result = append(result, Stats{Name: methodName, Label: method.name, Groups: stats})
	}
	return result, nil
}

func (b *Bookmark) Stat(method string) ([]Stat, error) {
	statMethod, ok := statMethods[method]
	if !ok {
		return nil, errors.New(method + " is not defined")
	}

	stats := statMethod.process(b.Entries())
	return stats, nil
}

func (b *Bookmark) Filter(ef *EntryFilter) *Bookmark {
	return &Bookmark{
		cachedEntries: ef.filter(b.Entries()),
	}
}
