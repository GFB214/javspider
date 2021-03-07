package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"time"
)

var (
	imgpath   = "covers"
	rooturl   = "https://www.javbus.com/"
	baseurl   = "https://www.javbus.com/genre/sub/"
	magneturl = "ajax/uncledatoolsbyajax.php"
	current   = 1
	limit     = false
	gidregexp *regexp.Regexp
	ucregexp  *regexp.Regexp
	randg     *rand.Rand
)

func init() {
	var err error
	_, err = os.Stat(imgpath)
	if err != nil && os.IsNotExist(err) {
		os.Mkdir(imgpath, 0777)
	}

	gidregexp, err = regexp.Compile(`var gid.*;`)
	if err != nil {
		fmt.Println(err)
	}

	ucregexp, err = regexp.Compile(`var uc.*;`)
	if err != nil {
		fmt.Println(err)
	}

	randg = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func test() {
	httpClient := getProxyClient()
	res, err := httpClient.Get(baseurl + string(rune(current)))
	if err != nil {
		fmt.Println(err)
		// panic(err)
	}
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	c, _ := ioutil.ReadAll(res.Body)
	ioutil.WriteFile("b.html", c, 0777)
	defer res.Body.Close()

}

func getProxyClient() *http.Client {
	// proxyurl := "http://127.0.0.1:7890"
	// proxy, err := url.Parse(proxyurl)
	// if err != nil {
	// 	panic(err)
	// }
	netTransport := &http.Transport{
		// Proxy:                 http.ProxyURL(proxy),
		MaxConnsPerHost:       10,
		ResponseHeaderTimeout: time.Second * time.Duration(5),
	}
	httpClient := &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
	return httpClient
}

func main() {
	var works []Work
	fmt.Println("清除工作")
	for true {
		works = getWorks(10)
		if len(works) > 0 {
			for _, work := range works {
				done := getDetail(work.URL)
				if done {
					deleteWorkByID(work.ID)
				}
				time.Sleep(5 * time.Second)
			}
		} else {
			fmt.Println("已清除")
			break
		}
	}
	for true {
		if limit {
			fmt.Println("爬取详情")
			works = getWorks(10)
			if len(works) > 0 {
				for _, work := range works {
					done := getDetail(work.URL)
					if done {
						deleteWorkByID(work.ID)
					}
					time.Sleep(5 * time.Second)
				}
			} else {
				limit = false
				current = 1
				getPage()
			}
		} else {
			fmt.Println("爬取页")
			getPage()
			time.Sleep(5 * time.Second)
		}
	}
}
