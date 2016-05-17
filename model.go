package main

import "time"

type( File struct {
	Id 		  int		`json:"id"`
	Name      string    `json:"name"`
    Completed bool      `json:"completed"`
    Due       time.Time `json:"due"`
}

Data struct {
		files []*File
	}
)

func (al *Data) append(file *File) {
	al.files = append(al.files, file)
}