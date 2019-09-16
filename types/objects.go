package types

import "github.com/gregdhill/go-openrpc/util"

type BasicType struct {
	Desc string
	Name string
	Type string
}

type FieldMap struct {
	fields map[string]BasicType
	keys   []string
}

func NewFieldMap() *FieldMap {
	return &FieldMap{
		fields: make(map[string]BasicType, 0),
		keys:   make([]string, 0),
	}
}

func (fm *FieldMap) Set(key string, value BasicType) {
	key = util.CamelCase(key)
	value.Name = util.CamelCase(value.Name)
	_, exists := fm.fields[key]
	fm.fields[key] = value
	if !exists {
		fm.keys = append(fm.keys, key)
	}
}

func (fm *FieldMap) Get(key string) BasicType {
	return fm.fields[key]
}

func (fm *FieldMap) GetKeys() []string {
	return fm.keys
}

type ObjectMap struct {
	objects map[string]*FieldMap
	keys    []string
}

func NewObjectMap() *ObjectMap {
	return &ObjectMap{
		objects: make(map[string]*FieldMap, 0),
		keys:    make([]string, 0),
	}
}

func (om *ObjectMap) Set(key string, value BasicType) {
	if key == "" {
		return
	} else if value.Name == "" {
		return
	} else if util.CamelCase(key) == value.Type {
		return
	}
	key = util.CamelCase(key)
	_, exists := om.objects[key]
	if !exists {
		if om.objects[key] == nil {
			om.objects[key] = NewFieldMap()
		}
		om.keys = append(om.keys, key)
	}
	om.objects[key].Set(value.Name, value)
}

func (om *ObjectMap) Get(key string) *FieldMap {
	return om.objects[key]
}

func (om *ObjectMap) GetKeys() []string {
	return om.keys
}
