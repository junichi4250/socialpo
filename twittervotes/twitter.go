package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// 接続
var conn net.Conn

func dial(netw, addr string) (net.Conn, error) {
	// 接続が閉じられているか確認
	if conn != nil {
		conn.Close()
		conn = nil
	}
	netc, err := net.DialTimeout(netw, addr, 5*time.Second)
	if err != nil {
		return nil, err
	}
	conn = netc
	return netc, nil
}

var reader io.ReadCloser

func closeConn() {
	if conn != nil {
		conn.Close()
	}
	if reader != nil {
		reader.Close()
	}
}

type tweet struct {
	Text string
}

// 選択肢を読み込む
func readFromTwitter(votes chan<- string) {
	// 全ての投票での選択肢を取得
	options, err := loadOptions()
	if err != nil {
		log.Println("選択肢の読み込みに失敗しました", err)
		return
	}
	u, err := url.Parse("https://stream.twitter.com/1.1/statuses/filter.json")
	if err != nil {
		log.Println("URLの解析に失敗しました:", err)
		return
	}

	query := make(url.Values)
	query.Set("track", strings.Join(options, ","))
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(query.Encode()))
	if err != nil {
		log.Println("検索のリクエスト作成に失敗しました", err)
		return
	}
	// resp, err := makeRequest(req, query)
	// if err != nil {
	// 	log.Println("検索のリクエストに失敗しました:", err)
	// 	return
	// }
	// reader := resp.Body
	// decoder := json.NewDecoder(reader)

	for {
		var tweet tweet
		if err := decoder.Decode(&tweet); err != nil {
			break
		}
		for _, option := range options {
			if strings.Contains(
				strings.ToLower(tweet.Text),
				strings.ToLower(option),
			) {
				log.Println("投票:", option)
				votes <- option
			}
		}
	}

}
