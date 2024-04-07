package main

import (
	"flag"
	"strconv"
	"testing"
)

func Test_run(t *testing.T) {
	tests := []struct {
		name string
		w    string
		m    int
	}{
		{name: "input0", w: "../../wiser-20140928/files/jawiki-latest-pages-articles.xml", m: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine.Set("w", tt.w)               // -target=iと指定したかの様に設定できる
			flag.CommandLine.Set("m", strconv.Itoa(tt.m)) // Convert tt.m to string before passing it as an argument
			main()
		})
	}
}
