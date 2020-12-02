package main

import (
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	mx sync.Mutex
)

func createFilePath(fileName string, fileGroup groupFiles) (string, string) {
	if strings.TrimSpace(root) == "" {
		setRoot()
	}
	dir := ""
	switch fileGroup {
	case avatarGroupFiles:
		dir = strings.Join([]string{
			root,
			"/storage/avatar/",
		}, "")
		fileName = strings.Join([]string{
			dir,
			fileName,
		}, "")
		break
	case advertGroupFiles:
		dir = strings.Join([]string{
			root,
			"/storage/advert/",
		}, "")
		fileName = strings.Join([]string{
			dir,
			fileName,
		}, "")
		break
	}
	return fileName, dir
}

func addFilenameExtension(fileName string, fileType typeFiles) string {
	switch fileType {
	case jpegFileType:
		fileName = strings.Join([]string{
			fileName,
			"jpeg",
		}, ".")
		break
	case pngFileType:
		fileName = strings.Join([]string{
			fileName,
			"png",
		}, ".")
		break
	}
	return fileName
}

func createFileName(fileGroup groupFiles, id uint64) string {
	fileName := ""
	switch fileGroup {
	case avatarGroupFiles:
		fileName = strings.Join([]string{
			"avatar",
			"id",
			strconv.Itoa(int(id)),
		}, "_")
		break
	case advertGroupFiles:
		fileName = strings.Join([]string{
			"advert",
			"id",
			strconv.Itoa(int(id)),
		}, "_")
		break
	}
	return fileName
}

func saveFile(reader io.Reader, fileGroup groupFiles, fileType typeFiles, viewModel *UploadViewModel) error {
	fileName := createFileName(fileGroup, viewModel.ID)
	fileName = addFilenameExtension(fileName, fileType)
	filePath, _ := createFilePath(fileName, fileGroup)
	mx.Lock()
	defer mx.Unlock()
	if fileExists(filePath) {
		if err := os.Remove(filePath); err != nil {
			return err
		}
	}
	switch fileType {
	case jpegFileType:
		if err := jpegSaveAndCompress(filePath, reader); err != nil {
			return err
		}
		break
	case pngFileType:
		if err := pngSaveAndCompress(filePath, reader); err != nil {
			return err
		}
		break
	}
	return nil
}

func pngSaveAndCompress(filePath string, reader io.Reader) error {
	compress, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer compress.Close()
	img, err := png.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	encoder := png.Encoder{CompressionLevel: png.BestCompression}
	if err := encoder.Encode(compress, img); err != nil {
		return err
	}
	return nil
}

func jpegSaveAndCompress(filePath string, reader io.Reader) error {
	compress, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer compress.Close()
	img, err := jpeg.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	if err := jpeg.Encode(compress, img, &jpeg.Options{Quality: 50}); err != nil {
		return err
	}
	return nil
}

func fileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
