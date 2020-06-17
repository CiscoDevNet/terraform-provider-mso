package models

type Model interface {
	ToMap() (map[string]interface{}, error)
}
