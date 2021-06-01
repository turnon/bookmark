package bookmark

import "errors"

type Stat struct {
	Group   string
	Entries []Entry
}

type Stats struct {
	Name   string
	Label  string
	Groups []Stat
}

type statMethod struct {
	name    string
	groupBy func(e *Entry) string
	onlyDup bool
}

var statMethods = map[string]*statMethod{
	"dupName": {
		name:    "名字重复",
		groupBy: func(e *Entry) string { return e.Name },
		onlyDup: true,
	},
	"dupURL": {
		name:    "URL重复",
		groupBy: func(e *Entry) string { return e.URL },
		onlyDup: true,
	},
	"hosts": {
		name:    "网站统计",
		groupBy: func(e *Entry) string { return e.Host() },
	},
	"folders": {
		name:    "目录统计",
		groupBy: func(e *Entry) string { return e.Folder() },
	},
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
