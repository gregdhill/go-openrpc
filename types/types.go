package types

import (
	"encoding/json"

	"github.com/go-openapi/spec"
)

type Contact struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Email string `json:"email"`
}

type License struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Info struct {
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	TermsOfService string  `json:"termsOfService"`
	Contact        Contact `json:"contact"`
	License        License `json:"license"`
	Version        string  `json:"version"`
}

type ServerVariable struct {
	Enum        []string `json:"enum"`
	Default     string   `json:"default"`
	Description string   `json:"description"`
}

type Server struct {
	Name        string                    `json:"name"`
	URL         string                    `json:"url"`
	Summary     string                    `json:"summary"`
	Description string                    `json:"description"`
	Variables   map[string]ServerVariable `json:"variables`
}

type ExternalDocs struct {
	Description string `json:"description"`
	URL         string `json:"url"`
}

type Tag struct {
	Name         string       `json:"name"`
	Summary      string       `json:"summary"`
	Description  string       `json:"description"`
	ExternalDocs ExternalDocs `json:"externalDocs"`
}

type Content struct {
	Name        string      `json:"name"`
	Summary     string      `json:"summary"`
	Description string      `json:"description"`
	Required    bool        `json:"required"`
	Deprecated  bool        `json:"deprecated"`
	Schema      spec.Schema `json:"schema"`
}

type ContentDescriptor struct {
	Content
}

func (cd *ContentDescriptor) UnmarshalJSON(data []byte) error {
	cont := new(Content)
	err := json.Unmarshal(data, cont)
	if err != nil {
		return err
	}
	cd.Content = *cont

	params := make(map[string]interface{})
	err = json.Unmarshal(data, &params)
	if err != nil {
		return err
	}

	if _, ok := params["$ref"]; ok {
		sch := new(spec.Schema)
		err = json.Unmarshal(data, sch)
		if err != nil {
			return err
		}
		cd.Schema = *sch
	}

	return nil

}

// https://www.jsonrpc.org/specification#error_object
type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Link struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Summary     string                 `json:"summary"`
	Method      string                 `json:"method"`
	Params      map[string]interface{} `json:"params"`
	Server      Server                 `json:"server"`
}

type Example struct {
	Name          string      `json:"name"`
	Summary       string      `json:"summary"`
	Description   string      `json:"description"`
	Value         interface{} `json:"value"`
	ExternalValue string      `json:"externalValue"`
}

type ExamplePairing struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Summary     string    `json:"summary"`
	Params      []Example `json:"params"`
	Result      Example   `json:"result"`
}

type Method struct {
	Name           string               `json:"name"`
	Tags           []Tag                `json:"tags"`
	Summary        string               `json:"summary"`
	Description    string               `json:"description"`
	ExternalDocs   ExternalDocs         `json:"externalDocs"`
	Params         []*ContentDescriptor `json:"params"`
	Result         *ContentDescriptor   `json:"result"`
	Deprecated     bool                 `json:"deprecated"`
	Servers        []Server             `json:"servers"`
	Errors         []Error              `json:"errors"`
	Links          []Link               `json:"links"`
	ParamStructure string               `json:"paramStructure"`
	Examples       []ExamplePairing     `json:"examples"`
}

type Components struct {
	ContentDescriptors    map[string]*ContentDescriptor `json:"contentDescriptors"`
	Schemas               map[string]spec.Schema        `json:"schemas"`
	Examples              map[string]Example            `json:"examples"`
	Links                 map[string]Link               `json:"links"`
	Errors                map[string]Error              `json:"errors"`
	ExamplePairingObjects map[string]ExamplePairing     `json:"examplePairingObjects"`
	Tags                  map[string]Tag                `json:"tags"`
}

type OpenRPCSpec1 struct {
	OpenRPC      string       `json:"openrpc"`
	Info         Info         `json:"info"`
	Servers      []Server     `json:"servers"`
	Methods      []Method     `json:"methods"`
	Components   Components   `json:"components"`
	ExternalDocs ExternalDocs `json:"externalDocs"`

	Objects *ObjectMap `json:"-"`
}

func NewOpenRPCSpec1() *OpenRPCSpec1 {
	return &OpenRPCSpec1{
		Servers: make([]Server, 0),
		Methods: make([]Method, 0),

		Objects: NewObjectMap(),
	}
}
