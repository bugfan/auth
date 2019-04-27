package goini

import (
	"github.com/go-ini/ini"
)

/*
*	Read ENV from config file
*   Use a configuration Config globally
**/

func LoadConfig(filename string) error {
	return Config.Load(filename)
}

var Config *config = new(config)

type config struct {
	conf *ini.File
}

func (self *config) Load(filename string) (err error) {
	self.conf, err = ini.Load(filename)
	return nil
}

func (self *config) GetString(key string) string {
	return self.GetSectionString("", key)
}

func (self *config) GetInt(key string) int64 {
	return self.GetSectionInt("", key)
}

func (self *config) GetFloat(key string) float64 {
	return self.GetSectionFloat("", key)
}

// read value depends on section token like '=' or ' '
func (self *config) GetSectionString(section string, key string) string {
	if self.conf == nil {
		return ""
	}
	s := self.conf.Section(section)
	return s.Key(key).String()
}

func (self *config) GetSectionInt(section string, key string) int64 {
	if self.conf == nil {
		return 0
	}
	s := self.conf.Section(section)
	v, _ := s.Key(key).Int64()
	return v
}

func (self *config) GetSectionFloat(section string, key string) float64 {
	if self.conf == nil {
		return 0
	}
	s := self.conf.Section(section)
	v, _ := s.Key(key).Float64()
	return v
}
