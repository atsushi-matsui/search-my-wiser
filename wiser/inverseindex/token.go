package inverseindex

import (
	"fmt"
	"localhost/search-my-wiser/analyzer"
	"localhost/search-my-wiser/common"
	"localhost/search-my-wiser/database"
	"strings"
)

func TextToPostingsLists(docId int, text string, postings *InverseIndex) error {
	texts := strings.Split(text, "")
	bufferPostings := &InverseIndex{}

	position := 0
	for {
		token := analyzer.TextAnalyze(&texts)
		tokenStr := strings.Join(token, "")
		if len(token) <= 0 {
			break
		}
		position++ //トークンの出現位置
		err := textToPostingsList(docId, tokenStr, position, bufferPostings)
		if err != nil {
			fmt.Printf("Error text to postings list: %v\n", err)
			return err
		}
	}

	if postings != nil {
		// ミニ転置インデックスにマージする
		mergeInvertedIndex(postings, bufferPostings)
	} else {
		postings = bufferPostings
	}

	return nil
}

func QueryToPostingList(query string, ii *InverseIndexList) error {
	texts := strings.Split(query, "")

	if len(texts) < int(common.EnvParams.Ngram) {
		return fmt.Errorf("query is too short. query: %s", query)
	}

	for {
		token := analyzer.TextAnalyze(&texts)
		tokenStr := strings.Join(token, "")
		// Nに満たないタームは利用しない
		if len(token) > 0 && len(token) < common.EnvParams.Ngram {
			*ii = append(*ii, nil)
			continue
		}
		if len(token) <= 0 {
			break
		}
		// トークンからDBのtokensを取得する
		t, err := database.GetToken(tokenStr, false)
		if err != nil {
			return err
		}
		// ポスティングリストをデコードする
		pl, err := decodePostings(t.Postings)
		if err != nil {
			return err
		}
		*ii = append(*ii, NewInverseIndex(t.Id, pl.PostingsList, pl.DocCount, pl.PostingsCount))
	}

	return nil
}

/**
 * テキストをポスティングリストに変換する
 */
func textToPostingsList(docId int, tokenStr string, position int, bufferPostings *InverseIndex) error {
	var p *Postings
	var iiVal *InverseIndexValue

	token, err := database.GetToken(tokenStr, docId > 0)
	if err != nil {
		return err
	}
	if val, ok := (*bufferPostings)[token.Id]; ok {
		iiVal = val
	} else {
		iiVal = nil
	}
	if iiVal != nil { // ミニ転置インデックス内にすでに対象のポスティングリストが存在する場合は出現回数を加算する
		p = (*iiVal.PostingsList)[0]
		p.PostingsCount++
	} else {
		// 新規でポスティングリストを作成
		if docId > 0 { // ストレージ内にトークンが存在しない場合は完全に新規作成
			iiVal = NewInverseIndex(token.Id, nil, 1, 0)
		} else { // 存在する場合はストレージから取得したドキュメント数を設定
			iiVal = NewInverseIndex(token.Id, nil, token.DocsCount, 0)
		}

		// ミニ転置インデックスへ追加
		(*bufferPostings)[token.Id] = iiVal
		// ポスティングリストを作成
		pl := NewPostingsList(docId)
		iiVal.PostingsList = pl
		p = (*pl)[0]
	}
	// 出現位置を追加
	p.Postings = append(p.Postings, position)
	// 転置インデックス全体でのトークンの総出現回数も追加
	iiVal.PostingsCount++

	return nil
}
