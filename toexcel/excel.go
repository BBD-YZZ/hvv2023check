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

	//保存文件
	// if !strings.Contains(fileName, ".") {
	// 	fileName = "./output/" + fileName + ".xlsx"
	// } else if strings.Split(fileName, ".")[1] != "xlsx" {
	// 	fileName = "./output/" + fmt.Sprintf("%v.xlsx", strings.Split(fileName, ".")[0])
	// } else {
	// 	fileName = "./output/" + fileName
	// }
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

	tableContent := "<tr><th style=\"border: 2px solid black;\" bgcolor=\"#66FFFF\"><font color=\"black\">序号</font></th><th style=\"border: 2px solid black;\" bgcolor=\"#66FFFF\"><font color=\"black\">结果列表</font></th></tr>"
	for i, rs := range contents {
		row := fmt.Sprintf("<tr><td style=\"border: 2px solid black; text-align: center;\">%d</td><td style=\"border: 2px solid black;\">%s</td></tr>", i+1, rs)
		tableContent += row
	}

	htmlConent := fmt.Sprintf("<table style=\"border-collapse: collapse; border: 1px solid black;\">%s</table>", tableContent)
	//result := strings.Join(contents, "<br>")
	if err := os.WriteFile(fileName, []byte(htmlConent), 0644); err != nil {
		return err
	}
	s := fmt.Sprintf("[+] 扫描结果已保存在%v文件中", fileName)
	colorOutput.Colorful.WithFrontColor("green").Println(s)

	return nil

}
