package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type artiryRange struct {
	min int
	max int
}

var commandRanges = map[string]artiryRange{
	"PING":    {min: 1, max: 2},
	"ECHO":    {min: 2, max: 2},
	"COMMAND": {min: 1, max: 2},
}

func checkArity(cmd string, args []string) error {
	limits, exists := commandRanges[strings.ToUpper(cmd)]
	if !exists {
		return nil
	}
	argCount := len(args)
	if argCount < limits.min || argCount > limits.max {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", cmd)
	}
	return nil
}

func eb(s string, ok bool) string {
	if !ok { return "$-1\r\n" }
	return fmt.Sprintf("$%d\r\n%s\r\n", len(s), s)
}
func es(s string) string { return fmt.Sprintf("+%s\r\n", s) }
func ee(m string) string { return fmt.Sprintf("-%s\r\n", m) }

func handleCommand(args []string) string {
	cmd := strings.ToUpper(args[0])

	switch cmd {
	case "PING":
		// TODO: Return "+PONG\r\n" for no args
		if len(args) == 1 {
			return "+PONG\r\n"
		}
		// TODO: Return bulk string for PING <message>
		if len(args) == 2 {
			return encodeBulkString(args[1])
		}
		if err := checkArity(cmd, args); err != nil {
			return ee(err.Error())
		}
	case "ECHO":
		if len(args) == 2 {
			return eb(args[1], true)
		}
		if err := checkArity(cmd, args); err != nil {
			return ee(err.Error())
		}
	case "COMMAND":
		if len(args) > 1 && strings.ToUpper(args[1]) == "DOCS" {
			return es("OK")
		}
		if err := checkArity(cmd, args); err != nil {
			return ee(err.Error())
		}
		// if es(args[1]) == "DOCS"{
		// 	return es("OK")
		// }
			
	}

	return ee("ERR unknown command '" + args[0] + "'")
}

func encodeBulkString(s string) string {
	return fmt.Sprintf("$%d\r\n%s\r\n", len(s), s)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		args := parseArgs(line)
		response := handleCommand(args)
		fmt.Print(response)
	}
}

func parseArgs(line string) []string {
	var args []string
	var current strings.Builder
	inQuotes := false
	for _, ch := range line {
		switch {
		case ch == '"' && !inQuotes:
			inQuotes = true
		case ch == '"' && inQuotes:
			inQuotes = false
		case ch == ' ' && !inQuotes:
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(ch)
		}
	}
	if current.Len() > 0 {
		args = append(args, current.String())
	}
	return args
}
