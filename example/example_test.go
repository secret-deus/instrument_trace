/*
 * @Author: xhang 1263403710@qq.com
 * @Date: 2025-04-01 11:16:01
 * @LastEditors: xhang 1263403710@qq.com
 * @LastEditTime: 2025-04-01 11:21:23
 * @FilePath: /github.com/instrument_trace/example/example_test.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package trace_test

import (
	trace "github.com/secret-deus/instrument_trace"
)

func a() {
	defer trace.Trace()()
	b()
}

func b() {
	defer trace.Trace()()
	c()
}

func c() {
	defer trace.Trace()()
	d()
}

func d() {
	defer trace.Trace()()
}

func ExampleTrace() {
	a()
	// Output:
	// g[00001]:    ->github.com/secret-deus/instrument_trace/example_test.a
	// g[00001]:        ->github.com/secret-deus/instrument_trace/example_test.b
	// g[00001]:            ->github.com/secret-deus/instrument_trace/example_test.c
	// g[00001]:                ->github.com/secret-deus/instrument_trace/example_test.d
	// g[00001]:                <-github.com/secret-deus/instrument_trace/example_test.d
	// g[00001]:            <-github.com/secret-deus/instrument_trace/example_test.c
	// g[00001]:        <-github.com/secret-deus/instrument_trace/example_test.b
	// g[00001]:    <-github.com/secret-deus/instrument_trace/example_test.a
}
