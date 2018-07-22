package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

func main() {
	// create listener
	listener, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		panic(err)
	}

	fmt.Println("Server is running")

	for {
		// Acceptはいつでも行えるように常時実行
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go func() {
			// リクエストの処理はノンブロッキング
			fmt.Printf("Accept %v\n", conn.RemoteAddr())
			for {
				// タイムアウトの設定
				conn.SetReadDeadline(time.Now().Add(5 * time.Second))
				// リクエストの読み込み
				request, err := http.ReadRequest(bufio.NewReader(conn))
				if err != nil {
					// タイムアウトもしくはソケットクローズ時は終了
					neterr, ok := err.(net.Error)
					if ok && neterr.Timeout() {
						fmt.Println("Timeout")
						break
					} else if err == io.EOF {
						fmt.Println("EOF")
						break
					}
					// それ以外はエラー
					panic(err)
				}
				// リクエスト表示
				dump, err := httputil.DumpRequest(request, true)
				if err != nil {
					panic(err)
				}
				fmt.Println(string(dump))

				content := "Hello World\n"
				// レスポンスの書き込み
				// HTTP/1.1 (major1, minor1)
				// ContentLengthの指定必須
				response := http.Response{
					StatusCode:    200,
					ProtoMajor:    1,
					ProtoMinor:    1,
					ContentLength: int64(len(content)),
					Body:          ioutil.NopCloser(strings.NewReader(content)),
				}

				response.Write(conn)
			}
			conn.Close()
		}()
	}
}
