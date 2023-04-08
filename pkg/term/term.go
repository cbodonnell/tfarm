package term

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

// StringPrompt asks for a string value using the label
func StringPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}

// PasswordPrompt asks for a string value using the label.
// The entered value will not be displayed on the screen
// while typing.
func PasswordPrompt(label string) string {
	var s string
	fmt.Fprint(os.Stderr, label+" ")
	b, _ := term.ReadPassword(int(syscall.Stdin))
	s = string(b)
	fmt.Println()
	return s
}

func PrettyJSON(v any) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("error marshaling info to json: %s", err)
	}
	var out bytes.Buffer
	if err := json.Indent(&out, b, "", "  "); err != nil {
		return nil, fmt.Errorf("error indenting json: %s", err)
	}
	return out.Bytes(), nil
}
