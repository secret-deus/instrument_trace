/*
 * @Author: xhang 1263403710@qq.com
 * @Date: 2025-04-01 11:22:09
 * @LastEditors: xhang 1263403710@qq.com
 * @LastEditTime: 2025-04-01 11:23:51
 * @FilePath: /github.com/instrument_trace/instrumenter/instrumenter.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package instrumenter

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/secret-deus/instrument_trace/instrumenter/ast"
)

func Instrument(path string, pkg string, wrote bool) error {
	// 如果path是目录，则遍历目录下的所有.go文件
	finfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat %s error: %v", path, err)
	}

	if finfo.IsDir() {
		return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// 跳过目录
			if info.IsDir() {
				return nil
			}

			// 只处理.go文件
			if filepath.Ext(path) != ".go" {
				return nil
			}

			return instrumentFile(path, pkg, wrote)
		})
	}

	return instrumentFile(path, pkg, wrote)
}

func instrumentFile(src string, pkg string, wrote bool) error {
	instrumenter := ast.NewInstrumenter(src, pkg)
	if err := instrumenter.Instrument(); err != nil {
		return err
	}
	return instrumenter.Write(wrote)
}
