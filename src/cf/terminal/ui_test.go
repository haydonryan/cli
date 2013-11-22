package terminal

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestSayWithStringOnly(t *testing.T) {
	ui := new(terminalUI)
	out := captureOutput(func() {
		ui.Say("Hello")
	})

	assert.Equal(t, "Hello\n", out)
}

func TestSayWithStringWithFormat(t *testing.T) {
	ui := new(terminalUI)
	out := captureOutput(func() {
		ui.Say("Hello %s", "World!")
	})

	assert.Equal(t, "Hello World!\n", out)
}

func TestConfirmYes(t *testing.T) {
	simulateStdin("y\n", func() {
		ui := new(terminalUI)

		var result bool
		out := captureOutput(func() {
			result = ui.Confirm("Hello %s", "World?")
		})

		assert.True(t, result)
		assert.Contains(t, out, "Hello World?")
	})
}

func TestConfirmNo(t *testing.T) {
	simulateStdin("wat\n", func() {
		ui := new(terminalUI)

		var result bool
		out := captureOutput(func() {
			result = ui.Confirm("Hello %s", "World?")
		})

		assert.False(t, result)
		assert.Contains(t, out, "Hello World?")
	})
}

func simulateStdin(input string, block func()) {
	defer func() {
		stdin = os.Stdin
	}()

	stdinReader, stdinWriter := io.Pipe()
	stdin = stdinReader

	go func() {
		stdinWriter.Write([]byte(input))
		defer stdinWriter.Close()
	}()

	block()
}

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	outC := make(chan string)

	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	w.Close()
	os.Stdout = old
	return <-outC
}
