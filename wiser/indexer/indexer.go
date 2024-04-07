package indexer

import (
	"fmt"
	"localhost/search-my-wiser/common"
	"localhost/search-my-wiser/database"
	"localhost/search-my-wiser/inverseindex"
	"os"
)

func Indexing(title, text string, iiBuffer *inverseindex.InverseIndex) error {
	id, err := database.AddDocument(title, text)
	if err != nil {
		fmt.Printf("Error adding document: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("id: %d, title: %s\n", id, title)

	// テキストをNGramに分解してバッファに追加する
	inverseindex.TextToPostingsLists(id, text, iiBuffer)
	common.EnvParams.IiBufferCount++

	// バッファに所定の文書数が溜まったらストレージに保存する
	if iiBuffer != nil && common.EnvParams.IiBufferCount >= common.EnvParams.IiBufferUpdateThreshold {
		FlashBuffer(iiBuffer)
	}

	return nil
}

func FlashBuffer(iiBuffer *inverseindex.InverseIndex) {
	// バッファをストレージに保存する
	for _, iiVal := range *iiBuffer {
		inverseindex.UpdatePostings(iiVal)
	}
	// バッファをクリアする
	iiBuffer.Clear()
}
