package dbase

import (
	"database/sql"
	"fmt"
)

type PostInfo struct {
	UserId int    `json:"userId"`
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type CommentInfo struct {
	PostId int    `json:"postId"`
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}

type HomePageStruct struct {
	UserId  int    `json:"userId"`
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	Comments []CommentInfo
}

var Db *sql.DB

const (
	DBName    = "db_user_7"
	PortDB    = "3306"
	PostTb    = "info_posts"
	CommentTb = "info_comments"
	Pass      = "passwordIK"
	UsName    = "root"
	Driver    = "mysql"
)

var DbSN = fmt.Sprintf("%s:%s@tcp(127.0.0.1:%v)/%s?charset=utf8", UsName, Pass, PortDB, DBName)

func CheckError(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}

}
