package wiki

import (
	"os"

	"github.com/dustin/go-wikiparse"
)

// https://meta.wikimedia.org/wiki/Data_dumps
func LoadWikiDump(file *os.File, maxIndexCount int) (wikiparse.Parser, error) {
	return wikiparse.NewParser(file)
}
