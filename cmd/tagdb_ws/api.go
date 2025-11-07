package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"dev.azure.com/trayport/Hackathon/_git/Q/internal/logger"
	"dev.azure.com/trayport/Hackathon/_git/Q/internal/tagdb"
)

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type TagKey struct {
	Tag string `json:"tag"`
	Key string `json:"key"`
}

func List(w http.ResponseWriter, r *http.Request) {
	logger.Infof("%s %s", r.Method, r.URL.String())

	// Read query string.

	queryString := r.URL.Query()
	var tags []string
	if rawTags := queryString.Get("tags"); rawTags != "" {
		tags = strings.Split(queryString.Get("tags"), ",")
	}

	// Connect to database.
	conn, err := tagdb.Connect()
	if err != nil {
		err = logger.Errorf("cannot connected to database because %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get result.
	items, err := conn.List(tags)
	if err != nil {
		err = logger.Errorf("cannot list tags `%v` because %s", tags, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Serialize.
	data, err := json.Marshal(&items)
	if err != nil {
		err = logger.Errorf("cannot serialize result because %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func Set(w http.ResponseWriter, r *http.Request) {
	logger.Infof("%s %s", r.Method, r.URL.String())

	// Read item.
	var kv KeyValue
	if err := json.NewDecoder(r.Body).Decode(&kv); err != nil {
		logger.Infof("cannot read body because %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Connect to db.
	conn, err := tagdb.Connect()
	if err != nil {
		err = logger.Errorf("cannot connected to database because %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create or update.
	if err := conn.Set(kv.Key, kv.Value); err != nil {
		err = logger.Errorf("cannot connected to database because %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Get(w http.ResponseWriter, r *http.Request) {
	logger.Infof("%s %s", r.Method, r.URL.String())

	// Read params.
	key := r.PathValue("key")
	if key == "" {
		msg := "cannot complete request because key not provided"
		logger.Info(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// Connect to db.
	conn, err := tagdb.Connect()
	if err != nil {
		err = logger.Errorf("cannot connected to database because %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get.
	item, found, err := conn.Get(key)
	if err != nil {
		err = logger.Errorf("cannot get from database because %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !found {
		logger.Infof("cannot find key %s", key)
		http.Error(w, "Resource not found", http.StatusNotFound)
		return
	}

	// Serialise.
	data, err := json.Marshal(&item)
	if err != nil {
		err = logger.Errorf("cannot serialize result because %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	logger.Infof("%s %s", r.Method, r.URL.String())

	// Read params.
	key := r.PathValue("key")
	if key == "" {
		msg := "cannot complete request because key not provided"
		logger.Info(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// Connect to db.
	conn, err := tagdb.Connect()
	if err != nil {
		err = logger.Errorf("cannot connected to database because %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set.
	if err := conn.Delete(key); err != nil {
		err = logger.Errorf("cannot delete from database because %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Tag(w http.ResponseWriter, r *http.Request) {
	logger.Infof("%s %s", r.Method, r.URL.String())

	// Read item.
	var tk TagKey
	if err := json.NewDecoder(r.Body).Decode(&tk); err != nil {
		logger.Infof("cannot read body because %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Connect to db.
	conn, err := tagdb.Connect()
	if err != nil {
		err = logger.Errorf("cannot connected to database because %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add tag.
	if err := conn.Tag(tk.Key, tk.Tag); err != nil {
		err = logger.Errorf("cannot add tag to database because %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Untag(w http.ResponseWriter, r *http.Request) {
	logger.Infof("%s %s", r.Method, r.URL.String())

	// Read params.
	var err error
	key := r.PathValue("key")
	if key == "" {
		err = errors.Join(err, errors.New("key is required"))
	}

	tag := r.PathValue("tag")
	if tag == "" {
		err = errors.Join(err, errors.New("tag is required"))
	}

	if err != nil {
		logger.Info(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Connect to db.
	conn, err := tagdb.Connect()
	if err != nil {
		err = logger.Errorf("cannot connected to database because %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Remove tag.
	if err := conn.Untag(key, tag); err != nil {
		err = logger.Errorf("cannot remove tag from database because %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
