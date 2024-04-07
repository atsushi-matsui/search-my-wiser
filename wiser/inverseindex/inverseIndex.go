package inverseindex

type InverseIndexValue struct {
	TokenId       int
	PostingsList  *PostingsList
	DocsCount     int //トークンを含むドキュメント数
	PostingsCount int //トークンの出現数
}

type InverseIndex map[int]*InverseIndexValue

type InverseIndexList []*InverseIndexValue

func NewInverseIndex(tokenId int, postingList *PostingsList, docCount, postingsCount int) *InverseIndexValue {
	return &InverseIndexValue{
		TokenId:       tokenId,
		PostingsList:  postingList,
		DocsCount:     docCount,
		PostingsCount: postingsCount,
	}
}

func MergeInvertedIndex(base, toBeAdded *InverseIndex) error {
	for _, addedI := range *toBeAdded {
		if baseI, ok := (*base)[addedI.TokenId]; ok {
			mergedPostings, err := MergePostings(*baseI.PostingsList, *addedI.PostingsList)
			if err != nil {
				return err
			}
			(*base)[addedI.TokenId].PostingsList = &mergedPostings
			(*base)[addedI.TokenId].DocsCount += addedI.DocsCount
		} else {
			(*base)[addedI.TokenId] = addedI
		}
	}

	return nil
}

func (ii *InverseIndex) Clear() {
	ii = &InverseIndex{}
}
