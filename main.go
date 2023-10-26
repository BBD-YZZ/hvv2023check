package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
	parse "xyz/argparse"
	"xyz/check_demo"
	"xyz/colorOutput"
	"xyz/http_client"
	"xyz/toexcel"
)

func main() {
	file, proxyURL, proxyType, username, password, outFile, threads := parse.Get_Parse()
	parse.Banner()

	colorOutput.Colorful.WithFrontColor("blue").Println("……………………………………………………………………………………………………请耐心等待扫描马上开始…………………………………………………………………………………………………………………")
	colorOutput.Colorful.WithFrontColor("blue").Println("============================================checking===================================================")

	start := time.Now()

	var urlFileChan = make(chan string)
	var rsChan = make(chan string)
	var exitChan = make(chan bool, threads)
	var rs []string
	var wg sync.WaitGroup
	var allowRedirect bool
	//var rsCounter int32
	var clientPool sync.Pool
	if proxyType == "" || proxyType == "none" {
		clientPool.New = func() interface{} {
			return http_client.GetClient()
		}
	} else if proxyType == "httpProxy" || proxyType == "socksProxy" && proxyURL != "" {
		clientPool.New = func() interface{} {
			client, err := http_client.SetProxy(proxyType, proxyURL, username, password)
			if err != nil {
				fmt.Println("设置代理异常，请检查")
				os.Exit(1)
			}
			return client
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	//go check_demo.Put_URL(urlFileChan, file, ctx)
	go func() {
		defer wg.Done()
		check_demo.Put_URL(http_client.GetClientFromPool(&clientPool), urlFileChan, file, ctx)
	}()

	for i := 0; i < threads; i++ {
		wg.Add(1)
		//go check_demo.Start(urlFileChan, rsChan, exitChan, client, ctx)
		go func() {
			defer wg.Done()
			check_demo.Start(urlFileChan, rsChan, exitChan, http_client.GetClientFromPool(&clientPool), ctx, allowRedirect)
		}()
	}

	for i := 0; i < threads; i++ {
		wg.Add(1)
		//go check_demo.PrintRS(rsChan, &wg)
		go func() {
			defer wg.Done()
			check_demo.PrintRS(rsChan, ctx, &rs)
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < threads; i++ {
			<-exitChan
		}
		close(rsChan)
	}()

	wg.Wait()
	cancel()

	if outFile != "" {
		ext := filepath.Ext(outFile)
		switch ext {
		case ".xlsx":
			var content = [][]interface{}{
				{"序号", "结果列表"},
			}
			for index, result := range rs {
				content = append(content, []interface{}{index + 1, result})
			}
			t := time.Now().Format("20060102")
			err := toexcel.SaveToExcel(t+outFile, "漏洞扫描结果", content, "A", "B", [2]string{"B", "B"}, 100)
			if err != nil {
				s := fmt.Sprintf("[*] 保存excel文件出错:%v", err)
				colorOutput.Colorful.WithFrontColor("red").Println(s)
			}
		case ".html":
			t := time.Now().Format("20060102")
			err := toexcel.SaveToHtml(t+outFile, rs)
			if err != nil {
				s := fmt.Sprintf("[*] 保存html文件出错:%v", err)
				colorOutput.Colorful.WithFrontColor("red").Println(s)
			}
		default:
			s := fmt.Sprintf("[*] 文件%v后缀错误，请填写正确后缀，目前只支持xlsx和html!", outFile)
			colorOutput.Colorful.WithFrontColor("red").Println(s)
		}
	}
	end := time.Now()
	if len(rs) == 0 {
		s := fmt.Sprintf("[x] %v", "本次扫描未发现任何问题！！！")
		colorOutput.Colorful.WithFrontColor("green").Println(s)
	} else {
		s := fmt.Sprintf("[Y] 本次扫描共发现%v个问题！！！", len(rs))
		colorOutput.Colorful.WithFrontColor("green").Println(s)
	}
	colorOutput.Colorful.WithFrontColor("blue").Println("==============================================end======================================================")
	s := fmt.Sprintf("[本次扫描完成，任务总用时:%v]", end.Sub(start))
	colorOutput.Colorful.WithFrontColor("purple").Println(s)
}
