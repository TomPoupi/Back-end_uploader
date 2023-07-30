package common

type Upload struct {
	Id           int
	Name         string
	Description  string
	Date         string
	Object_video Video
}

type Video struct {
	Video_id  string
	File_name string
	Path      string
}

// type News struct {
// 	Headline string
// 	Body     string
// }
