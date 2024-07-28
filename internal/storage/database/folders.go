package database

type Folder struct {
	Id       int
	Name     string
	Files    []*File
	Parent   *Folder
	Children []*Folder
}

func NewFolder(name string, Parent *Folder) *Folder {
	return &Folder{
		Name:     name,
		Files:    []*File{},
		Parent:   Parent,
		Children: []*Folder{},
	}
}
