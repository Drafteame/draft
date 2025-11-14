package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/alecthomas/chroma/quick"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-isatty"
)

// Package log provides simple, colored terminal logging for CLI tools.
// It supports leveled logging (Debug, Info, Success, Warn, Error),
// optional timestamps, and automatic color disablement when stdout/stderr
// are not TTYs or when NO_COLOR/LOG_NO_COLOR is set.

var (
	useColor     bool
	showTime     bool
	debugEnabled bool

	lblDebug   = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true) // blue
	lblInfo    = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true) // blue
	lblSuccess = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true) // green
	lblWarn    = lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true) // yellow
	lblError   = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)  // red

	msgStyle = lipgloss.NewStyle()

	// Optional style for code blocks (dim text when color enabled)
	codeStyle = lipgloss.NewStyle().Faint(true)
)

func init() {
	// Detect TTY for both stdout and stderr
	stdoutTTY := isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
	stderrTTY := isatty.IsTerminal(os.Stderr.Fd()) || isatty.IsCygwinTerminal(os.Stderr.Fd())

	// Honor NO_COLOR / LOG_NO_COLOR
	noColor := strings.TrimSpace(strings.ToLower(os.Getenv("NO_COLOR"))) != "" ||
		strings.TrimSpace(strings.ToLower(os.Getenv("LOG_NO_COLOR"))) != ""

	useColor = (stdoutTTY || stderrTTY) && !noColor

	// Timestamp and debug config via env
	showTime = strings.TrimSpace(strings.ToLower(os.Getenv("LOG_TS"))) == "1" ||
		strings.TrimSpace(strings.ToLower(os.Getenv("LOG_TS"))) == "true"
	debugEnabled = strings.TrimSpace(strings.ToLower(os.Getenv("DEBUG"))) == "1" ||
		strings.TrimSpace(strings.ToLower(os.Getenv("DEBUG"))) == "true"

	if !useColor {
		// If color is disabled, reset styles to no-op
		lblDebug = lipgloss.NewStyle()
		lblInfo = lipgloss.NewStyle()
		lblSuccess = lipgloss.NewStyle()
		lblWarn = lipgloss.NewStyle()
		lblError = lipgloss.NewStyle()
		msgStyle = lipgloss.NewStyle()
		codeStyle = lipgloss.NewStyle()
	}
}

// PrintCode prints a multi-line code snippet or file content to stdout with optional syntax highlighting.
// lang controls the syntax highlighter; pass "toml", "yaml", "json", "go", etc.
// If lang is empty or "auto", the lexer will be auto-detected. If highlighting fails or colors are disabled,
// the code is printed as plain text. Intended for TOML, YAML, JSON, Go, etc.
func PrintCode(code string, lang string) {
	// Normalize line endings and ensure trailing newline
	if !strings.HasSuffix(code, "\n") {
		code += "\n"
	}

	// Determine a formatter based on color availability
	formatter := "terminal16m" // truecolor ANSI for capable terminals
	if !useColor {
		formatter = "noop" // no ANSI escape sequences when color is disabled
	}

	// Normalize language input
	lexer := strings.TrimSpace(strings.ToLower(lang))
	if lexer == "" || lexer == "auto" {
		lexer = "" // let chroma auto-detect
	}

	var buf bytes.Buffer
	// Try to highlight; fall back to plain text on any error
	if err := quick.Highlight(&buf, code, lexer, formatter, "github"); err != nil {
		buf.Reset()
		if _, errWrite := buf.WriteString(code); errWrite != nil {
			panic(errWrite)
		}
	}

	rendered := buf.String()
	// Apply optional faint style when colors are enabled (noop otherwise)
	rendered = codeStyle.Render(rendered)

	// Put code on a new line after the label for clarity
	out(os.Stdout, "CODE", lblInfo, "\n"+rendered)
}

// Debug prints a debug-level message to stdout if DEBUG is enabled.
func Debug(args ...any) {
	if !debugEnabled {
		return
	}
	out(os.Stdout, "DEBUG", lblDebug, fmt.Sprint(args...))
}

// Debugf prints a formatted debug-level message to stdout if DEBUG is enabled.
func Debugf(format string, args ...any) {
	if !debugEnabled {
		return
	}
	out(os.Stdout, "DEBUG", lblDebug, fmt.Sprintf(format, args...))
}

// Info prints an info-level message to stdout.
func Info(args ...any) {
	out(os.Stdout, "INFO", lblInfo, fmt.Sprint(args...))
}

// Infof prints a formatted info-level message to stdout.
func Infof(format string, args ...any) {
	out(os.Stdout, "INFO", lblInfo, fmt.Sprintf(format, args...))
}

// Success prints a success-level message to stdout.
func Success(args ...any) {
	out(os.Stdout, "SUCCESS", lblSuccess, fmt.Sprint(args...))
}

// Successf prints a formatted success-level message to stdout.
func Successf(format string, args ...any) {
	out(os.Stdout, "SUCCESS", lblSuccess, fmt.Sprintf(format, args...))
}

// Warn prints a warning-level message to stdout.
func Warn(args ...any) {
	out(os.Stdout, "WARN", lblWarn, fmt.Sprint(args...))
}

// Warnf prints a formatted warning-level message to stdout.
func Warnf(format string, args ...any) {
	out(os.Stdout, "WARN", lblWarn, fmt.Sprintf(format, args...))
}

// Error prints an error-level message to stderr.
func Error(args ...any) {
	out(os.Stderr, "ERROR", lblError, fmt.Sprint(args...))
}

// Errorf prints a formatted error-level message to stderr.
func Errorf(format string, args ...any) {
	out(os.Stderr, "ERROR", lblError, fmt.Sprintf(format, args...))
}

func Exitf(code int, msg string, args ...any) {
	Errorf(msg, args...)
	os.Exit(code)
}

func Exit(code int, args ...any) {
	Error(args...)
	os.Exit(code)
}

func out(w io.Writer, level string, levelStyle lipgloss.Style, message string) {
	label := levelStyle.Render("[" + level + "]")
	if showTime {
		ts := time.Now().Format(time.RFC3339)
		_, _ = fmt.Fprintf(w, "%s %s %s\n", ts, label, message)
		return
	}
	_, _ = fmt.Fprintf(w, "%s %s\n", label, message)
}
