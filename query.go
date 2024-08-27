package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type queryCache struct {
	key     string
	data    []byte
	expired time.Time
}

var (
	SqlCount = make(map[string]int)
	QCache   sync.Map
)

func query(w http.ResponseWriter, r *http.Request, param *formData) {
	if len(param.Query) > 1 || param.Query == nil {
		Proxy.ServeHTTP(w, r)
		return
	}

	sql := param.Query[0]
	value, ok := QCache.Load(sql)
	if ok {
		if time.Now().Before(value.(queryCache).expired) {
			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Set("Content-Type", "text/plain; charset=x-user-defined")
			_, _ = w.Write(value.(queryCache).data)
			return
		}
	}

	ww := &writer{ResponseWriter: w}
	Proxy.ServeHTTP(ww, r)

	if _, ok := CacheSql[sql]; ok {
		QCache.Store(sql, queryCache{
			key:     sql,
			data:    ww.data,
			expired: time.Now().Add(5 * time.Minute),
		})
	}
	SqlCount[sql] = SqlCount[sql] + 1
}

func sql(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	data := make(map[string]int)

	for k, v := range SqlCount {
		if v < limit {
			continue
		}
		data[k] = v
	}

	bytes, _ := json.Marshal(data)
	_, _ = w.Write(bytes)
}
