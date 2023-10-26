package check_demo

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
	"xyz/colorOutput"
	"xyz/poc"

	"github.com/Knetic/govaluate"
)

type response struct {
	Text       string
	StatusCode int
}

func Start(rch chan string, wch chan string, bch chan bool, client *http.Client, ctx context.Context, allowRedirect bool) {
	var r response
	var params map[string]interface{}

	for url := range rch {
		select {
		case <-ctx.Done():
		default:
			poc, err := poc.GetYamlFile("./poc/poc.yaml")
			if err != nil {
				//fmt.Println(err)
				colorOutput.Colorful.WithFrontColor("red").Println("[*] 读取yaml文件出错")
				continue
			}
			for _, v := range poc.Poc_content {
				var req *http.Request
				var body io.Reader
				var err error
				if v.Rules.Rule0.Requests0.Method0 != "GET" && v.Rules.Rule0.Requests0.Method0 != "POST" {
					continue
				}
				if v.Rules.Rule0.Requests0.Method0 == "GET" || v.Rules.Rule0.Requests0.Method0 == "POST" {
					// allow_redirects := v.Rules.Rule.Requests.Allow_redirects
					path := v.Rules.Rule0.Requests0.Path0
					headers := v.Rules.Rule0.Requests0.Headers0
					allowRedirect = v.Rules.Rule0.Requests0.Allow_redirects0
					if v.Rules.Rule0.Requests0.Method0 == "GET" {
						req, err = http.NewRequest("GET", url+path, nil)
					} else {
						body = bytes.NewReader([]byte(v.Rules.Rule0.Requests0.Body0))
						req, err = http.NewRequest("POST", url+path, body)
					}

					if err != nil {
						colorOutput.Colorful.WithFrontColor("red").Println("[*] " + url + "http.NewRequest 请求错误,请检测网络是否可达!!!")
						continue
					}

					// 设置请求头部信息
					if headers.User_Agent0 != "" {
						req.Header.Set("User-Agent", headers.User_Agent0)
					}
					if headers.Content_Type0 != "" {
						req.Header.Set("Content-Type", headers.Content_Type0)
					}
					if headers.Cookie0 != "" {
						cookie := &http.Cookie{
							Name:  "Cookie",
							Value: headers.Cookie0,
						}
						req.AddCookie(cookie)
					}
					if headers.TestCMD0 != "" {
						req.Header.Set("TestCmd", headers.TestCMD0)
					}
				} else {
					continue
				}

				// 设置 client 的 CheckRedirect 字段
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					if !allowRedirect {
						return http.ErrUseLastResponse
					}
					return nil
				}

				resp, err := client.Do(req)
				if err != nil {
					//colorOutput.Colorful.WithFrontColor("red").Println("[*] " + url + " client.Do请求错误,请检测网络是否可达!!!")
					continue
				}
				defer resp.Body.Close()

				respBody, err := io.ReadAll(resp.Body)
				if err != nil {
					colorOutput.Colorful.WithFrontColor("red").Println("[*] io.ReadAll(resp.Body)出错！！！")
					continue
				}

				matchers0 := v.Rules.Rule0.Matchers0
				if matchers0 != "" {
					if !strings.Contains(matchers0, "Set-Cookie") {
						r = response{
							Text:       string(respBody),
							StatusCode: resp.StatusCode,
						}
						params = map[string]interface{}{
							"r.Text":       r.Text,
							"r.StatusCode": r.StatusCode,
						}
						result, err := evaluate(matchers0, params)
						if err != nil {
							colorOutput.Colorful.WithFrontColor("red").Println(err)
							continue
						}

						if result {
							rs := "[+] " + url + " | 存在" + v.Info.Name
							wch <- rs
						} else {
							continue
						}

					} else {
						new_cookie := resp.Header.Get("Set-Cookie")
						if v.Rules.Rule1.Requests1.Method1 == "GET" || v.Rules.Rule1.Requests1.Method1 == "POST" {
							path1 := v.Rules.Rule1.Requests1.Path1
							headers1 := v.Rules.Rule1.Requests1.Headers1
							if v.Rules.Rule1.Requests1.Method1 == "GET" {
								req, err = http.NewRequest("GET", url+path1, nil)
							} else {
								body = bytes.NewReader([]byte(v.Rules.Rule1.Requests1.Body1))
								req, err = http.NewRequest("POST", url+path1, body)
							}

							if err != nil {
								colorOutput.Colorful.WithFrontColor("red").Println("[*] " + url + "http.NewRequest 请求错误,请检测网络是否可达!!!")
								continue
							}
							if headers1.User_Agent1 != "" {
								req.Header.Set("User-Agent", headers1.User_Agent1)
							}
							if headers1.Content_Type1 != "" {
								req.Header.Set("Content-Type", headers1.Content_Type1)
							}

							req.Header.Set("Cookie", new_cookie)

							if headers1.TestCMD1 != "" {
								req.Header.Set("TestCmd", headers1.TestCMD1)
							}
							client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
								if !allowRedirect {
									return http.ErrUseLastResponse
								}
								return nil
							}

							resp1, err := client.Do(req)
							if err != nil {
								//colorOutput.Colorful.WithFrontColor("red").Println("[*] " + url + " client.Do请求错误,请检测网络是否可达!!!")
								continue
							}
							defer resp1.Body.Close()

							respBody1, err := io.ReadAll(resp1.Body)
							if err != nil {
								colorOutput.Colorful.WithFrontColor("red").Println("[*] io.ReadAll(resp1.Body)出错！！！")
								continue
							}

							matchers1 := v.Rules.Rule1.Matchers1
							if matchers1 != "" {
								r = response{
									Text:       string(respBody1),
									StatusCode: resp1.StatusCode,
								}
								params = map[string]interface{}{
									"r.Text":       r.Text,
									"r.StatusCode": r.StatusCode,
								}
								result, err := evaluate(matchers1, params)
								if err != nil {
									s := fmt.Sprintf("[*] %v/%v", url, err)
									colorOutput.Colorful.WithFrontColor("blue").Println(s)
									continue
								}

								if result {
									rs := "[+] " + url + " | 存在" + v.Info.Name
									wch <- rs
								} else {
									continue
								}
							}
						} else {
							continue
						}
					}
				} else {
					if v.Rules.Rule1.Requests1.Method1 == "GET" || v.Rules.Rule1.Requests1.Method1 == "POST" {
						path1 := v.Rules.Rule1.Requests1.Path1
						headers1 := v.Rules.Rule1.Requests1.Headers1
						if v.Rules.Rule1.Requests1.Method1 == "GET" {
							req, err = http.NewRequest("GET", url+path1, nil)
						} else {
							body = bytes.NewReader([]byte(v.Rules.Rule1.Requests1.Body1))
							req, err = http.NewRequest("POST", url+path1, body)
						}

						if err != nil {
							colorOutput.Colorful.WithFrontColor("red").Println("[*] " + url + "http.NewRequest 请求错误,请检测网络是否可达!!!")
							continue
						}
						if headers1.User_Agent1 != "" {
							req.Header.Set("User-Agent", headers1.User_Agent1)
						}
						if headers1.Content_Type1 != "" {
							req.Header.Set("Content-Type", headers1.Content_Type1)
						}

						if headers1.Cookie1 != "" {
							cookie := &http.Cookie{
								Name:  "Cookie",
								Value: headers1.Cookie1,
							}
							req.AddCookie(cookie)
						}

						if headers1.TestCMD1 != "" {
							req.Header.Set("TestCmd", headers1.TestCMD1)
						}
						// 设置 client 的 CheckRedirect 字段
						client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
							if !allowRedirect {
								return http.ErrUseLastResponse
							}
							return nil
						}

						resp1, err := client.Do(req)
						if err != nil {
							//colorOutput.Colorful.WithFrontColor("red").Println("[*] " + url + " client.Do请求错误,请检测网络是否可达!!!")
							continue
						}
						defer resp1.Body.Close()

						respBody1, err := io.ReadAll(resp1.Body)
						if err != nil {
							colorOutput.Colorful.WithFrontColor("red").Println("[*] io.ReadAll(resp1.Body)出错！！！")
							continue
						}

						matchers1 := v.Rules.Rule1.Matchers1
						if matchers1 != "" {
							r = response{
								Text:       string(respBody1),
								StatusCode: resp1.StatusCode,
							}
							params = map[string]interface{}{
								"r.Text":       r.Text,
								"r.StatusCode": r.StatusCode,
							}
							result, err := evaluate(matchers1, params)
							if err != nil {
								s := fmt.Sprintf("%v/%v", url, err)
								colorOutput.Colorful.WithFrontColor("red").Println(s)
								continue
							}

							if result {
								rs := "[+] " + url + " | 存在" + v.Info.Name
								wch <- rs
							} else {
								continue
							}
						}
					} else {
						continue
					}

				}

			}

		}
	}
	bch <- true
}

func readFile(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			urls = append(urls, line)
		}
		//urls = append(urls, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Println("读取文件失败，错误:", err)
		return nil
	}
	return urls
}

func checkUrl(client *http.Client, url string, timeout time.Duration) bool {
	// tr := &http.Transport{
	// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	// }

	// client := http.Client{
	// 	Transport: tr,
	// 	Timeout:   timeout,
	// }

	resp, err := client.Get(url)
	if err != nil {
		//fmt.Println("[*] " + url + "请求错误,请检测网络是否可达!")
		colorOutput.Colorful.WithFrontColor("red").Println("[-] " + url + "请求无法访问,请手动检查目标是否可正常访问!!!")
		return false
	}
	defer resp.Body.Close()
	return true
}

func Put_URL(client *http.Client, ch chan string, path string, ctx context.Context) {
	url := readFile(path)
	for _, v := range url {
		if !strings.HasPrefix(v, "http://") || !strings.HasPrefix(v, "https://") {
			v = strings.Join([]string{"http://", v}, "")
		}
		if checkUrl(client, v, time.Second*2) {
			select {
			case <-ctx.Done():
				return
			default:
				ch <- v
			}
		}
	}
	close(ch)
}

func PrintRS(ch chan string, ctx context.Context, rss *[]string) {
	for {
		select {
		case <-ctx.Done():
			return
		case rs, ok := <-ch:
			if !ok {
				return
			}
			colorOutput.Colorful.WithFrontColor("yellow").Println(rs)
			//fmt.Println(strings.Split(rs, "]")[1])
			*rss = append(*rss, strings.Split(rs, "]")[1])
		}
	}
}


func evaluate(expr string, params map[string]interface{}) (bool, error) {
	functions := map[string]govaluate.ExpressionFunction{
		"containsFunc": containsFunc,
		"regexFunc":    regexFunc,
		"equalsFunc":   equalsFunc,
		"lengthFunc":   lengthFunc,
	}

	expression, err := govaluate.NewEvaluableExpressionWithFunctions(expr, functions)
	if err != nil {
		return false, err
	}

	result, err := expression.Evaluate(params)
	if err != nil {
		return false, err
	}

	eval, ok := result.(bool)

	if !ok {
		return false, fmt.Errorf("Expression does not evaluate to a boolean result")
	}

	return eval, nil
}

func containsFunc(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("Invalid number of arguments for Contains function")
	}
	s1, ok1 := args[0].(string)
	s2, ok2 := args[1].(string)
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("Invalid argument types for Contains function")
	}
	return strings.Contains(s1, s2), nil
}

func regexFunc(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return false, fmt.Errorf("Invalid number of arguments for regexMatch function")
	}

	pattern, ok := args[0].(string)
	if !ok {
		return false, fmt.Errorf("Invalid pattern argument for regexMatch function")
	}

	input, ok := args[1].(string)
	if !ok {
		return false, fmt.Errorf("Invalid input argument for regexMatch function")
	}

	match, err := regexp.MatchString(pattern, input)
	if err != nil {
		return false, err
	}

	return match, nil
}

func equalsFunc(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return false, fmt.Errorf("Invalid number of arguments for equalsFunc function")
	}

	return args[0] == args[1], nil
}

func lengthFunc(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return false, fmt.Errorf("Invalid number of arguments for lengthFunc function")
	}

	str, ok1 := args[0].(string)
	length, ok2 := args[1].(float64)
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("Invalid argument types for Contains lengthFunc function")
	}
	return len(str) == int(length), nil
}
