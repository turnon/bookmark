package bookmark

import (
	"testing"
)

const testData = "./test_data.json"

var bookmark *Bookmark

func TestMain(m *testing.M) {
	bookmark = Load(testData)
	m.Run()
}
func TestEntries(t *testing.T) {
	entries := bookmark.Entries()
	if count := len(entries); count != 6 {
		t.Error("no entries !")
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
