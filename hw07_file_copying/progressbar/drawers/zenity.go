package drawers

import (
	"fmt"
	"io"
	"log"
	"os/exec"
)

type Zenity struct {
	w io.WriteCloser
}

func (d *Zenity) Draw(percent float64, text string) error {
	var err error

	if d.w == nil {
		cmd := exec.Command("zenity", "--progress", "--auto-close", "--auto-kill", "--no-cancel")
		if d.w, err = cmd.StdinPipe(); err != nil {
			log.Println(err)
			return err
		}
		if err := cmd.Start(); err != nil {
			log.Println(err)
			return err
		}
	}

	_, err = io.WriteString(d.w, fmt.Sprintf("%f\n#%s\n", percent, text))
	return err
}

func (d *Zenity) Close() error {
	if d.w != nil {
		return d.w.Close()
	}
	return nil
}
