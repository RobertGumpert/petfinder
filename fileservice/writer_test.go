package main

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func BenchmarkWritingFlow(b *testing.B) {
	dat, err := ioutil.ReadFile("C:/PetFinderRepos/petfinder/fileservice/test_jpeg.jpg")
	if err != nil {
		b.Fatal(err)
	}
	// r := bytes.NewReader(dat)
	v := &UploadViewModel{ID: 0}
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(dat)
		err = saveFile(
			r,
			avatarGroupFiles,
			jpegFileType,
			v,
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReadingFlow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := findFileByID(avatarGroupFiles, uint64(0), false)
		if err != nil {
			b.Fatal(err)
		}
	}
}
