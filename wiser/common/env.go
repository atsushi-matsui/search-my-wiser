package common

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Env struct {
	DbEnv DbEnv
	// 他にも必要なパラメータがあれば追加していく
	IiBufferUpdateThreshold int
	IiBufferCount           int
	IndexCount              int
	Ngram                   int
	CompressMethod          CompressMethod
	TextAnalyzeMethod       TextAnalyzeMethod
	ScoringMethod           ScoringMethod
}

type DbEnv struct {
	DbDriver   string
	DbUser     string
	DbPassword string
	DbName     string
	DbProtocol string
	DbHost     string
	DbPort     string
	DbSocket   string
}

const (
	CompressNone CompressMethod = iota
	CompressGolomb
)

type CompressMethod int

const (
	Ngram TextAnalyzeMethod = iota
)

type TextAnalyzeMethod int

const (
	TfIdf ScoringMethod = iota
	Bm25
)

type ScoringMethod int

var EnvParams *Env

func InitEnv() {
	err := godotenv.Load("./.env") // envファイルのパスを渡す。何も渡さないと、どうディレクトリにある、.envファイルを探す
	if err != nil {
		panic("Error loading .env file")
	}

	defaultIiBufferUpdateThreshold := 0
	if defaultIiBufferUpdateThreshold, err = strconv.Atoi(os.Getenv("DEFAULT_II_BUFFER_UPDATE_THRESHOLD")); err != nil {
		panic("DEFAULT_II_BUFFER_UPDATE_THRESHOLD is not integer")
	}

	compressMethodInteger := 0
	if compressMethodInteger, err = strconv.Atoi(os.Getenv("COMPRESS_METHOD")); err != nil {
		panic("COMPRESS_METHOD is not integer")
	}
	compressMethod := CompressMethod(compressMethodInteger)

	textAnalyzeMethodInteger := 0
	if textAnalyzeMethodInteger, err = strconv.Atoi(os.Getenv("TEXT_ANALYZE_METHOD")); err != nil {
		panic("TEXT_ANALYZE_METHOD is not integer")
	}
	textAnalyzeMethod := TextAnalyzeMethod(textAnalyzeMethodInteger)

	nGram := 0
	if nGram, err = strconv.Atoi(os.Getenv("N_GRAM")); err != nil {
		panic("N_GRAM is not integer")
	}

	scoringMethodInteger := 0
	if scoringMethodInteger, err = strconv.Atoi(os.Getenv("SCORING_METHOD")); err != nil {
		panic("SCORING_METHOD is not integer")
	}
	scoringMethod := ScoringMethod(scoringMethodInteger)

	EnvParams = &Env{
		DbEnv: DbEnv{
			DbDriver:   os.Getenv("DB_DRIVER"), // 読み込んだ後の使い方はいつも通り
			DbUser:     os.Getenv("DB_USER"),
			DbPassword: os.Getenv("DB_PASSWORD"),
			DbName:     os.Getenv("DB_NAME"),
			DbProtocol: os.Getenv("DB_PROTOCOL"),
			DbHost:     os.Getenv("DB_HOST"),
			DbPort:     os.Getenv("DB_PORT"),
			DbSocket:   os.Getenv("DB_SOCKET"),
		},
		IiBufferUpdateThreshold: defaultIiBufferUpdateThreshold,
		IiBufferCount:           0,
		CompressMethod:          compressMethod,
		IndexCount:              0,
		TextAnalyzeMethod:       textAnalyzeMethod,
		Ngram:                   nGram,
		ScoringMethod:           scoringMethod,
	}
}
