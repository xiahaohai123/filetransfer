package filetransfer

import "net/http"

type FileServer struct {
	http.Handler
}

func NewFileServer() *FileServer {
	fileServer := &FileServer{}
	router := http.NewServeMux()
	router.Handle("/file/upload/initialization", http.HandlerFunc(fileServer.uploadInitHandler))

	fileServer.Handler = router
	return fileServer
}

func (fs *FileServer) uploadInitHandler(w http.ResponseWriter, r *http.Request) {

}
