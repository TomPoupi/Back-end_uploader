package common

type VideoGene struct {
	Id           int
	Name         string
	Description  string
	Date         string
	Object_video VideoDetail
}

type VideoDetail struct {
	Video_id  string
	File_name string
	Path      string
	Size      int64
}

type Users struct {
	Id       int
	Username string
	Password string
}

// type News struct {
// 	Headline string
// 	Body     string
// }
