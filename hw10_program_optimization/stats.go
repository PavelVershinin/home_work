package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"bytes"
	"io"

	"github.com/valyala/fastjson"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	var parser fastjson.Parser

	scanner := bufio.NewScanner(r)
	result := make(DomainStat)
	bDomain := []byte("." + domain)
	bDog := []byte("@")

	for scanner.Scan() {
		v, err := parser.ParseBytes(scanner.Bytes())
		if err != nil {
			return nil, err
		}
		email := v.GetStringBytes("Email")
		if bytes.Contains(email, bDog) {
			email = bytes.ToLower(email)
			if bytes.HasSuffix(email, bDomain) {
				result[string(bytes.SplitN(email, bDog, 2)[1])]++
			}
		}
	}

	return result, nil
}
