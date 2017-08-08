package images

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/xianlubird/mydocker/util"
)

var (
	RootUrl          string = "/var/run/mydocker/images/diff/%s"
	defaultImagePath string = "/var/run/mydocker/images/config/%s"
)

type ImageInfo struct {
	Id          string `json:"id"`         //镜像Id
	Name        string `json:"name"`       //镜像名
	CreatedTime string `json:"createTime"` //创建时间
}

func LoadImages(path, tag string) error {
	tagInfo, err := GetImageInfoByTag(tag)
	if err != nil {
		log.Errorf("Check if tag used error %v", err)
		return err
	}
	finalPath := ""
	if check := strings.HasPrefix(path, "/"); !check {
		pwd, err := os.Getwd()
		if err != nil {
			log.Errorf("Get current location error %v", err)
			return err
		}
		finalPath = pwd + "/" + path
	} else {
		finalPath = path
	}
	exist, err := util.PathExists(finalPath)
	if err != nil {
		log.Infof("Fail to judge whether file %s exists. %v", finalPath, err)
		return err
	}
	if !exist {
		log.Errorf("loading image file not exist,%v", path)
		return errors.New("file not found")
	}
	imageId := util.RandStringBytes(10)
	imageRootUrl := fmt.Sprintf(RootUrl, imageId)
	exist, err := util.PathExists(imageRootUrl)
	if err != nil {
		log.Infof("Fail to judge whether dir %s exists. %v", imageRootUrl, err)
		return err
	}
	if !exist {
		if err := os.MkdirAll(imageRootUrl, 0622); err != nil {
			log.Errorf("Mkdir %s error %v", imageRootUrl, err)
			return err
		}

		if _, err := exec.Command("tar", "-xvf", finalPath, "-C", imageRootUrl).CombinedOutput(); err != nil {
			log.Errorf("Untar dir %s error %v", imageRootUrl, err)
			return err
		}
	}
	createTime := time.Now().Format("2006-01-02 15:04:05")
	imageInfo := &ImageInfo{
		Id:          imageId,
		Name:        tag,
		CreatedTime: createTime,
	}
	return nil
}

func GetImageInfoByTag(tag string) (ImageInfo, error) {

}

func 