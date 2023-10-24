package colorOutput

import (
	"fmt"
	"strings"

	"github.com/gogf/gf/container/garray"
)

// 前景 背景 颜色
// ---------------------------------------
// 30  40  黑色
// 31  41  红色
// 32  42  绿色
// 33  43  黄色
// 34  44  蓝色
// 35  45  紫红色
// 36  46  青蓝色
// 37  47  白色
//
// 代码 意义
// -------------------------
//  0  终端默认设置
//  1  高亮显示
//  4  使用下划线
//  5  闪烁
//  7  反白显示
//  8  不可见

const (
	FrontBlack = iota + 30
	FrontRed
	FrontGreen
	FrontYellow
	FrontBlue
	FrontPurple
	FrontCyan
	FrontWhite
)

const (
	BackBlack = iota + 40
	BackRed
	BackGreen
	BackYellow
	BackBlue
	BackPurple
	BackCyan
	BackWhite
)

const (
	ModeDefault   = 0
	ModeHighLight = 1
	ModeLine      = 4
	ModeFlash     = 5
	ModeReWhite   = 6
	ModeHidden    = 7
)

var modeArr = []int{0, 1, 4, 5, 6, 7}

type ColorOutput struct {
	frontColor int
	backColor  int
	mode       int
}

var Colorful ColorOutput
var frontMap map[string]int
var backMap map[string]int

func init() {
	Colorful = ColorOutput{frontColor: FrontGreen, backColor: BackBlack, mode: ModeDefault}

	frontMap = make(map[string]int)
	frontMap["black"] = FrontBlack
	frontMap["red"] = FrontRed
	frontMap["green"] = FrontGreen
	frontMap["yellow"] = FrontYellow
	frontMap["blue"] = FrontBlue
	frontMap["purple"] = FrontPurple
	frontMap["cyan"] = FrontCyan
	frontMap["white"] = FrontWhite

	backMap = make(map[string]int)
	backMap["black"] = BackBlack
	backMap["red"] = BackRed
	backMap["green"] = BackGreen
	backMap["yellow"] = BackYellow
	backMap["blue"] = BackBlue
	backMap["purple"] = BackPurple
	backMap["cyan"] = BackCyan
	backMap["white"] = BackWhite
}

// 其中0x1B是标记，[开始定义颜色，依次为：模式，背景色，前景色，0代表恢复默认颜色。
func (c ColorOutput) Println(str interface{}) {
	fmt.Println(fmt.Sprintf("%c[%d;%d;%dm%s%c[0m", 0x1B, c.mode, c.backColor, c.frontColor, str, 0x1B))
}

// 设置前景色
func (c ColorOutput) WithFrontColor(color string) ColorOutput {
	color = strings.ToLower(color)
	co, ok := frontMap[color]
	if ok {
		c.frontColor = co
	}

	return c
}

// 设置背景色
func (c ColorOutput) WithBackColor(color string) ColorOutput {
	color = strings.ToLower(color)
	co, ok := backMap[color]
	if ok {
		c.backColor = co
	}

	return c
}

// 设置模式
func (c ColorOutput) WithMode(mode int) ColorOutput {
	a := garray.NewIntArrayFrom(modeArr, true)
	bo := a.Contains(mode)
	if bo {
		c.mode = mode
	}

	return c
}
