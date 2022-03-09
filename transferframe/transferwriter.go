package transferframe

import (
	"errors"
	"io"
	"log"
)

var ReaderErr = errors.New("reader error")

const bufferSize = 1024

type TransferWriter interface {
	// BeforeTransfer 传输之前调用，若出现异常则踢出传输链，并立即调用ErrorTransfer
	BeforeTransfer() error
	// Write 每次从io.Reader中读取部分字节后，会调用传输链上的每一个Writer的Write([]byte)方法，传输过程。
	Write([]byte) error
	// AfterTransfer 若传输全程不出现问题，则会调用此方法以结束传输，可以不实现
	AfterTransfer()
	// ErrorTransfer 传输过程中出现任何error都会调用此方法，并将err传入
	ErrorTransfer(err error)
}

var NilParamErr = errors.New("got nil param")

type BasicWriter struct {
	writer io.Writer
}

func NewBasicWriter(writer io.Writer) (*BasicWriter, error) {
	if writer == nil {
		return nil, NilParamErr
	} else {
		return &BasicWriter{writer: writer}, nil
	}
}

func (b *BasicWriter) BeforeTransfer() error {
	// Do nothing
	return nil
}

func (b *BasicWriter) Write(bytes []byte) error {
	_, err := b.writer.Write(bytes)
	return err
}

func (b *BasicWriter) AfterTransfer() {
	// Do nothing
}

func (b *BasicWriter) ErrorTransfer(err error) {
	log.Printf("problem transfer: %v", err)
}

type TransferManager struct {
	reader  io.Reader
	writers []TransferWriter
}

// NewTransferManager 创建传输管理器
// io.Reader 已经open的输入流，管理器不负责关闭输入流和输出流
// error 传入参数有问题时会返回异常
func NewTransferManager(reader io.Reader) (*TransferManager, error) {
	if reader == nil {
		return nil, errors.New("got nil reader")
	} else {
		manager := &TransferManager{reader: reader, writers: []TransferWriter{}}
		return manager, nil
	}
}

// AddWriter 添加传输输入端
// error 传入参数为nil时会抛出异常
func (t *TransferManager) AddWriter(writer TransferWriter) error {
	if writer == nil {
		return errors.New("got nil writer")
	} else {
		t.writers = append(t.writers, writer)
		return nil
	}
}

// StartTransfer 开始传输
// 如果没有输出端，也会读完输入端
// error 输入端出现异常时返回该异常，在此之前会调用所有输出端的异常结束方法，并传入ReadErr
func (t *TransferManager) StartTransfer() error {
	t.callBeforeFunc()
	err := t.doTransfer()
	if err != nil {
		t.callReadErr()
		return err
	}
	t.callAfterFunc()
	return nil
}

// doTransfer 执行传输过程
// error 当出现读入端错误时会返回该错误，该方法不会在读入错误时调用ErrorTransfer方法
func (t *TransferManager) doTransfer() error {
	buf := make([]byte, bufferSize)
	reader := t.reader
	writers := t.writers
	for {
		readLen, err := reader.Read(buf)
		if readLen > 0 {
			i := 0
			for _, writer := range writers {
				writeErr := writer.Write(buf[:readLen])
				if writeErr != nil {
					writer.ErrorTransfer(writeErr)
				} else {
					writers[i] = writer
					i++
				}
			}
			if i != len(writers) {
				writers = writers[:i]
			}
		}
		if err == io.EOF {
			t.writers = writers
			return nil
		} else if err != nil {
			return err
		}
	}
}

func (t *TransferManager) callReadErr() {
	for _, writer := range t.writers {
		writer.ErrorTransfer(ReaderErr)
	}
}

func (t *TransferManager) callBeforeFunc() {
	i := 0
	writers := t.writers
	for _, writer := range writers {
		err := writer.BeforeTransfer()
		if err != nil {
			writer.ErrorTransfer(err)
		} else {
			writers[i] = writer
			i++
		}
	}
	t.writers = writers[:i]
}

func (t *TransferManager) callAfterFunc() {
	for _, writer := range t.writers {
		writer.AfterTransfer()
	}
}
