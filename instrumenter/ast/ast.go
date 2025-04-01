package ast

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
)

type instrumenter struct {
	fileSet    *token.FileSet
	src        string
	pkg        string
	file       *ast.File
	needImport bool
}

func NewInstrumenter(src string, pkg string) *instrumenter {
	return &instrumenter{
		fileSet: token.NewFileSet(),
		src:     src,
		pkg:     pkg,
	}
}

func (i *instrumenter) Instrument() error {
	// 解析Go源码文件
	file, err := parser.ParseFile(i.fileSet, i.src, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parsing %s: %v", i.src, err)
	}
	i.file = file

	// 如果指定了包名，且与源文件包名不同，则跳过
	if i.pkg != "" && i.pkg != file.Name.Name {
		return nil
	}

	// 遍历AST，在函数声明处注入追踪代码
	ast.Inspect(file, i.inspect)

	// 如果需要导入trace包
	if i.needImport {
		i.addImport()
	}

	return nil
}

func (i *instrumenter) inspect(node ast.Node) bool {
	funcDecl, ok := node.(*ast.FuncDecl)
	if !ok {
		return true
	}

	// 不处理main函数
	if funcDecl.Name.Name == "main" {
		return true
	}

	// 不处理已经有defer trace.Trace()()的函数
	for _, stmt := range funcDecl.Body.List {
		if deferStmt, ok := stmt.(*ast.DeferStmt); ok {
			if callExpr, ok := deferStmt.Call.Fun.(*ast.CallExpr); ok {
				if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
					if ident, ok := selExpr.X.(*ast.Ident); ok {
						if ident.Name == "trace" && selExpr.Sel.Name == "Trace" {
							return true
						}
					}
				}
			}
		}
	}

	// 注入defer trace.Trace()()
	i.needImport = true
	traceStmt := &ast.DeferStmt{
		Call: &ast.CallExpr{
			Fun: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("trace"),
					Sel: ast.NewIdent("Trace"),
				},
			},
		},
	}

	funcDecl.Body.List = append([]ast.Stmt{traceStmt}, funcDecl.Body.List...)
	return true
}

func (i *instrumenter) addImport() {
	hasTrace := false
	for _, imp := range i.file.Imports {
		if imp.Path.Value == `"github.com/bigwhite/instrument_trace"` {
			hasTrace = true
			break
		}
	}

	if !hasTrace {
		i.file.Imports = append(i.file.Imports, &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: `"github.com/bigwhite/instrument_trace"`,
			},
		})
	}
}

func (i *instrumenter) Write(wrote bool) error {
	if !i.needImport {
		return nil
	}

	if wrote {
		// 将修改后的AST写回文件
		f, err := os.Create(i.src)
		if err != nil {
			return fmt.Errorf("creating %s: %v", i.src, err)
		}
		defer f.Close()

		if err := format.Node(f, i.fileSet, i.file); err != nil {
			return fmt.Errorf("writing %s: %v", i.src, err)
		}
	} else {
		// 输出到标准输出
		if err := format.Node(os.Stdout, i.fileSet, i.file); err != nil {
			return fmt.Errorf("writing %s: %v", i.src, err)
		}
	}

	return nil
}
