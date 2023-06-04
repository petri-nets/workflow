package flow

import (
	"os"
	"testing"
)

// 测试包钩子
func TestMain(m *testing.M) {
	// 运行测试前的以前初始化工作

	exitCode := m.Run()
	// 运行测试后的一些清理工作

	// 退出
	os.Exit(exitCode)
}
