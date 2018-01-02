package bbld


type WriterText interface {
	Write(dest string, txt string) error
}