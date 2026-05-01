package security

import (
	"fmt"
	"regexp"
)

// validIdentifier matches safe SQL identifier characters only.
var validIdentifier = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// QuoteIdentifier safely quotes a SQL identifier (table or column name)
// to prevent SQL injection. Returns an error if the identifier is unsafe.
func QuoteIdentifier(name string) (string, error) {
	if !validIdentifier.MatchString(name) {
		return "", fmt.Errorf("unsafe SQL identifier: %q", name)
	}
	return fmt.Sprintf(`"%s"`, name), nil
}

// RedactDSN removes passwords from a DSN string for safe logging.
func RedactDSN(dsn string) string {
	// postgres: password=xxx  or  :password@
	pgPass := regexp.MustCompile(`password=[^ ]+`)
	dsn = pgPass.ReplaceAllString(dsn, "password=***")

	urlPass := regexp.MustCompile(`://([^:]+):([^@]+)@`)
	dsn = urlPass.ReplaceAllString(dsn, "://$1:***@")
	return dsn
}
