package main

import (
	"errors"
	"image"
	_ "image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	jsoniter "github.com/json-iterator/go"
)

// 1. 判断图片缩放比例
// 2. 调整数据
// 3. 保存为新的文件

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
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

const (
	clientHeight int = 636
	clientWidth  int = 888
)

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
	filePath := "./label_result.json"

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
		err = handelSingleLabel(&label)
		if err != nil {
			logger.Println(err)
		}
	}

	newFile, err := os.OpenFile("new_result.json", os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)

	if err != nil {
		logger.Fatal(err)
	}
	defer newFile.Close()
	enc := jsoniter.NewEncoder(newFile)
	if err = enc.Encode(&result); err != nil {
		logger.Println(err)
	}
}

// 处理没一张图片
func handelSingleLabel(label *Label) error {
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
	imgWidth := img.Width
	imgHeight := img.Height
	var (
		scaleX float32
		scaleY float32
	)

	logger.Printf("imgWidth:%d,clientWidth:%d,imgHeight:%d,clientHeight:%d", imgWidth, clientWidth, imgHeight, clientHeight)

	if imgWidth > clientWidth && imgHeight <= clientHeight {
		scaleX = float32(clientWidth) / float32(imgWidth)
		scaleY = scaleX
	} else if imgHeight > clientHeight && imgWidth <= clientWidth {
		scaleY = float32(clientHeight) / float32(imgHeight)
		scaleX = scaleY
	} else if imgHeight > clientHeight && imgWidth > clientWidth {
		scaleX = float32(clientWidth) / float32(imgWidth)
		scaleY = float32(clientHeight) / float32(imgHeight)
		if scaleX < scaleY {
			scaleX = scaleY
			// this.canvasElement.width = clientWidth
			// this.canvasElement.height = (imgHeight * clientWidth) / imgWidth
		} else {
			scaleY = scaleX
			// this.canvasElement.height = clientHeight
			// this.canvasElement.width = (imgWidth * clientHeight) / imgHeight
		}
	}

	if scaleX != 0 && scaleX != 0 {
		logger.Printf("scaleX:%d,scaleY:%d", scaleX, scaleY)
		for cmpIndex := range label.Component {
			component := label.Component[cmpIndex]
			for pointIndex := range component.PointList {
				point := &component.PointList[pointIndex]
				x := int(float32(point.X) / scaleX)
				point.X = x
				logger.Printf("x:%d,pointX:%d", x, point.X)
				point.Y = int(float32(point.Y) / scaleY)
			}
			logger.Printf("component:%v", component)

		}
	}

	return nil

}
