package toexcel

import (
	"fmt"
	"os"
	"path/filepath"
	"xyz/colorOutput"

	"github.com/xuri/excelize/v2"
)

func SaveToExcel(fileName, sheetName string, content [][]interface{}, A, B string, colwidth [2]string, colInt float64) error {
	eFile := excelize.NewFile()
	defer func() {
		if err := eFile.Close(); err != nil {
			return
		}
	}()

	// 设置sheet name
	if err := eFile.SetSheetName("Sheet1", sheetName); err != nil {
		return err
	}

	// 写入数据
	for index, row := range content {
		cell, err := excelize.CoordinatesToCellName(1, index+1)
		if err != nil {
			return err
		}
		if err := eFile.SetSheetRow(sheetName, cell, &row); err != nil {
			return err
		}
	}

	// 设置第一行样式
	style1, err := eFile.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center"},
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#00BFFF"}, Pattern: 1},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	if err != nil {
		return err
	}
	eFile.SetCellStyle(sheetName, A+"1", B+"1", style1)

	// 设置除第一行以外的样式
	style2, err := eFile.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "left"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return err
	}
	rows, _ := eFile.GetRows(sheetName)
	for i := 2; i <= len(rows); i++ {
		eFile.SetCellStyle(sheetName, fmt.Sprintf("%s%d", A, i), fmt.Sprintf("%s%d", B, i), style2)
	}

	if err := eFile.SetColWidth(sheetName, colwidth[0], colwidth[1], colInt); err != nil {
		return err
	}

	extension := filepath.Ext(fileName)
	if extension == "" {
		fileName = filepath.Join("./output", fileName+".xlsx")
	} else if extension != ".xlsx" {
		fileName = filepath.Join("./output", fmt.Sprintf("%v.xlsx", fileName[:len(fileName)-len(extension)]))
	} else {
		fileName = filepath.Join("./output", fileName)
	}

	if err := eFile.SaveAs(fileName); err != nil {
		return err
	}
	s := fmt.Sprintf("[+] 扫描结果已保存在%v文件中", fileName)
	colorOutput.Colorful.WithFrontColor("green").Println(s)

	return nil
}

func SaveToHtml(fileName string, contents []string) error {
	extension := filepath.Ext(fileName)
	if extension == "" {
		fileName = filepath.Join("./output", fileName+".html")
	} else if extension != ".html" {
		fileName = filepath.Join("./output", fmt.Sprintf("%v.html", fileName[:len(fileName)-len(extension)]))
	} else {
		fileName = filepath.Join("./output", fileName)
	}

	// 创建 HTML 表格的表头
	tableContent := "<tr><th style=\"border: 2px solid black;\" bgcolor=\"#66FFFF\"><font color=\"black\">序号</font></th><th style=\"border: 2px solid black;\" bgcolor=\"#66FFFF\"><font color=\"black\">结果列表</font></th></tr>"

	// 遍历内容切片，生成 HTML 表格的行
	for i, rs := range contents {
		row := fmt.Sprintf("<tr><td style=\"border: 2px solid black; text-align: center;\">%d</td><td style=\"border: 2px solid black;\">%s</td></tr>", i+1, rs)
		tableContent += row
	}

	// 生成最终的 HTML 内容，包含表格和样式
	htmlContent := fmt.Sprintf("<html><head><style>table {margin: 0 auto;}</style></head><body><table style=\"border-collapse: collapse; border: 1px solid black;\">%s</table></body></html>", tableContent)

	// 将 HTML 内容写入文件
	if err := os.WriteFile(fileName, []byte(htmlContent), 0644); err != nil {
		return err
	}

	// 打印保存成功的信息
	s := fmt.Sprintf("[+] 扫描结果已保存在 %v 文件中", fileName)
	colorOutput.Colorful.WithFrontColor("green").Println(s)

	return nil
}
