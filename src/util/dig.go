package util

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
