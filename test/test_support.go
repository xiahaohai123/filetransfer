package testutil

import (
	"github.com/modern-go/reflect2"
	"reflect"
	"testing"
)

func AssertIntEquals(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("want %d but got %d", want, got)
	}
}

func AssertStructEquals(t *testing.T, got, want interface{}) {
	t.Helper()
	if got == nil || want == nil {
		t.Errorf("unexpted nil pointer")
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %+v but got %+v", want, got)
	}
}

func AssertContent(t *testing.T, content string) {
	t.Helper()
	if content == "" {
		t.Errorf("want a content but got empty")
	}
}

func AssertStringEqual(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("want %s but got %s", want, got)
	}
}

func AssertTrue(t *testing.T, got bool) {
	t.Helper()
	if !got {
		t.Errorf("want true but got false")
	}
}

func AssertFalse(t *testing.T, got bool) {
	t.Helper()
	if got {
		t.Errorf("want false but got true")
	}
}

func AssertNotNil(t *testing.T, got interface{}) {
	t.Helper()
	if got == nil || reflect2.IsNil(got) {
		t.Errorf("want not nil but got")
	}
}

func AssertNil(t *testing.T, got interface{}) {
	t.Helper()
	if got != nil && !reflect2.IsNil(got) {
		t.Errorf("want nil but got other: %+v", got)
	}
}

func AssertErrEquals(t *testing.T, got, want error) {
	t.Helper()
	if got != want {
		t.Errorf("want %v but got %v", want, got)
	}
}
