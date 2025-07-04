package myglob

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
测试 glob 模式，不匹配隐藏文件

匹配所有 html 文件，不匹配隐藏文件
*/
func TestMatchGlobWithoutDot(t *testing.T) {

	os.MkdirAll("dist/pages", 0755)

	// 这个文件会被匹配到
	os.Create("dist/index.html")

	// 这些文件会被匹配到
	os.Create("dist/pages/page1.html")
	os.Create("dist/pages/page2.html")
	os.Create("dist/pages/page3.html")

	// 这些文件不会被匹配到
	os.MkdirAll("dist/.pages", 0755)
	os.Create("dist/.pages/page4.html")
	os.Create("dist/.pages/page5.html")

	matches, err := Match(os.DirFS("."), "**/*.html", Options{
		Dot:    false,
		NoCase: true,
	})

	for _, match := range matches {
		log.Println(match)
	}

	assert.NoError(t, err)
	assert.Equal(t, len(matches), 4)
	assert.Equal(t, matches[0], "dist/index.html")
	assert.Equal(t, matches[1], "dist/pages/page1.html")
	assert.Equal(t, matches[2], "dist/pages/page2.html")
	assert.Equal(t, matches[3], "dist/pages/page3.html")

	os.RemoveAll("dist")
}

/*
测试 glob 模式，匹配隐藏文件

匹配所有 html 文件，匹配隐藏文件
*/
func TestMatchGlobWithDot(t *testing.T) {
	os.MkdirAll("dist/pages", 0755)
	os.MkdirAll("dist/.pages", 0755)

	// 这个文件会被匹配到
	os.Create("dist/index.html")

	// 这些文件会被匹配到
	os.Create("dist/pages/page1.html")
	os.Create("dist/pages/page2.html")
	os.Create("dist/pages/page3.html")

	// 这些文件会被匹配到
	os.Create("dist/.pages/page4.html")
	os.Create("dist/.pages/page5.html")

	matches, err := Match(os.DirFS("."), "**/*.html", Options{
		Dot:    true,
		NoCase: true,
	})
	assert.NoError(t, err)
	assert.Equal(t, len(matches), 6)

	for _, match := range matches {
		log.Println(match)
	}

	assert.Equal(t, matches[0], "dist/.pages/page4.html")
	assert.Equal(t, matches[1], "dist/.pages/page5.html")
	assert.Equal(t, matches[2], "dist/index.html")
	assert.Equal(t, matches[3], "dist/pages/page1.html")
	assert.Equal(t, matches[4], "dist/pages/page2.html")
	assert.Equal(t, matches[5], "dist/pages/page3.html")

	os.RemoveAll("dist")
}

/*
测试 no glob 模式，不匹配隐藏文件

匹配所有 dist 目录下的文件，不匹配隐藏文件
*/
func TestMatchNoGlobWithoutDot(t *testing.T) {
	os.MkdirAll("dist", 0755)
	os.MkdirAll("dist/pages", 0755)
	os.MkdirAll("dist/.pages", 0755)

	// 这些文件会被匹配到
	os.Create("dist/index.html")
	os.Create("dist/pages/page1.html")
	os.Create("dist/pages/page2.html")
	os.Create("dist/pages/page3.html")

	// 这些文件不会被匹配到
	os.Create("dist/.pages/page4.html")
	os.Create("dist/.pages/page5.html")

	matches, err := Match(os.DirFS("."), "dist/**/*", Options{
		Dot:    false,
		NoCase: true,
	})
	for _, match := range matches {
		log.Println(match)
	}

	assert.NoError(t, err)
	assert.Equal(t, len(matches), 5)
	assert.Equal(t, matches[0], "dist/index.html")
	assert.Equal(t, matches[1], "dist/pages")
	assert.Equal(t, matches[2], "dist/pages/page1.html")
	assert.Equal(t, matches[3], "dist/pages/page2.html")
	assert.Equal(t, matches[4], "dist/pages/page3.html")

	os.RemoveAll("dist")
}
