package model

type Entry struct {
	Word    string
	Content string
}

type Resource struct {
	Audio []Audio `json:"audios"`
}

type Audio struct {
	Path     string `json:"path"`
	Language string `json:"language"`
}
