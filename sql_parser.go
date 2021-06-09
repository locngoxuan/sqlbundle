package sqlbundle

import (
	"bufio"
	"bytes"
	"os"
	"strings"
)

const (
	ParserStarted   int = iota // 0
	ParserUpBegin              // 1
	ParserDownBegin            // 3
	ParserClosed
)

func (pg SQLiteDialect) parseStatement(filePath string, up bool) (stmts []string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()
	var buf bytes.Buffer
	scanner := bufio.NewScanner(file)
	stateMachine := ParserStarted

	for scanner.Scan() {
		if stateMachine == ParserClosed {
			break
		}
		line := scanner.Text()
		if strings.HasPrefix(line, "--") {
			cmd := strings.TrimSpace(strings.TrimPrefix(line, "--"))
			if strings.HasPrefix(cmd, "+up BEGIN") && up {
				stateMachine = ParserUpBegin
			} else if strings.HasPrefix(cmd, "+up END") && stateMachine == ParserUpBegin {
				if buf.Len() > 0 {
					statement := buf.String()
					buf.Reset()
					if strings.TrimSpace(statement) != "" {
						stmts = append(stmts, statement)
					}
				}
				stateMachine = ParserClosed
			} else if strings.HasPrefix(cmd, "+down BEGIN") && !up {
				stateMachine = ParserDownBegin
			} else if strings.HasPrefix(cmd, "+down END") && stateMachine == ParserDownBegin {
				if buf.Len() > 0 {
					statement := buf.String()
					buf.Reset()
					if strings.TrimSpace(statement) != "" {
						stmts = append(stmts, statement)
					}
				}
				stateMachine = ParserClosed
			} else {
				// ignore comment
			}
			continue
		}

		if stateMachine == ParserStarted {
			//ignore line due to parser still in started state
			continue
		}

		line = strings.TrimSpace(line)
		if line == "" {
			//ignore empty line
			continue
		}

		if _, err = buf.WriteString(line + "\n"); err != nil {
			break
		}
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}
	// EOF
	return
}


func (pg PostgresDialect) parseStatement(filePath string, up bool) (stmts []string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()
	var buf bytes.Buffer
	scanner := bufio.NewScanner(file)
	stateMachine := ParserStarted

	for scanner.Scan() {
		if stateMachine == ParserClosed {
			break
		}
		line := scanner.Text()
		if strings.HasPrefix(line, "--") {
			cmd := strings.TrimSpace(strings.TrimPrefix(line, "--"))
			if strings.HasPrefix(cmd, "+up BEGIN") && up {
				stateMachine = ParserUpBegin
			} else if strings.HasPrefix(cmd, "+up END") && stateMachine == ParserUpBegin {
				if buf.Len() > 0 {
					statement := buf.String()
					buf.Reset()
					if strings.TrimSpace(statement) != "" {
						stmts = append(stmts, statement)
					}
				}
				stateMachine = ParserClosed
			} else if strings.HasPrefix(cmd, "+down BEGIN") && !up {
				stateMachine = ParserDownBegin
			} else if strings.HasPrefix(cmd, "+down END") && stateMachine == ParserDownBegin {
				if buf.Len() > 0 {
					statement := buf.String()
					buf.Reset()
					if strings.TrimSpace(statement) != "" {
						stmts = append(stmts, statement)
					}
				}
				stateMachine = ParserClosed
			} else {
				// ignore comment
			}
			continue
		}

		if stateMachine == ParserStarted {
			//ignore line due to parser still in started state
			continue
		}

		line = strings.TrimSpace(line)
		if line == "" {
			//ignore empty line
			continue
		}

		if _, err = buf.WriteString(line + "\n"); err != nil {
			break
		}
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}
	// EOF
	return
}

func (od OracleDialect) parseStatement(filePath string, up bool) (stmts []string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()
	var buf bytes.Buffer
	scanner := bufio.NewScanner(file)
	stateMachine := ParserStarted
	for scanner.Scan() {
		if stateMachine == ParserClosed {
			break
		}
		line := scanner.Text()
		if strings.HasPrefix(line, "--") {
			cmd := strings.TrimSpace(strings.TrimPrefix(line, "--"))
			if strings.HasPrefix(cmd, "+up BEGIN") && up {
				stateMachine = ParserUpBegin
			} else if strings.HasPrefix(cmd, "+up END") && stateMachine == ParserUpBegin {
				stateMachine = ParserClosed
			} else if strings.HasPrefix(cmd, "+down BEGIN") && !up {
				stateMachine = ParserDownBegin
			} else if strings.HasPrefix(cmd, "+down END") && stateMachine == ParserDownBegin {
				stateMachine = ParserClosed
			} else {
				// ignore comment
			}
			continue
		}

		if stateMachine == ParserStarted {
			//ignore line due to parser still in started state
			continue
		}

		line = strings.TrimSpace(line)
		if line == "" {
			//ignore empty line
			continue
		}

		if line == "/" {
			statement := buf.String()
			//statement = strings.TrimSuffix(statement, ";\n")
			buf.Reset()
			stmts = append(stmts, statement)
			continue
		}

		if _, err = buf.WriteString(line + "\n"); err != nil {
			break
		}
	}

	if buf.Len() > 0 {
		statement := buf.String()
		//statement = strings.TrimSuffix(statement, ";\n")
		buf.Reset()
		stmts = append(stmts, statement)
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}
	// EOF
	return
}
