package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main(){
	testGet()
}
func testGet() {

	url := "https://baidu.com"

	req, err := http.NewRequest("GET", url, nil)

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {

		panic(err)

	}

	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)

	fmt.Println("response Headers:", resp.Header)

	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("response Body:", string(body))

}