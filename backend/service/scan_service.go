package service

type ScanService interface {
	Scan(args []string) error
	CheckDependencies() error
}
