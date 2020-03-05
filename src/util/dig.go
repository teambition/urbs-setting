package util

// util 模块不要引入其它内部模块
import "go.uber.org/dig"

var globalDig = dig.New()

// DigInvoke ...
func DigInvoke(function interface{}, opts ...dig.InvokeOption) error {
	return globalDig.Invoke(function, opts...)
}

// DigProvide ...
func DigProvide(constructor interface{}, opts ...dig.ProvideOption) error {
	return globalDig.Provide(constructor, opts...)
}
