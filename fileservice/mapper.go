package main

type UploadViewModel struct {
	ID                       uint64 `json:"id"`
	AdditionalIdentification string `json:"additional_identification"`
}

type URLViewModel struct {
	URL string `json:"url"`
}

type DownloadViewModel struct {
	ID []uint64 `json:"id"`
}

type ListFile64ViewModel struct {
	Files []*File64ViewModel `json:"files"`
}

type File64ViewModel struct {
	ID     uint64 `json:"id"`
	File64 string `json:"file_64"`
}
