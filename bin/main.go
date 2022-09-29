package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/sirupsen/logrus"
)

var flagdir string
var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	flag.StringVar(&flagdir, "dir", "", "指定目录")
	flag.Parse()
	logrus.SetLevel(logrus.DebugLevel)
	rand.Seed(time.Now().UnixNano())
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func main() {
	logrus.Println("程序已启动", os.Args, flagdir)
	//var files []string
	root := strings.TrimRight(flagdir, "/")
	if strings.IndexAny(root, "~/") == 0 {
		homePath, err := os.UserHomeDir()
		if err != nil {
			logrus.Fatalln("homePath error: ", err)
		}
		root = homePath + strings.TrimLeft(root, "~")
	}
	logrus.Println("jjjj", root)

	loopDir(root)
}

func loopDir(rootPath string) {
	files, err := ioutil.ReadDir(rootPath)
	if err != nil {
		logrus.Fatalln("ReadDir error:", err)
	}
	for _, f := range files {
		fileExt := strings.ToLower(filepath.Ext(f.Name()))
		fullPath := rootPath + "/" + f.Name()
		logrus.Debugln("====: ", fullPath, f.Name(), fileExt)
		//if f.IsDir() {
		//	loopDir(fullPath)
		//	continue
		//}
		if fileExt == ".png" || fileExt == ".jpg" || fileExt == ".jpeg" ||
			fileExt == ".gif" {
			logrus.Debugln("++++: ", f.Name(), fileExt)
			ocrText := runTesseract(fullPath)
			logrus.Debugln("----: ", ocrText)
			if len(ocrText) > 0 {
				newPath := rootPath + "/emo-" + RandStringRunes(8) + "-" + ocrText + fileExt
				logrus.Debugln("move: ", newPath)
				err := os.Rename(fullPath, newPath)
				if err != nil {
					logrus.Errorln("Rename error: ", err, newPath)
				}
			}
		}

	}
}

func runTesseract(fullPath string) string {
	cmd := fmt.Sprintf("tesseract %s stdout -l chi_sim", fullPath)
	logrus.Debugln("cmd: ", cmd)
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		logrus.Errorln("执行命令出错: %s", cmd)
		return ""
	}
	return stripSpaces(string(out))
}

func stripSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) || (!unicode.Is(unicode.Han, r) && !unicode.IsLetter(r) &&
			!unicode.IsDigit(r)) {
			return -1
		}
		return r
	}, str)
}
