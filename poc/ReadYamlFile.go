package poc

import (
	"os"

	"gopkg.in/yaml.v3"
)

type PocSlice struct {
	Poc_content []Poc_content `yaml:"poc_content,flow"`
}

type Poc_content struct {
	Id   string `yaml:"id"`
	Info struct {
		Name        string `yaml:"name"`
		Author      string `yaml:"author"`
		Severity    string `yaml:"severity"`
		Verified    bool   `yaml:"verified"`
		Description string `yaml:"description"`
	} `yaml:"info"`
	Rules struct {
		Rule0 struct {
			Requests0 struct {
				Method0          string `yaml:"method"`
				Path0           string `yaml:"path"`
				Allow_redirects0 bool   `yaml:"allow_redirects"`
				Headers0         struct {
					Content_Type0 string `yaml:"content-type"`
					User_Agent0   string `yaml:"user-agent"`
					Cookie0       string `yaml:"cookie"`
					TestCMD0      string `yaml:"TestCmd"`
				} `yaml:"headers"`
				Body0 string `yaml:"body"`
			} `yaml:"requests"`
			Matchers0 string `yaml:"matchers"`
		} `yaml:"rule0"`
		Rule1 struct {
			Requests1 struct {
				Method1         string `yaml:"method"`
				Path1            string `yaml:"path"`
				Allow_redirects1 bool   `yaml:"allow_redirects"`
				Headers1         struct {
					Content_Type1 string `yaml:"content-type"`
					User_Agent1   string `yaml:"user-agent"`
					Cookie1       string `yaml:"cookie"`
					TestCMD1      string `yaml:"TestCmd"`
				} `yaml:"headers"`
				Body1 string `yaml:"body"`
			} `yaml:"requests"`
			Matchers1 string `yaml:"matchers"`
		} `yaml:"rule1"`
	} `yaml:"rules"`
}

func GetYamlFile(path string) (*PocSlice, error) {
	poc := new(PocSlice)
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	//fmt.Println(string(b))
	err = yaml.Unmarshal(b, poc)
	if err != nil {
		return nil, err
	}
	return poc, nil
}
