package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	setup(6)
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup(n int) {
	path := "teststorage"
	err := os.Mkdir(path, 0777)
	if err != nil {
		log.Fatal(err)
	}
	exts := []string{".jpg", ".png", ".gif"}
	for i := 0; i < n; i++ {
		_, err = os.Create(fmt.Sprintf("./%s/%s%s", path, strconv.Itoa(i), exts[i%3]))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func teardown() {
	path := "teststorage"
	err := os.RemoveAll(path)
	if err != nil {
		log.Fatal(err)
	}
}

func TestFileServer(t *testing.T) {
	r, err := http.NewRequest("GET", "/0.jpg", nil)
	assert.NoError(t, err)
	rec := httptest.NewRecorder()
	uploadHandler := UploadHandler{UploadDir: "teststorage"}
	dirToServe := http.Dir(uploadHandler.UploadDir)
	http.FileServer(dirToServe).ServeHTTP(rec, r)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestGetAllFiles(t *testing.T) {
	r, err := http.NewRequest("GET", "/get/", nil)
	assert.NoError(t, err)
	rec := httptest.NewRecorder()
	uh := UploadHandler{UploadDir: "teststorage"}
	uh.ServeHTTP(rec, r)
	require.Equal(t, http.StatusOK, rec.Code)
	var res []fileInfo
	err = json.Unmarshal(rec.Body.Bytes(), &res)
	require.NoError(t, err)
	require.Equal(t, 6, len(res))
	require.Equal(t, "jpg", res[0].Extension)
}

func TestGetFilesWithExt(t *testing.T) {
	r, err := http.NewRequest("GET", "/get/?ext=jpg", nil)
	assert.NoError(t, err)
	rec := httptest.NewRecorder()
	uh := UploadHandler{UploadDir: "teststorage"}
	uh.ServeHTTP(rec, r)
	require.Equal(t, http.StatusOK, rec.Code)
	var res []fileInfo
	err = json.Unmarshal(rec.Body.Bytes(), &res)
	require.NoError(t, err)
	require.Equal(t, 2, len(res))
	require.Equal(t, "jpg", res[0].Extension)
	require.Equal(t, "jpg", res[1].Extension)
}

func TestGetNoFilesWithExt(t *testing.T) {
	r, err := http.NewRequest("GET", "/get/?ext=txt", nil)
	assert.NoError(t, err)
	rec := httptest.NewRecorder()
	uh := UploadHandler{UploadDir: "teststorage"}
	uh.ServeHTTP(rec, r)
	require.Equal(t, http.StatusOK, rec.Code)
	var res []fileInfo
	err = json.Unmarshal(rec.Body.Bytes(), &res)
	require.NoError(t, err)
	require.Equal(t, 0, len(res))
}
