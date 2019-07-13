package test_main

import (
    "testing"
    "net/http"
    "net/http/httptest"
    . "github.com/jelgar/login"
)

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder{
    req, _ := http.NewRequest(method, path, nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    return w
}

func Testping(t *test.T) {
    body := gin.H{
        "world": "Hello",
    }
   

    db := NewTestDB()
    env := &Env{db: db}
    router := SetupRouter(env)
}
