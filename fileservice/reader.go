package main

import (
	"bufio"
	"encoding/base64"
	"fileservice/pckg/runtimeinfo"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func findFileByID(fileGroup groupFiles, id uint64) (*os.File, error) {
	var file *os.File = nil
	var group sync.WaitGroup
	var pool = make(chan struct{}, 12)
	//
	fileName := createFileName(fileGroup, id)
	filePath, dir := createFilePath(fileName, fileGroup)
	filesInDir, err := ioutil.ReadDir(dir)
	//
	if err != nil {
		log.Println(runtimeinfo.Runtime(1), "ERROR=[", err, "]")
		return nil, err
	}
	for _, fileInfo := range filesInDir {
		group.Add(1)
		go func(fileInfo os.FileInfo, fileName, filePath string, group *sync.WaitGroup) {
			//
			if file != nil {
				return
			}
			// lock
			pool <- struct{}{}
			defer func() {
				<-pool
				group.Done()
			}()
			//
			ext := filepath.Ext(fileInfo.Name())
			if fileName == strings.Split(fileInfo.Name(), ext)[0] {
				filePath = strings.Join([]string{
					filePath,
					ext,
				}, "")
				openFile, err := os.Open(filePath)
				if err != nil {
					return
				}
				file = openFile
			}
			return
		}(fileInfo, fileName, filePath, &group)
	}
	group.Wait()
	close(pool)
	return file, nil
}

func readFilesByID(fileGroup groupFiles, viewModel *DownloadViewModel) map[uint64]*os.File {
	files := make(map[uint64]*os.File)
	//
	var wg sync.WaitGroup
	var mx = new(sync.Mutex)
	var filesInDir []os.FileInfo
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
			filePath, dir := createFilePath(fileName, fileGroup)
			if filesInDir == nil || len(filesInDir) == 0 {
				list, err := ioutil.ReadDir(dir)
				if err != nil {
					log.Println(runtimeinfo.Runtime(1), "ERROR=[", err, "]")
					return
				}
				filesInDir = list
			}
			for _, f := range filesInDir {
				ext := filepath.Ext(f.Name())
				if fileName == strings.Split(f.Name(), ext)[0] {
					filePath = strings.Join([]string{
						filePath,
						ext,
					}, "")
					break
				}
			}
			file, err := os.Open(filePath)
			if err != nil {
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
				_ = file.Close()
				return
			}
			files[id] = file
			return
		}(id, &wg)
	}
	wg.Wait()
	return files
}

func filesToBase64(files map[uint64]*os.File) *ListFile64ViewModel {
	fileListViewModel := &ListFile64ViewModel{Files: make([]*File64ViewModel, 0)}
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
