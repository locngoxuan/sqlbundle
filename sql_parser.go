package sqlbundle

import (
	"bufio"
	"bytes"
	"os"
	"strings"
)

func parseStatements(filePath string, up bool) (stmts []string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()
	var buf bytes.Buffer
	scanner := bufio.NewScanner(file)
	stateMachine := PARSER_START

	for scanner.Scan() {
		if (stateMachine == PARSER_UP_END && up) || (stateMachine == PARSER_DOWN_END && !up) {
			if buf.Len() > 0 {
				statement := buf.String()
				buf.Reset()
				if strings.TrimSpace(statement) != "" {
					stmts = append(stmts, statement)
				}
			}
			break
		}
		line := scanner.Text()
		if strings.HasPrefix(line, "--") {
			cmd := strings.TrimSpace(strings.TrimPrefix(line, "--"))
			if strings.HasPrefix(cmd, "+up BEGIN") && up {
				stateMachine = PARSER_UP_BEGIN
			} else if strings.HasPrefix(cmd, "+up END") && stateMachine == PARSER_UP_BEGIN {
				stateMachine = PARSER_UP_END
			} else if strings.HasPrefix(cmd, "+down BEGIN") && !up {
				stateMachine = PARSER_DOWN_BEGIN
			} else if strings.HasPrefix(cmd, "+down END") && stateMachine == PARSER_DOWN_BEGIN {
				stateMachine = PARSER_DOWN_END
				if buf.Len() > 0 {
					statement := buf.String()
					buf.Reset()
					if strings.TrimSpace(statement) != "" {
						stmts = append(stmts, statement)
					}
				}
			} else {
				// ignore comment
			}
			continue
		}

		if stateMachine == PARSER_START {
			//ignore line due to parser still in started state
			continue
		}

		line = strings.TrimSpace(line)
		if strings.TrimSpace(line) == "" {
			//ignore empty line
			continue
		}

		if _, err = buf.WriteString(line + "\n"); err != nil {
			break
			//return nil, false, errors.Wrap(err, "failed to write to buf")
		}
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}
	// EOF
	return
}
