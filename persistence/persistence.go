package persistence

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dimanevelev/aqua/model"
	"path/filepath"
)

const table string = "fileStats"

type Client struct {
	Client *sql.DB
}

func (persis *Client) Store(file model.File) error {
	stmtIns, err := persis.Client.Prepare("INSERT INTO " + table + " (path, extension, size, info) VALUES( ?, ?, ?, ?)")
	defer stmtIns.Close()
	if err != nil {
		return err
	}
	info, err := json.Marshal(file.FileInfo)
	if err != nil {
		return err
	}
	_, err = stmtIns.Exec(file.Path, filepath.Ext(file.Path), file.FileInfo.Size, info)
	return err
}

type MySqlConf struct {
	Port     string
	URL      string
	DBName   string
	User     string
	Password string
}

func (persis *Client) CountRows() (uint64, error) {
	var result uint64
	err := persis.Client.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&result)
	return result, err
}

func (persis *Client) MaxFileSize() (int64, string, error) {
	var size int64
	var path string
	err := persis.Client.QueryRow(
		"SELECT size, path "+
			"FROM "+table+" "+
			"ORDER BY size DESC "+
			"LIMIT 1").Scan(&size, &path)
	if err == sql.ErrNoRows {
		// empty table, can be ignored
		err = nil
	}
	return size, path, err
}

func (persis *Client) AVGFileSize() (float64, error) {
	var result sql.NullFloat64
	err := persis.Client.QueryRow("SELECT AVG(size) FROM " + table).Scan(&result)
	if result.Valid {
		return result.Float64, nil
	}
	return 0, err
}

func (persis *Client) ExtensionsList() ([]string, error) {
	var result []string
	results, err := persis.Client.Query("SELECT DISTINCT extension FROM " + table)
	if err != nil {
		return nil, err
	}
	for results.Next() {
		var extension string
		err := results.Scan(&extension)
		if err != nil {
			return nil, err
		}
		result = append(result, extension)
	}
	return result, nil
}

func (persis *Client) MostCommonExt() (string, error) {
	var result string
	var ignore int64
	err := persis.Client.QueryRow(""+
		"SELECT extension, COUNT(extension) as extCount "+
		"FROM "+table+" "+
		"GROUP BY extension "+
		"ORDER BY COUNT(extension) DESC "+
		"LIMIT 1 "+
		"OFFSET 1").Scan(&result, &ignore)
	if err == sql.ErrNoRows {
		// empty table, can be ignored
		err = nil
	}
	return result, err
}

func InitClient(conf MySqlConf) (*sql.DB, error) {
	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/",
			conf.User,
			conf.Password,
			conf.URL,
			conf.Port))

	defer db.Close()

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + conf.DBName)
	if err != nil {
		panic(err)
	}
	db.Close()

	db, err = sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			conf.User,
			conf.Password,
			conf.URL,
			conf.Port,
			conf.DBName))
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS " + table + " ( id integer AUTO_INCREMENT, path varchar(4096), extension varchar(32), size integer, info varchar(4096), primary key (id))")
	if err != nil {
		panic(err)
	}
	return db, err
}
