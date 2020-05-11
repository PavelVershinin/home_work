package main

import (
	"flag"
	"log"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	if len(from) == 0 {
		log.Fatalln("-from must contains the path to the source file")
	}

	if len(to) == 0 {
		log.Fatalln("-to must contains the path to the destination file")
	}

	if err := Copy(from, to, offset, limit); err != nil {
		log.Fatalln(err)
	}
}
