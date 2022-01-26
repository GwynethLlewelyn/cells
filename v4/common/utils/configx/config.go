package configx

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cast"

	json "github.com/pydio/cells/v4/common/utils/jsonx"
)

var (
	ErrNoMarshallerDefined  = errors.New("no marshaller defined")
	ErrNoUnmarshalerDefined = errors.New("no unmarshaler defined")
)

type Scanner interface {
	Scan(interface{}) error
}

type Watcher interface {
	Watch(path ...string) (Receiver, error)
}

type Receiver interface {
	Next() (Values, error)
	Stop()
}

// TODO - we should be returning a Value
type KV struct {
	Key   string
	Value interface{}
}

type Key interface{}

type Value interface {
	Default(interface{}) Value

	Bool() bool
	Bytes() []byte
	Int() int
	Int64() int64
	Duration() time.Duration
	String() string
	StringMap() map[string]string
	StringArray() []string
	Slice() []interface{}
	Map() map[string]interface{}

	Scanner
}

type KVStore interface {
	Get() Value
	Set(value interface{}) error
	Del() error
}

type Entrypoint interface {
	KVStore
	Val(path ...string) Values
}

type Values interface {
	Entrypoint
	Value
}

type Ref interface {
	Get() string
}

type Source interface {
	Entrypoint
	Watcher
}

func NewFrom(i interface{}) Values {
	c := New()
	c.Set(i)
	return c
}

// config is standard
type config struct {
	v    interface{}
	d    interface{} // Default
	r    *config     // Root
	k    []string    // Reference to key for re-assignment
	opts Options
}

func New(opts ...Option) Values {
	options := Options{}

	for _, o := range opts {
		o(&options)
	}

	return &config{
		opts: options,
	}
}

func (c *config) get() interface{} {
	if c == nil {
		return nil
	}

	if c.v != nil {
		useDefault := false

		switch vv := c.v.(type) {
		case map[interface{}]interface{}:
			if ref, ok := vv["$ref"]; ok {
				vvv := c.r.Val(ref.(string)).Get()
				switch vvvv := vvv.(type) {
				case *config:
					return vvvv.get()
				default:
					return vvvv
				}
			}
		case map[string]interface{}:
			if ref, ok := vv["$ref"]; ok {
				vvv := c.r.Val(ref.(string)).Get()
				switch vvvv := vvv.(type) {
				case *config:
					return vvvv.get()
				default:
					return vvvv
				}
			}
		case string:
			if vv == "default" {
				useDefault = true
			}
		}

		if !useDefault {
			str, ok := c.v.(string)
			if ok {
				if d := c.opts.Decrypter; d != nil {
					b, err := d.Decrypt(str)
					if err != nil {
						return c.v
					}
					return string(b)
				}
			}
			return c.v
		}
	}

	if c.d != nil {
		switch vv := c.d.(type) {
		case map[string]interface{}:
			if ref, ok := vv["$ref"]; ok {
				vvv := c.r.Val(ref.(string)).Get()
				switch vvvv := vvv.(type) {
				case *config:
					return vvvv.get()
				default:
					return vvvv
				}
			}
		case *ref:
			vvv := c.r.Val(vv.Get()).Get()
			switch vvvv := vvv.(type) {
			case *config:
				return vvvv.get()
			default:
				return vvvv
			}
		}
		return c.d
	}

	return nil
}

// Get retrieve interface
func (c *config) Get() Value {
	if mtx := c.opts.RWMutex; mtx != nil {
		mtx.RLock()
		defer mtx.RUnlock()
	}

	if c.v == nil && c.d == nil {
		return nil
	}

	switch vv := c.v.(type) {
	case map[string]interface{}:
		if ref, ok := vv["$ref"]; ok {
			return c.r.Val(ref.(string)).Get()
		}
	case *ref:
		return c.r.Val(vv.Get()).Get()
	}

	return c
}

// Default value set
func (c *config) Default(i interface{}) Value {
	if c.d == nil {
		c.d = i
	}

	switch vv := c.v.(type) {
	case string:
		if vv == "default" {
			c.v = nil
		}
	}

	return c.Get()
}

// Set data in interface
func (c *config) Set(data interface{}) error {

	if c == nil {
		return fmt.Errorf("value doesn't exist")
	}

	if c.opts.Unmarshaler != nil {
		switch vv := data.(type) {
		case []byte:
			if len(vv) > 0 {
				if err := c.opts.Unmarshaler.Unmarshal(vv, &data); err != nil {
					return err
				}
			}
		}
	}

	switch d := data.(type) {
	case *config:
		data = d.v
	}

	if len(c.k) == 0 {
		c.v = data
		return nil
	}

	k := c.k[len(c.k)-1]
	pk := c.k[0 : len(c.k)-1]

	// Retrieve parent value
	p := c.r.Val(pk...)
	m := p.Map()
	if data == nil {
		delete(m, k)
	} else {
		if enc := c.opts.Encrypter; enc != nil {
			switch vv := data.(type) {
			case []byte:
				// Encrypting value
				str, err := enc.Encrypt(vv)
				if err != nil {
					return err
				}

				data = str
			case string:
				// Encrypting value
				str, err := enc.Encrypt([]byte(vv))
				if err != nil {
					return err
				}

				data = str
			}
		}

		m[k] = data
	}

	p.Set(m)

	if mtx := c.opts.RWMutex; mtx != nil {
		mtx.Lock()
		defer mtx.Unlock()
	}

	c.v = data

	return nil
}

func (c *config) Del() error {
	if c == nil {
		return fmt.Errorf("value doesn't exist")
	}

	return c.Set(nil)
}

// Val values cannot retrieve lower values as it is final
func (c *config) Val(s ...string) Values {
	if mtx := c.opts.RWMutex; mtx != nil {
		mtx.RLock()
		defer mtx.RUnlock()
	}

	keys := StringToKeys(s...)

	// Need to do something for reference
	if len(keys) == 1 && keys[0] == "#" {
		if c.r != nil {
			return c.r
		}
		return c
	} else if len(keys) > 0 && keys[0] == "#" {
		keys = keys[1:]
	} else {
		keys = append(c.k, keys...)
	}

	root := c.r
	if c.r == nil {
		root = c
	}

	if len(keys) == 0 {
		return c
	}

	pk := keys

	// Looking for the specific key
	var current interface{} = root.Map()

	for _, pkk := range pk {
		switch cv := current.(type) {
		case map[interface{}]interface{}:
			cvv, ok := cv[pkk]
			if !ok {
				// The parent doesn't actually exist here, we return the nil value
				return &config{nil, nil, root, keys, c.opts}
			}

			current = cvv
		case map[string]string:
			cvv, ok := cv[pkk]
			if !ok {
				// The parent doesn't actually exist here, we return the nil value
				return &config{nil, nil, root, keys, c.opts}
			}

			current = cvv
		case map[string]interface{}:
			cvv, ok := cv[pkk]
			if !ok {
				// The parent doesn't actually exist here, we return the nil value
				return &config{nil, nil, root, keys, c.opts}
			}

			current = cvv
		case []interface{}:
			i, err := strconv.Atoi(pkk)
			if err != nil || i < 0 || i >= len(cv) {
				return &config{nil, nil, root, keys, c.opts}
			}

			cvv := cv[i]

			current = cvv
		default:
			return &config{nil, nil, root, keys, c.opts}
		}
	}

	return &config{current, nil, root, keys, c.opts}
}

// Scan to interface
func (c *config) Scan(val interface{}) error {
	if mtx := c.opts.RWMutex; mtx != nil {
		mtx.RLock()
		defer mtx.RUnlock()
	}

	v := c.get()
	if v == nil {
		return nil
	}

	marshaller := c.opts.Marshaller
	if marshaller == nil {
		return ErrNoMarshallerDefined
	}

	str, err := marshaller.Marshal(v)
	if err != nil {
		return err
	}

	unmarshaler := c.opts.Unmarshaler
	if unmarshaler == nil {
		return ErrNoUnmarshalerDefined
	}

	return unmarshaler.Unmarshal(str, val)
}

func (c *config) Bool() bool {
	if mtx := c.opts.RWMutex; mtx != nil {
		mtx.RLock()
		defer mtx.RUnlock()
	}

	v := c.get()
	if v == nil {
		return false
	}
	return cast.ToBool(v)
}
func (c *config) Bytes() []byte {
	if mtx := c.opts.RWMutex; mtx != nil {
		mtx.RLock()
		defer mtx.RUnlock()
	}

	v := c.get()
	if v == nil {
		return []byte{}
	}
	switch v := c.v.(type) {
	case []interface{}, map[string]interface{}:
		if m := c.opts.Marshaller; m != nil {
			data, err := m.Marshal(v)
			if err != nil {
				return []byte{}
			}

			return data
		}

		return []byte{}
	case string:
		// Need to handle it differently
		if v == "default" {
			c.v = nil
		}
	}
	return []byte(cast.ToString(v))
}
func (c *config) Int() int {
	if mtx := c.opts.RWMutex; mtx != nil {
		mtx.RLock()
		defer mtx.RUnlock()
	}

	v := c.get()
	if v == nil {
		return 0
	}
	return cast.ToInt(v)
}
func (c *config) Int64() int64 {
	if mtx := c.opts.RWMutex; mtx != nil {
		mtx.RLock()
		defer mtx.RUnlock()
	}

	v := c.get()
	if v == nil {
		return 0
	}
	return cast.ToInt64(v)
}
func (c *config) Duration() time.Duration {
	if mtx := c.opts.RWMutex; mtx != nil {
		mtx.RLock()
		defer mtx.RUnlock()
	}

	v := c.get()
	if v == nil {
		return 0 * time.Second
	}
	return cast.ToDuration(v)
}
func (c *config) String() string {
	if mtx := c.opts.RWMutex; mtx != nil {
		mtx.RLock()
		defer mtx.RUnlock()
	}

	v := c.get()

	switch v := c.v.(type) {
	case []interface{}, map[string]interface{}:
		if m := c.opts.Marshaller; m != nil {
			data, err := m.Marshal(v)
			if err != nil {
				return ""
			}

			return string(data)
		}

		return ""
	case string:
		// Need to handle it differently
		if v == "default" {
			c.v = nil
		}
	}

	return cast.ToString(v)
}
func (c *config) StringMap() map[string]string {
	if mtx := c.opts.RWMutex; mtx != nil {
		mtx.RLock()
		defer mtx.RUnlock()
	}

	v := c.get()
	if v == nil {
		return map[string]string{}
	}
	return cast.ToStringMapString(v)
}
func (c *config) StringArray() []string {
	if mtx := c.opts.RWMutex; mtx != nil {
		mtx.RLock()
		defer mtx.RUnlock()
	}

	v := c.get()
	if v == nil {
		return []string{}
	}
	return cast.ToStringSlice(c.get())
}
func (c *config) Slice() []interface{} {
	if mtx := c.opts.RWMutex; mtx != nil {
		mtx.RLock()
		defer mtx.RUnlock()
	}

	v := c.get()
	if v == nil {
		return []interface{}{}
	}
	return cast.ToSlice(c.get())
}
func (c *config) Map() map[string]interface{} {
	if mtx := c.opts.RWMutex; mtx != nil {
		mtx.RLock()
		defer mtx.RUnlock()
	}

	v := c.get()
	if v == nil {
		return map[string]interface{}{}
	}
	r, _ := cast.ToStringMapE(v)
	return r
}
func (c *config) UnmarshalJSON(data []byte) error {
	var m map[string]interface{}

	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	c.v = m

	return nil
}

func (c *config) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.v)
}
