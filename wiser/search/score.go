package search

import (
	"fmt"
	"localhost/search-my-wiser/common"
	"math"
	"os"
)

func Score(searchCandidates searchCandidates) float64 {
	switch common.EnvParams.ScoringMethod {
	case common.TfIdf:
		return tfIdf(searchCandidates)
	case common.Bm25:
		return bm25(searchCandidates)
	default:
		return 0.0
	}
}

func tfIdf(searchCandidates searchCandidates) float64 {
	score := 0.0
	for _, searchCandidate := range searchCandidates {
		postings := searchCandidate.searchCursors[searchCandidate.current]
		tf := (float64(postings.PostingsCount) / float64(searchCandidate.postingsCount))            // 文書中の特定の単語の出現数/文書中の単語総数
		idf := math.Log2(float64(common.EnvParams.IndexCount) / float64(searchCandidate.docsCount)) // 総文書数/検索対象の単語が含まれる文書数
		score += float64(tf) * idf                                                                  // Convert the result back to int
	}

	return score
}

func bm25(_ searchCandidates) float64 {
	fmt.Printf("undefined function")
	os.Exit(1)

	return 0.0
}
