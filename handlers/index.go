package handlers

import (
	"KeyValueDB/db"
	"KeyValueDB/util"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func IndexHandler(d db.IDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getHandler(d, w, r)
		case http.MethodPut:
			putHandler(d, w, r)
		case http.MethodDelete:
			deleteHandler(d, w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func getHandler(d db.IDatabase, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		v, err := d.GetAllKeys()
		if err != nil {
			http.Error(w, "error - getting all keys", http.StatusInternalServerError)
			fmt.Println("error - getting all keys: ", err)
			return
		}

		err = json.NewEncoder(w).Encode(v)
		if err != nil {
			http.Error(w, "error - encoding response", http.StatusInternalServerError)
			fmt.Println("error - encoding response: ", err)
			return
		}
		return
	}

	key := r.URL.Path[1:]

	v, err := d.Get(key)
	if err != nil {
		http.Error(w, "error - getting key", http.StatusInternalServerError)
		fmt.Printf("error - getting key %s: %s\n", key, err)
		return
	}

	if v == nil {
		w.WriteHeader(404)
		return
	}

	err = json.NewEncoder(w).Encode(v)
	if err != nil {
		http.Error(w, "error - encoding response", http.StatusInternalServerError)
		fmt.Println("error - encoding response: ", err)
		return
	}
}

func putHandler(d db.IDatabase, w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:]

	if len(key) == 0 {
		http.Error(w, "error - no key provided", http.StatusBadRequest)
		return
	}

	b, err := util.StreamToByte(r.Body)

	//Encase not valid JSON, stream to []byte instead of decoder to preserve original stream.
	reader := bytes.NewReader(b)

	data := make(map[string]interface{})
	err = json.NewDecoder(reader).Decode(&data)
	if err != nil {
		err = d.Set(key, string(b))
		if err != nil {
			http.Error(w, "error - putting kv pair", http.StatusInternalServerError)
			//This would all be logged away with splunk or similar, not logging directly to console.
			fmt.Printf("error - putting kvpair %s: %s\n", key, err)
			return
		}
		return
	}

	err = d.Set(key, data)
	if err != nil {
		http.Error(w, "error - putting json kv pair", http.StatusInternalServerError)
		fmt.Printf("error - putting json kvpair %s: %s\n", key, err)
		return
	}
}

func deleteHandler(d db.IDatabase, w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:]

	if len(key) == 0 {
		http.Error(w, "error - no key provided", http.StatusBadRequest)
		return
	}
	/*
			The following code that checking the existence of the key purely for a 404, is not necessary.
		    The data is to be discarded anyway, and so we don't need to see if it exists. Just delete it.
		    Removing the existence check is an optimization and simplification that I would push for.

		    I am including it here, as the spec specifically requires it.
	*/
	v, err := d.Get(key)
	if err != nil {
		http.Error(w, "error - getting key", http.StatusInternalServerError)
		fmt.Printf("error - getting key %s: %s\n", key, err)
		return
	}

	if v == nil {
		w.WriteHeader(404)
		return
	}

	err = d.Delete(key)
	if err != nil {
		http.Error(w, "error - deleting key", http.StatusInternalServerError)
		fmt.Printf("error - deleting key %s: %s\n", key, err)
		return
	}
}
