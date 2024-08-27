package main

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type connCache struct {
	key     string
	data    []byte
	expired time.Time
}

var CCache sync.Map

func connect(w http.ResponseWriter, r *http.Request, param *formData) {
	value, ok := CCache.Load(param.Key)

	if ok {
		if time.Now().Before(value.(connCache).expired) {
			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Set("Content-Type", "text/plain; charset=x-user-defined")
			_, _ = w.Write(value.(connCache).data)
			return
		}
	}
	log.Println("conn cache miss", param.Key)

	ww := &writer{ResponseWriter: w}
	Proxy.ServeHTTP(ww, r)

	CCache.Store(param.Key, connCache{
		key:     param.Key,
		data:    ww.data,
		expired: time.Now().Add(5 * time.Minute),
	})
}
