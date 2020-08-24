package v3

import (
	"fmt"
	"net/url"
	"testing"
)

func Test_Encode(t *testing.T) {
	params := url.Values{}
	params.Add("name", "@Rajeev")
	params.Add("phone", "+919999999999")
	res := "hello" + params.Encode()
	fmt.Println(res)
}

func Test_Encode1(t *testing.T) {
	params := url.Values{}
	params.Add("name", "@Rajeev")
	params.Add("phone", "+919999999999")
	res := "hello" + params.Encode()
	fmt.Println(res)
}
