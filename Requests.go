package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func startServer(addr string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "getHandler: incoming request %#v\n", r)
		fmt.Fprintf(w, "getHandler: URL %#v\n", r.URL)
	})
	http.HandleFunc("/raw_body", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		fmt.Fprintf(w, "postHandler: raw body %s\n", string(body))
	})
	fmt.Printf("Started server at %s\n", addr)
	err := http.ListenAndServe(addr, nil)

	if err != nil {
		panic(err)
	}
}

func runGet() {
	url := "http://127.0.0.1:8080/?param=123&param2=test"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error happened", err)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	fmt.Printf("http.Get body %#v\n\n\n\n", string(respBody))
}

func runGetFullRequest() {
	req := &http.Request{
		Method: http.MethodGet,
		Header: http.Header{
			"User-Agent": {"niko/baby"},
		},
	}

	req.URL, _ = url.Parse("http://127.0.0.1:8080/?id=42")
	req.URL.Query().Set("user", "niko")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("error happened", err)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	fmt.Printf("getFullReq response: %#v\n\n\n\n", string(respBody))
}

func runTransportAndPost() {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport,
	}

	data := `{"id": 52, "user": "niko"}`
	body := bytes.NewBufferString(data)

	url := "http://127.0.0.1:8080/raw_body"
	req, _ := http.NewRequest(http.MethodPost, url, body)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(data)))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error happened", err)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	fmt.Printf("runTransportPost %#v\n\n\n", string(respBody))
}

func main() {
	go startServer(":8080")

	time.Sleep(100 * time.Millisecond)

	runGet()
	runGetFullRequest()
	runTransportAndPost()
}
