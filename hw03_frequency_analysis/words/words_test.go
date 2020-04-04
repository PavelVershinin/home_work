package words_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/PavelVershinin/home_work/hw03_frequency_analysis/words"
)

func TestCounter(t *testing.T) {
	var testTable = []struct {
		Text          string
		ExpectedList  []string
		ExpectedCount int
	}{
		{
			"Мама мыла раму",
			[]string{"мама", "мыла", "раму"},
			3,
		},
		{
			"Мама - мыла раму!",
			[]string{"мама", "мыла", "раму"},
			3,
		},
		{
			"Я должнa отомстить! Я должнa отомстить... Дa! Я нaпою его в хлaм, отвезу домой, рaздену, уложу нa кровaть, нaкрою одеялом и уйду, остaвив зaписку: «Ты был прэвaсходэн. Цэлую твой, Гиви!»",
			[]string{"я", "должнa", "отомстить"},
			26,
		},
	}
	for i, ts := range testTable {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var counter = &words.Counter{}
			counter.AddText(ts.Text)
			require.Equal(t, ts.ExpectedCount, counter.Count())
			require.Subset(t, ts.ExpectedList, counter.MostCommon(3))
		})
	}
}

func TestCounter_AddWord(t *testing.T) {
	var testTable = []struct {
		Word          string
		ExpectedList  []string
		ExpectedError error
	}{
		{"Мама", []string{"мама"}, nil},
		{"мыла раму", []string{}, words.ErrorMoreThanOneWord},
		{"!", []string{}, words.ErrorBadWord},
		{"", []string{}, words.ErrorBadWord},
		{"-", []string{}, words.ErrorBadWord},
		{"Римский-Корсаков", []string{"римский-корсаков"}, nil},
	}

	for i, ts := range testTable {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var counter = &words.Counter{}
			require.Equal(t, ts.ExpectedError, counter.AddWord(ts.Word))
			require.Equal(t, ts.ExpectedList, counter.MostCommon(-1))
		})
	}
}
