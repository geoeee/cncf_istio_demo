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
	adminToken = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTA4MDE0MTUsImlzcyI6IjE5Mi4xNjguOTkuMTAwIiwidXNlcl9yb2xlcyI6WyJhZG1pbiJdfQ.dWqlPbzAS2ScdsN38E3GCFAOU_e47Jn2a1SUaIRWoCmhzqsKZn6af0Fju1Va7rZrRV8UVX2lFR7QhEKCY5SzrcEdvk1Kiw3MVk1jkYYtUZup84IPq4-QqhHE1EcNRCvO0rQfMngtM-jfS7K98WlPs22PJTRvQopt8LloGBSp2ype2Vq4bxlJFNgG2E73y1GPEno7kWNdtgSA_bh4ObsTiJuzWmW10kZnu1ElzEP9t4Y8FD0NkYrASaOJ4lrymtPbjIuGTklMPPgVJ7jxjKa5NoVIAita6k0S-7drcpCEom8eCSW18Gqosg2XFXTBY3B3UMwy2tsc54A0f9FD9gGTtc4ZFM3avlARvrNMrk8ImSb-HYO9VKrimSZSqk4L1NmarYzILD_Qo36ziRNTzh1VoB2NmLavCnlAguhIXmWPsQp2BZI3tZ1H-GibnOzu3PnTRC73Z8WnwBQ2QPJidpKE3TNR6FSw999Wbn4FEIiWR5wKQTH4MyhyAMZ2-kp0WWs045I7_a2iDGjxAAKb8qRv0PJ3XtHKWzoZRj5z8ic9qJKus2ST3Q_0JkTiWJj2nuWUXw9MK3T93K_BSmrbjv-rwfzehwc0-kHTt4VRCHJPZaCR4ktILY-5bnVhlMEQLBEeBdxz-VQ6oZuho9bVESpiuYakCEOcWJts0GH9a21qKVA"
	userToken = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTA4MDE0MTUsImlzcyI6IjE5Mi4xNjguOTkuMTAwIiwidXNlcl9yb2xlcyI6WyJ1c2VyIl19.VCr_ing7kOn6n4ljq60aEid7ALu-jJFCoRgY89WXnHb-DIEZGjkz6--SRND5dl5J4jNXntWsYch54-pJ4SVaiYyD50yLx4mCt9xct7kLz91tlwNeFoZMxO64vZ5VkgOunWCIaP5VEZOB9GpsjedDgUuLu60lLdYy3-KTSl6X0jDBdGJdGbT2SiznZW_MpwzXu3_9FrUqMZye-IENxa4Gwn_8sgWmhRwnnc5NDC5VXbf7_QDWWp4wy5o4N5UKnbd9uv377oEGsfMEH5NuJBXg2_Lm2FeTEL_pFr0o7B__fB7ewyFEgfa0Kp-h-Wx5cUAgzC8n9gloQMCVz6jcM_cB3CSYcZwGUBQ62SI3gjUzomldOIzrRq6Om520Uk6OHF3pKWSc330RVEv9bf6_v4vKIGbDflVO5UZThMRg2RC5wWnKFDsYoh5vrP6122sqtJIjAf3XXaFWvOf2F1V7117BqaWQ2ELL6wIc7nlUZbmV2foAG-uhKCelhlCt2LpOqRrrlZde_UdIIqnCtHDPn0gB7l5ta-PcKMM1Gsjgkhu4-7bHVaG01Mk8z1lvNGUtQmyBEmnUkWFciB5xBpDpxF_PuMGvRs5AK0zu10n_7HNtNXMZLhSV7XU5S0_yB7x1JR3kKY5N3TmEKy2v1Jiw3YHaGILQPBw-fgGZJd7a_hXbOUI"
)

func getCompanies(auth string) {
	url := "http://192.168.99.100:31380/company/api/v1/companies"
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Cache-Control", "no-cache")

	if auth == "admin" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", adminToken))
	}
	if auth == "user" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", userToken))
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