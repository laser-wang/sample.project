package rest

import (
	"errors"
	"testing"
)

func Test_checkErr_1(t *testing.T) {
	err := errors.New("test sample")
	checkErr(err, CHECK_FLAG_LOG)
	t.Log("checkErr_1测试通过")

}

func Test_checkErr_2(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Log("checkErr_2测试通过")
		} else {
			t.Error("checkErr_2测试没有通过!!!")
		}

	}()
	err := errors.New("test sample")
	checkErr(err, CHECK_FLAG_EXIT)

}
