package bookmark

import (
	"regexp"
	"testing"
	"time"
)

const testData = "./test_data.json"

var bookmark *Bookmark

func TestMain(m *testing.M) {
	bookmark = Load(testData)
	m.Run()
}

func TestUnixTime(t *testing.T) {
	timing := time.Unix((-11644473600 + 13267972139374695/1000000), 0)
	t.Log(timing)
}

func TestEntries(t *testing.T) {
	entries := bookmark.Entries()
	if count := len(entries); count != 7 {
		t.Error(count)
	}
}

func TestDupName(t *testing.T) {
	stats, err := bookmark.Stat("dupName")
	if err != nil {
		t.Error(err)
	}
	if count := len(stats); count != 1 {
		t.Error(count, stats)
	}
}

func TestDupURL(t *testing.T) {
	stats, err := bookmark.Stat("dupURL")
	if err != nil {
		t.Error(err)
	}
	if count := len(stats); count != 1 {
		t.Error(count, stats)
	}
}

func TestHosts(t *testing.T) {
	stats, err := bookmark.Stat("hosts")
	if err != nil {
		t.Error(err)
	}
	if count := len(stats); count != 4 {
		t.Error(count, stats)
	}
}

func TestFolders(t *testing.T) {
	stats, err := bookmark.Stat("folders")
	if err != nil {
		t.Error(err)
	}
	if count := len(stats); count != 2 {
		t.Error(count, stats)
	}
}

func TestOrderByCount(t *testing.T) {
	stats, err := bookmark.Stat("hosts")
	if err != nil {
		t.Error(err)
	}
	if stats[0].Group != "coolshell.cn" {
		t.Error(stats[0].Group)
	}
}

func TestOrderBynName(t *testing.T) {
	stats, err := bookmark.Stat("dirs")
	if err != nil {
		t.Error(err)
	}
	if stats[0].Group != "/amtinfo/collection" || stats[1].Group != "/amtinfo/notes" {
		t.Error(stats[0].Group)
	}
}

func TestNameFilter(t *testing.T) {
	ef := EntryFilter{Name: "壳"}
	re := regexp.MustCompile("壳")
	entries := ef.filter(bookmark.Entries())
	if count := len(entries); count != 3 {
		t.Error(count, entries)
	}
	for _, e := range entries {
		if !re.Match([]byte(e.Name)) {
			t.Error(entries)
		}
	}
}

func TestURLFilter(t *testing.T) {
	ef := EntryFilter{URL: "shell"}
	re := regexp.MustCompile("shell")
	entries := ef.filter(bookmark.Entries())
	if count := len(entries); count != 3 {
		t.Error(count, entries)
	}
	for _, e := range entries {
		if !re.Match([]byte(e.URL)) {
			t.Error(entries)
		}
	}
}

func TestFolderFilter(t *testing.T) {
	ef := EntryFilter{Folder: "notes"}
	re := regexp.MustCompile("notes")
	entries := ef.filter(bookmark.Entries())
	if count := len(entries); count != 5 {
		t.Error(count, entries)
	}
	for _, e := range entries {
		if !re.Match([]byte(e.Folder())) {
			t.Error(entries)
		}
	}
}

func TestFilter(t *testing.T) {
	ef := EntryFilter{Folder: "notes", Name: "壳", URL: "featured"}
	b := bookmark.Filter(&ef)
	if count := len(b.Entries()); count != 1 {
		t.Error(b)
	}
}
