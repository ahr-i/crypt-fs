package logIPFS

import log2 "github.com/ipfs/go-log/v2"

type logger struct {
	System *log2.ZapEventLogger
}
