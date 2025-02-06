package history

import "log"

type History struct {
	Filename string
	Log      *log.Logger // Json忽略Log欄位
}

func New(filename string) *History {
	return &History{Filename: filename, Log: log.Default()}
}
