package main

import (
	"github.com/zhangzhoujian/istio-demo/echo/internal/schema"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type comResponse struct {
	Meta     schema.Meta
	Companies []*schema.Company `json:"elements"`
}
var (
	prodV1Count int
	prodV2Count int
)

func getCompanies(xHeader string) {
	url := "http://192.168.99.100:31380/company/api/v1/companies"
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("x-request-platform", xHeader)
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
	for _, c := range er.Companies {
		for _, p := range c.Products {
			fmt.Printf("{ x-request-platform: [%s] } ==> product version [%s], pod name [%s]\n",xHeader, p.Meta.Version, p.Meta.PodName)
		}
	}
}

func main() {
	argNum := len(os.Args)
	// fmt.Printf("the num of input is %d\n",argNum)

	// fmt.Printf("%v", os.Args)
	if argNum < 2 {
		for i := 0; i < 10; i++ {
			getCompanies("")
		}
		return
	} 

	for i := 0; i < 10; i++ {
		getCompanies(os.Args[1])
	}
}