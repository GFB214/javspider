package main

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

//Video 影片
type Video struct {
	ID         int    `db:"id"`
	Code       string `db:"code"`
	Title      string `db:"title"`
	Date       string `db:"date"`
	Magnet     string `db:"magnet"`
	Cover      string `db:"cover"`
	Downloaded int    `db:"downloaded"`
	Like       int    `db:"like"`
}

//Work 未完成的任务
type Work struct {
	ID  int    `db:"id"`
	URL string `db:"url"`
}

//Db 数据库
var Db *sqlx.DB

func init() {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, pwd, dbURL, database)
	db, err := sqlx.Open(sqltype, dataSourceName)
	if err != nil {
		fmt.Println(err)
		fmt.Println("数据库初始化出错，退出")
		os.Exit(1)
	}
	db.SetMaxOpenConns(3)
	db.SetMaxIdleConns(2)

	Db = db
}

func exist(code string) bool {
	var record int
	err := Db.Get(&record, "SELECT count(*) FROM video WHERE code=?", code)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if record == 1 {
		return true
	}
	return false
}

func insert(video Video) bool {
	res, err := Db.Exec("INSERT INTO video (code, title, date, magnet, cover) VALUES (?, ?, ?, ?, ?)",
		video.Code, video.Title, video.Date, video.Magnet, video.Cover)
	if err != nil {
		fmt.Println(err)
		return false
	}

	row, err := res.RowsAffected()
	if err != nil {
		return false
	}
	return row >= 1
}

func getWorks(count int) []Work {
	var works []Work
	err := Db.Select(&works, "SELECT * FROM work LIMIT ?", count)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return works
}

func deleteWorkByID(id int) bool {
	res, err := Db.Exec("DELETE FROM work WHERE id=?", id)
	if err != nil {
		return false
	}
	row, err := res.RowsAffected()
	if err != nil {
		return false
	}
	return row >= 1
}

func deleteWorkByURL(url string) bool {
	res, err := Db.Exec("DELETE FROM work WHERE url=?", url)
	if err != nil {
		return false
	}
	row, err := res.RowsAffected()
	if err != nil {
		return false
	}
	return row >= 1
}

func insertWork(url string) bool {
	res, err := Db.Exec("INSERT INTO work (url) VALUES (?)", url)
	if err != nil {
		return false
	}

	row, err := res.RowsAffected()
	if err != nil {
		return false
	}
	return row >= 1
}
