package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/gorilla/mux"
    "strconv"
    "io"
    "io/ioutil"
)

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "RESTful api test /data : /data/{:id}")
}

func DataShowAll(w http.ResponseWriter, r *http.Request) {

    RepoFindAllFiles()

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(data.files); err != nil {
        panic(err)
    }
}

func DataShow(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    dataId := vars["dataId"]
    indx, _ :=strconv.Atoi(dataId)

    t := RepoFindFile(indx)

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(t); err != nil {
        panic(err)
    }
}

func DataCreate(w http.ResponseWriter, r *http.Request) {
    var file File
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        panic(err)
    }
    if err := r.Body.Close(); err != nil {
        panic(err)
    }
    if err := json.Unmarshal(body, &file); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422)
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
    }

    t := RepoCreateFile(&file)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(t); err != nil {
        panic(err)
    }
}