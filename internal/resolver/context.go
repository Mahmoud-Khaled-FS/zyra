package resolver

// TODO (MAHMOUD) - Default Values {{TOKEN|guest}}
// TODO (MAHMOUD) - Nested resolution {{API_{{ENV}}}}
// TODO (MAHMOUD) - Escaping \{{NOT_A_VAR}}

type Context struct {
	values map[string]string
}

func NewContext() *Context {
	return &Context{
		values: make(map[string]string),
	}
}

func (c *Context) Set(key string, value string) {
	c.values[key] = value
}

func (c *Context) Get(key string) (string, bool) {
	v, ok := c.values[key]
	return v, ok
}

func (c *Context) Merge(resource map[string]string) {
	for k, v := range resource {
		c.values[k] = v
	}
}
