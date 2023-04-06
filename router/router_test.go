package router

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetupRouter(t *testing.T) {

	testCases := []struct {
		name        string
		endpoint    string
		method      string
		contendType string
		expect      int
	}{
		{"TestUploadFileSingle", "/uploadFile", http.MethodPost, "multipart/form-data", http.StatusBadRequest},
		{"TestUploadFiles", "/uploadFiles", http.MethodPost, "multipart/form-data", http.StatusBadRequest},
		{"TestDownloadFile", "/downloadFile/single/myfile.txt", "", http.MethodGet, http.StatusNotFound},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := setupRouter()
			w := httptest.NewRecorder()
			resp, err := http.NewRequest(tc.method, tc.endpoint, nil)
			resp.Header.Set("Content-Type", tc.contendType)
			router.ServeHTTP(w, resp)
			assert.Nil(t, err)
			assert.Equal(t, tc.expect, w.Code)
		})
	}

}
