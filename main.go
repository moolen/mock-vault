package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
)

type vaultData struct {
	User  string `json:"user"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

func main() {
	port := os.Getenv("VAULT_PORT")
	vaultPath := os.Getenv("VAULT_PATH")
	if len(port) == 0 || len(vaultPath) == 0 {
		panic(errors.New("VAULT_PORT || VAULT_PATH missing"))
	}
	err := os.MkdirAll(vaultPath, os.ModePerm)
	if err != nil {
		panic(err)
	}
	log.SetFormatter(&log.JSONFormatter{})
	http.HandleFunc("/vault", func(res http.ResponseWriter, req *http.Request) {
		log.WithFields(log.Fields{
			"method": req.Method,
			"uri":    req.RequestURI,
		}).Infoln("recv request")
		if req.Method == http.MethodPost {
			create(vaultPath, res, req)
		}
		if req.Method == http.MethodGet {
			get(vaultPath, res, req)
		}
	})
	log.Infof("listening on :%s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

func handleError(err error, res http.ResponseWriter) {
	res.WriteHeader(500)
	fmt.Fprintf(res, "error: %s", err)
	log.Errorf("error: %s", err)
}

func create(vaultPath string, res http.ResponseWriter, req *http.Request) {
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		handleError(err, res)
		return
	}
	defer req.Body.Close()
	var vault vaultData
	err = json.Unmarshal(data, &vault)
	if err != nil {
		handleError(err, res)
		return
	}
	if vault.User == "" {
		handleError(errors.New("missing user"), res)
		return
	}
	if vault.Value == "" {
		handleError(errors.New("missing value LOL"), res)
		return
	}
	fileBasename := path.Base(vault.User)
	if err != nil {
		handleError(err, res)
		return
	}
	if len(vault.Key) == 0 {
		handleError(errors.New("empty key not allowed"), res)
		return
	}
	cipher, err := encrypt(vault.Key, vault.Value)
	if err != nil {
		handleError(err, res)
		return
	}
	err = ioutil.WriteFile(path.Join(vaultPath, fmt.Sprintf("vault-%s", fileBasename)), []byte(cipher), os.ModePerm)
	if err != nil {
		handleError(err, res)
		return
	}
	res.WriteHeader(201)
}

func get(vaultPath string, res http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	user := q["user"]
	key := q["key"]
	if len(user) == 0 || len(user[0]) == 0 {
		handleError(errors.New("no user provided"), res)
		return
	}
	if len(key) == 0 || len(key[0]) == 0 {
		handleError(errors.New("no key provided"), res)
		return
	}
	fileBasename := path.Base(user[0])
	requestedFile := path.Join(vaultPath, fmt.Sprintf("vault-%s", fileBasename))
	cipher, err := ioutil.ReadFile(requestedFile)
	if err != nil {
		handleError(err, res)
		return
	}

	plaintext, err := decrypt(key[0], string(cipher))
	if err != nil {
		handleError(err, res)
		return
	}
	log.Infof("plain: %#v", plaintext)
	fmt.Fprintf(res, "%s", plaintext)
}
