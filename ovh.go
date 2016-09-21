package logrusOVH

import "github.com/Sirupsen/logrus"

// Protocol define available transfert proto
type Protocol uint8

const (
	// GELFUDP for Gelf + UDP
	GELFUDP Protocol = 1 + iota
	// GELFTCP for Gelf + TCP
	GELFTCP
	// GELFTLS for Gelf + TLS
	GELFTLS
	// CAPNPROTOUDP for Cap'n proto + UDP
	CAPNPROTOUDP
	// CAPNPROTOTCP for Cap'n proto + TCP
	CAPNPROTOTCP
	// CAPNPROTOTLS for Cap'n proto + TLS
	CAPNPROTOTLS
)

// TODO reverse map
func (p Protocol) String() string {
	switch p {
	case GELFTCP:
		return "GELFTCP"
	case GELFUDP:
		return "GELFUDP"
	case GELFTLS:
		return "GELFTLS"
	case CAPNPROTOUDP:
		return "CAPNPROTOUDP"
	case CAPNPROTOTCP:
		return "CAPNPROTOTCP"
	case CAPNPROTOTLS:
		return "CAPNPROTOTLS"
	default:
		return "UNKNOW"
	}
}

// CompressAlgo the compression algorithm used
type CompressAlgo uint8

const (
	// COMPRESSNONE No compression
	COMPRESSNONE = 1 + iota
	// COMPRESSGZIP GZIP compression
	COMPRESSGZIP
	// COMPRESSZLIB ZLIB compression
	COMPRESSZLIB
)

const (
	UDPCHUNKMAXSIZE           = 8192
	UDP_CHUNK_MAX_SIZE_NOFRAG = 1472
	UDP_CHUNK_MAX_SIZE_FRAG   = 8192
	//UDP_CHUNK_MAX_SIZE        = 8164 // 8192 - (IP header) - (UDP header)
	//UDP_CHUNK_MAX_DATA_SIZE   = 8144 // UDP_CHUNK_MAX_SIZE - ( 2 + 8 + 1 + 1)
	UDP_CHUNK_MAX_SIZE      = 1420
	UDP_CHUNK_MAX_DATA_SIZE = 1348 // UDP_CHUNK_MAX_SIZE - ( 2 + 8 + 1 + 1)
)

var (
	GELF_CHUNK_MAGIC_BYTES = []byte{0x1e, 0x0f}
)

// OvhHook represents an OVH PAAS Log
type OvhHook struct {
	async       bool
	token       string
	levels      []logrus.Level
	proto       Protocol
	compression CompressAlgo
}

// NewOvhHook returns a sync Hook
func NewOvhHook(ovhToken string, proto Protocol) (*OvhHook, error) {
	return newOvhHook(ovhToken, proto, false)
}

// NewAsyncOvhHook returns a async hook
func NewAsyncOvhHook(ovhToken string, proto Protocol) (*OvhHook, error) {
	return newOvhHook(ovhToken, proto, true)
}

// generic (ooops)
func newOvhHook(ovhToken string, proto Protocol, async bool) (*OvhHook, error) {
	hook := OvhHook{
		async:       async,
		token:       ovhToken,
		proto:       proto,
		levels:      logrus.AllLevels,
		compression: COMPRESSNONE,
	}
	return &hook, nil
}

// SetCompression set compression algorithm
func (hook *OvhHook) SetCompression(algo CompressAlgo) {
	hook.compression = algo
}

// TODO SetLevels

// Fire is called when a log event is fired.
func (hook *OvhHook) Fire(logrusEntry *logrus.Entry) error {
	e := Entry{
		entry:    logrusEntry,
		ovhToken: hook.token,
	}
	return e.send(hook.proto, hook.compression)
}

// Levels returns the available logging levels (interface impl)
func (hook *OvhHook) Levels() []logrus.Level {
	return hook.levels
}
