package filetransfer

import (
	"errors"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"path/filepath"
	"time"
)

var DownloadDir = errors.New("can not download directory")

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
	return f.createUploadSftpChannel(*uploadData)
}

func (f *FileTranDataAdapter) IsDownloadTaskExist(taskId string) bool {
	return f.dataStore.IsDownloadTaskExist(taskId)
}

func (f *FileTranDataAdapter) GetDownloadChannelFilename(taskId string) (io.ReadCloser, string, error) {
	downloadData := f.dataStore.GetDownloadDataRemove(taskId)
	channel, err := f.createSftpDownloadChannel(downloadData.Resource, downloadData.Path)
	if err != nil {
		if err == DownloadDir {
			return nil, "", err
		} else {
			return nil, "", fmt.Errorf("problem create channel: %v", err)
		}
	}
	filename := filepath.Base(downloadData.Path)
	return channel, filename, nil
}

func (f *FileTranDataAdapter) SaveDownloadData(taskId string, downloadData DownloadData) {
	f.dataStore.SaveDownloadData(taskId, downloadData)
}

func (f *FileTranDataAdapter) createUploadSftpChannel(data UploadData) (WriteCloseRollback, error) {
	sftpClient, err := f.createSftpClient(data.Resource)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	filePath := sftp.Join(data.Path, data.Filename)
	transferChannel, err := sftpClient.Create(filePath)
	if err != nil {
		_ = sftpClient.Close()
		return nil, fmt.Errorf("problem create upload channel: %v", err)
	}
	channel := &SftpUploadChannel{sftpClient.sshClient, sftpClient.Client, transferChannel, filePath}

	return channel, nil
}

func (f *FileTranDataAdapter) createSftpDownloadChannel(resource Resource, path string) (io.ReadCloser, error) {
	sftpClient, err := f.createSftpClient(resource)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	fileInfo, err := sftpClient.Stat(path)
	if err != nil {
		_ = sftpClient.Close()
		return nil, fmt.Errorf("problem while search file %v", err)
	}
	if fileInfo.IsDir() {
		_ = sftpClient.Close()
		return nil, DownloadDir
	}
	file, err := sftpClient.Open(path)
	if err != nil {
		_ = sftpClient.Close()
		return nil, fmt.Errorf("problem open file %v", err)
	}
	return &sftpDownloadChannel{sftpClient, file}, nil
}

func (f *FileTranDataAdapter) createShhConfig(account Account) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User: account.Name,
		Auth: []ssh.AuthMethod{
			ssh.Password(account.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		ClientVersion:   "",
		Timeout:         10 * time.Second,
	}
}

func (f *FileTranDataAdapter) createSftpClient(resource Resource) (*ClientPackage, error) {
	sshConfig := f.createShhConfig(resource.Account)
	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", resource.Address, resource.Port), sshConfig)
	if err != nil {
		return nil, fmt.Errorf("problem dial target resource: %v", err)
	}
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		closeWithErrLog(sshClient)
		return nil, fmt.Errorf("problem create sftp client: %v", err)
	}
	c := &ClientPackage{sshClient: sshClient}
	c.Client = sftpClient
	return c, nil
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

type ClientPackage struct {
	*sftp.Client
	sshClient io.Closer
}

func (c *ClientPackage) Close() error {
	closeWithErrLog(c.Client)
	closeWithErrLog(c.sshClient)
	return nil
}

type sftpDownloadChannel struct {
	client *ClientPackage
	io.ReadCloser
}

func (sf *sftpDownloadChannel) Close() error {
	closeWithErrLog(sf.ReadCloser)
	_ = sf.client.Close()
	return nil
}

func closeWithErrLog(closer io.Closer) {
	err := closer.Close()
	if err != nil {
		log.Printf("problem close io: %v", err)
	}
}
