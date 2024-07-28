package database

type File struct {
	Name string
}

func NewFile(name string) *File {
	return &File{
		Name: name,
	}
}
