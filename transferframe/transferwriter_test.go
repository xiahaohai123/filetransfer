package transferframe_test

import (
	"bytes"
	"errors"
	"strings"
	testutil "summersea.top/filetransfer/test"
	"summersea.top/filetransfer/transferframe"
	"testing"
)

var stubBeforeErr = errors.New("stub before err")
var stubWriteErr = errors.New("stub write err")

const testInput = "test input"

type stubTransferWriter struct {
	gotErr          error
	stringBuf       bytes.Buffer
	beforeCall      int
	shouldBeforeErr bool
	writeCall       int
	shouldWriteErr  bool
	afterCall       int
}

func (s *stubTransferWriter) BeforeTransfer() error {
	s.beforeCall++
	if s.shouldBeforeErr {
		return stubBeforeErr
	}
	return nil
}

func (s *stubTransferWriter) Write(bytes []byte) error {
	s.writeCall++
	if s.shouldWriteErr {
		return stubWriteErr
	}
	_, err := s.stringBuf.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

func (s *stubTransferWriter) AfterTransfer() {
	s.afterCall++
}

func (s *stubTransferWriter) ErrorTransfer(err error) {
	s.gotErr = err
}

type StubReader struct{}

var stubReadErr = errors.New("stub read error")

func (StubReader) Read(p []byte) (n int, err error) {
	return 0, stubReadErr
}

func TestNewTransferManager(t *testing.T) {
	t.Run("input nil", func(t *testing.T) {
		manager, err := transferframe.NewTransferManager(nil)
		testutil.AssertNotNil(t, err)
		testutil.AssertNil(t, manager)
	})

	t.Run("input reader", func(t *testing.T) {
		manager, err := transferframe.NewTransferManager(StubReader{})
		testutil.AssertNotNil(t, manager)
		testutil.AssertNil(t, err)
	})
}

func TestTransferManager_AddWriter(t *testing.T) {
	t.Run("input nil", func(t *testing.T) {
		manager := createCommonManager()
		err := manager.AddWriter(nil)
		testutil.AssertNotNil(t, err)
	})

	t.Run("common test", func(t *testing.T) {
		manager := createCommonManager()
		writer := &stubTransferWriter{}
		err := manager.AddWriter(writer)
		testutil.AssertNil(t, err)
	})
}

func TestTransferManager_StartTransfer(t *testing.T) {
	t.Run("read err", func(t *testing.T) {
		manager, _ := transferframe.NewTransferManager(StubReader{})
		writer := &stubTransferWriter{}
		_ = manager.AddWriter(writer)
		err := manager.StartTransfer()
		testutil.AssertErrEquals(t, err, stubReadErr)
		testutil.AssertIntEquals(t, writer.beforeCall, 1)
		testutil.AssertIntEquals(t, writer.writeCall, 0)
		testutil.AssertErrEquals(t, writer.gotErr, transferframe.ReaderErr)
	})

	t.Run("before err", func(t *testing.T) {
		writer := &stubTransferWriter{shouldBeforeErr: true}
		manager := createManagerWithWriter(writer)
		err := manager.StartTransfer()
		testutil.AssertNil(t, err)
		testutil.AssertErrEquals(t, writer.gotErr, stubBeforeErr)
	})

	t.Run("write err", func(t *testing.T) {
		writer := &stubTransferWriter{shouldWriteErr: true}
		manager := createManagerWithWriter(writer)
		err := manager.StartTransfer()
		testutil.AssertNil(t, err)
		testutil.AssertIntEquals(t, writer.writeCall, 1)
		testutil.AssertErrEquals(t, writer.gotErr, stubWriteErr)
	})

	t.Run("common write", func(t *testing.T) {
		writer := &stubTransferWriter{}
		manager := createManagerWithWriter(writer)
		err := manager.StartTransfer()
		testutil.AssertNil(t, err)
		testutil.AssertIntEquals(t, writer.afterCall, 1)
		testutil.AssertNil(t, writer.gotErr)
		s := writer.stringBuf.String()
		testutil.AssertStringEqual(t, s, testInput)
	})
}

func createManagerWithWriter(writer transferframe.TransferWriter) *transferframe.TransferManager {
	manager := createCommonManager()
	_ = manager.AddWriter(writer)
	return manager
}

func createCommonManager() *transferframe.TransferManager {
	manager, _ := transferframe.NewTransferManager(strings.NewReader(testInput))
	return manager
}

func TestNewBasicWriter(t *testing.T) {
	t.Run("input nil", func(t *testing.T) {
		writer, err := transferframe.NewBasicWriter(nil)
		testutil.AssertErrEquals(t, err, transferframe.NilParamErr)
		testutil.AssertNil(t, writer)
	})

	t.Run("input reader", func(t *testing.T) {
		writer, err := transferframe.NewBasicWriter(&bytes.Buffer{})
		testutil.AssertNotNil(t, writer)
		testutil.AssertNil(t, err)
	})
}

func TestBasicWriter_Write(t *testing.T) {
	buf := &bytes.Buffer{}
	writer, _ := transferframe.NewBasicWriter(buf)
	err := writer.Write([]byte(testInput))
	testutil.AssertNil(t, err)
	testutil.AssertStringEqual(t, buf.String(), testInput)
}
