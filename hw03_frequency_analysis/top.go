package hw03_frequency_analysis //nolint:golint,stylecheck

import (
	"regexp"
	"sort"
)

const topSize int = 10

type wordStats struct {
	word  string
	count int
}

func Top10(s string) []string {
	if len(s) == 0 {
		return nil
	}

	chart := map[string]int{}

	for _, w := range regexp.MustCompile(`[\s\t\r\n]+`).Split(s, -1) {
		chart[w]++
	}

	chartStats := []wordStats{}

	for w, c := range chart {
		chartStats = append(chartStats, wordStats{w, c})
	}

	sort.Slice(chartStats, func(i, j int) bool {
		return chartStats[i].count > chartStats[j].count
	})

	res := []string{}

	limit := topSize
	if limit > len(chartStats) {
		limit = len(chartStats)
	}

	for _, cs := range chartStats[:limit] {
		res = append(res, cs.word)
	}

	return res
}
