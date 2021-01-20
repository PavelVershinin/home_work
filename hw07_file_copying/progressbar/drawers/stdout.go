package drawers

import (
	"fmt"
	"os"
)

type Stdout struct{}

func (d *Stdout) Draw(percent float64, text string) error {
	_, err := fmt.Fprintf(os.Stdout, "\r%s", text)
	return err
}

func (d *Stdout) Close() error {
	_, err := fmt.Fprint(os.Stdout, "\n")
	return err
}
