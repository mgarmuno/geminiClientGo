package coms

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
)

const (
	geminiPrefix = "gemini://"
	geminiSufix  = "/"
)

func Request(u string) string {
	u = formatURL(u)
	parsed, err := url.Parse(u)
	conn, err := tls.Dial("tcp", parsed.Host+":1965", &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		fmt.Println("Falied to connect: " + err.Error())
		return err.Error()
	}
	defer conn.Close()

	conn.Write([]byte(u + "\r\n"))

	reader := bufio.NewReader(conn)
	responseHeader, err := reader.ReadString('\n')
	parts := strings.Fields(responseHeader)
	status, err := strconv.Atoi(parts[0][0:1])
	switch status {
	case 1, 3, 6:
		fmt.Print("Unsupported feature!")
	case 2:
		bodyBytes, err := ioutil.ReadAll(reader)
		if err != nil {
			fmt.Println("Error reading body")
			return err.Error()
		}
		body := string(bodyBytes)
		return body
	}
	return ""
}

func formatURL(u string) string {
	if !strings.HasPrefix(u, geminiPrefix) {
		u = geminiPrefix + u
	}
	if !strings.HasSuffix(u, geminiSufix) {
		u = u + geminiSufix
	}
	return u
}
