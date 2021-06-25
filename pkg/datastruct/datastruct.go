package datastruct

import "github.com/Igor-Koniukhov/rest_api/dbase"

type Data struct{
	Int int
	Int2 int
	Ints []int
	Str string
	Strs []string
	StructHomePage *[]dbase.HomePageStruct
	Comment dbase.CommentInfo
	Comments []dbase.CommentInfo
	Post dbase.PostInfo
	Posts []dbase.PostInfo

}
