package configx

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
	"sync"

	json "github.com/pydio/cells/v4/common/utils/jsonx"
)

type Unmarshaler interface {
	Unmarshal([]byte, interface{}) error
}

type Marshaller interface {
	Marshal(interface{}) ([]byte, error)
}

type Encrypter interface {
	Encrypt([]byte) (string, error)
}

type Decrypter interface {
	Decrypt(string) ([]byte, error)
}

type Option func(*Options)

type Options struct {
	*sync.RWMutex
	Unmarshaler
	Marshaller
	Encrypter
	Decrypter

	AutoUpdate bool
}

type jsonReader struct{}

func (j *jsonReader) Unmarshal(data []byte, out interface{}) error {
	var err error

	switch v := out.(type) {
	case proto.Message:
		err = protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(data, v)
	default:
		err = json.Unmarshal(data, v)
	}

	return err
}

type jsonWriter struct{}

func (j *jsonWriter) Marshal(in interface{}) ([]byte, error) {
	return json.Marshal(in)
}

func WithJSON() Option {
	return func(o *Options) {
		o.Unmarshaler = &jsonReader{}
		o.Marshaller = &jsonWriter{}
	}
}

type yamlReader struct{}

func (j *yamlReader) Unmarshal(data []byte, out interface{}) error {
	return yaml.Unmarshal(data, out)
}

type yamlWriter struct{}

func (j *yamlWriter) Marshal(in interface{}) ([]byte, error) {
	return yaml.Marshal(in)
}

func WithYAML() Option {
	return func(o *Options) {
		o.Unmarshaler = &yamlReader{}
		o.Marshaller = &yamlWriter{}
	}
}

func WithEncrypt(e Encrypter) Option {
	return func(o *Options) {
		o.Encrypter = e
	}
}

func WithDecrypt(d Decrypter) Option {
	return func(o *Options) {
		o.Decrypter = d
	}
}

func WithLock(m *sync.RWMutex) Option {
	return func(o *Options) {
		o.RWMutex = m
	}
}
