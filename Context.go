package clux

import (
	"log"
	"sync"

	"github.com/hkloudou/clux/binding"
	"github.com/hkloudou/clux/render"
)

type Context struct {
	// This mutex protect Keys map
	Writer  render.ResponseWriter
	mu      sync.RWMutex
	Request binding.RequestTransportData
	Keys    map[string]interface{}
	abort   bool
	// err     error
}

func (c *Context) JSON(obj interface{}) {
	c.Render(render.JSON{Data: obj}, nil)
}
func (c *Context) JSONError(obj interface{}, err error) {
	c.Render(render.JSON{Data: obj}, err)
}

func (c *Context) String(format string, values ...interface{}) {
	// c.Render(render.JSON{Data: obj}, err)
	c.Render(render.String{Format: format, Data: values}, nil)
}

func (c *Context) StringError(format string, err error) {
	// c.Render(render.JSON{Data: obj}, err)
	c.Render(render.String{Format: format, Data: nil}, err)
}

func (c *Context) Error(err error) {
	c.Render(render.Data{}, err)
}

// Data writes some data into the body stream and updates the HTTP code.
func (c *Context) Data(contentType string, data []byte) {
	c.Render(render.Data{
		ContentType: contentType,
		Data:        data,
	}, nil)
}

func (c *Context) Render(r render.Render, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.abort {
		return
	}
	c.abort = true
	if err != nil {
		c.Writer.Header().Set("clux-err", err.Error())
		r.WriteContentType(c.Writer)
		c.Writer.Write([]byte{}) //现在就flush
		return
	}
	r.WriteContentType(c.Writer)
	if err := r.Render(c.Writer); err != nil {
		log.Println("Render.Error:", err)
	}
}
