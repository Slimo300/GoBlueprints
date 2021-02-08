package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joeshaw/envdecode"

	"github.com/matryer/go-oauth/oauth"
)

var conn net.Conn

func dial(netw, addr string) (net.Conn, error) {
	if conn != nil {
		conn.Close()
		conn = nil
	}

	netC, err := net.DialTimeout(netw, addr, 5*time.Second)
	if err != nil {
		return nil, err
	}
	conn = netC
	return netC, nil
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

var (
	authClient *oauth.Client
	creds      *oauth.Credentials
)

func setupTwitterAuth() {
	var ts struct {
		ConsumerKey    string `env:"SP_TWITTER_KEY,required"`
		ConsumerSecret string `env:"SP_TWITTER_SECRET,required"`
		AccessToken    string `env:"SP_TWITTER_ACCESS_KEY,required"`
		AccessSecret   string `env:"SP_TWITTER_ACCESS_SECRET,required"`
	}

	if err := envdecode.Decode(&ts); err != nil {
		log.Fatalln(err)
	}

	creds = &oauth.Credentials{
		Token:  ts.AccessToken,
		Secret: ts.AccessSecret,
	}

	authClient = &oauth.Client{
		Credentials: oauth.Credentials{
			Token:  ts.ConsumerKey,
			Secret: ts.ConsumerSecret,
		},
	}
}

var (
	authSetupOnce sync.Once
	httpClient    *http.Client
)

func makeRequest(req *http.Request, params url.Values) (*http.Response, error) {
	authSetupOnce.Do(func() {
		setupTwitterAuth()
		httpClient = &http.Client{
			Transport: &http.Transport{
				Dial: dial,
			},
		}
	})

	formEnc := params.Encode()
	// Zadziałało po zmianie kolejności ustawiania wartości headera HTTP :(((((
	req.Header.Set("Authorization", authClient.AuthorizationHeader(creds, "POST", req.URL, params))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-length", strconv.Itoa(len(formEnc)))
	return httpClient.Do(req)
}

type tweet struct {
	Text string
}

func readFromTwitter(votes chan<- string) {

	options, err := loadOptions()
	if err != nil {
		log.Println("failed tp load options:", err)
		return
	}

	hashtags := make([]string, len(options))
	for i := range options {
		hashtags[i] = "#" + strings.ToLower(options[i])
	}

	u, err := url.Parse("https://stream.twitter.com/1.1/statuses/filter.json")
	if err != nil {
		log.Println("creating filter request failed:", err)
		return
	}

	query := make(url.Values)

	query.Set("track", strings.Join(hashtags, ","))

	req, err := http.NewRequest("POST", u.String(), strings.NewReader(query.Encode()))
	if err != nil {
		log.Println("creating filter request failed:", err)
		return
	}

	resp, err := makeRequest(req, query)
	if err != nil {
		log.Println("making request failed:", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		// this is a nice way to see what the error actually is:
		s := bufio.NewScanner(resp.Body)
		s.Scan()
		log.Println(s.Text())
		log.Println(hashtags)
		log.Println("StatusCode =", resp.StatusCode)
		return
	}

	reader := resp.Body
	decoder := json.NewDecoder(reader)
	for {
		var t tweet
		if err := decoder.Decode(&t); err != nil {
			break
		}

		for _, option := range options {
			if strings.Contains(
				strings.ToLower(t.Text),
				strings.ToLower(option),
			) {
				log.Println("vote:", option)
				votes <- option
			}
		}
	}
}

func startTwitterStream(stopchan <-chan struct{}, votes chan<- string) <-chan struct{} {

	stoppedchan := make(chan struct{}, 1)

	go func() {
		defer func() {
			stoppedchan <- struct{}{}
		}()
		for {
			select {
			case <-stopchan:
				log.Println("stopping Twitter...")
				return
			default:
				log.Println("Querying Twitter...")
				readFromTwitter(votes)
				log.Println("  (waiting)")
				time.Sleep(10 * time.Second) // wait before reconnecting
			}
		}
	}()
	return stoppedchan

}