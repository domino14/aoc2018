package main

import (
	"bufio"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

func init() {
	addSolutions(4, (*problemContext).problem4)
}

func (ctx *problemContext) problem4() {
	var records []*guardRecord
	scanner := bufio.NewScanner(ctx.f)
	for scanner.Scan() {
		gr, err := parseGuardRecord(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		records = append(records, gr)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	ctx.reportLoad()
	sort.Slice(records, func(i, j int) bool {
		return records[i].t.Before(records[j].t)
	})

	sleep := make(map[int]*sleepRecord)

	var sr *sleepRecord
	var start time.Time
	for _, gr := range records {
		switch gr.action {
		case actionBeginShift:
			var ok bool
			sr, ok = sleep[gr.guardID]
			if !ok {
				sr = new(sleepRecord)
				sleep[gr.guardID] = sr
			}
		case actionFallAsleep:
			if sr == nil {
				log.Fatal("fall asleep before first shift")
			}
			if !start.IsZero() {
				log.Fatal("fall asleep while already asleep")
			}
			start = gr.t
		case actionWakeUp:
			if sr == nil {
				log.Fatal("wake up before first shift")
			}
			if start.IsZero() {
				log.Fatal("wake up while already awake")
			}
			sr.countAsleep(start, gr.t)
			start = time.Time{}
		}
	}

	var part1 int
	var maxSleeper *sleepRecord
	for id, sr := range sleep {
		if maxSleeper == nil || sr.total > maxSleeper.total {
			var maxMin int
			var maxCount int64
			for m, n := range sr.minutes {
				if n > maxCount {
					maxMin = m
					maxCount = n
				}
			}
			part1 = id * maxMin
			maxSleeper = sr
		}
	}
	ctx.reportPart1(part1)

	var part2 int
	var max int64
	for id, sr := range sleep {
		for m, n := range sr.minutes {
			if n > max {
				max = n
				part2 = id * m
			}
		}
	}
	ctx.reportPart2(part2)
}

type guardAction int

const (
	actionBeginShift guardAction = iota
	actionFallAsleep
	actionWakeUp
)

type guardRecord struct {
	t       time.Time
	action  guardAction
	guardID int // only for actionBeginShift
}

func parseGuardRecord(s string) (*guardRecord, error) {
	parts := strings.Split(s, "] ")
	timestamp := strings.TrimPrefix(parts[0], "[")
	t, err := time.Parse("2006-01-02 15:04", timestamp)
	if err != nil {
		return nil, err
	}

	gr := &guardRecord{t: t}
	switch parts[1] {
	case "falls asleep":
		gr.action = actionFallAsleep
	case "wakes up":
		gr.action = actionWakeUp
	default:
		id := strings.TrimPrefix(parts[1], "Guard #")
		id = strings.TrimSuffix(id, " begins shift")
		gr.guardID, err = strconv.Atoi(id)
		if err != nil {
			return nil, err
		}
	}
	return gr, nil
}

type sleepRecord struct {
	minutes [60]int64
	total   int64
}

func (sr *sleepRecord) countAsleep(start, end time.Time) {
	m0 := start.Minute()
	m1 := end.Minute()
	if m1 <= m0 {
		panic(fmt.Sprintf("countAsleep with start=%s; end=%s", start, end))
	}
	for m := m0; m < m1; m++ {
		sr.minutes[m]++
		sr.total++
	}
}
