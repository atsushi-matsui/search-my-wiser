package database

import (
	"database/sql"
	"fmt"
	"localhost/search-my-wiser/common"

	"github.com/go-sql-driver/mysql"
)

var (
	db                        *sql.DB
	prepareCreateToken        *sql.Stmt
	prepareCreateDocument     *sql.Stmt
	prepareGetTokenById       *sql.Stmt
	prepareGetTokenByToken    *sql.Stmt
	prepareGetDocumentByTitle *sql.Stmt
	prepareUpdateTokenById    *sql.Stmt
	prepareGetDocumentById    *sql.Stmt
	prepareUpdateDocumentById *sql.Stmt
	prepareGetDocumentsCount  *sql.Stmt
)

type Document struct {
	Id    int
	Title string
	Body  string
}

type Token struct {
	Id        int
	Token     string
	DocsCount int
	Postings  []byte
}

func Init(env common.DbEnv) error {
	var err error
	con := setDbConfig(env)
	db, err = sql.Open(env.DbDriver, con.FormatDSN())
	if err != nil {
		return err
	}

	//db.SetMaxOpenConns(100)
	//db.SetMaxIdleConns(100)

	err = db.Ping()
	if err != nil {
		return err
	}
	fmt.Printf("Connected to %s\n", env.DbName)

	setPrepares(db)

	return nil
}

func setDbConfig(env common.DbEnv) *mysql.Config {
	con := mysql.NewConfig()
	con.User = env.DbUser
	con.Passwd = env.DbPassword
	con.Net = env.DbProtocol
	switch env.DbProtocol {
	case "tcp":
		con.Addr = fmt.Sprintf("%s:%s", env.DbHost, env.DbPort)
	case "unix":
		con.Addr = env.DbSocket
	}
	con.DBName = env.DbName
	con.CheckConnLiveness = true

	return con
}

/*
 * https://golang.shop/post/go-databasesql-06-prepared-ja/
 * SQL文をプリペアすると、プール内のコネクション上でプリペアされます。
 */
func setPrepares(db *sql.DB) error {
	var err error

	prepareCreateToken, err = db.Prepare("INSERT INTO tokens (token, docs_count, postings) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	prepareCreateDocument, err = db.Prepare("INSERT INTO documents (title, body) VALUES (?, ?)")
	if err != nil {
		return err
	}
	prepareGetTokenById, err = db.Prepare("SELECT * FROM tokens WHERE id = ?")
	if err != nil {
		return err
	}
	prepareGetTokenByToken, err = db.Prepare("SELECT * FROM tokens WHERE token = ?")
	if err != nil {
		return err
	}
	prepareGetDocumentById, err = db.Prepare("SELECT id, title FROM documents WHERE id = ?")
	if err != nil {
		return err
	}
	prepareGetDocumentByTitle, err = db.Prepare("SELECT * FROM documents WHERE title = ?")
	if err != nil {
		return err
	}
	prepareUpdateTokenById, err = db.Prepare("UPDATE tokens SET docs_count = ?, postings = ? WHERE id = ?")
	if err != nil {
		return err
	}
	prepareUpdateDocumentById, err = db.Prepare("UPDATE documents SET title = ?, body = ? WHERE id = ?")
	if err != nil {
		return err
	}
	prepareGetDocumentsCount, err = db.Prepare("SELECT count(*) FROM documents")
	if err != nil {
		return err
	}

	return nil
}

func AddDocument(title string, text string) (int, error) {
	doc, err := getDocumentByTitle(title)
	if err != nil && err != sql.ErrNoRows {
		fmt.Printf("Error getting document: %v\n", err)
		return -1, err
	}

	id := -1
	if doc == nil {
		id, err = createDocument(title, text)
		if err != nil {
			fmt.Printf("Error creating document: %v\n", err)
			return id, err
		}
	} else {
		err = updateDocumentById(doc.Id, title, text)
		if err != nil {
			fmt.Printf("Error updating document: %v\n", err)
			return id, err
		}
	}

	return id, nil
}

func GetToken(tokenStr string, isInsert bool) (*Token, error) {
	var token *Token

	if isInsert {
		token, err := getTokenByToken(tokenStr)
		if err != nil {
			if err == sql.ErrNoRows {
				tokenId, err := createToken(tokenStr, 0, []byte{})
				if err != nil {
					return nil, err
				}
				token = &Token{
					Id:        tokenId,
					Token:     tokenStr,
					DocsCount: 0,
					Postings:  []byte{},
				}
			} else {
				return nil, err
			}
		}

		return token, nil
	} else {
		var err error
		token, err = getTokenByToken(tokenStr)
		if err != nil {
			return nil, err
		}
	}

	return token, nil
}

func GetPosting(tokenId int) (*Token, error) {
	token, err := getTokenById(tokenId)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func UpdatePostings(tokenId int, docCount int, postings []byte) error {
	err := updateTokenById(tokenId, docCount, postings)
	if err != nil {
		fmt.Printf("Error updating postings: %v\n", err)
		return err
	}

	return nil
}

func GetIndexCount() (int, error) {
	return getDocumentCount()
}

func GetDocumentTitle(id int) (string, error) {
	doc, err := getDocumentById(id)
	if err != nil {
		return "", err
	}

	return doc.Title, nil
}

func createToken(token string, docsCount int, postings []byte) (int, error) {
	res, err := prepareCreateToken.Exec(token, docsCount, postings)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(id), nil
}

func createDocument(title string, text string) (int, error) {
	res, err := prepareCreateDocument.Exec(title, text)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(id), nil
}

func getTokenById(tokenStr int) (*Token, error) {
	rows := prepareGetTokenById.QueryRow(tokenStr)

	token := &Token{}
	err := rows.Scan(&token.Id, &token.Token, &token.DocsCount, &token.Postings)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func getTokenByToken(tokenStr string) (*Token, error) {
	token := &Token{}
	err := prepareGetTokenByToken.QueryRow(tokenStr).Scan(&token.Id, &token.Token, &token.DocsCount, &token.Postings)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func getDocumentById(id int) (*Document, error) {
	document := &Document{}
	err := prepareGetDocumentById.QueryRow(id).Scan(&document.Id, &document.Title)
	if err != nil {
		return nil, err
	}

	return document, nil
}

func getDocumentByTitle(title string) (*Document, error) {
	document := &Document{}
	err := prepareGetDocumentByTitle.QueryRow(title).Scan(&document.Id, &document.Title, &document.Body)
	if err != nil {
		return nil, err
	}

	return document, nil
}

func updateTokenById(id int, docsCount int, postings []byte) error {
	_, err := prepareUpdateTokenById.Exec(docsCount, postings, id)
	if err != nil {
		return err
	}

	return nil
}

func updateDocumentById(id int, title string, text string) error {
	_, err := prepareUpdateDocumentById.Exec(title, text, id)
	if err != nil {
		return err
	}

	return nil
}

func getDocumentCount() (int, error) {
	count := 0
	err := prepareGetDocumentsCount.QueryRow().Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
