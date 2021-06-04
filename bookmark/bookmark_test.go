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

func TestOrder(t *testing.T) {
	stats, err := bookmark.Stat("hosts")
	if err != nil {
		t.Error(err)
	}
	if stats[0].Group != "coolshell.cn" {
		t.Error(stats[0].Group)
	}
}
