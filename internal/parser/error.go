package parser

import "fmt"

func (p *parser) error(msg string) error {
	return fmt.Errorf("line %d: %s", p.current().Num, msg)
}
