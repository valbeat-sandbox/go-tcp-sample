package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
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
		go func()() {()
			fmt.Printf("Accept %v\n", conn.RemoteAddr())
			// リクエストの処理はノンブロッキング
			// レスポンスの読み込み
			request, err := http.ReadRequest(bufio.NewReader(conn))
			if err != nil {
				panic(err)
			}
			// dumpは便利なデバッグ用関数
			dump, err := httputil.DumpRequest(request, true)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(dump))

			// レスポンスの書き込み
			response := http.Response{
				StatusCode: 200,
				ProtoMajor: 1,
				ProtoMinor: 0,
				Body:       ioutil.NopCloser(strings.NewReader("Hello World\n")),
			}

			response.Write(conn)
			conn.Close()
		}()
	}

}
