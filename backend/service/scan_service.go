package service

type ScanService interface {
	CheckDependencies() error
	GetDefaultScanArgs() []string
	Scan(args []string) error
	ScanStream(args []string) error
}
