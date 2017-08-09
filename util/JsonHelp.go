package util

import (
	"encoding/json"
	"os"
	"path"

	"github.com/Sirupsen/logrus"
)

type JsonHelp struct {
	Id string
}

func (self *JsonHelp) dump(dumpPath string) error {
	if _, err := os.Stat(dumpPath); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(dumpPath, 0644)
		} else {
			return err
		}
	}

	jsonPath := path.Join(dumpPath, self.Id)
	jsonFile, err := os.OpenFile(jsonPath, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		logrus.Errorf("dump to file %s error:%v", jsonPath, err)
		return err
	}
	defer jsonFile.Close()

	JsonStr, err := json.Marshal(self)
	if err != nil {
		logrus.Errorf("dump to file %s error:%v", jsonPath, err)
		return err
	}

	_, err = jsonFile.Write(JsonStr)
	if err != nil {
		logrus.Errorf("dump to file %s error:%v", jsonPath, err)
		return err
	}
	return nil
}

func (self *JsonHelp) remove(dumpPath string) error {
	if _, err := os.Stat(path.Join(dumpPath, self.Id)); err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	} else {
		return os.Remove(path.Join(dumpPath, self.Id))
	}
}

func (self *JsonHelp) load(dumpPath string) error {
	ConfigFile, err := os.Open(dumpPath)
	defer ConfigFile.Close()
	if err != nil {
		return err
	}
	JsonStr := make([]byte, 2000)
	n, err := ConfigFile.Read(JsonStr)
	if err != nil {
		return err
	}

	err = json.Unmarshal(JsonStr[:n], self)
	if err != nil {
		logrus.Errorf("load file %s error:%v", ConfigFile, err)
		return err
	}
	return nil
}
