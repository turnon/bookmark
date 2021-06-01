package bookmark

import "errors"

type Stat struct {
	Group   string
	Entries []Entry
}

type statMethod struct {
	groupBy func(e *Entry) string
	onlyDup bool
}

var statMethods = map[string]*statMethod{
	"dupName": {
		groupBy: func(e *Entry) string { return e.Name },
		onlyDup: true,
	},
	"dupURL": {
		groupBy: func(e *Entry) string { return e.URL },
		onlyDup: true,
	},
}

func (b *Bookmark) Stat(method string) ([]Stat, error) {
	statMethod, ok := statMethods[method]
	if !ok {
		return nil, errors.New(method + " is not defined")
	}

	// grouping into map
	m := make(map[string][]Entry)
	for _, e := range b.Entries() {
		attr := statMethod.groupBy(&e)
		m[attr] = append(m[attr], e)
	}

	// turn groups into array
	var stats []Stat
	if statMethod.onlyDup {
		stats = []Stat{}
	} else {
		stats = make([]Stat, 0, len(m))
	}
	for group, entries := range m {
		if statMethod.onlyDup && len(entries) <= 1 {
			continue
		}
		st := Stat{Group: group, Entries: entries}
		stats = append(stats, st)
	}

	return stats, nil
}
