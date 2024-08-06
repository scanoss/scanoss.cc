package domain

//
//type StagingDirectory struct {
//	rootPath string
//	files    []*StagedFile
//}
//
//func NewGitStagingDirectory(rootPath string) *StagingDirectory {
//	return &StagingDirectory{
//		files:    []*StagedFile{},
//		rootPath: rootPath,
//	}
//}
//
//func (gsd *StagingDirectory) AddFile(path string) error {
//	f, err := NewGitFile(path)
//	if err != nil {
//		return err
//	}
//	gsd.files = append(gsd.files, f)
//	return nil
//}
//
//func (gsd *StagingDirectory) GetRootPath() string {
//	return gsd.rootPath
//}
