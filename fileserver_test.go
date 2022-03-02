package filetransfer_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"summersea.top/filetransfer"
	"testing"
)

func TestUploadFileInitialise(t *testing.T) {
	url := "/file/upload/initialization"

	t.Run("return status when input some param", func(t *testing.T) {
		testTables := []struct {
			uploadInitReqBody  filetransfer.UploadInitReqBody
			wantResponseStatus int
		}{
			{
				uploadInitReqBody:  filetransfer.UploadInitReqBody{},
				wantResponseStatus: http.StatusBadRequest,
			},
			{
				uploadInitReqBody: filetransfer.UploadInitReqBody{
					Resource: filetransfer.Resource{
						Address: "",
						Port:    255,
						Account: filetransfer.Account{
							Name:     "a",
							Password: "pwd",
						},
					},
					Path: "/root",
				},
				wantResponseStatus: http.StatusBadRequest,
			},
			{
				uploadInitReqBody: filetransfer.UploadInitReqBody{
					Resource: filetransfer.Resource{
						Address: "10.12.1.12",
						Port:    -1,
						Account: filetransfer.Account{
							Name:     "a",
							Password: "pwd",
						},
					},
					Path: "/root",
				},
				wantResponseStatus: http.StatusBadRequest,
			},
			{
				uploadInitReqBody: filetransfer.UploadInitReqBody{
					Resource: filetransfer.Resource{
						Address: "10.12.1.12",
						Port:    65536,
						Account: filetransfer.Account{
							Name:     "a",
							Password: "pwd",
						},
					},
					Path: "/root",
				},
				wantResponseStatus: http.StatusBadRequest,
			},
			{
				uploadInitReqBody: filetransfer.UploadInitReqBody{
					Resource: filetransfer.Resource{
						Address: "10.12.1.12",
						Port:    255,
						Account: filetransfer.Account{
							Name:     "",
							Password: "pwd",
						},
					},
					Path: "/root",
				},
				wantResponseStatus: http.StatusBadRequest,
			},
			{
				uploadInitReqBody: filetransfer.UploadInitReqBody{
					Resource: filetransfer.Resource{
						Address: "10.12.1.12",
						Port:    256,
						Account: filetransfer.Account{
							Name:     "a",
							Password: "",
						},
					},
					Path: "/root",
				},
				wantResponseStatus: http.StatusBadRequest,
			},
			{
				uploadInitReqBody: filetransfer.UploadInitReqBody{
					Resource: filetransfer.Resource{
						Address: "10.12.1.12",
						Port:    256,
						Account: filetransfer.Account{
							Name:     "a",
							Password: "pwddd",
						},
					},
					Path: "",
				},
				wantResponseStatus: http.StatusBadRequest,
			},
			{
				uploadInitReqBody: filetransfer.UploadInitReqBody{
					Resource: filetransfer.Resource{
						Address: "10.12.1.12",
						Port:    256,
						Account: filetransfer.Account{
							Name:     "a",
							Password: "pwddd",
						},
					},
					Path: "pwd",
				},
				wantResponseStatus: http.StatusBadRequest,
			},
			{
				uploadInitReqBody: filetransfer.UploadInitReqBody{
					Resource: filetransfer.Resource{
						Address: "10.12.1.12",
						Port:    256,
						Account: filetransfer.Account{
							Name:     "a",
							Password: "pwddd",
						},
					},
					Path: "/root/pwd",
				},
				wantResponseStatus: http.StatusOK,
			},
		}

		for _, test := range testTables {
			testHttpStatus(t, url, test.uploadInitReqBody, test.wantResponseStatus)
		}
	})
}

func testHttpStatus(t *testing.T, url string, body filetransfer.UploadInitReqBody, wantStatus int) {
	t.Helper()
	requestBody := new(bytes.Buffer)
	err := json.NewEncoder(requestBody).Encode(body)
	if err != nil {
		t.Fatalf("wrong body: %+v", body)
	}
	request, err := http.NewRequest(http.MethodPost, url, requestBody)
	if err != nil {
		t.Fatalf("problem new request %v", err)
	}
	response := httptest.NewRecorder()

	fileServer := filetransfer.NewFileServer()
	fileServer.ServeHTTP(response, request)

	if response.Code != wantStatus {
		t.Errorf("test case: %+v \n got %d but want %d from response", body, response.Code, wantStatus)
	}
}
