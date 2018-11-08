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
	Meta     schema.Meta
	Companies []*schema.Company `json:"elements"`
}
var (
	prodV1Count int
	prodV2Count int
	okToken = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTA4MDE0MTUsImlzcyI6IjE5Mi4xNjguOTkuMTAwIn0.naXDV8R6o4JpyLbN058xc_5M6EQa6FxqlAC_imrEj97lMxN0ZrPdT0e3PaRBRaPMXsCVSOQcIpMEVhna3tRoViIKtPT3IhgkxDOEbx4wbGOYCMM0hpdyHoV3g15mV70pDXP2R6oCGPy5A9x-cpOn0iUAAPNCe1WpWv4FMoev9rWAt8TDLuxCpzMd0-wBBNZ-WWcHEa1Qwo598We_T15xHXSsIBsVe_vkmeFGp7yg9t2kJqokgDXOB7rr1k2w-ATjB2x9Lt7w93wjGcTgX57RKL8sV2jjb0XbJ9vB5eY3yoJw5u21nFBBmLoA2JkllZDFCZ7COWHaiZfqtxa-3d56BmP_LCZGc9_aY44MRObjVo38k9oXuENU8nG0NuaYMDM320lhehG6aq8ozRmh3e1zvkDpsnmm1LvQzahVxIkb_2-dvCBrrXFdsJhzc2xYHhU270oiwnMYhTgkJZ4VZGy_odgzHxuC1QErkGhB_MMoxMAKBtsNEPIVJUNtJZG0MBlEUGVCVdFM9I5f47yhlh6tSggG0WOOwImJiVj_SQhrGyd7BOySbfQYEFE_3TXKXRrytAZj_GGXTrAQS8AOqfy4eGv-kMuCc10K-8_ewK31DL8pe_RYpGOPtWTmQ-sHx0C8c_hBqEBFoM-GjRnmzo-9g1FZ-jo78Mqt47fWU2yszIM"
	wrongToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
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