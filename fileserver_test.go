package filetransfer_test

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/kirinlabs/utils/str"
	uuid "github.com/satori/go.uuid"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"summersea.top/filetransfer"
	"testing"
)

const correctJson = `{"resource":{"address":"summersea1.top","port":22,"account":{"name":"ccc","password":"pwd"}},"path":"/root","filename":"test.txt"}`
const uploadUrl = "/file/upload"
const downloadUrl = "/file/download"

type initTestCase struct {
	requestBody        interface{}
	wantResponseStatus int
	wantResponseBody   interface{}
}

func TestUploadFileInitialise(t *testing.T) {
	url := "/file/upload/initialization"
	fileServer := filetransfer.NewFileServer(&StubAdapter{})

	t.Run("return status when input some param", func(t *testing.T) {
		errResponseBody := getInvalidErrBody()
		testTables := []initTestCase{
			{
				requestBody:        filetransfer.UploadInitReqBody{},
				wantResponseStatus: http.StatusBadRequest,
				wantResponseBody:   errResponseBody,
			},
			{
				requestBody: filetransfer.UploadInitReqBody{
					Resource: filetransfer.Resource{
						Address: "",
						Port:    255,
						Account: filetransfer.Account{
							Name:     "a",
							Password: "pwd",
						},
					},
					Path:     "/home/test",
					Filename: "testFile.txt",
				},
				wantResponseStatus: http.StatusBadRequest,
				wantResponseBody:   errResponseBody,
			},
			{
				requestBody: filetransfer.UploadInitReqBody{
					Resource: filetransfer.Resource{
						Address: "10.12.1.12",
						Port:    -1,
						Account: filetransfer.Account{
							Name:     "a",
							Password: "pwd",
						},
					},
					Path:     "/root",
					Filename: "testFile.txt",
				},
				wantResponseStatus: http.StatusBadRequest,
				wantResponseBody:   errResponseBody,
			},
			{
				requestBody: filetransfer.UploadInitReqBody{
					Resource: filetransfer.Resource{
						Address: "10.12.1.12",
						Port:    65536,
						Account: filetransfer.Account{
							Name:     "a",
							Password: "pwd",
						},
					},
					Path:     "/root",
					Filename: "testFile.txt",
				},
				wantResponseStatus: http.StatusBadRequest,
				wantResponseBody:   errResponseBody,
			},
			{
				requestBody: filetransfer.UploadInitReqBody{
					Resource: filetransfer.Resource{
						Address: "10.12.1.12",
						Port:    255,
						Account: filetransfer.Account{
							Name:     "",
							Password: "pwd",
						},
					},
					Path:     "/root",
					Filename: "testFile.txt",
				},
				wantResponseStatus: http.StatusBadRequest,
				wantResponseBody:   errResponseBody,
			},
			{
				requestBody: filetransfer.UploadInitReqBody{
					Resource: filetransfer.Resource{
						Address: "10.12.1.12",
						Port:    256,
						Account: filetransfer.Account{
							Name:     "a",
							Password: "",
						},
					},
					Path:     "/root",
					Filename: "testFile.txt",
				},
				wantResponseStatus: http.StatusBadRequest,
				wantResponseBody:   errResponseBody,
			},
			{
				requestBody: filetransfer.UploadInitReqBody{
					Resource: filetransfer.Resource{
						Address: "10.12.1.12",
						Port:    256,
						Account: filetransfer.Account{
							Name:     "a",
							Password: "pwddd",
						},
					},
					Path:     "",
					Filename: "testFile.txt",
				},
				wantResponseStatus: http.StatusBadRequest,
				wantResponseBody:   errResponseBody,
			},
			{
				requestBody: filetransfer.UploadInitReqBody{
					Resource: filetransfer.Resource{
						Address: "10.12.1.12",
						Port:    256,
						Account: filetransfer.Account{
							Name:     "a",
							Password: "pwddd",
						},
					},
					Path:     "pwd",
					Filename: "testFile.txt",
				},
				wantResponseStatus: http.StatusBadRequest,
				wantResponseBody:   errResponseBody,
			},
			{
				requestBody: filetransfer.UploadInitReqBody{
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
				wantResponseStatus: http.StatusBadRequest,
				wantResponseBody:   errResponseBody,
			},
			{
				requestBody: filetransfer.UploadInitReqBody{
					Resource: filetransfer.Resource{
						Address: "10.12.1.12",
						Port:    256,
						Account: filetransfer.Account{
							Name:     "a",
							Password: "pwddd",
						},
					},
					Path:     "/root/pwd",
					Filename: "testFile.txt",
				},
				wantResponseStatus: http.StatusOK,
			},
		}

		for _, test := range testTables {
			response := testCase(t, test, url, fileServer)
			if response.Code != http.StatusOK {
				var gotErrorBody filetransfer.ErrorBody
				_ = json.NewDecoder(response.Body).Decode(&gotErrorBody)
				assertStructEquals(t, gotErrorBody, test.wantResponseBody)
			} else {
				okBody := extractOkBody(response.Body)
				_, exist := okBody.Data["taskId"]
				assertTrue(t, exist)
			}
		}
	})

	t.Run("input wrong type", func(t *testing.T) {

		reader := strings.NewReader(`"aa":"aa"`)
		request := newPostRequestReader(url, reader)
		response := httptest.NewRecorder()
		fileServer.ServeHTTP(response, request)
		assertIntEquals(t, response.Code, http.StatusBadRequest)

		request1 := newPostRequestReader(url, nil)
		response1 := httptest.NewRecorder()
		fileServer.ServeHTTP(response1, request1)
		assertIntEquals(t, response1.Code, http.StatusBadRequest)
	})

	t.Run("using wrong method", func(t *testing.T) {
		request := newGetRequest(url)
		response := httptest.NewRecorder()
		fileServer.ServeHTTP(response, request)
		assertIntEquals(t, response.Code, http.StatusNotFound)
	})

	t.Run("get task id when correctly use method", func(t *testing.T) {
		request := newPostRequestReader(url, strings.NewReader(correctJson))
		response := httptest.NewRecorder()
		fileServer.ServeHTTP(response, request)
		assertIntEquals(t, response.Code, http.StatusOK)

		okBody := extractOkBody(response.Body)
		taskId := okBody.Data["taskId"].(string)
		if len(strings.Split(taskId, "-")) != 5 {
			t.Errorf("want uuid but got other '%s'", taskId)
		}
		assertContent(t, taskId)
	})
}

type StubAdapter struct {
	uploadTaskId   string
	filename       string
	downloadTaskId string
}

type fileRollback struct {
	*os.File
}

func (f *fileRollback) RollBack() error {
	return os.Remove(f.Name())
}

func (s *StubAdapter) GetUploadChannel(taskId string) (filetransfer.WriteCloseRollback, error) {
	if s.uploadTaskId == taskId {
		rollback := fileRollback{}
		file, _ := os.OpenFile(s.filename, os.O_RDWR|os.O_CREATE, 0777)
		rollback.File = file
		return &rollback, nil
	}
	return nil, nil
}

func (s *StubAdapter) SaveUploadData(taskId string, uploadData filetransfer.UploadData) {
	s.uploadTaskId = taskId
}

func (s *StubAdapter) IsUploadTaskExist(taskId string) bool {
	return s.uploadTaskId == taskId
}

func (s *StubAdapter) IsDownloadTaskExist(taskId string) bool {
	return s.downloadTaskId == taskId
}

func (s *StubAdapter) GetDownloadChannelFilename(taskId string) (io.ReadCloser, string, error) {
	if s.downloadTaskId == taskId {
		file, _ := os.OpenFile(s.filename, os.O_RDWR, 0666)
		return file, s.filename, nil
	}
	return nil, "", nil
}

func TestUploadFile(t *testing.T) {
	url := uploadUrl
	fileServer := filetransfer.NewFileServer(&StubAdapter{uploadTaskId: uuid.NewV4().String()})

	t.Run("api exists", func(t *testing.T) {
		request := newGetRequest(url)
		response := httptest.NewRecorder()
		fileServer.ServeHTTP(response, request)
		assertIntEquals(t, response.Code, http.StatusNotFound)
	})

	t.Run("upload", func(t *testing.T) {
		taskId := uuid.NewV4().String()
		contentFilename, deleteContentFile := createTempFileWithContent(t)
		defer deleteContentFile()
		dstFilename := createRandomFilename("tempFile", ".txt")
		fileServer := filetransfer.NewFileServer(&StubAdapter{uploadTaskId: taskId, filename: dstFilename})
		contentFile, _ := os.Open(contentFilename)
		defer contentFile.Close()
		uploadUrl := fmt.Sprintf("%s?taskId=%s", url, taskId)
		request := newPostRequestReader(uploadUrl, contentFile)
		response := httptest.NewRecorder()
		fileServer.ServeHTTP(response, request)
		assertIntEquals(t, response.Code, http.StatusNoContent)
		assertFileContentEquals(t, contentFilename, dstFilename)
		_ = os.Remove(dstFilename)
	})
}

func TestUploadByIntegration(t *testing.T) {
	urlInit := "/file/upload/initialization"
	urlUpload := "/file/upload"
	dstFilename := createRandomFilename("tempFile", ".txt")
	fileServer := filetransfer.NewFileServer(&StubAdapter{filename: dstFilename})
	request := newPostRequestReader(urlInit, strings.NewReader(correctJson))
	response := httptest.NewRecorder()
	fileServer.ServeHTTP(response, request)
	assertIntEquals(t, response.Code, http.StatusOK)
	okBody := extractOkBody(response.Body)
	taskId := okBody.Data["taskId"]
	uploadUrl := fmt.Sprintf("%s?taskId=%s", urlUpload, taskId)

	contentFilename, deleteContentFile := createTempFileWithContent(t)
	defer deleteContentFile()
	contentFile, _ := os.Open(contentFilename)
	defer contentFile.Close()
	uploadReq := newPostRequestReader(uploadUrl, contentFile)
	uploadResponse := httptest.NewRecorder()
	fileServer.ServeHTTP(uploadResponse, uploadReq)
	assertIntEquals(t, uploadResponse.Code, http.StatusNoContent)
	assertFileContentEquals(t, contentFilename, dstFilename)
	_ = os.Remove(dstFilename)
}

func TestDownloadFileInit(t *testing.T) {
	url := "/file/download/initialization"
	fileServer := filetransfer.NewFileServer(&StubAdapter{})
	t.Run("test bad request", func(t *testing.T) {
		errResponseBody := getInvalidErrBody()
		testCases := []initTestCase{
			{
				filetransfer.DownloadInitReqBody{},
				http.StatusBadRequest,
				errResponseBody,
			},
			{
				filetransfer.DownloadInitReqBody{
					Resource: filetransfer.Resource{
						Address: "",
						Port:    255,
						Account: filetransfer.Account{
							Name:     "a",
							Password: "pwd",
						},
					},
					Path: "/home/test",
				},
				http.StatusBadRequest,
				errResponseBody,
			},
			{
				filetransfer.DownloadInitReqBody{
					Resource: filetransfer.Resource{
						Address: "addr",
						Port:    0,
						Account: filetransfer.Account{
							Name:     "a",
							Password: "pwd",
						},
					},
					Path: "/home/test",
				},
				http.StatusBadRequest,
				errResponseBody,
			},
			{
				filetransfer.DownloadInitReqBody{
					Resource: filetransfer.Resource{
						Address: "addr",
						Port:    22,
						Account: filetransfer.Account{
							Name:     "",
							Password: "pwd",
						},
					},
					Path: "/home/test",
				},
				http.StatusBadRequest,
				errResponseBody,
			},
			{
				filetransfer.DownloadInitReqBody{
					Resource: filetransfer.Resource{
						Address: "addr",
						Port:    22,
						Account: filetransfer.Account{
							Name:     "a",
							Password: "",
						},
					},
					Path: "/home/test",
				},
				http.StatusBadRequest,
				errResponseBody,
			},
			{
				filetransfer.DownloadInitReqBody{
					Resource: filetransfer.Resource{
						Address: "addr",
						Port:    22,
						Account: filetransfer.Account{
							Name:     "a",
							Password: "pwd",
						},
					},
					Path: "",
				},
				http.StatusBadRequest,
				errResponseBody,
			},
			{
				filetransfer.DownloadInitReqBody{
					Resource: filetransfer.Resource{
						Address: "addr",
						Port:    22,
						Account: filetransfer.Account{
							Name:     "a",
							Password: "pwd",
						},
					},
					Path: "/home/test/",
				},
				http.StatusBadRequest,
				errResponseBody,
			},
		}
		for _, test := range testCases {
			response := testCase(t, test, url, fileServer)
			var gotErrorBody filetransfer.ErrorBody
			_ = json.NewDecoder(response.Body).Decode(&gotErrorBody)
			assertStructEquals(t, gotErrorBody, test.wantResponseBody)
		}
	})

	t.Run("test get Task Id", func(t *testing.T) {
		downloadInitReqBody := filetransfer.DownloadInitReqBody{
			Resource: filetransfer.Resource{
				Address: "addr",
				Port:    22,
				Account: filetransfer.Account{
					Name:     "test",
					Password: "pwd",
				},
			},
			Path: "/home/test",
		}
		request := newPostReqBody(t, url, downloadInitReqBody)
		response := httptest.NewRecorder()
		fileServer.ServeHTTP(response, request)
		assertIntEquals(t, response.Code, http.StatusOK)
		okBody := extractOkBody(response.Body)
		taskId := okBody.Data["taskId"].(string)
		assertContent(t, taskId)
	})
}

func TestDownloadFile(t *testing.T) {
	url := downloadUrl

	taskId := uuid.NewV4().String()
	contentFilename, deleteContentFile := createTempFileWithContent(t)
	defer deleteContentFile()
	fileServer := filetransfer.NewFileServer(&StubAdapter{downloadTaskId: taskId, filename: contentFilename})
	requestUrl := fmt.Sprintf("%s?taskId=%s", url, taskId)
	request := newGetRequest(requestUrl)
	response := httptest.NewRecorder()
	fileServer.ServeHTTP(response, request)
	assertIntEquals(t, response.Code, http.StatusOK)
	contentDisposition := response.Header().Get("Content-Disposition")
	filenamePrefix := "attachment; filename="
	if !str.StartsWith(contentDisposition, filenamePrefix) {
		t.Errorf("got uncorrect Content-Disposition")
	}
	gotFilename := contentDisposition[len(filenamePrefix):]
	assertDirectlyEqual(t, gotFilename, contentFilename)
	downloadFilename := "download-" + gotFilename
	downloadFile(t, downloadFilename, response.Body)
	assertFileContentEquals(t, contentFilename, gotFilename)
	_ = os.Remove(downloadFilename)
}

func TestTaskNotFound(t *testing.T) {
	test := func(url string, fn func(requestUrl string) *http.Request) {
		fileServer := filetransfer.NewFileServer(&StubAdapter{})
		taskId := uuid.NewV4().String()
		requestUrl := fmt.Sprintf("%s?taskId=%s", url, taskId)
		request := fn(requestUrl)
		response := httptest.NewRecorder()
		fileServer.ServeHTTP(response, request)
		assertIntEquals(t, response.Code, http.StatusBadRequest)

		wantErrorBody := getTaskNotFoundBody()
		var gotErrorBody filetransfer.ErrorBody
		_ = json.NewDecoder(response.Body).Decode(&gotErrorBody)
		assertStructEquals(t, gotErrorBody, wantErrorBody)
	}

	test(uploadUrl, func(requestUrl string) *http.Request {
		return newPostRequestReader(requestUrl, nil)
	})
	test(downloadUrl, func(requestUrl string) *http.Request {
		return newGetRequest(requestUrl)
	})
}

func testHttpStatus(t *testing.T, requestBody interface{}, got, wantStatus int) {
	t.Helper()
	if got != wantStatus {
		t.Errorf("test case: %+v \n got %d but want %d from response", requestBody, got, wantStatus)
	}
}

func newPostRequestReader(url string, requestBody io.Reader) *http.Request {
	request, _ := http.NewRequest(http.MethodPost, url, requestBody)
	return request
}

func newPostReqBody(t *testing.T, url string, requestBody interface{}) *http.Request {
	requestBodyReader := new(bytes.Buffer)
	err := json.NewEncoder(requestBodyReader).Encode(requestBody)
	if err != nil {
		t.Fatalf("wrong body: %+v", requestBody)
	}
	return newPostRequestReader(url, requestBodyReader)
}

func newGetRequest(url string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	return request
}

func assertFileContentEquals(t *testing.T, filename1, filename2 string) {
	t.Helper()
	if fileMd5Hash(t, filename1) != fileMd5Hash(t, filename2) {
		t.Errorf("file hash not equals: '%s' , '%s'", filename1, filename2)
	}
}

func fileMd5Hash(t *testing.T, filename string) string {
	t.Helper()
	file, _ := os.Open(filename)
	defer file.Close()
	buffer := make([]byte, 256)
	hash := md5.New()
	for {
		readLen, _ := file.Read(buffer)
		if readLen == 0 {
			break
		}
		hash.Write(buffer[:readLen])
	}
	sum := hash.Sum(nil)
	return fmt.Sprintf("%x", sum)
}

func downloadFile(t *testing.T, filename string, reader io.Reader) {
	t.Helper()

	tempFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		t.Fatalf("could not create file %v", err)
	}
	defer tempFile.Close()
	buffer := make([]byte, 256)
	for {
		readLen, _ := reader.Read(buffer)
		if readLen == 0 {
			break
		}
		tempFile.Write(buffer[:readLen])
	}
}

func createTempFileWithContent(t *testing.T) (string, func()) {
	t.Helper()

	tempFile, err := os.OpenFile(createRandomFilename("tempFileWithContent", ".txt"), os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}
	content := uuid.NewV4().String() + "中文测试"
	_, err = tempFile.Write([]byte(content))
	if err != nil {
		t.Fatalf("could not input content to file %v", err)
	}
	_ = tempFile.Close()
	removeFile := func() {
		_ = os.Remove(tempFile.Name())
	}
	return tempFile.Name(), removeFile
}

func createRandomFilename(pattern, suffix string) string {
	return fmt.Sprintf("%s%s%s", pattern, uuid.NewV4().String(), suffix)
}

func extractOkBody(reader io.Reader) filetransfer.OkBody {
	var okBody filetransfer.OkBody
	_ = json.NewDecoder(reader).Decode(&okBody)
	return okBody
}

func testCase(t *testing.T, test initTestCase, url string, fileServer http.Handler) *httptest.ResponseRecorder {
	request := newPostReqBody(t, url, test.requestBody)
	response := httptest.NewRecorder()
	fileServer.ServeHTTP(response, request)
	testHttpStatus(t, test.requestBody, response.Code, test.wantResponseStatus)
	return response
}

func getInvalidErrBody() filetransfer.ErrorBody {
	errResponseBody := filetransfer.ErrorBody{Error: filetransfer.ErrorContent{
		Message: filetransfer.ErrorContentInvalidParam, Code: filetransfer.ErrorCodeInvalidParam}}
	return errResponseBody
}

func getTaskNotFoundBody() filetransfer.ErrorBody {
	wantErrorBody := filetransfer.ErrorBody{
		Error: filetransfer.ErrorContent{
			Message: filetransfer.ErrorContentTaskNotFound,
			Code:    filetransfer.ErrorCodeResourceNotFound,
		},
	}
	return wantErrorBody
}
