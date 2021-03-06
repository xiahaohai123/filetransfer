package filetransfer_test

import (
	"summersea.top/filetransfer"
	testutil "summersea.top/filetransfer/test"
	"testing"
)

// 该测试使用外部环境进行测试
// docker run --redis-test -dp 6379:6379 redis
func TestNewRedisStore(t *testing.T) {
	t.Run("common", func(t *testing.T) {
		store, err := filetransfer.NewRedisStore("localhost:6379", "", 0)
		testutil.AssertNil(t, err)
		testutil.AssertNotNil(t, store)
	})

	t.Run("wrong message", func(t *testing.T) {
		store, err := filetransfer.NewRedisStore("localhost:6381", "", 0)
		testutil.AssertNotNil(t, err)
		testutil.AssertNil(t, store)
	})
}
