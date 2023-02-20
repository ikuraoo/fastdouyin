package util

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/ikuraoo/fastdouyin/entity"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"log"
	"os"
)

func SnapShotFromVideo(videoPath, snapShotPath string, frameNum int) (err error) {
	buf := bytes.NewBuffer(nil)
	err = ffmpeg.Input(videoPath).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return err
	}

	img, err := imaging.Decode(buf)
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return err
	}

	err = imaging.Save(img, snapShotPath)
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return err
	}
	return nil
}

func NewFileCoverName(userId int64) string {
	var count int64

	err := entity.NewVideoDaoInstance().QueryVideoCountByUserId(userId, &count)
	if err != nil {
		log.Println(err)
	}
	return fmt.Sprintf("%d-%d.jpg", userId, count)
}
