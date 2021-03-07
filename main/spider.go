package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getPage() {
	url := baseurl + strconv.Itoa(current)
	fmt.Println(url)

	defer func() {
		if !limit {
			current++
		}
	}()

	details := make(map[string]string)
	httpClient := getProxyClient()
	res, err := httpClient.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	if res.StatusCode != 200 {
		fmt.Printf("status code error: %d %s \n", res.StatusCode, res.Status)
	}
	c, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	dom, err := goquery.NewDocumentFromReader(bytes.NewReader(c))
	if err != nil {
		log.Fatalln("dom.err", err)
	}

	dom.Find(".movie-box").EachWithBreak(func(i int, s *goquery.Selection) bool {
		code := s.Find("date:nth-of-type(1)").Text()
		date := s.Find("date:nth-of-type(2)").Text()
		link, _ := s.Attr("href")

		year, _ := strconv.Atoi(date[0:4])
		if year < 2017 {
			limit = true
			return false
		}
		fmt.Println(code, ": ", link, ", ", date[0:4])
		details[code] = link
		return true
	})

	for key := range details {
		if exist(key) {
			delete(details, key)
		} else {
			insertWork(details[key])
		}
	}

}

func getDetail(url string) bool {
	fmt.Println(url)

	httpClient := getProxyClient()
	res, err := httpClient.Get(url)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if res.StatusCode != 200 {
		fmt.Printf("status code error: %d %s \n", res.StatusCode, res.Status)
		return false
	}
	c, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	dom, err := goquery.NewDocumentFromReader(bytes.NewReader(c))
	if err != nil {
		log.Fatalln("dom.err", err)
	}
	var gid, uc string
	title, _ := dom.Find(".screencap").Find("img").Attr("title")
	cover, _ := dom.Find(".screencap").Find("img").Attr("src")
	code := dom.Find(".info").Find("p").First().Find("span").Last().Text()
	date := strings.Trim(dom.Find(".info").Find("p").First().Next().After("span").Text(), " ")
	dom.Find("script").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if strings.Index(s.Text(), "gid") > -1 {
			gid = replace(gidregexp.FindString(s.Text()))
			uc = replace(ucregexp.FindString(s.Text()))
		}
		return true
	})
	magnet := getMagnet(gid, uc)
	if magnet == "" {
		deleteWorkByURL(url)
		return false
	}
	saveCover(code, cover)
	video := Video{
		Code:   code,
		Title:  title,
		Cover:  cover,
		Date:   date,
		Magnet: magnet,
	}
	insert(video)
	return true
}

func getMagnet(gid, uc string) (link string) {
	url := rooturl + magneturl + "?" + gid + "&" + uc + "&" + "floor=" + fmt.Sprint(randg.Intn(1e3))
	httpClient := getProxyClient()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Referer", baseurl)
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
		return ""
	}
	c, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var hd, sub, gotlink bool
	dom, err := goquery.NewDocumentFromReader(bytes.NewReader(c))
	if err != nil {
		log.Fatalln("dom.err", err)
	}
	dom.Find("a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if hd && sub {
			link, _ = s.Attr("href")
			gotlink = true
			return false
		}
		if strings.Contains(s.Text(), "HD") {
			hd = true
		} else if strings.Contains(s.Text(), "SUB") {
			sub = true
		} else {
			hd = false
			sub = false
		}
		return true
	})

	sub = false
	hd = false
	if !gotlink {
		dom.Find("a").EachWithBreak(func(i int, s *goquery.Selection) bool {
			if sub {
				link, _ = s.Attr("href")
				gotlink = true
				return false
			}
			if strings.Contains(s.Text(), "SUB") {
				sub = true
			} else {
				sub = false
			}
			return true
		})
	}

	sub = false
	hd = false
	if !gotlink {
		dom.Find("a").EachWithBreak(func(i int, s *goquery.Selection) bool {
			if hd {
				link, _ = s.Attr("href")
				gotlink = true
				return false
			}
			if strings.Contains(s.Text(), "HD") {
				hd = true
			} else {
				hd = false
			}
			return true
		})
	}

	if link == "" {
		dom.Find("a").EachWithBreak(func(i int, s *goquery.Selection) bool {
			var exist bool
			link, exist = s.Attr("href")
			if !exist {
				return true
			}
			gotlink = true
			return false
		})
	}

	return link
}
