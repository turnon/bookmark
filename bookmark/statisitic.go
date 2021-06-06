package bookmark

import (
	"sort"
)

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
	order   int
}

const (
	unorder = iota
	desc
	asc
	atoz
	ztoa
)

var statMethods = map[string]*statMethod{
	"dupName": {
		name:    "名字重复",
		groupBy: func(e *Entry) string { return e.Name },
		onlyDup: true,
		order:   desc,
	},
	"dupURL": {
		name:    "URL重复",
		groupBy: func(e *Entry) string { return e.URL },
		onlyDup: true,
		order:   desc,
	},
	"hosts": {
		name:    "网站统计",
		groupBy: func(e *Entry) string { return e.Host() },
		order:   desc,
	},
	"folders": {
		name:    "目录大小",
		groupBy: func(e *Entry) string { return e.Folder() },
		order:   desc,
	},
	"dirs": {
		name:    "目录列表",
		groupBy: func(e *Entry) string { return e.Folder() },
		order:   atoz,
	},
}

func (stm *statMethod) process(entries []Entry) []Stat {
	// grouping into map
	m := make(map[string][]Entry)
	for _, e := range entries {
		attr := stm.groupBy(&e)
		m[attr] = append(m[attr], e)
	}

	// turn groups into array
	var stats []Stat
	if stm.onlyDup {
		stats = []Stat{}
	} else {
		stats = make([]Stat, 0, len(m))
	}
	for group, entries := range m {
		if stm.onlyDup && len(entries) <= 1 {
			continue
		}
		st := Stat{Group: group, Entries: entries}
		stats = append(stats, st)
	}

	// sort
	stm.sort(stats)

	return stats
}

func (stm *statMethod) sort(stats []Stat) {
	if stm.order == unorder {
		return
	}

	var less func(i int, j int) bool
	switch stm.order {
	case desc:
		less = func(i, j int) bool {
			return len(stats[i].Entries) > len(stats[j].Entries)
		}
	case asc:
		less = func(i, j int) bool {
			return len(stats[i].Entries) < len(stats[j].Entries)
		}
	case atoz:
		less = func(i, j int) bool {
			return stats[i].Group < stats[j].Group
		}
	case ztoa:
		less = func(i, j int) bool {
			return stats[i].Group > stats[j].Group
		}
	}

	sort.Slice(stats, less)
}

func (st *Stat) Count() int {
	return len(st.Entries)
}

func (sts *Stats) Count() int {
	c := 0
	for _, st := range sts.Groups {
		c += st.Count()
	}
	return c
}
