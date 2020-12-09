package main

import (
	"errors"
	"os"
	"path"

	"github.com/zxdvd/go-libs/std-helper/M"
	"github.com/zxdvd/go-libs/std-helper/encode"
	"github.com/zxdvd/go-libs/std-helper/fs"
)

type LoadConf struct {
	Filename     string
	FallbackHome bool
	config       map[string]interface{}
}

func (c *LoadConf) FindFile() (p string, err error) {
	p = c.Filename
	if fs.Exists(p) {
		return p, nil
	}
	if !c.FallbackHome {
		return "", errors.New("no config file found")
	}
	homedir := os.Getenv("HOME")
	p = path.Join(homedir, c.Filename)
	if fs.Exists(p) {
		return p, nil
	}
	return "", errors.New("no config file found")
}

func (c *LoadConf) Decode() error {
	p, err := c.FindFile()
	if err != nil {
		return err
	}
	return encode.DecodeJSONFile(p, &c.config)
}

func (c *LoadConf) Get(p string) (interface{}, error) {
	if c.config == nil {
		c.config = map[string]interface{}{}
		if err := c.Decode(); err != nil {
			return nil, err
		}
	}
	return M.GetPath(c.config, p)
}

func Get(path, key string) (interface{}, error) {
	conf := &LoadConf{
		Filename:     path,
		FallbackHome: true,
	}
	return conf.Get(key)
}
