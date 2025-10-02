package api

import "sync/atomic"

type ApiConfig struct {
	serverHits atomic.Int32
}
