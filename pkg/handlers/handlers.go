package handlers

import (
	"fmt"
	"github.com/Igor-Koniukhov/rest_api/dbase"
	"github.com/Igor-Koniukhov/rest_api/pkg/config"
	"github.com/Igor-Koniukhov/rest_api/pkg/datastruct"
	"github.com/Igor-Koniukhov/rest_api/pkg/render"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	data       dbase.HomePageStruct
	DataStruct dbase.HomePageStruct
	DataArr    []dbase.HomePageStruct
	wg         sync.WaitGroup
	err        error
	Repo       *Repository
	userNumber = 7
)

type Repository struct {
	App *config.AppConfig
}

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func newData(ds []dbase.HomePageStruct, d dbase.HomePageStruct) {
	DataArr = ds
	DataStruct = d
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	var datas []dbase.HomePageStruct

	sqlStmt := fmt.Sprintf("SELECT * FROM %s ", dbase.PostTb)

	rows, err := dbase.Db.Query(sqlStmt)
	defer rows.Close()

	dbase.CheckError(err)

	for rows.Next() {
		_ = rows.Scan(
			&data.UserId,
			&data.Id,
			&data.Title,
			&data.Body)

		id := data.Id
		wg.Add(1)
		go getComments(id)
		wg.Wait()
		datas = append(datas, data)
	}
	newData(datas, data)


	render.RenderTemplate(w, "home.page.tmpl", datastruct.Data{StructHomePage: &datas})

}

func (m *Repository) GetPost(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	var post dbase.PostInfo

	sqlStmt := fmt.Sprintf("SELECT * FROM %s WHERE id=%v", dbase.PostTb, id)
	row := dbase.Db.QueryRow(sqlStmt)

	err := row.Scan(
		&post.UserId,
		&post.Id,
		&post.Title,
		&post.Body)
	dbase.CheckError(err)

	render.RenderTemplate(w, "post.page.tmpl", datastruct.Data{Post: post,})

}

func (m *Repository) GetForUpdatePost(w http.ResponseWriter, r *http.Request) {
	var post dbase.PostInfo

	id := r.FormValue("id")
	sqlStmt := fmt.Sprintf("SELECT * FROM %s WHERE id=%v", dbase.PostTb, id)
	row := dbase.Db.QueryRow(sqlStmt)

	err = row.Scan(
		&post.UserId,
		&post.Id,
		&post.Title,
		&post.Body)
	dbase.CheckError(err)

	render.RenderTemplate(w, "updatepost.page.tmpl", datastruct.Data{Post: post,})
}

func (m *Repository) UpdatePost(w http.ResponseWriter, r *http.Request) {
	var post dbase.PostInfo

	userId := r.FormValue("userId")
	id := r.FormValue("id")
	title := r.FormValue("title")
	body := r.FormValue("body")

	sqlStmt := fmt.Sprintf("UPDATE %s SET userId=?, id=?, title=?, body=? WHERE id=%v ", dbase.PostTb, id)
	stmt, err := dbase.Db.Prepare(sqlStmt)
	dbase.CheckError(err)

	_, err = stmt.Exec(
		userId,
		id,
		title,
		body)
	dbase.CheckError(err)

	render.RenderTemplate(w, "massage.page.tmpl", datastruct.Data{Post: post})

}

func (m *Repository) PostCreate(w http.ResponseWriter, r *http.Request) {
	var post dbase.PostInfo
	userId := r.FormValue("userId")
	id := r.FormValue("id")
	title := r.FormValue("title")
	body := r.FormValue("body")
	idl := lastIdPost()

	sqlStmt := fmt.Sprintf("INSERT INTO %s (`userId`, `id`, `title`, `body`) VALUES (?, ?, ?, ?) ", dbase.PostTb)
	stmt, err := dbase.Db.Prepare(sqlStmt)
	dbase.CheckError(err)


	dbase.CheckError(err)

	_, err = stmt.Exec(
		userId,
		id,
		title,
		body)
	dbase.CheckError(err)

	render.RenderTemplate(w, "createpost.page.tmpl", datastruct.Data{
		Post: post,
		Int:  idl,
		Int2: userNumber,
	})
}

func (m *Repository) DeletePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	dbase.CheckError(err)

	delete(dbase.PostTb, "id", id)

	 delete(dbase.CommentTb, "postId", id)

	render.RenderTemplate(w, "delete.page.tmpl", datastruct.Data{Int: id})
}

func delete(table, name string, id int) {
	sqlStmt2 := fmt.Sprintf("DELETE FROM %s WHERE %s=%v", table, name, id)
	_, err := dbase.Db.Exec(sqlStmt2)
	dbase.CheckError(err)
}

func (m *Repository) GetComment(w http.ResponseWriter, r *http.Request) {
	fields := strings.Split(r.URL.String(), "/")
	id := fields[len(fields)-1]
	var comment dbase.CommentInfo

	sqlStmt := fmt.Sprintf("SELECT * FROM %s WHERE id=%v", dbase.CommentTb, id)
	row := dbase.Db.QueryRow(sqlStmt)
	err := row.Scan(
		&comment.PostId,
		&comment.Id,
		&comment.Name,
		&comment.Email,
		&comment.Body)
	dbase.CheckError(err)

	render.RenderTemplate(w, "comment.page.tmpl", datastruct.Data{Comment: comment,})

}
func (m *Repository) GetForUpdateComment(w http.ResponseWriter, r *http.Request) {
	var comment dbase.CommentInfo
	id := r.FormValue("id")
	sqlStmt := fmt.Sprintf("SELECT * FROM %s WHERE id=%v", dbase.CommentTb, id)
	row := dbase.Db.QueryRow(sqlStmt)
	err = row.Scan(
		&comment.PostId,
		&comment.Id,
		&comment.Name,
		&comment.Email,
		&comment.Body)
	dbase.CheckError(err)
	render.RenderTemplate(w, "updatemassage.page.tmpl", datastruct.Data{Comment: comment,})

}

func (m *Repository) UpdateComment(w http.ResponseWriter, r *http.Request) {
	var comment dbase.CommentInfo
	postId := r.FormValue("postId")
	id := r.FormValue("id")
	name := r.FormValue("name")
	email := r.FormValue("email")
	body := r.FormValue("body")

	sqlStmt := fmt.Sprintf("UPDATE %s SET postId=?,id=?, name=?, email=?, body=? WHERE id=%v ", dbase.CommentTb, id)
	stmt, err := dbase.Db.Prepare(sqlStmt)
	dbase.CheckError(err)
	_, err = stmt.Exec(
		postId,
		id,
		name,
		email,
		body)
	dbase.CheckError(err)

	render.RenderTemplate(w, "massage.page.tmpl", datastruct.Data{Comment: comment})

}

func (m *Repository) DeleteComment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	dbase.CheckError(err)

	delete(dbase.CommentTb, "id", id)

	render.RenderTemplate(w, "massage.page.tmpl", datastruct.Data{Int: id})
}

func getComments(id int) []dbase.CommentInfo {
	defer wg.Done()
	var comment dbase.CommentInfo
	var comments []dbase.CommentInfo
	data.Comments = comments

	sqlStmt := fmt.Sprintf("SELECT * FROM %s WHERE postId=%v", dbase.CommentTb, id)
	rows, err := dbase.Db.Query(sqlStmt)
	defer rows.Close()
	dbase.CheckError(err)
	for rows.Next() {
		_ = rows.Scan(
			&comment.PostId,
			&comment.Id,
			&comment.Name,
			&comment.Email,
			&comment.Body)
		data.Comments = append(data.Comments, comment)
	}
	time.Sleep(time.Millisecond * 50)

	return data.Comments

}

func lastIdComment(postId int) int {
	var id int
	for _, d := range DataArr {
		if d.Id == postId {
			for _, comment := range d.Comments {
				id = comment.Id
			}
		}
	}
	return id
}

func lastIdPost() int {
	var id int
	for _, data := range DataArr {
		id = data.Id
	}
	newId := id + 1

	return newId

}

func (m *Repository) WriteComment(w http.ResponseWriter, r *http.Request) {
	var comment dbase.CommentInfo
	var n int
	var post dbase.PostInfo
	urlId := strings.Split(r.URL.String(), "/")
	idPref := strings.Split(urlId[len(urlId)-1], "=")
	idPost, err := strconv.Atoi(idPref[len(idPref)-1])
	dbase.CheckError(err)

	n = lastIdComment(idPost) + 1

	postId := r.FormValue("postId")
	id := r.FormValue("id")
	name := r.FormValue("name")
	email := r.FormValue("email")
	body := r.FormValue("body")

	sqlStmt := fmt.Sprintf("INSERT INTO %s VALUES(?,?,?,?,?)", dbase.CommentTb)
	stmt, err := dbase.Db.Prepare(sqlStmt)
	dbase.CheckError(err)
	_, err = stmt.Exec(postId, id, name, email, body)

	if err != nil {
		fmt.Println(err, " sql execution error")
	}

	render.RenderTemplate(w, "writeComment.page.tmpl", datastruct.Data{
		Comment: comment,
		Post:    post,
		Int:     idPost,
		Int2:    n,
	})

}
