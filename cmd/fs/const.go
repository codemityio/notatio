package fs

import "time"

const (
	port              = 8080
	readBufSize       = 32 * 1024
	readHeaderTimeout = 5 * time.Second
	rfc3339Len        = 20
	datePosYYYY       = 4
	datePosYYYYMM     = 7
	datePosYYYYMMDD   = 10
)
