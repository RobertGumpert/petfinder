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

func openFile(filePath, ext string) (*os.File, error) {
	filePath = strings.Join([]string{
		filePath,
		ext,
	}, "")
	openFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return openFile, nil
}

func findFile(fileName, filePath, dir string, filesInDir []os.FileInfo) *os.File {
	var file *os.File = nil
	for _, fileInfo := range filesInDir {
		ext := filepath.Ext(fileInfo.Name())
		if fileName == strings.Split(fileInfo.Name(), ext)[0] {
			open, err := openFile(filePath, ext)
			if err == nil {
				file = open
			}
			break
		}
	}
	if file == nil {
		open, err := openFile(dir, "base.jpg")
		if err != nil {
			return nil
		}
		file = open
	}
	return file
}

func asyncFindFile(fileName, filePath, dir string, filesInDir []os.FileInfo) *os.File {
	var file *os.File = nil
	var baseFile *os.File = nil
	var group sync.WaitGroup
	var poolDescriptorHolders = make(chan struct{}, 100)
	var resultSearchFile = make(chan bool)
	var countCallers = 0
	//
	for _, fileInfo := range filesInDir {
		if file != nil {
			break
		}
		countCallers++
		go func(fileInfo os.FileInfo, fileName, filePath string) {
			poolDescriptorHolders <- struct{}{}
			defer func() {
				<-poolDescriptorHolders
			}()
			if file != nil {
				resultSearchFile <- false
				return
			}
			ext := filepath.Ext(fileInfo.Name())
			if fileName == strings.Split(fileInfo.Name(), ext)[0] {
				open, err := openFile(filePath, ext)
				if err == nil {
					file = open
					resultSearchFile <- true
					return
				}
			}
			resultSearchFile <- false
		}(fileInfo, fileName, filePath)
	}
	group.Add(2)
	go func(dir string, group *sync.WaitGroup) {
		defer group.Done()
		open, err := openFile(dir, "base.jpg")
		if err == nil {
			baseFile = open
		}
	}(dir, &group)
	go func() {
		defer group.Done()
		var count = 0
		for flag := range resultSearchFile {
			count++
			if countCallers == count {
				return
			}
			if flag {
				return
			}
		}
	}()
	group.Wait()
	if file == nil {
		file = baseFile
	}
	return file
}

func findFileByID(fileGroup groupFiles, id uint64, async bool) (*os.File, error) {
	var file *os.File = nil
	fileName := createFileName(fileGroup, id)
	filePath, dir := createFilePath(fileName, fileGroup)
	filesInDir, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println(runtimeinfo.Runtime(1), "ERROR=[", err, "]")
		return nil, err
	}
	if async {
		file = asyncFindFile(fileName, filePath, dir, filesInDir)
	} else {
		file = findFile(fileName, filePath, dir, filesInDir)
	}
	return file, nil
}

func findListFilesByID(fileGroup groupFiles, ids []uint64, async bool) map[uint64]*os.File {
	files := make(map[uint64]*os.File)
	var filesInDir []os.FileInfo
	for _, id := range ids {
		fileName := createFileName(fileGroup, id)
		filePath, dir := createFilePath(fileName, fileGroup)
		if filesInDir == nil || len(filesInDir) == 0 {
			list, err := ioutil.ReadDir(dir)
			if err != nil {
				log.Println(runtimeinfo.Runtime(1), "ERROR=[", err, "]")
				break
			}
			filesInDir = list
		}
		if async {
			files[id] = asyncFindFile(fileName, filePath, dir, filesInDir)
		} else {
			files[id] = findFile(fileName, filePath, dir, filesInDir)
		}
	}
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
