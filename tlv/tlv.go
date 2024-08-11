package tlv

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	logger "github.com/sirupsen/logrus"

)

type fragments map[int][][]byte

func (f fragments) Add(tag int, buf []byte) {
	f[tag] = append(f[tag], buf)
}

func (f fragments) Get(tag int) ([][]byte, bool) {
	ret, t := f[tag]
	return ret, t
}

func Unmarshal(b []byte, v interface{}) error {
	return decodeValue(b, v)
}

func isNumber(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	}
	return false
}

func isRefType(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer,
		reflect.Interface, reflect.Slice:
		return true
	default:
		return false
	}
}

func hasValue(value reflect.Value) bool {
	if isRefType(value.Type()) {
		return !value.IsNil()
	} else {
		return value.IsValid()
	}
}

func decodeValue(b []byte, v interface{}) error {
	value := reflect.ValueOf(v)

	if unmarshaler, ok := value.Interface().(encoding.BinaryUnmarshaler); ok {
		err := unmarshaler.UnmarshalBinary(b)
		return err
	}

	value = reflect.Indirect(value)
	valueType := reflect.TypeOf(value.Interface())
	switch value.Kind() {
	case reflect.Int8:
		tmp := int64(int8(b[0]))
		value.SetInt(tmp)
	case reflect.Int16:
		tmp := int64(int16(binary.BigEndian.Uint16(b)))
		value.SetInt(tmp)
	case reflect.Int32:
		tmp := int64(int32(binary.BigEndian.Uint32(b)))
		value.SetInt(tmp)
	case reflect.Int64:
		tmp := int64(binary.BigEndian.Uint64(b))
		value.SetInt(tmp)
	case reflect.Int:
		tmp := int64(binary.BigEndian.Uint64(b))
		value.SetInt(tmp)
	case reflect.Uint8:
		tmp := uint64(b[0])
		value.SetUint(tmp)
	case reflect.Uint16:
		tmp := uint64(binary.BigEndian.Uint16(b))
		value.SetUint(tmp)
	case reflect.Uint32:
		tmp := uint64(binary.BigEndian.Uint32(b))
		value.SetUint(tmp)
	case reflect.Uint64:
		tmp := binary.BigEndian.Uint64(b)
		value.SetUint(tmp)
	case reflect.Uint:
		tmp := binary.BigEndian.Uint64(b)
		value.SetUint(tmp)
	case reflect.String:
		value.SetString(string(b))
	case reflect.Ptr:
		if value.IsNil() {
			value.Set(reflect.New(value.Type().Elem()))
		}
		if err := decodeValue(b, value.Interface()); err != nil {
			return err
		}
	case reflect.Struct:
		var tlvFragment fragments
		if tlvFragmentTmp, err := parseTLV(b); err != nil {
			return err
		} else {
			tlvFragment = tlvFragmentTmp
		}
		for i := 0; i < value.NumField(); i++ {
			fieldValue := value.Field(i)
			fieldType := valueType.Field(i)

			tag, hasTLV := fieldType.Tag.Lookup("tlv")
			if !hasTLV {
				return errors.New("field " + fieldType.Name + " need tag `tlv`")
			}

			tagVal, err := strconv.Atoi(tag)
			if err != nil {
				return fmt.Errorf("invalid tlv tag \"%s\", need to be decimal number", tag)
			}

			if len(tlvFragment[tagVal]) == 0 {
				continue
			}

			if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
				fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
			} else if fieldValue.Kind() == reflect.Slice && fieldValue.IsNil() {
				fieldValue.Set(reflect.MakeSlice(fieldValue.Type(), 0, 1))
			}

			for _, buf := range tlvFragment[tagVal] {
				if fieldValue.Kind() != reflect.Ptr {
					fieldValue = fieldValue.Addr()
				}
				err = decodeValue(buf, fieldValue.Interface())
				if err != nil {
					return err
				}
			}
		}
	case reflect.Slice:
		if value.IsNil() {
			value.Set(reflect.MakeSlice(value.Type(), 0, 1))
		}
		if valueType.Elem().Kind() == reflect.Uint8 {
			value.SetBytes(b)
		} else if valueType.Elem().Kind() == reflect.Ptr || valueType.Elem().Kind() == reflect.Struct ||
			isNumber(valueType.Elem()) {
			elemValue := reflect.New(valueType.Elem())
			if err := decodeValue(b, elemValue.Interface()); err != nil {
				return err
			}
			value.Set(reflect.Append(value, elemValue.Elem()))
		} else {
			return errors.New("value type `Slice of " + valueType.String() + "` is not support decode")
		}
	}
	return nil
}

func parseTLV(b []byte) (fragments, error) {
	tlvFragment := make(fragments)
	buffer := bytes.NewBuffer(b)

	var tag uint16
	var length uint16
	for {
		if err := binary.Read(buffer, binary.BigEndian, &tag); err != nil {
			fmt.Printf("Binary Read error: %v", err)
		}
		if err := binary.Read(buffer, binary.BigEndian, &length); err != nil {
			fmt.Printf("Binary Read error: %v", err)
		}
		value := make([]byte, length)
		if _, err := buffer.Read(value); err != nil {
			return nil, err
		}
		tlvFragment.Add(int(tag), value)
		if buffer.Len() == 0 {
			break
		}
	}
	return tlvFragment, nil
}

func Marshal(v interface{}) ([]byte, error) {
	if reflect.TypeOf(v).Kind() != reflect.Struct && reflect.TypeOf(v).Kind() != reflect.Ptr {
		return nil, errors.New("tlv need struct value to encode")
	}
	return buildTLV(0, v)
}

func makeTLV(tag int, value []byte) []byte {
	buf := new(bytes.Buffer)
	if tag != 0 {
		if err := binary.Write(buf, binary.BigEndian, uint16(tag)); err != nil {
			logger.Printf("makeTLV binary write type error: %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, uint16(len(value))); err != nil {
			logger.Printf("makeTLV binary write value error: %+v", err)
		}
	}

	if err := binary.Write(buf, binary.BigEndian, value); err != nil {
		logger.Printf("makeTLV write value error %+v", err)
	}

	return buf.Bytes()
}

func buildTLV(tag int, v interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	value := reflect.ValueOf(v)

	if marshaler, ok := value.Interface().(encoding.BinaryMarshaler); ok {
		bin, err := marshaler.MarshalBinary()
		if err != nil {
			return nil, err
		}

		return makeTLV(tag, bin), nil
	}

	value = reflect.Indirect(value)
	switch value.Kind() {
	case reflect.Int8:
		if err := binary.Write(buf, binary.BigEndian, uint16(tag)); err != nil {
			logger.Printf("tag write error: %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, uint16(1)); err != nil {
			logger.Printf("len(1) write error: %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, v); err != nil {
			logger.Printf("value write error: %+v", err)
		}
		return buf.Bytes(), nil
	case reflect.Int16:
		if err := binary.Write(buf, binary.BigEndian, uint16(tag)); err != nil {
			logger.Printf("tag write error: %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, uint16(2)); err != nil {
			logger.Printf("len(2) write error %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, v); err != nil {
			logger.Printf("value write error %+v", err)
		}
		return buf.Bytes(), nil
	case reflect.Int32:
		if err := binary.Write(buf, binary.BigEndian, uint16(tag)); err != nil {
			logger.Printf("tag write error: %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, uint16(4)); err != nil {
			logger.Printf("len(4) write error %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, v); err != nil {
			logger.Printf("value write error %+v", err)
		}
		return buf.Bytes(), nil
	case reflect.Int64:
		if err := binary.Write(buf, binary.BigEndian, uint16(tag)); err != nil {
			logger.Printf("tag write error: %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, uint16(8)); err != nil {
			logger.Printf("len(8) write error: %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, v); err != nil {
			logger.Printf("value write error: %+v", err)
		}
		return buf.Bytes(), nil
	case reflect.Uint8:
		if err := binary.Write(buf, binary.BigEndian, uint16(tag)); err != nil {
			logger.Printf("tag write error: %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, uint16(1)); err != nil {
			logger.Printf("len(1) write error: %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, v); err != nil {
			logger.Printf("value write error: %+v", err)
		}
		return buf.Bytes(), nil
	case reflect.Uint16:
		if err := binary.Write(buf, binary.BigEndian, uint16(tag)); err != nil {
			logger.Printf("tag write error: %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, uint16(2)); err != nil {
			logger.Printf("len(2) write error: %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, v); err != nil {
			logger.Printf("value write error: %+v", err)
		}
		return buf.Bytes(), nil
	case reflect.Uint32:
		if err := binary.Write(buf, binary.BigEndian, uint16(tag)); err != nil {
			logger.Printf("tag write error: %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, uint16(4)); err != nil {
			logger.Printf("len(4) write error: %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, v); err != nil {
			logger.Printf("value write error: %+v", err)
		}
		return buf.Bytes(), nil
	case reflect.Uint64:
		if err := binary.Write(buf, binary.BigEndian, uint16(tag)); err != nil {
			logger.Printf("tag write error: %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, uint16(8)); err != nil {
			logger.Printf("len(8) write error: %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, v); err != nil {
			logger.Printf("value write error: %+v", err)
		}
		return buf.Bytes(), nil
	case reflect.String:
		str := v.(string)
		if err := binary.Write(buf, binary.BigEndian, uint16(tag)); err != nil {
			logger.Printf("tag write error: %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, uint16(len(str))); err != nil {
			logger.Printf("len(str) write error: %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, []byte(str)); err != nil {
			logger.Printf("value write error: %+v", err)
		}
		return buf.Bytes(), nil
	case reflect.Struct:
		for i := 0; i < value.Type().NumField(); i++ {
			field := value.Field(i)
			if !hasValue(field) {
				continue
			}
			structField := value.Type().Field(i)
			tlvTag, hasTLV := structField.Tag.Lookup("tlv")
			if !hasTLV {
				return nil, errors.New("field " + structField.Name + " need tag `tlv`")
			}
			tagVal, err := strconv.Atoi(tlvTag)
			if err != nil {
				return nil, err
			}
			subValue, err := buildTLV(tagVal, field.Interface())
			if err != nil {
				return nil, err
			}
			if err := binary.Write(buf, binary.BigEndian, subValue); err != nil {
				return nil, err
			}
		}

		return makeTLV(tag, buf.Bytes()), nil
	case reflect.Slice:
		if value.Type().Elem().Kind() == reflect.Uint8 {
			if err := binary.Write(buf, binary.BigEndian, uint16(tag)); err != nil {
				logger.Printf("Binary write error: %+v", err)
			}
			if err := binary.Write(buf, binary.BigEndian, uint16(value.Len())); err != nil {
				logger.Printf("Binary write error: %+v", err)
			}
			if err := binary.Write(buf, binary.BigEndian, v); err != nil {
				logger.Printf("Binary write error: %+v", err)
			}
			return buf.Bytes(), nil
		} else if value.Type().Elem().Kind() == reflect.Ptr || value.Type().Elem().Kind() == reflect.Struct ||
			isNumber(value.Type().Elem()) {
			for i := 0; i < value.Len(); i++ {
				elem := value.Index(i)
				if value.Type().Elem().Kind() == reflect.Struct {
					elem = elem.Addr()
				}
				elemBuf, err := buildTLV(tag, elem.Interface())
				if err != nil {
					return nil, err
				}
				buf.Write(elemBuf)
			}
			return buf.Bytes(), nil
		} else {
			return nil, errors.New("value type `Slice of " + value.Type().Elem().Name() + "` is not support TLV encode")
		}
	default:
		return nil, errors.New("value type " + value.Type().String() + " is not support TLV encode")
	}
}