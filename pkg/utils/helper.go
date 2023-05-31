package utils

import (
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

// SetTable 设置table的样式
func SetTable(table *tablewriter.Table) {
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)
}

// YamlFile2Struct 从指定yaml文件解析内容到struct
func YamlFile2Struct(path string, v any) error {
	content, err := ReadFile(path)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal([]byte(content), v); err != nil {
		return err
	}

	return nil
}

// ReadFile 从文件读取内容
func ReadFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	data, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
