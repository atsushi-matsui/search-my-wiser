package search

import (
	"fmt"
	"localhost/search-my-wiser/database"
	"localhost/search-my-wiser/inverseindex"
	"sort"
)

type searchCandidate struct {
	tokenId       int
	docsCount     int
	postingsCount int
	base          int
	current       int
	searchCursors inverseindex.PostingsList
}

type searchCandidates []searchCandidate

type phraseCandidate struct {
	tokenId       int
	base          int
	postings      []int
	current       int
	postingsCount int
}

type phraseCandidates []phraseCandidate

type SearchResult struct {
	DocumentId int
	Title      string
	Score      float64
}
type SearchResults []SearchResult

func Search(query string) (*SearchResults, error) {
	fmt.Println("start search.")
	ii := &inverseindex.InverseIndexList{}
	// DBからトークンに紐づくポスティングリストを取得
	err := inverseindex.QueryToPostingList(query, ii)
	if err != nil {
		return nil, err
	}
	if len(*ii) == 0 {
		return nil, fmt.Errorf("no search result. query: %s", query)
	}

	// ドキュメント内でのトークンの出現回数で昇順にソート
	var searchCandidates searchCandidates
	for i, iiVal := range *ii {
		if iiVal == nil {
			continue
		}
		searchCandidates = append(searchCandidates, searchCandidate{
			tokenId:       iiVal.TokenId,
			docsCount:     iiVal.DocsCount,
			postingsCount: iiVal.PostingsCount,
			base:          i,
			current:       0,
			searchCursors: *iiVal.PostingsList,
		})
	}
	sort.Slice(searchCandidates, func(i, j int) bool {
		return searchCandidates[i].docsCount < searchCandidates[j].docsCount
	})

	return searchDocs(searchCandidates)
}

func searchDocs(searchCandidates searchCandidates) (*SearchResults, error) {
	res := &SearchResults{}
	// 最小トークンの現在位置
	baseCursor := 0
	for baseCursor < searchCandidates[0].docsCount {
		// 最小トークンの現在位置でのドキュメントID
		baseDocId := searchCandidates[0].searchCursors[baseCursor].DocId
		// 次候補のドキュメントID
		nextDocId := 0
		// 最小トークンのドキュメントIDを基準に、他のトークンのポスティングリスト内に同じドキュメントIDがあるか確認する
		for i := 1; i < len(searchCandidates); i++ {
			for searchCandidates[i].docsCount > searchCandidates[i].current && baseDocId > searchCandidates[i].searchCursors[searchCandidates[i].current].DocId {
				searchCandidates[i].current++
			}
			if searchCandidates[i].docsCount <= searchCandidates[i].current {
				goto EXIT
			}
			if baseDocId < searchCandidates[i].searchCursors[searchCandidates[i].current].DocId {
				nextDocId = searchCandidates[i].searchCursors[searchCandidates[i].current].DocId
				break
			}
		}
		if nextDocId > 0 {
			for searchCandidates[0].docsCount > searchCandidates[0].current && searchCandidates[0].searchCursors[searchCandidates[0].current].DocId < nextDocId {
				searchCandidates[0].current++
			}
		} else { // 全てのトークンで共通のドキュメントが見つかった場合
			// フレーズ検索
			phraseCount := searchPhrase(searchCandidates)
			// スコアリング
			if phraseCount > 0 {
				score := Score(searchCandidates)
				title, err := database.GetDocumentTitle(baseDocId)
				if err != nil {
					return nil, err
				}
				*res = append(*res, SearchResult{
					DocumentId: baseDocId,
					Title:      title,
					Score:      score})
			}
			searchCandidates[0].current++
		}

		// 次の文書を検索
		baseCursor = searchCandidates[0].current
	}

	return res, nil

EXIT:
	return res, nil
}

func searchPhrase(searchCandidates searchCandidates) int {
	var phraseCandidates phraseCandidates
	for _, candidate := range searchCandidates {
		phraseCandidates = append(phraseCandidates, phraseCandidate{
			tokenId:       candidate.tokenId,
			base:          candidate.base,
			postings:      candidate.searchCursors[candidate.current].Postings,
			current:       0,
			postingsCount: candidate.searchCursors[candidate.current].PostingsCount,
		})
	}

	sort.Slice(phraseCandidates, func(i, j int) bool {
		return phraseCandidates[i].postingsCount < phraseCandidates[j].postingsCount
	})

	baseCursor := 0
	phraseCount := 0
	for baseCursor < phraseCandidates[0].postingsCount {
		basePos := phraseCandidates[0].postings[baseCursor] - phraseCandidates[0].base
		nextPos := 0
		for i := 1; i < len(phraseCandidates); i++ {
			for phraseCandidates[i].current < phraseCandidates[i].postingsCount && basePos > phraseCandidates[i].postings[phraseCandidates[i].current]-phraseCandidates[i].base {
				phraseCandidates[i].current++
			}
			if phraseCandidates[i].postingsCount <= phraseCandidates[i].current {
				goto EXIT
			}
			if basePos < phraseCandidates[i].postings[phraseCandidates[i].current]-phraseCandidates[i].base {
				nextPos = phraseCandidates[i].postings[phraseCandidates[i].current] - phraseCandidates[i].base
				break
			}
		}
		if nextPos > 0 {
			for phraseCandidates[0].postingsCount > phraseCandidates[0].current && phraseCandidates[0].postings[phraseCandidates[0].current]-phraseCandidates[0].base < nextPos {
				phraseCandidates[0].current++
			}
		} else {
			phraseCount++
			phraseCandidates[0].current++
		}
		baseCursor = phraseCandidates[0].current
	}

	return phraseCount

EXIT:
	return phraseCount
}
