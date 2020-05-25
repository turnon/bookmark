package bookmark

func (b *Bookmark) DupName() map[string]int {
	m := groupBy(b.Entries(), func(e *Entry) string {
		return e.Name
	})
	return onlyCount(onlyDup(m))
}

func (b *Bookmark) DupURL() map[string][]Entry {
	m := groupBy(b.Entries(), func(e *Entry) string {
		return e.URL
	})
	return onlyDup(m)
}

func (b *Bookmark) Hosts() map[string]int {
	m := groupBy(b.Entries(), func(e *Entry) string {
		return e.Host()
	})
	return onlyCount(m)
}

func (b *Bookmark) Folders() map[string]int {
	m := groupBy(b.Entries(), func(e *Entry) string {
		return e.Folder()
	})
	return onlyCount(m)
}

func onlyDup(m map[string][]Entry) map[string][]Entry {
	new_m := make(map[string][]Entry)
	for attr, es := range m {
		if len(es) == 1 {
			continue
		}
		new_m[attr] = es
	}
	return new_m
}

func onlyCount(m map[string][]Entry) map[string]int {
	new_m := make(map[string]int)
	for attr, es := range m {
		new_m[attr] = len(es)
	}
	return new_m
}

func groupBy(es []Entry, fn func(*Entry) string) map[string][]Entry {
	m := make(map[string][]Entry)
	for _, e := range es {
		attr := fn(&e)
		m[attr] = append(m[attr], e)
	}
	return m
}
