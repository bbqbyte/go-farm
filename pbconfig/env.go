package pbconfig

import (
	"os/user"
	"errors"
	"path/filepath"
	"os"
	"keywea.com/cloud/pblib/pbvalidator"
	"path"
	"fmt"
	"strings"
	"unsafe"
)

const INT_SIZE int = int(unsafe.Sizeof(0))

// Dir returns the home directory for the executing user.
// An error is returned if a home directory cannot be detected.
func HomeDir() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}
	if currentUser.HomeDir == "" {
		return "", errors.New("cannot find user-specific home dir")
	}

	return currentUser.HomeDir, nil
}

// Expand expands the path to include the home directory if the path
// is prefixed with `~`. If it isn't prefixed with `~`, the path is
// returned as-is.
func Expand(path string) (string, error) {
	if len(path) == 0 {
		return path, nil
	}

	if path[0] != '~' {
		return path, nil
	}

	if len(path) > 1 && path[1] != '/' && path[1] != '\\' {
		return "", errors.New("cannot expand user-specific home dir")
	}

	dir, err := HomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, path[1:]), nil
}

// ExpandValueEnv returns value of convert with environment variable.
// Return default value if environment variable is empty or not exist.
// Must begin with $.
//
// It accept value formats "$env" , "${env}" , "${env||defaultValue}".
// Examples:
//	v1 := config.ExpandValueEnv("$the result: $GOPATH")			// return the GOPATH environment variable.
//	v2 := config.ExpandValueEnv("$the result: ${GOPATH}")			// return the GOPATH environment variable.
//	v3 := config.ExpandValueEnv("$the result: ${GOPATHX||/usr/local/go}")	// return the default value "/usr/local/go/".
func ExpandValueEnv(s string) string {
	if len(s) < 2 || s[0] != '$' {
		return s
	}
	buf := make([]byte, 0, 2*len(s)-2)
	i := 1
	for j := 1; j < len(s); j++ {
		if s[j] == '$' && j+1 < len(s) {
			buf = append(buf, s[i:j]...)
			name, defv, w := getShellName(s[j+1:])
			v := os.Getenv(name)
			if v == "" {
				v = defv
			}
			buf = append(buf, v...)
			j += w
			i = j + 1
		}
	}
	return string(buf) + s[i:]
}

// return shellName, defaultValue, pos
func getShellName(s string) (string, string, int) {
	key := ""
	defval := ""
	pos := 1
	vlen := 0
	slen := len(s)
	if s[0] == '{' {
		for i := 1; i < slen; i++ {
			if s[i] == '}' {
				pos = i + 1
				key = s[1:i]
				vlen = len(key)
				if pbvalidator.IsAlphanumeric(key) {
					break
				}
				for j := 0; j < vlen; j++ {
					if key[j] == '|' && (j+1 < vlen && key[j+1] == '|') {
						defval = key[j+2: vlen]
						key = key[0:j]
						break
					}
				}
			}
		}
		return key, defval, pos // Bad syntax; just eat the brace.
	}
	// Scan alphanumerics.
	var i int
	for i = 0; i < slen && isAlphaNum(s[i]); i++ {
	}
	return s[:i], defval, i
}

func isAlphaNum(c uint8) bool {
	return c == '_' || '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}

func getConfigurationFileWithENVPrefix(file, env string) (string, error) {
	var (
		envFile string
		extname = path.Ext(file)
	)

	if extname == "" {
		envFile = fmt.Sprintf("%v.%v", file, env)
	} else {
		envFile = fmt.Sprintf("%v.%v%v", strings.TrimSuffix(file, extname), env, extname)
	}

	if fileInfo, err := os.Stat(envFile); err == nil && fileInfo.Mode().IsRegular() {
		return envFile, nil
	}
	return "", fmt.Errorf("failed to find file %v", file)
}

// system byte order
func IsBigEndian() bool {
	var i int = 0x1
	bs := (*[INT_SIZE]byte)(unsafe.Pointer(&i))
	if bs[0] == 0 { // little endian
		return false
	}
	// big endian
	return true
}
