package domain

type ScanResult struct {
	content []byte
}

func NewScanResult(content []byte) ScanResult {
	return ScanResult{
		content,
	}
}

func (sr *ScanResult) GetContent() []byte {
	return sr.content
}
