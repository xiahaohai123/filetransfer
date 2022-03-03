package filetransfer

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
}

type ErrorBody struct {
	Error ErrorContent `json:"error"`
}

type ErrorContent struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}
