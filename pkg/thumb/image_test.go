package thumb

import (
	"fmt"
	"github.com/HFO4/cloudreve/pkg/cache"
	"github.com/HFO4/cloudreve/pkg/util"
	"github.com/stretchr/testify/assert"
	"image"
	"image/jpeg"
	"os"
	"testing"
)

func CreateTestImage() *os.File {
	file, err := os.Create("TestNewThumbFromFile.jpeg")
	alpha := image.NewAlpha(image.Rect(0, 0, 500, 200))
	jpeg.Encode(file, alpha, nil)
	if err != nil {
		fmt.Println(err)
	}
	_, _ = file.Seek(0, 0)
	return file
}

func TestNewThumbFromFile(t *testing.T) {
	asserts := assert.New(t)
	file := CreateTestImage()
	defer file.Close()

	// 无扩展名时
	{
		thumb, err := NewThumbFromFile(file, "123")
		asserts.Error(err)
		asserts.Nil(thumb)
	}

	{
		thumb, err := NewThumbFromFile(file, "123.jpg")
		asserts.NoError(err)
		asserts.NotNil(thumb)
	}
	{
		thumb, err := NewThumbFromFile(file, "123.jpeg")
		asserts.Error(err)
		asserts.Nil(thumb)
	}
	{
		thumb, err := NewThumbFromFile(file, "123.png")
		asserts.Error(err)
		asserts.Nil(thumb)
	}
	{
		thumb, err := NewThumbFromFile(file, "123.gif")
		asserts.Error(err)
		asserts.Nil(thumb)
	}
	{
		thumb, err := NewThumbFromFile(file, "123.3211")
		asserts.Error(err)
		asserts.Nil(thumb)
	}
}

func TestThumb_GetSize(t *testing.T) {
	asserts := assert.New(t)
	file := CreateTestImage()
	defer file.Close()
	thumb, err := NewThumbFromFile(file, "123.jpg")
	asserts.NoError(err)

	w, h := thumb.GetSize()
	asserts.Equal(500, w)
	asserts.Equal(200, h)
}

func TestThumb_GetThumb(t *testing.T) {
	asserts := assert.New(t)
	file := CreateTestImage()
	defer file.Close()
	thumb, err := NewThumbFromFile(file, "123.jpg")
	asserts.NoError(err)

	asserts.NotPanics(func() {
		thumb.GetThumb(10, 10)
	})
}

func TestThumb_Save(t *testing.T) {
	asserts := assert.New(t)
	file := CreateTestImage()
	defer file.Close()
	thumb, err := NewThumbFromFile(file, "123.jpg")
	asserts.NoError(err)

	err = thumb.Save("/:noteexist/")
	asserts.Error(err)

	err = thumb.Save("TestThumb_Save.png")
	asserts.NoError(err)
	asserts.True(util.Exists("TestThumb_Save.png"))

}

func TestThumb_CreateAvatar(t *testing.T) {
	asserts := assert.New(t)
	file := CreateTestImage()
	defer file.Close()

	thumb, err := NewThumbFromFile(file, "123.jpg")
	asserts.NoError(err)

	cache.Set("setting_avatar_path", "tests", 0)
	cache.Set("setting_avatar_size_s", "50", 0)
	cache.Set("setting_avatar_size_m", "130", 0)
	cache.Set("setting_avatar_size_l", "200", 0)

	asserts.NoError(thumb.CreateAvatar(1))
	asserts.True(util.Exists(util.RelativePath("tests/avatar_1_1.png")))
	asserts.True(util.Exists(util.RelativePath("tests/avatar_1_2.png")))
	asserts.True(util.Exists(util.RelativePath("tests/avatar_1_0.png")))
}
