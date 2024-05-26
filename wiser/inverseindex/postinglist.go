package inverseindex

import (
	"fmt"
	"localhost/search-my-wiser/common"
	"localhost/search-my-wiser/database"
	"localhost/search-my-wiser/util"
	"os"
	"sort"
)

type Postings struct {
	DocId         int
	Postings      []int //ドキュメント内のトークンの出現位置
	PostingsCount int   //ドキュメント内のトークンの出現数
}
type PostingsList []*Postings

type decodedPostingsList struct {
	PostingsList  *PostingsList
	DocCount      int
	PostingsCount int
}

func NewPostingsList(docId int) *PostingsList {
	pl := append(PostingsList{}, NewPostings(docId, []int{}, 1))

	return &pl
}

func NewPostings(docId int, p []int, pc int) *Postings {
	return &Postings{
		DocId:         docId,
		Postings:      p,
		PostingsCount: pc,
	}
}

// 同一のドキュメントIDを持つ場合は動作を保証しない
func mergePostings(base, toBeAdded PostingsList) (PostingsList, error) {
	merged := PostingsList{}

	for len(base) > 0 || len(toBeAdded) > 0 {
		var p Postings
		if len(toBeAdded) == 0 || (len(base) > 0 && base[0].DocId <= toBeAdded[0].DocId) {
			p = *base[0]
			base = base[1:]
		} else if len(base) == 0 || base[0].DocId >= toBeAdded[0].DocId {
			p = *toBeAdded[0]
			toBeAdded = toBeAdded[1:]
		} else {
			panic("Invalid condition")
		}

		merged = append(merged, &p)
	}

	return merged, nil
}

func UpdatePostings(iiVal *InverseIndexValue) error {
	//dbからポスティグリストを取得して、ポスティングリストをデコードする
	pl, err := fetchPostings(iiVal.TokenId)
	if err != nil {
		return err
	}
	//取得したポスティングリストとマージする
	mergedPl, err := mergePostings(*iiVal.PostingsList, *pl.PostingsList)
	if err != nil {
		return err
	}
	iiVal.PostingsList = &mergedPl
	iiVal.DocsCount += pl.DocCount
	//ポスティングリストをエンコードする
	eMergedPl, err := encodePostings(mergedPl)
	if err != nil {
		return err
	}
	//dbへ最新のポスティングリストを保存
	err = database.UpdatePostings(iiVal.TokenId, iiVal.DocsCount, eMergedPl)
	if err != nil {
		return err
	}

	return nil
}

func fetchPostings(tokenId int) (*decodedPostingsList, error) {
	//dbからポスティングリストを取得
	token, err := database.GetPosting(tokenId)
	if err != nil {
		return nil, err
	}

	pl, err := decodePostings(token.Postings)
	if err != nil {
		return nil, err
	}

	return pl, nil
}

func decodePostings(ePl []byte) (*decodedPostingsList, error) {
	//デコード処理
	var dpl *decodedPostingsList
	var err error
	switch common.EnvParams.CompressMethod {
	case common.CompressNone:
		dpl, err = decodePostingsNone(ePl)

	case common.CompressGolomb:
		dpl, err = decodePostingsGolomb(ePl)
	}

	if err != nil {
		return nil, err
	}

	return dpl, nil
}

func decodePostingsNone(ePl []byte) (*decodedPostingsList, error) {
	pl := PostingsList{}
	docCount := 0 //文書の出現数
	pCount := 0   //タームの出現数
	cursor := 0
	for cursor < len(ePl) {
		docID := util.ByteToInt(ePl[cursor : cursor+4])
		cursor += 4
		pc := util.ByteToInt(ePl[cursor : cursor+4])
		cursor += 4

		docCount++
		var posl []int
		for i := cursor; i < cursor+pc*4; i += 4 {
			p := util.ByteToInt(ePl[i : i+4])
			posl = append(posl, p)
			pCount++
		}
		cursor += pc * 4

		pl = append(pl, NewPostings(docID, posl, pc))
	}

	return &decodedPostingsList{
		PostingsList:  &pl,
		DocCount:      docCount,
		PostingsCount: pCount}, nil
}

func decodePostingsGolomb(_ []byte) (*decodedPostingsList, error) {
	fmt.Println("undefined function.")
	os.Exit(1)

	return &decodedPostingsList{}, nil
}

func encodePostings(pl PostingsList) ([]byte, error) {
	var ePl []byte
	var err error
	switch common.EnvParams.CompressMethod {
	case common.CompressNone:
		ePl, err = encodePostingsNone(pl)

	case common.CompressGolomb:
		ePl, err = encodePostingsGolomb(pl)
	}

	if err != nil {
		return nil, err
	}

	return ePl, nil
}

func encodePostingsNone(pl PostingsList) ([]byte, error) {
	var ePl []byte
	for _, p := range pl {
		ePl = append(ePl, util.IntToBetes(p.DocId)...)
		ePl = append(ePl, util.IntToBetes(p.PostingsCount)...)

		// ドキュメントの出現位置をソートしておく
		sort.Ints(p.Postings)
		for _, pos := range p.Postings {
			ePl = append(ePl, util.IntToBetes(pos)...)
		}
	}

	return ePl, nil
}

func encodePostingsGolomb(_ PostingsList) ([]byte, error) {
	fmt.Println("undefined function.")
	os.Exit(1)

	return nil, nil
}
