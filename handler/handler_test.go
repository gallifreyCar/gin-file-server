package handler

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/uploadFile", uploadFileSingle)
	r.POST("/uploadFiles", uploadFiles)
	return r
}

func TestUploadFileSingle(t *testing.T) {

	// Create a temporary file to upload
	file, err := os.CreateTemp("", "example*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	// Write some data to the file
	file.WriteString("testing file upload")
	if err := file.Sync(); err != nil {
		t.Fatal(err)
	}
	// Create a new request with the file as the body of the request
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		t.Fatal(err)
	}

	// Create a cleanup function to remove the temporary file after the test finishes
	cleanup := func() {
		err := os.Remove(file.Name())
		if err != nil {
			t.Fatal(err)
		}
	}
	if _, err = file.Seek(0, 0); err != nil {
		cleanup()
		t.Fatal(err)
	}
	if _, err = io.Copy(part, file); err != nil {
		cleanup()
		t.Fatal(err)
	}
	if err = writer.Close(); err != nil {
		cleanup()
		t.Fatal(err)
	}
	// Create a test router
	router := setupRouter()
	// Create a response recorder to record the response
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/uploadFile", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	//Serve the request using the test router
	router.ServeHTTP(w, req)

	// Check the response status code
	if w.Code != http.StatusOK {
		cleanup()
		t.Errorf("unexpected status code %d", w.Code)
		return
	}

}

func TestUploadFiles(t *testing.T) {
	fieldName := "testfiles"

	// Create temporary files to upload
	var files []*os.File
	for i := 0; i < 3; i++ {
		tempFiles, err := os.CreateTemp("", "example*.txt")
		// Write some data to the file
		tempFiles.WriteString(fmt.Sprintf("testing file upload %v", i))

		if err != nil {
			t.Fatal(err)
		}
		files = append(files, tempFiles)
	}

	//Remove files from os
	defer func() {
		for _, file := range files {
			os.Remove(file.Name())
		}
	}()

	// Create a new request with the files as the body of the request
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for _, file := range files {

		part, err := writer.CreateFormFile(fieldName, filepath.Base(file.Name()))
		if err != nil {
			t.Fatal(err)
		}

		if _, err = file.Seek(0, 0); err != nil {
			t.Fatal(err)
		}
		if _, err = io.Copy(part, file); err != nil {
			t.Fatal(err)
		}
	}

	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	// Create a test router
	router := setupRouter()
	// Create a response recorder to record the response
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/uploadFiles", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.PostForm = url.Values{"fieldName": {fieldName}}
	// Serve the request using the test router
	router.ServeHTTP(w, req)

	// Check the response status code
	if w.Code != http.StatusOK {
		t.Errorf("unexpected status code %d", w.Code)
		return
	}

}
