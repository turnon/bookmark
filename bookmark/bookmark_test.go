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
	if count := len(entries); count != 5 {
		t.Error("no entries !")
	}
}
