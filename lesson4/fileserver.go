package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type UploadHandler struct {
	HostAddr  string
	UploadDir string
}

type fileInfo struct {
	Name      string `json:"name"`
	Extension string `json:"extension"`
	Size      int64  `json:"size"`
}

func (uh *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var response []byte
	var err error
	ext := r.FormValue("ext")
	response, err = uh.HandleGetFiles(ctx, r, ext)
	if err != nil {
		http.Error(w, "can't handle request: "+err.Error(), http.StatusInternalServerError)
	}
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(response)
	if err != nil {
		http.Error(w, "can't handle request: "+err.Error(), http.StatusInternalServerError)
	}

}

func (uh *UploadHandler) HandleGetFiles(ctx context.Context, r *http.Request, ext string) ([]byte, error) {
	dir, err := os.ReadDir(uh.UploadDir)
	if err != nil {
		return nil, fmt.Errorf("can't open dir %s", err)
	}
	files := make([]fileInfo, 0)
	fis := fileInfo{}
	for _, f := range dir {
		fis.Name = f.Name()
		fis.Extension = strings.Split(fis.Name, ".")[1]
		info, errr := f.Info()
		if errr != nil {
			log.Printf("error while reading file: %s\n", err)
		}
		fis.Size = info.Size()
		if ext == "" || ext == fis.Extension {
			files = append(files, fis)
		}
	}
	// add filter
	js, err := json.Marshal(files)
	if err != nil {
		return nil, fmt.Errorf("can't encode json: %s", err)
	}
	return js, nil
}

func main() {
	uploadHandler := UploadHandler{UploadDir: "storage"}
	dirToServe := http.Dir(uploadHandler.UploadDir)
	http.Handle("/", http.FileServer(dirToServe))
	http.Handle("/get/", &uploadHandler)

	fs := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err := fs.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
