package goini

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

/*
*	Read System ENV
**/

var Env MyEnv

func NewEnv(path ...string) MyEnv {
	if len(path) < 1 {
		return &myEnv{}
	}
	m, err := ReadFile(path...)
	if err != nil {
		log.Fatal("read env file error:", err)
	}
	return &myEnv{useFile: true, env: m}
}

type MyEnv interface {
	Load(path ...string) error
	Getenv(string) string
	Getenvd(key, def string) string
	GetAll() map[string]string
}

type myEnv struct {
	env     map[string]string
	m       sync.RWMutex
	useFile bool
}

func (s *myEnv) Load(path ...string) error {
	s.m.RLock()
	defer s.m.RUnlock()
	m, err := ReadFile(path...)
	if err != nil {
		log.Fatal("read env file error:", err)
		return err
	}
	s.useFile = true
	s.env = m
	return nil
}

func (s *myEnv) Getenv(key string) string {
	return strings.TrimSpace(s.getEnv(key))
}

func (s *myEnv) getEnv(key string) string {
	s.m.RLock()
	defer s.m.RUnlock()
	if s.useFile {
		return s.env[key]
	}
	return os.Getenv(key)
}

func (s *myEnv) Getenvd(key, def string) string {
	return strings.TrimSpace(s.getEnvd(key, def))
}

func (s *myEnv) getEnvd(key, def string) (val string) {
	val = s.Getenv(key)
	if val == "" {
		val = def
	}
	return
}

func (s *myEnv) GetAll() map[string]string {
	s.m.RLock()
	defer s.m.RUnlock()
	if s.useFile {
		m := s.env
		return m
	}
	m := make(map[string]string)
	for _, v := range os.Environ() {
		strs := strings.Split(v, "=")
		if len(strs) > 1 {
			// m[strs[0]] = strs[1]
			m[strings.TrimSpace(strs[0])] = strings.TrimSpace(strs[1])
		}
	}
	return m
}

// read one file
func ReadFile(path ...string) (m map[string]string, err error) {
	return ReadFiles(path...)
}

// Read all env (with same file loading semantics as Load) but return values as
// a map rather than automatically writing values into env
func ReadFiles(filenames ...string) (envMap map[string]string, err error) {
	filenames = filenamesOrDefault(filenames)
	envMap = make(map[string]string)

	for _, filename := range filenames {
		individualEnvMap, individualErr := readFile(filename)

		if individualErr != nil {
			err = individualErr
			return // return early on a spazout
		}

		for key, value := range individualEnvMap {
			envMap[key] = value
		}
	}

	return
}

// parse reads an env file from io.Reader, returning a map of keys and values.
func parse(r io.Reader) (envMap map[string]string, err error) {
	envMap = make(map[string]string)

	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		return
	}

	for _, fullLine := range lines {
		if !isIgnoredLine(fullLine) {
			var key, value string
			key, value, err = parseLine(fullLine, envMap)

			if err != nil {
				return
			}
			envMap[key] = value
		}
	}
	return
}
func filenamesOrDefault(filenames []string) []string {
	if len(filenames) == 0 {
		return []string{".env"}
	}
	return filenames
}

func loadFile(filename string, overload bool) error {
	envMap, err := readFile(filename)
	if err != nil {
		return err
	}

	currentEnv := map[string]bool{}
	rawEnv := os.Environ()
	for _, rawEnvLine := range rawEnv {
		key := strings.Split(rawEnvLine, "=")[0]
		currentEnv[key] = true
	}

	for key, value := range envMap {
		if !currentEnv[key] || overload {
			os.Setenv(key, value)
		}
	}

	return nil
}

func readFile(filename string) (envMap map[string]string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	return parse(file)
}

func parseLine(line string, envMap map[string]string) (key string, value string, err error) {
	if len(line) == 0 {
		err = errors.New("zero length string")
		return
	}

	// ditch the comments (but keep quoted hashes)
	if strings.Contains(line, "#") {
		segmentsBetweenHashes := strings.Split(line, "#")
		quotesAreOpen := false
		var segmentsToKeep []string
		for _, segment := range segmentsBetweenHashes {
			if strings.Count(segment, "\"") == 1 || strings.Count(segment, "'") == 1 {
				if quotesAreOpen {
					quotesAreOpen = false
					segmentsToKeep = append(segmentsToKeep, segment)
				} else {
					quotesAreOpen = true
				}
			}

			if len(segmentsToKeep) == 0 || quotesAreOpen {
				segmentsToKeep = append(segmentsToKeep, segment)
			}
		}

		line = strings.Join(segmentsToKeep, "#")
	}

	firstEquals := strings.Index(line, "=")
	firstColon := strings.Index(line, ":")
	splitString := strings.SplitN(line, "=", 2)
	if firstColon != -1 && (firstColon < firstEquals || firstEquals == -1) {
		//this is a yaml-style line
		splitString = strings.SplitN(line, ":", 2)
	}

	if len(splitString) != 2 {
		err = errors.New("Can't separate key from value")
		return
	}

	// parse the key
	key = splitString[0]
	if strings.HasPrefix(key, "export") {
		key = strings.TrimPrefix(key, "export")
	}
	key = strings.Trim(key, " ")

	// parse the value
	value = parseValue(splitString[1], envMap)
	return
}

func parseValue(value string, envMap map[string]string) string {

	// trim
	value = strings.Trim(value, " ")

	// check if we've got quoted values or possible escapes
	if len(value) > 1 {
		rs := regexp.MustCompile(`\A'(.*)'\z`)
		singleQuotes := rs.FindStringSubmatch(value)

		rd := regexp.MustCompile(`\A"(.*)"\z`)
		doubleQuotes := rd.FindStringSubmatch(value)

		if singleQuotes != nil || doubleQuotes != nil {
			// pull the quotes off the edges
			value = value[1 : len(value)-1]
		}

		if doubleQuotes != nil {
			// expand newlines
			escapeRegex := regexp.MustCompile(`\\.`)
			value = escapeRegex.ReplaceAllStringFunc(value, func(match string) string {
				c := strings.TrimPrefix(match, `\`)
				switch c {
				case "n":
					return "\n"
				case "r":
					return "\r"
				default:
					return match
				}
			})
			// unescape characters
			e := regexp.MustCompile(`\\([^$])`)
			value = e.ReplaceAllString(value, "$1")
		}

		if singleQuotes == nil {
			value = expandVariables(value, envMap)
		}
	}

	return value
}

func expandVariables(v string, m map[string]string) string {
	r := regexp.MustCompile(`(\\)?(\$)(\()?\{?([A-Z0-9_]+)?\}?`)

	return r.ReplaceAllStringFunc(v, func(s string) string {
		submatch := r.FindStringSubmatch(s)

		if submatch == nil {
			return s
		}
		if submatch[1] == "\\" || submatch[2] == "(" {
			return submatch[0][1:]
		} else if submatch[4] != "" {
			return m[submatch[4]]
		}
		return s
	})
}

func isIgnoredLine(line string) bool {
	trimmedLine := strings.Trim(line, " \n\t")
	return len(trimmedLine) == 0 || strings.HasPrefix(trimmedLine, "#")
}

/*
* old struct
 */

type env struct {
}

func (s *env) GetString(key string, backup string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v != "" {
		return v
	}
	return backup
}
func (s *env) GetInt(key string, backup int64) int64 {
	v := strings.TrimSpace(os.Getenv(key))
	if v != "" {
		iv, _ := strconv.Atoi(v)
		return int64(iv)
	}
	return backup
}
func (s *env) GetFloat(key string, backup float64) float64 {
	v := strings.TrimSpace(os.Getenv(key))
	if v != "" {
		fv, _ := strconv.ParseFloat(v, 64)
		return fv
	}
	return backup
}
func (s *env) GetBool(key string, backup bool) bool {
	v := strings.TrimSpace(os.Getenv(key))
	if v != "" {
		iv, _ := strconv.ParseBool(v)
		return iv
	}
	return backup
}
