package models

type Model interface {
	ToMap() (map[string]string, error)
}
