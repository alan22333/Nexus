// nexus/pkg/utils/code.go
package utils

import (
	"fmt"
	"math/rand"
)

// GenerateRandomCode 生成一个指定位数的随机数字验证码
func GenerateRandomCode(width int) string {
	// 断言：确保位数在合理范围内
	if width <= 0 {
		width = 6 // 默认为6位
	}

	// 使用纳秒级时间戳作为随机数种子，确保每次生成的随机数都不同
	// 在 Go 1.20+ 中，rand.Seed() 已被弃用，因为默认种子已经足够随机。
	// 但为了兼容旧版本和更明确的意图，这里依然保留。
	// 若使用 Go 1.20+，可以省略 rand.Seed 这一行。
	// rand.Seed(time.Now().UnixNano())

	// format 字符串，例如 "%06v" 表示如果数字不足6位，前面用0补齐
	format := fmt.Sprintf("%%0%dv", width)

	// 计算最大值，例如6位就是 10^6 - 1 = 999999
	max := int(pow10(width)) - 1

	// 生成一个 [0, max] 范围内的随机整数
	code := rand.Intn(max)

	return fmt.Sprintf(format, code)
}

// pow10 是一个辅助函数，计算 10 的 n 次方
func pow10(n int) int64 {
	var result int64 = 1
	for range n {
		result *= 10
	}
	return result
}
