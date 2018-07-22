package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
)

func main() {
	sendMessages := []string{
		"メッセージ1",
		"メッセージ2",
		"メッセージ3",
	}
	current := 0
	var conn net.Conn
	var err error
	for {
		// リトライ用にループ内でコネクション
		if conn == nil {
			conn, err = net.Dial("tcp", "localhost:8888")
			if err != nil {
				panic(err)
			}
			fmt.Printf("Access: %d\n", current)
		}
		// リクエストを作成
		request, err := http.NewRequest(
			"GET",
			"http://localhost:8888",
			strings.NewReader(sendMessages[current]))
		if err != nil {
			panic(err)
		}
		err = request.Write(conn)
		if err != nil {
			panic(err)
		}
		// レスポンスを読み込み
		response, err := http.ReadResponse(bufio.NewReader(conn), request)
		if err != nil {
			fmt.Println("Retry")
			conn = nil
			continue
		}
		// 結果を表示
		dump, err := httputil.DumpResponse(response, true)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(dump))
		// 全部送信済み
		current++
		if current == len(sendMessages) {
			break
		}
	}
}
