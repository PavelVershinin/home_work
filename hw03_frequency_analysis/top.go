package hw03_frequency_analysis //nolint:golint,stylecheck

import "github.com/PavelVershinin/home_work/hw03_frequency_analysis/words"

func Top10(s string) []string {
	var topNumber = 10
	var counter = &words.Counter{}

	counter.AddText(s)

	return counter.SortedList(topNumber)
}
