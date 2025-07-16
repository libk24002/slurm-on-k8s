package utils

import (
	"log"
	"strconv"
	"strings"
)

// ParseMemory 将内存字符串（如"4Gi"）转换为整数（MB为单位）
func parse_ram_str(memoryStr string) int {
	// 移除空格
	memoryStr = strings.TrimSpace(memoryStr)
	// 如果字符串为空，返回默认值
	if memoryStr == "" {
		return 1024 // 默认1Gi = 1024MB
	}
	// 提取数字部分和单位部分
	var value float64
	var unit string
	// 查找第一个非数字字符的位置
	var i int
	for i = 0; i < len(memoryStr); i++ {
		if memoryStr[i] < '0' || memoryStr[i] > '9' {
			if memoryStr[i] == '.' {
				continue
			}
			break
		}
	}
	// 解析数字部分
	if i > 0 {
		valueStr := memoryStr[:i]
		if i < len(memoryStr) {
			unit = memoryStr[i:]
		}
		var err error
		value, err = strconv.ParseFloat(valueStr, 64)
		if err != nil {
			log.Printf("Error parsing memory value: %v", err)
			return 1024 // 默认值
		}
	} else {
		// 如果没有数字部分，返回默认值
		return 1024
	}
	// 根据单位转换为MB
	switch strings.ToLower(unit) {
	case "e", "ei", "eib", "exbi", "exbibyte":
		return int(value * 1024 * 1024 * 1024 * 1024 * 1024 * 1024)
	case "p", "pi", "pib", "pebi", "pebibyte":
		return int(value * 1024 * 1024 * 1024 * 1024 * 1024)
	case "t", "ti", "tib", "tebi", "tebibyte":
		return int(value * 1024 * 1024 * 1024 * 1024)
	case "g", "gi", "gib", "gibi", "gibibyte":
		return int(value * 1024)
	case "m", "mi", "mib", "mebi", "mebibyte":
		return int(value)
	case "k", "ki", "kib", "kibi", "kibibyte":
		return int(value / 1024)
	case "eb":
		return int(value * 1000 * 1000 * 1000 * 1000 * 1000 * 1000 / 1024 / 1024)
	case "pb":
		return int(value * 1000 * 1000 * 1000 * 1000 * 1000 / 1024 / 1024)
	case "tb":
		return int(value * 1000 * 1000 * 1000 * 1000 / 1024 / 1024)
	case "gb":
		return int(value * 1000 * 1000 * 1000 / 1024 / 1024)
	case "mb":
		return int(value * 1000 * 1000 / 1024 / 1024)
	case "kb":
		return int(value * 1000 / 1024 / 1024)
	default:
		// 如果没有单位或单位不识别，假设是MB
		return int(value)
	}
}
