package domain

type File struct {
	name    string
	path    string
	content string
}

func NewFile() *File {
	return &File{}
}

func (f *File) SetName(name string) {
	f.name = name
}

func (f *File) SetPath(path string) {
	f.path = path
}

func (f *File) SetContent(content string) {
	f.content = content
}

func (f *File) GetName() string {
	return f.name
}

func (f *File) GetPath() string {
	return f.path
}

func (f *File) GetContent() string {
	return f.content
}
