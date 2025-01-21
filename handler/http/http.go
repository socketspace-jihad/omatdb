package http

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/socketspace-jihad/omatdb/engine"
)

type KVHttpHandler struct {
	*engine.KVStorage
	*http.ServeMux
}

type kvBodyHandler struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

func (k *kvBodyHandler) isValid() bool {
	if k.Key == "" {
		return false
	}
	if k.Value == nil {
		return false
	}
	return true
}

func NewKVHandler(kvs *engine.KVStorage) *KVHttpHandler {
	return &KVHttpHandler{
		KVStorage: kvs,
		ServeMux:  &http.ServeMux{},
	}
}

func (k *KVHttpHandler) Run() {

	k.ServeMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello from omatdb"))
	})

	k.ServeMux.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query().Get("key")
		if params == "" {
			http.Error(w, "query params 'key' must be defined", http.StatusBadRequest)
			return
		}
		val, err := k.KVStorage.Get(params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp := kvBodyHandler{
			Key:   params,
			Value: val,
		}
		json.NewEncoder(w).Encode(resp)
	})

	k.ServeMux.HandleFunc("/store", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var data kvBodyHandler
		if err := json.Unmarshal(body, &data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !data.isValid() {
			http.Error(w, "kv body is not valid", http.StatusBadRequest)
			return
		}
		if err := k.KVStorage.Store(data.Key, data.Value); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(data)
	})

	k.ServeMux.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var data kvBodyHandler
		if err := json.Unmarshal(body, &data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !data.isValid() {
			http.Error(w, "kv body is not valid", http.StatusBadRequest)
			return
		}
		if err := k.KVStorage.Update(data.Key, data.Value); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(data)
	})

	k.ServeMux.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var data kvBodyHandler
		if err := json.Unmarshal(body, &data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if data.Key == "" {
			http.Error(w, "key in body must exists", http.StatusBadRequest)
			return
		}
		if err := k.KVStorage.Delete(data.Key); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(data)
	})

	log.Println("omatdb listening on port: 3300")
	if err := http.ListenAndServe(":3300", k.ServeMux); err != nil {
		log.Fatalln(err.Error())
	}

}
