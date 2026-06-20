package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type Format string

const (
	FormatText    Format = "text"
	FormatJSON    Format = "json"
	FormatVerbose Format = "verbose"
)

func ResolveFormat(flag string) Format {
	switch strings.ToLower(flag) {
	case "json":
		return FormatJSON
	case "verbose":
		return FormatVerbose
	default:
		return FormatText
	}
}

func Print(v any, format Format) {
	PrintTo(os.Stdout, v, format)
}

func PrintTo(w io.Writer, v any, format Format) {
	switch format {
	case FormatJSON:
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		if err := enc.Encode(v); err != nil {
			fmt.Fprintf(os.Stderr, "Error encoding JSON: %s\n", err)
		}
	case FormatVerbose:
		printVerbose(w, v)
	case FormatText:
		printText(w, v)
	}
}

func printText(w io.Writer, v any) {
	switch val := v.(type) {
	case string:
		fmt.Fprintln(w, val)
	case fmt.Stringer:
		fmt.Fprintln(w, val.String())
	default:
		b, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error formatting output: %s\n", err)
			return
		}
		fmt.Fprintln(w, string(b))
	}
}

func printVerbose(w io.Writer, v any) {
	switch val := v.(type) {
	case string:
		fmt.Fprintln(w, val)
	default:
		b, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error formatting output: %s\n", err)
			return
		}
		fmt.Fprintln(w, string(b))
	}
}

func PrintError(msg string) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", msg)
}

func PrintHint(hint string) {
	if hint != "" {
		fmt.Fprintf(os.Stderr, "Hint:  %s\n", hint)
	}
}

func PrintTry(action string) {
	if action != "" {
		fmt.Fprintf(os.Stderr, "Try:   %s\n", action)
	}
}

func MaskKey(key string) string {
	if len(key) <= 12 {
		return key
	}
	return key[:12] + "..." + key[len(key)-4:]
}
