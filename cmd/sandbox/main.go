package main

import (
	"github.com/imroc/req/v3"
)

func main() {
	testClient := req.C()
	testClient.SetProxyURL("sock5://hello:1415456@165.154.3.235:4532")

	res, err := testClient.R().Get("https://www.google.com")
	if err != nil {
		panic(err)
	}
	print(res.Status)
}
