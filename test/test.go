package test

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Knetic/govaluate"
)

type Response struct {
	Text       string
	StatusCode int
}

func evaluate(expr string, params map[string]interface{}) (bool, error) {
	functions := map[string]govaluate.ExpressionFunction{
		"containsFunc": containsFunc,
		"regexFunc":    regexFunc,
		"equalsFunc":   equalsFunc,
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
		return false, fmt.Errorf("Invalid number of arguments for regexMatch function")
	}

	return args[0] == args[1], nil
}

func Test() {
	expr := `containsFunc(r\.Text, "bootstrapProperties") && containsFunc(r\.Text, "profiles") && equalsFunc(r\.StatusCode, 200)`
	r := Response{
		Text:       "bootstrapProperties--profiles",
		StatusCode: 200,
	}
	fmt.Println(r.Text, r.StatusCode)
	params := map[string]interface{}{
		"r.Text":       r.Text,
		"r.StatusCode": r.StatusCode,
	}

	b, err := evaluate(expr, params)
	if err != nil {
		fmt.Println(err)
	}

	if b {
		fmt.Println("success")
	}
}

func main02() {
	// c, _ := config.GetConfig()
	// fmt.Println(&c.Mysql_Config)
	// data, _ := ceye.Get_Ceye_RS(http_client.GetClient())
	// fmt.Println(data)

	// poc, err := poc.GetYamlFile("./poc/poc.yaml")
	// if err != nil {
	// 	fmt.Println("err:", err)
	// 	return
	// }
	// fmt.Println(&poc.Poc_content)

	// test.Test()
}
