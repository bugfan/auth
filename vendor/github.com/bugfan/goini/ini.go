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

	"github.com/go-ini/ini"
)

var CONFIG *config = new(config)

// 读取纯文本配置文件代码 支持 ‘=’ ‘ ’ 等
type config struct {
	conf *ini.File
}

func (self *config) Load(filename string) error {
	conf, err := ini.Load(filename)
	if err != nil {
		log.Fatal(err)
		return err
	}
	self.conf = conf
	return nil
}

func (self *config) GetString(key string) string {
	return self.GetSectionString("", key)
}

func (self *config) GetInt64(key string) int64 {
	return self.GetSectionInt64("", key)
}

// 根据指定的分隔符读取配置
func (self *config) GetSectionString(section string, key string) string {
	if self.conf == nil {
		return ""
	}
	s := self.conf.Section(section)
	return s.Key(key).String()
}

func (self *config) GetSectionInt64(section string, key string) int64 {
	if self.conf == nil {
		return 0
	}
	s := self.conf.Section(section)
	v, _ := s.Key(key).Int64()
	return v
}

// 根据分隔符写配置文件
// func (s *config) Append(key, value, section string) error {
// 	if s.conf == nil {
// 		return errors.New("File not open!")
// 	}
// 	if section == "" {
// 		section = " "
// 	}
// 	data := fmt.Sprintf("%s %s %s", key, section, value)

// 	return s.write([]byte(data))
// }

// func (s *config) AppendTo(key, value, section string) error {
// 	data := fmt.Sprintf("%s %s %s", key, section, value)
// 	return s.conf.Append(s.conf, data)
// }

// func (s *config) write(data []byte) error {
// 	s.conf.Append(s.conf, data)
// 	return nil
// }

func LoadConfig(filename string) {
	CONFIG.Load(filename)
}

/*
*	直接读取环境变量，可以直接拿到对应的类型，支持默认类型
**/

var ENV *env = new(env)

// 读取环境变量相关代码 #所有方法入参都带一个默认值，即没有此环境变量就用默认值
type env struct {
}

func (s *env) GetString(key string, backup string) string {
	v := os.Getenv(key)
	if v != "" {
		return v
	}
	return backup
}
func (s *env) GetInt(key string, backup int64) int64 {
	v := os.Getenv(key)
	if v != "" {
		iv, _ := strconv.Atoi(v)
		return int64(iv)
	}
	return backup
}
func (s *env) GetFloat(key string, backup float64) float64 {
	v := os.Getenv(key)
	if v != "" {
		fv, _ := strconv.ParseFloat(v, 64)
		return fv
	}
	return backup
}
func (s *env) GetBool(key string, backup bool) bool {
	v := os.Getenv(key)
	if v != "" {
		iv, _ := strconv.ParseBool(v)
		return iv
	}
	return backup
}

type MyEnv interface {
	Getenv(string) string
	Load(path string) error
	Getenvd(key, def string) string
	GetAllenv() map[string]string
}

var Env MyEnv

type myEnv struct {
	env     map[string]string
	m       sync.RWMutex
	useFile bool
}

func (s *myEnv) Load(path string) error {
	s.m.RLock()
	defer s.m.RUnlock()
	s.useFile = true
	m, err := ReadFile(path)
	if err != nil {
		log.Fatal("Env file not exists!", err)
		return err
	}
	s.env = m
	return nil
}
func (s *myEnv) Getenv(key string) string {
	s.m.RLock()
	defer s.m.RUnlock()
	if s.useFile {
		return s.env[key]
	}
	return os.Getenv(key)
}
func (s *myEnv) Getenvd(key, def string) (val string) {
	val = s.Getenv(key)
	if val == "" {
		val = def
	}
	return
}
func (s *myEnv) GetAllenv() map[string]string {
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
func NewMyEnv(path ...string) MyEnv {
	if len(path) < 1 {
		return &myEnv{}
	}
	m, err := ReadFile(path[0])
	if err != nil {
		log.Fatal("Env file not exists!", err)
	}
	return &myEnv{useFile: true, env: m}
}

// read one file
func ReadFile(path string) (m map[string]string, err error) {
	if m, err = ReadFiles(path); err != nil {
		return nil, err
	}
	return m, nil
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

// Parse reads an env file from io.Reader, returning a map of keys and values.
func Parse(r io.Reader) (envMap map[string]string, err error) {
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

	return Parse(file)
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

	// Parse the key
	key = splitString[0]
	if strings.HasPrefix(key, "export") {
		key = strings.TrimPrefix(key, "export")
	}
	key = strings.Trim(key, " ")

	// Parse the value
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
