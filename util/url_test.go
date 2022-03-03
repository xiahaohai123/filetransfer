package util_test

import (
	"summersea.top/filetransfer/util"
	"testing"
)

func TestExtractUrlParam(t *testing.T) {
	util.ExtractUrlParam("")
	util.ExtractUrlParam("aaa")
	param := util.ExtractUrlParam("/aa?aa")
	if param["aa"] != "" {
		t.Errorf("want empty string bug got '%s'", param["aa"])
	}
	util.ExtractUrlParam("/aa?=aa")
	if len(param) != 0 {
		t.Errorf("want empty map bug got %+v", param)
	}
}
