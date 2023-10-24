package parse

import (
	"flag"
	"fmt"
	"xyz/colorOutput"
	"xyz/poc"
)

func Get_Parse() (file, proxy, proxyType, username, password, outFile string, thread int) {
	//flag.StringVar(&url, "u", "", "目标url")
	flag.StringVar(&file, "f", "", "目标文件")
	flag.StringVar(&proxy, "i", "", "代理地址(127.0.0.1:8080)")
	flag.StringVar(&proxyType, "p", "", "代理类型(none, httpProxy, socksProxy)")
	flag.StringVar(&username, "u", "", "代理用户名")
	flag.StringVar(&password, "w", "", "代理密码")
	flag.StringVar(&outFile, "o", "", "保存结果(excel or html) .xlsx .html")
	flag.IntVar(&thread, "t", 3, "go程：默认为3个go程,可根据电脑性能增加go程")

	flag.Usage = func() {
		fmt.Println("Usage [-f file] [-i proxy] [-p proxyType] [-u username] [-w password] [-t thread] Options: ")
		flag.PrintDefaults()
	}
	flag.Parse()

	return
}

func Banner() {
	// 	str := `
	// ################################################
	// #       #### #     #  ######  #### #   #       #
	// #      #     #     #  #      #     #  #        #
	// #      #     #######  ###### #     # #         #
	// #      #     #     #  #      #     #  #        #
	// #       #### #     #  ######  #### #    #      #
	// ################################################
	// Author: 小燕子
	// @@Version: 1.0
	// Explain: 只适用于辅助扫描！！！
	// QQ: 786474326 (有问题请及时沟通！！！)
	// Attention: 仅供安全测试使用，请勿非法使用！！！
	// POC:%d个
	// -------------------------------------------------------------------------------------------------------
	// `
	poc, _ := poc.GetYamlFile("./poc/poc.yaml")
	s := fmt.Sprintf(`
################################################
#       #### #     #  ######  #### #   #       #
#      #     #     #  #      #     #  #        #
#      #     #######  ###### #     # #         #
#      #     #     #  #      #     #  #        #
#       #### #     #  ######  #### #    #      #
################################################
Author: 小燕子
@@Version: 1.0
Explain: 只适用于辅助扫描！！！
QQ: 786474326 (有问题请及时沟通！！！)
Attention: 仅供安全测试使用，请勿非法使用！！！
POC: 现poc总共%d个！可根据模板自行添加(ps: 注意yaml文件格式,特殊字符的转义)~
-------------------------------------------------------------------------------------------------------
`, len(poc.Poc_content)) //len(poc.Poc_content)
	colorOutput.Colorful.WithFrontColor("green").Println(s)
}
