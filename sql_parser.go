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
			statement := buf.String()
			buf.Reset()
			stmts = append(stmts, statement)
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

		//if strings.HasSuffix(line, ";") {
		//	statement := buf.String()
		//	statement = strings.TrimSuffix(statement, ";\n")
		//	buf.Reset()
		//	stmts = append(stmts, statement)
		//}
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}
	// EOF
	return
}
