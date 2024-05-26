package main

import (
	"flag"
	"strconv"
	"testing"
)


func TestIndex(t *testing.T) {
	tests := []struct {
		name string
		w    string
		m    int
	}{
		{name: "5件登録", w: "../files/jawiki-latest-pages-articles.xml", m: 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine.Set("w", tt.w)               // 登録するwikiデータの指定
			flag.CommandLine.Set("m", strconv.Itoa(tt.m)) // 登録するwikiデータの最大件数
			main()
		})
	}
}


func TestSearch(t *testing.T) {
	tests := []struct {
		name string
		q    string
	}{
		{name: "言語を検索", q: "言語"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine.Set("q", tt.q) // 検索クエリを設定
			main()
		})
	}
}
