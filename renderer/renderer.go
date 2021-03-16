package renderer

import (
	"io"
)

type Renderer interface {
	Render(io.Writer, string, map[string]interface{}, string) (int, error)
}
