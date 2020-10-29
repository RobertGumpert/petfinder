package main

type UploadViewModel struct {
	ID                       uint64 `json:"id"`
	AdditionalIdentification string `json:"additional_identification"`
}

type DownloadViewModel struct {
	ID []uint64 `json:"id"`
}

type File64ListViewModel struct {
	Files []*File64ViewModel `json:"files"`
}

type FileNameListViewModel struct {
	Files []*FileNameViewModel `json:"files"`
}

type File64ViewModel struct {
	ID     uint64 `json:"id"`
	File64 string `json:"file_64"`
}

type FileNameViewModel struct {
	ID       uint64 `json:"id"`
	FileName string `json:"file_name"`
}
