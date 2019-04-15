package main

import (
	"image"
	_ "image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/kataras/iris/core/errors"
)

// 1. 判断图片缩放比例
// 2. 调整数据
// 3. 保存为新的文件

type Point struct {
	x int `json:"x"`
	y int `json:"y"`
}

type Component_detail struct {
	PointList  []Point `json:"pointList"`
	Tag_values string  `json:"tag_values"`
}
type Component struct {
	Component_id     int    `json:"component_id"`
	Component_type   string `json:"component_type"`
	Component_detail `json:"component_detail"`
}

type Label struct {
	Label_ID  int         `json:"label_ID"`
	Image     string      `json:"image"`
	Component []Component `json:"component"`
}

type Result struct {
	Label_template string  `json:"label_template"`
	Task_ID        string  `json:"task_ID"`
	Label          []Label `json:"label"`
}

var logger *log.Logger

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)
	// file := "./" + time.Now().Format("20180102150405") + ".log"
	// logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	// if err != nil {
	// 	fmt.Println("日志文件创建失败")
	// }
	// logger = log.New(logFile, "前缀", log.Ldate|log.Ltime|log.Lshortfile)
	logger = log.New(os.Stdout, "前缀", log.Ldate|log.Ltime|log.Lshortfile)
}
func main() {
	filePath := "/Users/herrdu/tmp/label_result.json"

	jsonFile, err := os.Open(filePath)
	if err != nil {
		logger.Fatal(err)
	}
	defer jsonFile.Close()
	var result Result

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = jsoniter.Unmarshal(byteValue, &result)
	if err != nil {
		logger.Fatal(err)
	}

	for index := range result.Label {
		label := result.Label[index]
		err = handelSingleLabel(label)
		if err != nil {
			logger.Println(err)
		}
	}

}

// 处理没一张图片
func handelSingleLabel(label Label) error {
	logger.Printf("%#v", label.Image)
	res, err := http.Get(label.Image)
	if err != nil {
		logger.Println(err)
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		logger.Printf("status code error: %d %s", res.StatusCode, res.Status)
		return errors.New("error in get image file")
	}
	img, _, err := image.DecodeConfig(res.Body)

	if err != nil {
		logger.Println(err)
		return err
	}

	logger.Println(img.Width)

	return nil

}
