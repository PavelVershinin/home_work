package words

import (
	"bytes"
	"errors"
	"sort"
	"strings"
	"sync"
	"unicode"
)

var (
	ErrorMoreThanOneWord = errors.New("more than one word")
	ErrorBadWord         = errors.New("bad word")
)

type Counter struct {
	mu    sync.Mutex
	words map[string]int
}

// Count Вернёт количество слов
func (c *Counter) Count() int {
	c.mu.Lock()
	cnt := len(c.words)
	c.mu.Unlock()
	return cnt
}

// MostCommon Вернёт список слов, отсортированный по частотности от большего к меньшему
func (c *Counter) MostCommon(topNumber int) []string {
	c.mu.Lock()
	list := make([]string, 0, len(c.words))
	for w := range c.words {
		list = append(list, w)
	}
	sort.Slice(list, func(i, j int) bool {
		return c.words[list[i]] > c.words[list[j]]
	})
	c.mu.Unlock()

	if topNumber < 0 || len(list) <= topNumber {
		return list
	}

	return list[:topNumber]
}

// AddText Добавит текст, ошибки не возвращаются, в тексте вполне могут быть проблемные слова и это нормально
func (c *Counter) AddText(s string) {
	for _, word := range strings.Fields(s) {
		_ = c.AddWord(word)
	}
}

// AddWord Добавит одно слово. Вернёт ошибку, если слово неправильное, например пустая строка или более одного слова
func (c *Counter) AddWord(word string) error {
	word, err := clear(word)
	if err != nil {
		return err
	}
	c.mu.Lock()
	if c.words == nil {
		c.words = make(map[string]int)
	}
	c.words[word]++
	c.mu.Unlock()
	return nil
}

// clear Вернёт слово очищенное от лишних символов в нижнем регистре
func clear(word string) (string, error) {
	buff := bytes.Buffer{}

	for _, r := range strings.ToLower(word) {
		if r == '-' || unicode.IsLetter(r) {
			buff.WriteRune(r)
		} else if unicode.IsSpace(r) {
			return "", ErrorMoreThanOneWord
		}
	}

	if buff.Len() == 0 || buff.String() == "-" {
		return "", ErrorBadWord
	}

	return buff.String(), nil
}
