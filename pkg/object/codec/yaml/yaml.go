package yaml

import (
	"fmt"
	"github.com/Aptomi/aptomi/pkg/object"
	"github.com/Aptomi/aptomi/pkg/object/codec"
	"gopkg.in/yaml.v2"
)

// NewCodec creates a new YAML codec to marshal and unmarshal objects
func NewCodec(catalog *object.Catalog) codec.MarshallerUnmarshaller {
	return &yamlCodec{catalog}
}

type yamlCodec struct {
	catalog *object.Catalog
}

func (c *yamlCodec) GetName() string {
	return "yaml"
}

func (c *yamlCodec) MarshalOne(object object.Base) ([]byte, error) {
	return yaml.Marshal(&object)
}

func (c *yamlCodec) MarshalMany(objects []object.Base) ([]byte, error) {
	return yaml.Marshal(&objects)
}

func (c *yamlCodec) UnmarshalOne(data []byte) (object.Base, error) {
	objects, err := c.unmarshalOneOrMany(data, true)
	if err != nil {
		return nil, err
	}

	return objects[0], nil
}

func (c *yamlCodec) UnmarshalOneOrMany(data []byte) ([]object.Base, error) {
	return c.unmarshalOneOrMany(data, false)
}

func (c *yamlCodec) unmarshalOneOrMany(data []byte, strictOne bool) ([]object.Base, error) {
	raw := new(interface{})
	err := yaml.Unmarshal(data, raw)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshalling data to raw interface{}: %s", err)
	}

	result := make([]object.Base, 0)

	if elem, ok := (*raw).(map[interface{}]interface{}); ok { // if it's a single object (map)
		obj, err := c.unmarshalRaw(elem, data)
		if err != nil {
			return nil, fmt.Errorf("error while unmarshalling single object: %s", err)
		}

		result = append(result, obj)
	} else if strictOne { // if single object strictly required
		return nil, fmt.Errorf("single object expected")
	} else if slice, ok := (*raw).([]interface{}); ok { // if it's an object slice
		for idx, rawElem := range slice {
			sliceElem, ok := rawElem.(map[interface{}]interface{}) // each slice elem should be map
			if !ok {
				return nil, fmt.Errorf("element #%d isn't an object", idx)
			}

			elemData, err := yaml.Marshal(sliceElem) // get []byte for current elem only
			if err != nil {
				return nil, fmt.Errorf("error while unmarshalling element #%d (marshal step): %s", idx, err)
			}

			obj, err := c.unmarshalRaw(sliceElem, elemData) // unmarshal to kind type
			if err != nil {
				return nil, fmt.Errorf("error while unmarshalling element #%d (final step): %s", idx, err)
			}

			result = append(result, obj)
		}
	} else { // if it's not an object or object slice
		return nil, fmt.Errorf("unmarshalable data (not an object or object slice): %T", raw)
	}

	return result, nil
}

func (c *yamlCodec) unmarshalRaw(single map[interface{}]interface{}, data []byte) (object.Base, error) {
	metaField, ok := single["metadata"]
	if !ok {
		return nil, fmt.Errorf("can't find metadata field inside object: %v", single)
	}

	meta, ok := metaField.(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("metadata field isn't a map: %v", single)
	}

	kindField, ok := meta["kind"]
	if !ok {
		return nil, fmt.Errorf("can't find kind field in metadata: %v", single)
	}

	kind, ok := kindField.(string)
	if !ok {
		return nil, fmt.Errorf("kind field in metadata isn't a string: %v", single)
	}

	if len(kind) == 0 {
		return nil, fmt.Errorf("empty kind")
	}

	objectInfo := c.catalog.Get(kind)
	if objectInfo == nil {
		return nil, fmt.Errorf("unknown kind: %s", kind)
	}

	obj := objectInfo.New()
	err := yaml.Unmarshal(data, obj)
	if err != nil {
		return nil, err
	}

	return obj, nil
}
