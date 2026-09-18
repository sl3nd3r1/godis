package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type artiryRange struct {
	min int
	max int
}

var store = map[string]string{}

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
	case "SET":
		if err := checkArity(cmd, args); err != nil {
			return ee(err.Error())
		}
		store[args[1]] = args[2]
		return es("OK")
	case "GET":
		if err := checkArity(cmd, args); err != nil {
			return ee(err.Error())
		}
		value, ok := store[args[1]]
		return eb(value, ok)	
			
	}

	return ee("ERR unknown command '" + args[0] + "'")
}

func encodeBulkString(s string) string {
	return fmt.Sprintf("$%d\r\n%s\r\n", len(s), s)
}

func main() {
	r := bufio.NewReader(os.Stdin)
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	for {
		args, err := parseArgs(r)
		if err != nil {
			return
		}
		w.WriteString(handleCommand(args))
		w.Flush()
	}
}


// parseRequest reads one RESP array from r and returns its arg list.
func parseArgs(r *bufio.Reader) ([]string, error) {
	// Read the *N\r\n line
	line, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}
	if !strings.HasPrefix(line, "*") {
		return nil, fmt.Errorf("expected '*', got %q", line)
	}
	// TODO: 1. Parse N from line[1:] (trim \r\n)
	n, err := strconv.Atoi(strings.TrimSpace(line[1:]))
	if err != nil {
		return nil, fmt.Errorf("invalid array length: %q", line)
	}
	args := make([]string, 0, n)
	for i := 0; i < n; i++ {
		line, err = r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		strLen, err := strconv.Atoi(strings.TrimSpace(line[1:]))
		if err != nil {
			return nil, fmt.Errorf("invalid bulk string length: %q", line)
		}
		buf := make([]byte, strLen)
		_, err = io.ReadFull(r, buf)
		if err != nil {
			return nil, err
		}
		args = append(args, string(buf))
		// Read the trailing \r\n after the bulk string
		// Option1: Use ReadString to read until '\n' and discard the trailing \r\n
		// _, err = r.ReadString('\n')
		// if err != nil {
		// 	return nil, err
		// }
		// Option2: Use Discard to skip the next 2 bytes (\r\n)
		// _, err =r.Discard(2)
		// if err != nil {
		// 	return nil, err
		// }
		// Option3: Use io.ReadFull to read exactly 2 bytes and check if they are \r\n
		var crlf [2]byte
		_, err = io.ReadFull(r, crlf[:])
		if err != nil {
			return nil, err
		}
		if crlf[0] != '\r' || crlf[1] != '\n'{
			return nil, fmt.Errorf("expected CRLF after bulk string, got %q", crlf)
		}
	}

	//=========================DONE=============================
	// TODO: 2. Loop N times: read $len\r\n, then read exactly len bytes + trailing \r\n
	// TODO: 3. Use io.ReadFull(r, buf) — DO NOT use ReadString once you know the byte length,
	//          because bulk bodies may contain \r\n.
	// TODO: 4. Return the assembled []string
	//=========================DONE=============================
	_ = strconv.Atoi
	_ = io.ReadFull
	return args, nil
}
