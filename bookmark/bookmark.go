package bookmark

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"path/filepath"
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

func (entry *Entry) ToJson() string {
	bytes, err := json.Marshal(entry)
	if err != nil {
		panic(err)
	}
	return string(bytes)
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
