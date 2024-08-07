package domain

type File struct {
	name          string
	path          string
	localContent  string
	remoteContent string
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

func (f *File) SetLocalContent(content string) {
	f.localContent = content
}

func (f *File) SetRemoteContent(content string) {
	f.remoteContent = content
}

func (f *File) GetName() string {
	return f.name
}

func (f *File) GetPath() string {
	return f.path
}

func (f *File) GetLocalContent() string {
	return f.localContent
}

func (f *File) GetRemoteContent() string {
	return f.remoteContent
}
