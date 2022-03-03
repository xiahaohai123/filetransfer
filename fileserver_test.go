package filetransfer_test

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"summersea.top/filetransfer"
	"testing"
)

func TestUploadFileInitialise(t *testing.T) {
	url := "/file/upload/initialization"
	fileServer := filetransfer.NewFileServer(&StubStore{})
	correctJson := `{"resource":{"address":"summersea1.top","port":22,"account":{"name":"ccc","password":"pwd"}},"path":"/root"}`

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
			requestBody := new(bytes.Buffer)
			err := json.NewEncoder(requestBody).Encode(test.uploadInitReqBody)
			if err != nil {
				t.Fatalf("wrong body: %+v", test.uploadInitReqBody)
			}
			request := newPostRequest(url, requestBody)
			response := httptest.NewRecorder()
			fileServer.ServeHTTP(response, request)
			testHttpStatus(t, requestBody, response.Code, test.wantResponseStatus)
		}
	})

	t.Run("input wrong type", func(t *testing.T) {

		reader := strings.NewReader(`"aa":"aa"`)
		request := newPostRequest(url, reader)
		response := httptest.NewRecorder()
		fileServer.ServeHTTP(response, request)
		assertIntEquals(t, response.Code, http.StatusBadRequest)

		request1 := newPostRequest(url, nil)
		response1 := httptest.NewRecorder()
		fileServer.ServeHTTP(response1, request1)
		assertIntEquals(t, response1.Code, http.StatusBadRequest)
	})

	t.Run("using wrong method", func(t *testing.T) {
		request := newGetRequest(url)
		response := httptest.NewRecorder()
		fileServer.ServeHTTP(response, request)
		assertIntEquals(t, response.Code, http.StatusForbidden)
	})

	t.Run("get task id when correctly use method", func(t *testing.T) {
		request := newPostRequest(url, strings.NewReader(correctJson))
		response := httptest.NewRecorder()
		fileServer.ServeHTTP(response, request)
		assertIntEquals(t, response.Code, http.StatusOK)

		sc := bufio.NewScanner(response.Body)
		sc.Scan()
		taskId := sc.Text()
		if len(strings.Split(taskId, "-")) != 5 {
			t.Errorf("want uuid but got other '%s'", taskId)
		}
		assertContent(t, taskId)
	})
}

type StubStore struct {
	taskId string
}

func (s *StubStore) GetUploadData(taskId string) *filetransfer.UploadData {
	if s.taskId == taskId {
		return &filetransfer.UploadData{}
	}
	return nil
}

func TestUploadFile(t *testing.T) {
	url := "/file/upload"
	fileServer := filetransfer.NewFileServer(&StubStore{})

	t.Run("api exists", func(t *testing.T) {
		request := newPostRequest(url, nil)
		response := httptest.NewRecorder()
		fileServer.ServeHTTP(response, request)
		assertIntNotEquals(t, response.Code, http.StatusNotFound)

		request = newGetRequest(url)
		response = httptest.NewRecorder()
		fileServer.ServeHTTP(response, request)
		assertIntEquals(t, response.Code, http.StatusForbidden)
	})

	t.Run("can not find task id in system", func(t *testing.T) {
		taskId := "a55b14b2-fb55-40b8-8311-6d1e7d949fb5"
		uploadUrl := fmt.Sprintf("%s?taskId=%s", url, taskId)
		request := newPostRequest(uploadUrl, nil)
		response := httptest.NewRecorder()
		fileServer.ServeHTTP(response, request)
		assertIntEquals(t, response.Code, http.StatusBadRequest)
		assertStringEquals(t, response.Header().Get("Content-Type"), filetransfer.ContentTypeJsonValue)

		wantErrorBody := filetransfer.ErrorBody{
			Error: filetransfer.ErrorContent{
				Message: "The task id is not found.",
				Code:    "ResourceNotFound",
			},
		}
		var gotErrorBody filetransfer.ErrorBody
		_ = json.NewDecoder(response.Body).Decode(&gotErrorBody)
		assertStructEquals(t, gotErrorBody, wantErrorBody)
	})

	t.Run("can find task id in system", func(t *testing.T) {
		taskId := uuid.NewV4().String()
		fileServer := filetransfer.NewFileServer(&StubStore{taskId: taskId})
		uploadUrl := fmt.Sprintf("%s?taskId=%s", url, taskId)
		request := newPostRequest(uploadUrl, nil)
		response := httptest.NewRecorder()
		fileServer.ServeHTTP(response, request)
		assertIntEquals(t, response.Code, http.StatusNoContent)
	})
}

func testHttpStatus(t *testing.T, requestBody io.Reader, got, wantStatus int) {
	t.Helper()
	if got != wantStatus {
		t.Errorf("test case: %v \n got %d but want %d from response", requestBody, got, wantStatus)
	}
}

func newPostRequest(url string, requestBody io.Reader) *http.Request {
	request, _ := http.NewRequest(http.MethodPost, url, requestBody)
	return request
}

func newGetRequest(url string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	return request
}

func assertIntEquals(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("want %d but got %d", want, got)
	}
}

func assertIntNotEquals(t *testing.T, got, notWant int) {
	t.Helper()
	if got == notWant {
		t.Errorf("don't want %d bug got", notWant)
	}
}

func assertStringEquals(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("want '%s' bug got '%s'", want, got)
	}
}

func assertStructEquals(t *testing.T, got, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %+v but got %+v", want, got)
	}
}

func assertContent(t *testing.T, content string) {
	t.Helper()
	if content == "" {
		t.Errorf("want a content but got empty")
	}
}
