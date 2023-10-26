package http_client

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/proxy"
)

func GetClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 10,
	}

	return client
}

func GetClientFromPool(pool *sync.Pool) *http.Client {
	return pool.Get().(*http.Client)
}

func PutClientToPool(pool *sync.Pool, client *http.Client) {
	pool.Put(client)
}

func SetProxy(proxyType, proxyURL, username, password string) (client *http.Client, err error) {
	switch proxyType {
	case "httpProxy":
		if !strings.Contains(proxyURL, "http") {
			proxyURL = fmt.Sprintf("http://%s", proxyURL)
		}
		proxyURLStr, err := url.Parse(proxyURL)
		if err != nil {
			fmt.Println("无法解析代理URL:", err)
			os.Exit(1)
		}
		if username != "" && password != "" {
			proxyURLStr.User = url.UserPassword(username, password)
		}

		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyURL(proxyURLStr),
		}
		//client.Transport = transport
		client = &http.Client{
			Transport: transport,
			Timeout:   time.Second * 10,
		}

	case "socksProxy":
		var dialer proxy.Dialer
		if username != "" && password != "" {
			proxyURLStr := fmt.Sprintf("socks5://%v:%v@%s", username, password, proxyURL)
			var pURL *url.URL
			pURL, err = url.Parse(proxyURLStr)
			if err != nil {
				fmt.Println("解析代理url失败:", err)
				return
			}
			dialer, err = proxy.SOCKS5("tcp", pURL.Host, &proxy.Auth{
				User:     username,
				Password: password,
			}, proxy.Direct)
			if err != nil {
				fmt.Println("创建SOCKS代理拨号器失败:", err)
				return
			}
		} else {
			dialer, err = proxy.SOCKS5("tcp", proxyURL, nil, proxy.Direct)
			if err != nil {
				fmt.Println("创建SOCKS代理拨号器失败:", err)
				return
			}
		}
		transport := &http.Transport{
			Dial: dialer.Dial,
		}
		// client.Transport = transport
		client = &http.Client{
			Transport: transport,
			Timeout:   time.Second * 10,
		}
	default:
		//不做任何处理
	}
	return
}
