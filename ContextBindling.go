package clux

import "github.com/hkloudou/clux/binding"

func (c *Context) ShouldBindHeader(obj interface{}) error {
	return c.ShouldBindWith(obj, binding.Header)
}

func (c *Context) ShouldBindHeaderRaw(obj interface{}) error {
	return c.ShouldBindWith(obj, binding.HeaderRaw)
}
func (c *Context) ShouldBindJson(obj interface{}) error {
	return c.ShouldBindWith(obj, binding.JSON)
}

func (c *Context) ShouldBindWith(obj interface{}, b binding.Binding) error {
	return b.Bind(c.Request, obj)
}
