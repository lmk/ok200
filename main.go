package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
)

var port string
var charset string

func main() {

	port = "8888"                                            // 기본 포트 설정
	charset = "utf8"                                         // 기본 charset 설정
	flag.StringVar(&port, "p", port, "port")                 // 플래그 설정
	flag.StringVar(&charset, "c", charset, "output charset") // 플래그 설정
	flag.Parse()                                             // 플래그 파싱
	fmt.Println("Port", port)
	fmt.Println("Charset", charset)

	handler := http.HandlerFunc(handleRequest)
	http.Handle("/", handler)
	http.ListenAndServe(":"+port, nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// r.Header를 문자열로 변환
	headerString := fmt.Sprintf("%v", r.Header)
	fmt.Println("Header", convertCharset(headerString, charset)) // charset 변환 후 출력

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("Error reading body. Err: %s", err)
	}
	fmt.Println("Body", convertCharset(string(body), charset)) // charset 변환 후 출력
	resp := make(map[string]string)
	resp["message"] = "Status OK"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func convertCharset(input string, charset string) string {
	var encoder *encoding.Encoder
	if charset == "euckr" {
		encoder = korean.EUCKR.NewEncoder()
		// 지원되지 않는 문자가 있는지 확인
		if _, err := encoder.String(input); err != nil {
			return input // 지원되지 않는 경우 기본적으로 utf8로 반환
		}
	} else if charset == "utf8" {
		return input // 기본적으로 utf8로 반환
	} else {
		log.Fatalf("Unsupported charset: %s", charset) // 지원하지 않는 charset 처리
	}

	// 변환
	reader := transform.NewReader(strings.NewReader(input), encoder)
	output, err := io.ReadAll(reader)
	if err != nil {
		log.Fatalf("Error converting charset. Err: %s", err)
	}
	return string(output)
}
