package images

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/tabwriter"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/xianyouq/mydocker/util"
)

var (
	rootUrl          string = "/var/lib/mydocker/images/diff"
	defaultImagePath string = "/var/lib/mydocker/images/config"
	images                  = map[string]*ImageInfo{}
)

type ImageInfo struct {
	util.JsonHelp
	Name        string `json:"name"`       //镜像名
	CreatedTime string `json:"createTime"` //创建时间
}

func LoadImages(outImagePath, tag string) error {
	if tag == "" {
		tag = path.Base(outImagePath)
	}
	_, err := GetImageInfoByTag(tag)
	if err != nil {
		log.Errorf("Check if tag used error %v", err)
		return err
	}
	finalPath := ""
	if check := strings.HasPrefix(outImagePath, "/"); !check {
		pwd, err := os.Getwd()
		if err != nil {
			log.Errorf("Get current location error %v", err)
			return err
		}
		finalPath = pwd + "/" + outImagePath
	} else {
		finalPath = outImagePath
	}
	exist, err := util.PathExists(finalPath)
	if err != nil {
		log.Infof("Fail to judge whether file %s exists. %v", finalPath, err)
		return err
	}
	if !exist {
		log.Errorf("loading image file not exist,%v", outImagePath)
		return errors.New("file not found")
	}
	imageId := util.RandStringBytes(10)
	imageRootUrl := path.Join(rootUrl, imageId)
	exist, err = util.PathExists(imageRootUrl)
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
	imageInfo.dump(defaultImagePath)
	images[tag] = imageInfo
	return nil
}

func ListImages() {
	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprint(w, "NAME\tId\tCreateTime\n")
	for _, image := range images {
		fmt.Fprintf(w, "%s\t%s\t%s\n",
			image.Name,
			image.Id,
			image.CreatedTime,
		)
	}
	if err := w.Flush(); err != nil {
		log.Errorf("Flush error %v", err)
		return
	}
}

func init() {
	if file, err := os.Stat(defaultImagePath); err != nil {
		if os.IsNotExist(err) {
			log.Warnf("default imageConfigPath %v not exist", defaultImagePath)
			return
		}
		if !file.IsDir() {
			log.Warnf("default imageConfigPath %v is not Dir", defaultImagePath)
			return
		}
	}
	dirList, err := ioutil.ReadDir(defaultImagePath)
	if err != nil {
		log.Errorf("init Images error on list file:%v", err)
		return
	}
	for _, file := range dirList {
		if !file.IsDir() {
			imageInfo := &ImageInfo{
				Id: file.Name(),
			}
			err := imageInfo.load(defaultImagePath)
			if err != nil {
				log.Errorf("error occur when load image config :%v", err)
				continue
			}
			images[imageInfo.Name] = imageInfo
		}
	}
}
func GetImageInfoByTag(tag string) (*ImageInfo, error) {
	for _, imageInfo := range images {
		if imageInfo.Name == tag {
			return imageInfo, nil
		}
	}
	imageNotFound := errors.New("image not found")
	return &ImageInfo{}, imageNotFound
}

func GetImagePathByTag(tag string) (string, error) {
	imageInfo, err := GetImageInfoByTag(tag)
	if err != nil {
		return "", nil
	}
	imagePath := path.Join(rootUrl, imageInfo.Id)
	return imagePath, nil
}
