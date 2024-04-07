package analyzer

import (
	"localhost/search-my-wiser/common"
	"os"
)

func TextAnalyze(text *[]string) []string {
	switch common.EnvParams.TextAnalyzeMethod {
	case common.Ngram:
		return nGram(text)
	default:
		os.Exit(1)
		return nil
	}
}

func nGram(text *[]string) []string {
	terms := []string{}

	for len(*text) > 0 && ignoreWord((*text)[0]) {
		*text = (*text)[1:]
	}

	for i := 0; i < common.EnvParams.Ngram && len(*text) > i && !ignoreWord((*text)[i]); i++ {
		terms = append(terms, (*text)[i])
	}

	if len(*text) > 0 {
		*text = (*text)[1:]
	}

	return terms
}

func ignoreWord(word string) bool {
	//今回は4byte以下しか渡さないので一度しかループしないはず
	for _, char := range word {
		switch char {
		case ' ', '\f', '\n', '\r', '\t', '\v',
			'!', '"', '#', '$', '%', '&',
			'\'', '(', ')', '*', '+', ',',
			'-', '.', '/', ':', ';', '<',
			'=', '>', '?', '@',
			'[', '\\', ']', '^', '_', '`',
			'{', '|', '}', '~',
			'\u3000', // 全角スペース
			'\u3001', // 、
			'\u3002', // 。
			'\uFF08', // （
			'\uFF09': // ）
			return true
		}
	}
	return false
}
