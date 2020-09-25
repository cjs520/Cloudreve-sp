package thumb

import (
	"errors"
	"fmt"
	model "github.com/HFO4/cloudreve/models"
	"github.com/HFO4/cloudreve/pkg/util"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

// Thumb 缩略图
type Thumb struct {
	src image.Image
	ext string
}

// NewThumbFromFile 从文件数据获取新的Thumb对象，
// 尝试通过文件名name解码图像
func NewThumbFromFile(file io.Reader, name string) (*Thumb, error) {
	ext := strings.ToLower(filepath.Ext(name))
	// 无扩展名时
	if len(ext) == 0 {
		return nil, errors.New("未知的图像类型")
	}

	var err error
	var img image.Image
	switch ext[1:] {
	case "jpg":
		img, err = jpeg.Decode(file)
	case "jpeg":
		img, err = jpeg.Decode(file)
	case "gif":
		img, err = gif.Decode(file)
	case "png":
		img, err = png.Decode(file)
	default:
		return nil, errors.New("未知的图像类型")
	}
	if err != nil {
		return nil, err
	}

	return &Thumb{
		src: img,
		ext: ext[1:],
	}, nil
}

// GetThumb 生成给定最大尺寸的缩略图
func (image *Thumb) GetThumb(width, height uint) {
	image.src = resize.Thumbnail(width, height, image.src, resize.Lanczos3)
}

// GetSize 获取图像尺寸
func (image *Thumb) GetSize() (int, int) {
	b := image.src.Bounds()
	return b.Max.X, b.Max.Y
}

// Save 保存图像到给定路径
func (image *Thumb) Save(path string) (err error) {
	out, err := util.CreatNestedFile(path)

	if err != nil {
		return err
	}
	defer out.Close()

	err = png.Encode(out, image.src)
	return err

}

// CreateAvatar 创建头像
func (image *Thumb) CreateAvatar(uid uint) error {
	// 读取头像相关设定
	savePath := util.RelativePath(model.GetSettingByName("avatar_path"))
	s := model.GetIntSetting("avatar_size_s", 50)
	m := model.GetIntSetting("avatar_size_m", 130)
	l := model.GetIntSetting("avatar_size_l", 200)

	// 生成头像缩略图
	src := image.src
	for k, size := range []int{s, m, l} {
		image.src = resize.Resize(uint(size), uint(size), src, resize.Lanczos3)
		err := image.Save(filepath.Join(savePath, fmt.Sprintf("avatar_%d_%d.png", uid, k)))
		if err != nil {
			return err
		}
	}

	return nil

}
