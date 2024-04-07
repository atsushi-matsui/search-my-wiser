package main

import (
	"flag"
	"fmt"
	"localhost/search-my-wiser/common"
	"localhost/search-my-wiser/database"
	"localhost/search-my-wiser/indexer"
	"localhost/search-my-wiser/inverseindex"
	"localhost/search-my-wiser/search"
	"localhost/search-my-wiser/wiki"
	"os"
	"sort"
)

var (
	compressMethodStr       string
	wikipediaDumpFile       string
	query                   string
	maxIndexCount           int
	iiBufferUpdateThreshold int

	iiBuffer *inverseindex.InverseIndex
)

func init() {
	flag.StringVar(&compressMethodStr, "c", "", "compress method for postings list")
	flag.StringVar(&wikipediaDumpFile, "w", "", "wikipedia dump file")
	flag.StringVar(&query, "q", "", "query")
	flag.IntVar(&maxIndexCount, "m", -1, "max index count")

	iiBuffer = &inverseindex.InverseIndex{}

	common.InitEnv()
	database.Init(common.EnvParams.DbEnv)
}

func main() {
	flag.Parse()

	if wikipediaDumpFile != "" {
		file, err := os.Open(wikipediaDumpFile)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		parser, err := wiki.LoadWikiDump(file, maxIndexCount)

		if err != nil {
			fmt.Printf("Error wiki dump: %v\n", err)
			os.Exit(1)
		}
		wikiCount := 0
		for err == nil {
			wikiCount++
			page, err := parser.Next()
			if err == nil {
				indexer.Indexing(page.Title, page.Revisions[0].Text, iiBuffer)
			}
			if maxIndexCount > 0 && wikiCount >= maxIndexCount {
				break
			}
		}
		indexer.FlashBuffer(iiBuffer)
	}

	if query != "" {
		indexCount, err := database.GetIndexCount()
		if err != nil {
			fmt.Printf("not index count: %v\n", err)
			os.Exit(1)
		}
		common.EnvParams.IndexCount = indexCount

		res, err := search.Search(query)
		if err != nil {
			fmt.Printf("search error: %v\n", err)
			os.Exit(1)
		}
		sort.Slice(*res, func(i, j int) bool {
			return (*res)[i].Score > (*res)[j].Score
		})
		for _, r := range *res {
			fmt.Printf("docId: %d, title: %s, score: %f\n", r.DocumentId, r.Title, r.Score)
		}
	}
}
