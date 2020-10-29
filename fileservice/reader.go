package main

import (
	"bufio"
	"encoding/base64"
	"fileservice/pckg/runtimeinfo"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

func readFiles(fileGroup groupFiles, viewModel *DownloadViewModel) map[uint64]*os.File {
	files := make(map[uint64]*os.File)
	//
	var wg sync.WaitGroup
	var mx = new(sync.Mutex)
	//
	for _, id := range viewModel.ID {
		wg.Add(1)
		go func(id uint64, wg *sync.WaitGroup) {
			mx.Lock()
			defer func() {
				wg.Done()
				mx.Unlock()
			}()
			fileName := createFileName(fileGroup, id)
			fileName = addFilenameExtension(fileName, jpegFileType)
			filePath, dir := createFilePath(fileName, fileGroup)
			file, err := os.Open(filePath)
			if err != nil {
				log.Println(runtimeinfo.Runtime(1), "ERROR=[", err, "]")
				filePath = strings.Join([]string{
					dir,
					"base.jpg",
				}, "")
				baseFile, err := os.Open(filePath)
				if err != nil {
					log.Println(runtimeinfo.Runtime(1), "ERROR=[", err, "]")
				} else {
					files[id] = baseFile
				}
				return
			}
			files[id] = file
			return
		}(id, &wg)
	}
	wg.Wait()
	return files
}

func filesToBase64(files map[uint64]*os.File) *File64ListViewModel {
	fileListViewModel := &File64ListViewModel{Files: make([]*File64ViewModel, 0)}
	for id, file := range files {
		fileViewModel := new(File64ViewModel)
		fileViewModel.ID = id
		reader := bufio.NewReader(file)
		content, err := ioutil.ReadAll(reader)
		if err != nil {
			continue
		}
		fileViewModel.File64 = base64.StdEncoding.EncodeToString(content)
		fileListViewModel.Files = append(
			fileListViewModel.Files,
			fileViewModel,
		)
		_ = file.Close()
	}
	return fileListViewModel
}
