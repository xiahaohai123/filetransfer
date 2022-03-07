package filetransfer

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"time"
)

type UploadData UploadInitReqBody

type DownloadData DownloadInitReqBody

type FileTranDataAdapter struct {
	dataStore DataStore
}

func NewFileTranDataAdapter(store DataStore) *FileTranDataAdapter {
	return &FileTranDataAdapter{store}
}

func (f *FileTranDataAdapter) SaveUploadData(taskId string, uploadData UploadData) {
	f.dataStore.SaveUploadData(taskId, uploadData)
}

func (f *FileTranDataAdapter) IsUploadTaskExist(taskId string) bool {
	return f.dataStore.IsUploadTaskExist(taskId)
}

func (f *FileTranDataAdapter) GetUploadChannel(taskId string) (WriteCloseRollback, error) {
	uploadData := f.dataStore.GetUploadDataRemove(taskId)
	return f.createSftpChannel(*uploadData)
}

func (f *FileTranDataAdapter) IsDownloadTaskExist(taskId string) bool {
	panic("implement me")
}

func (f *FileTranDataAdapter) GetDownloadChannelFilename(taskId string) (io.ReadCloser, string, error) {
	panic("implement me")
}

func (f *FileTranDataAdapter) SaveDownloadData(taskId string, downloadData DownloadData) {
	panic("implement me")
}

func (f *FileTranDataAdapter) createSftpChannel(data UploadData) (WriteCloseRollback, error) {
	resource := data.Resource
	sshConfig := &ssh.ClientConfig{
		User: resource.Account.Name,
		Auth: []ssh.AuthMethod{
			ssh.Password(resource.Account.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		ClientVersion:   "",
		Timeout:         10 * time.Second,
	}

	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", resource.Address, resource.Port), sshConfig)
	if err != nil {
		return nil, fmt.Errorf("problem dial target resource: %v", err)
	}
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		closeWithErrLog(sshClient)
		return nil, fmt.Errorf("problem create sftp client: %v", err)
	}
	filePath := sftp.Join(data.Path, data.Filename)
	transferChannel, err := sftpClient.Create(filePath)
	if err != nil {
		closeWithErrLog(sftpClient)
		closeWithErrLog(sshClient)
		return nil, fmt.Errorf("problem create upload channel: %v", err)
	}
	channel := &SftpUploadChannel{sshClient, sftpClient, transferChannel, filePath}

	return channel, nil
}

type DataStore interface {
	SaveUploadData(taskId string, data UploadData)
	GetUploadDataRemove(taskId string) *UploadData
	IsUploadTaskExist(taskId string) bool
	SaveDownloadData(taskId string, data DownloadData)
	GetDownloadDataRemove(taskId string) *DownloadData
	IsDownloadTaskExist(taskId string) bool
}

type WriteCloseRollback interface {
	io.WriteCloser
	RollBack() error
}

type SftpUploadChannel struct {
	sshClient  io.Closer
	sftpClient *sftp.Client
	io.WriteCloser
	filePath string
}

func (s *SftpUploadChannel) Close() error {
	closeWithErrLog(s.WriteCloser)
	closeWithErrLog(s.sftpClient)
	closeWithErrLog(s.sshClient)
	return nil
}

func (s *SftpUploadChannel) RollBack() error {
	return s.sftpClient.Remove(s.filePath)
}

func closeWithErrLog(closer io.Closer) {
	err := closer.Close()
	if err != nil {
		log.Printf("problem close io: %v", err)
	}
}
