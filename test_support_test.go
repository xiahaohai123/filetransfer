package filetransfer_test

import (
	"reflect"
	"testing"
)

func assertIntEquals(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("want %d but got %d", want, got)
	}
}

func assertStructEquals(t *testing.T, got, want interface{}) {
	t.Helper()
	if got == nil || want == nil {
		t.Errorf("unexpted nil pointer")
		return
	}
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

func assertDirectlyEqual(t *testing.T, got, want interface{}) {
	t.Helper()
	if got != want {
		t.Errorf("want %d but got %d", want, got)
	}
}
