package http_client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
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

func GetClientProxy(proxyStr string) *http.Client {
	proxyURL, _ := url.Parse(proxyStr)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy:           http.ProxyURL(proxyURL),
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 2,
	}

	return client
}

func HttpProxy(proxyStr, user, pwd string) *http.Client {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "http", err)
		}
	}()

	urli := url.URL{}

	if !strings.Contains(proxyStr, "http") {
		proxyStr = fmt.Sprintf("http://%s", proxyStr)
	}

	//proxyURL, _ := url.Parse(proxyStr)
	proxyURL, _ := urli.Parse(proxyStr)

	if user != "" && pwd != "" {
		proxyURL.User = url.UserPassword(user, pwd)
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
		Timeout: time.Second * 2,
	}

	return client
}

func Socks5Proxy(proxyStr, user, pwd string) *http.Client {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "http", err)
		}
	}()

	var userAuth proxy.Auth
	if user != "" && pwd != "" {
		userAuth.User = user
		userAuth.Password = pwd
	}

	dialer, err := proxy.SOCKS5("tcp", proxyStr, &userAuth, proxy.Direct)
	if err != nil {
		fmt.Println("proxy.SOCKS5 err:", err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, err error) {
				return dialer.Dial(network, addr)
			},
		},
		Timeout: time.Second * 2,
	}

	return client
}

func GetNetIP() {
	resp, err := http.Get("https://api.myip.la/en?json")
	if err != nil {
		fmt.Println("获取本机ip失败：", err)
		return
	}
	defer resp.Body.Close()

	body, err2 := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err2.Error())
		return
	}

	//fmt.Println(string(body))
	fmt.Println(strings.Split(strings.Split(string(body), ",")[0], "{")[1])
	//fmt.Println(strings.Split(strings.Split(strings.Split(string(body), ",")[0],"{")[1],":")[1])
}

func proxyWithAuth(username, password, proxyStr string) *url.URL {
	//创建代理
	if !strings.Contains(proxyStr, "http") {
		proxyStr = fmt.Sprintf("http://%s", proxyStr)
	}
	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		fmt.Println("无法解析代理URL:", err)
		os.Exit(1)
	}

	//解析代理的鉴权信息
	proxyURL.User = url.UserPassword(username, password)
	return proxyURL
}

func ClientWithAuth(username, password, proxyStr string) *http.Client {
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyWithAuth(username, password, proxyStr)),
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second * 5,
	}

	return client
}

func ResultMap(response *http.Response) (map[string]interface{}, error) {
	var result map[string]interface{}
	body, err := io.ReadAll(response.Body)
	if err == nil {
		err = json.Unmarshal(body, &result)
	}
	return result, err
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
