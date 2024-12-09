package service

type ScanService interface {
	Scan(dirPath string, args []string) error
	CheckDependencies() error
}
