package main

import (
	"github.com/zhangzhoujian/istio-demo/echo/internal/schema"
	// "encoding/json"
	"fmt"
	// "io/ioutil"
	"net/http"
	"os"
	"time"
)

type comResponse struct {
	Meta      schema.Meta
	Companies []*schema.Company `json:"elements"`
}

var (
	prodV1Count int
	prodV2Count int
	okToken     = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTA4MDE0MTUsImlzcyI6IjE5Mi4xNjguOTkuMTAwIiwiYXVkIjpbImNvbXBhbnkiXX0.tBUUGFY7JXHpjiOUkaW2QJ4r5i0IgxorasR7P6aQ6m7hWoMnOesBGMAHs48bdqFi2WQIFQ8x1E7WV-wPu31t05cE1YvhQvXd1L5cQbrtu6hsgfr9DU9FMDArHmJ2h2_inHpRI7xSKETU_uOD3hn8kSuzHetVlRJbIsM9M-7MTg3gAqSVbx6LDYIJVKA0Crox-bG1DMoeGdUBQ8RUprQ7K-1uEugtiYY8VjmDLQywx_8p07YZnxbkeii2IdEY9242hsBdbc_K3dHjB3yswXTwdBkFpAh0TdS9GfWr83OgSSRxNygSZJfWlJ3RLFL8lYbOcj_BAeTf6qKPVCiVukrW7TnPTKG9gS9ESHq2xOZg1VRh4FDUZsLkeiO8lWcL864J-xXKR93CodocGV5FOgsp33GKibXGQJsXxdTEqUHBOHT0hQJ4LOZXVxBkIcFjqG-vvEPCPdpCS-Axr1gb9llBtmOdvzv8CS1hQgBI2c6Ljqp6RwY9_vVsRRs7fqYEqRTR_HOuYTZHraj7_T-pm0Bk0SNShExXF0IyLdMqk6b4qeI_xRGbn7iL4B_J0of5evIOPtAoFqQOYAWes_tEWVZuo4Ul5LDWu2fwrTsEEhd7geGEWustVZ7QRMIAo1zjhqkqE1KqClTkGFJgDJfC5SXxXBt5jf9bsz0vEHVvWWeyg4Q"
	wrongToken  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
)

func getCompanies(auth string) {
	url := "http://192.168.99.100:31380/company/api/v1/companies"
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Cache-Control", "no-cache")

	if auth == "RS256" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", okToken))
	}
	if auth == "HS256" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", wrongToken))
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("[%s] resp status code:[%s]\n", time.Now(), res.Status)
	defer res.Body.Close()
}

func main() {
	argNum := len(os.Args)
	if argNum < 2 {
		fmt.Println("Usage: required args [RS256 | HS256]")
		return
	}

	for i := 0; i < 2000; i++ {
		getCompanies(os.Args[1])
		time.Sleep(1 * time.Second)
	}
}
