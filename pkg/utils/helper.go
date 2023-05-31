package utils

import (
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"time"
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

// NsToTime 解析纳秒为日期
func NsToTime(ns int64) string {
	tm := time.Unix(0, 1685526846768616411)
	return tm.Format("2006-01-02 15:04:05")
}
