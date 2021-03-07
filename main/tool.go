package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path"
	"runtime/debug"
	"strings"
)

func saveCover(code, url string) {
	httpClient := getProxyClient()
	fileSuffix := path.Ext(url)
	filename := code + fileSuffix
	resp, err := httpClient.Get(url)

	if err != nil {
		debug.PrintStack()
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		debug.PrintStack()
		return
	}

	out, err := os.Create(imgpath + "/" + filename)
	if err != nil {
		debug.PrintStack()
		return
	}
	io.Copy(out, bytes.NewReader(body))
}

func replace(s string) string {
	s = strings.Replace(s, "var", "", -1)
	s = strings.Replace(s, ";", "", -1)
	s = strings.Replace(s, " ", "", -1)
	s = strings.Replace(s, "'", "", -1)
	return s
}
