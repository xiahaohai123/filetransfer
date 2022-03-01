package filetransfer_test

import (
	"net/http"
	"net/http/httptest"
	"summersea.top/filetransfer"
	"testing"
)

func TestUploadFile(t *testing.T) {
	t.Run("return status 200 from /file/upload/initialization", func(t *testing.T) {

		request, err := http.NewRequest(http.MethodPost, "/file/upload/initialization", nil)
		if err != nil {
			t.Fatalf("problem new request %v", err)
		}
		response := httptest.NewRecorder()

		fileServer := filetransfer.NewFileServer()
		fileServer.ServeHTTP(response, request)

		if response.Code != http.StatusOK {
			t.Errorf("got %d but want %d from response", response.Code, http.StatusOK)
		}
	})
}
