package filetransfer

const ErrorCodeInvalidParam = "InvalidParam"
const ErrorContentInvalidParam = "Invalid Parameter"
const ErrorCodeResourceNotFound = "ResourceNotFound"
const ErrorContentTaskNotFound = "The task id is not found"

type Resource struct {
	Address string  `json:"address"`
	Port    int     `json:"port"`
	Account Account `json:"account"`
}

type Account struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UploadInitReqBody struct {
	Resource Resource `json:"resource"`
	Path     string   `json:"path"`
	Filename string   `json:"filename"`
}

type OkBody struct {
	Data Data `json:"data"`
}

type Data map[string]interface{}

type ErrorBody struct {
	Error ErrorContent `json:"error"`
}

type ErrorContent struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func getTaskNotFoundErr() ErrorBody {
	return NewErrorBody(ErrorCodeResourceNotFound, ErrorContentTaskNotFound)
}

func getInvalidParamErr() ErrorBody {
	return NewErrorBody(ErrorCodeInvalidParam, ErrorContentInvalidParam)
}

func NewErrorBody(code, message string) ErrorBody {
	return ErrorBody{
		Error: ErrorContent{
			Message: message,
			Code:    code,
		},
	}
}
