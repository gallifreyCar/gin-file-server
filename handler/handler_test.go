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
	"path"
	"path/filepath"
	"testing"
)

//todo Extract code to generate files.

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/uploadFile", UploadFileSingle)
	r.POST("/uploadFiles", UploadFiles)
	r.GET("/DownloadFile/:folder/:file_name", DownloadFile)
	return r
}

func TestUploadFileSingle(t *testing.T) {

	// Create a temporary file to upload
	file, err := os.CreateTemp("", "example*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	// Write some data to the file
	_, err = file.WriteString("testing file upload")
	if err != nil {
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
		if err != nil {
			t.Fatal(err)
		}

		// Write some data to the file
		_, err = tempFiles.WriteString(fmt.Sprintf("testing file upload %v", i))

		if err != nil {
			t.Fatal(err)
		}
		files = append(files, tempFiles)
	}

	//Remove files from os
	defer func() {

		for _, file := range files {
			err := file.Close()
			if err != nil {
				t.Fatal(err)
			}
			err = os.Remove(file.Name())
			if err != nil {
				t.Fatal(err)
			}
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

// todo use testCase struct to rewrite this unit test
func TestDownloadFile(t *testing.T) {
	folder := "single"
	fileName := "example617170769.txt"
	folder2 := "multiple"
	fileName2 := "example357758603.txt"
	downloadAndSave(folder, fileName, t)
	downloadAndSave(folder2, fileName2, t)
}

func downloadAndSave(folder, fileName string, t *testing.T) {
	// Create a test router
	router := setupRouter()
	// Create a response recorder to record the response
	w := httptest.NewRecorder()
	requestUrl := path.Join("/DownloadFile", folder, fileName)
	req, _ := http.NewRequest("GET", requestUrl, nil)

	// Serve the request using the test router
	router.ServeHTTP(w, req)

	// Check the response status code
	if w.Code == http.StatusNotFound {
		t.Logf("file not found %d", w.Code)
		return
	}
	if w.Code != http.StatusOK {
		t.Errorf("unexpected status code %d", w.Code)
		return
	}

	// Create the download directory if it doesn't exist
	downloadDirPath := path.Join("..", "target", "download", folder)
	// Save the files in the local
	saveFilePath := path.Join(downloadDirPath, fileName)

	err := os.MkdirAll(downloadDirPath, os.ModePerm)
	if err != nil {
		t.Errorf("error creating download directory: %v", err)
		return
	}

	// Create a new file to save the downloaded content
	out, err := os.Create(saveFilePath)
	if err != nil {
		t.Errorf("error creating file: %v", err)
		return
	}
	defer func() {
		err = out.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	// Write the downloaded content to the new file
	_, err = io.Copy(out, w.Body)
	if err != nil {
		t.Errorf("error writing file: %v", err)
		return
	}

}
