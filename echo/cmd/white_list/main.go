package main

import (
	"github.com/zhangzhoujian/istio-demo/echo/internal/schema"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type comResponse struct {
	Meta     schema.Meta
	Companies []*schema.Company `json:"elements"`
}
var (
	prodV1Count int
	prodV2Count int
)

func getCompanies(auth string) {
	url := "http://192.168.99.100:31380/company/api/v1/companies"
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Cache-Control", "no-cache")
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	er := &comResponse{}
	err = json.Unmarshal(body, er)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(er.Companies) > 1 {
		c := er.Companies[0]
	
		var productOutput string
		if c.Products == nil {
			productOutput = fmt.Sprintf("[%s] products is {null}", time.Now())
		} else {
			productOutput = fmt.Sprintf("[%s] get products ok: num:{%d}", time.Now(), len(c.Products))
		}
		fmt.Println(productOutput)
	} else {
		fmt.Println("get companies err companies is 0 len")
	}

}

func main() {
	argNum := len(os.Args)
	if argNum < 2 {
		for i := 0; i < 2000; i++ {
			getCompanies("")
			time.Sleep(1 * time.Second)
		}
		return
	}

	for i := 0; i < 2000; i++ {
		getCompanies(os.Args[1])
		time.Sleep(1 * time.Second)
	}
}