package tiny_gltf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type AccessorComponentTypeKind int

const (
	AccessorComponentTypeKind_UNDEFINED AccessorComponentTypeKind = iota
	AccessorComponentTypeKind_BYTE
	AccessorComponentTypeKind_UNSIGNED_BYTE
	AccessorComponentTypeKind_SHORT
	AccessorComponentTypeKind_UNSIGNED_SHORT
	AccessorComponentTypeKind_UNSIGNED_INT
	AccessorComponentTypeKind_FLOAT
)

type AccessorTypeKind int

const (
	AccessorTypeKind_UNDEFINED AccessorTypeKind = iota
	AccessorTypeKind_SCALAR
	AccessorTypeKind_VEC2
	AccessorTypeKind_VEC3
	AccessorTypeKind_VEC4
	AccessorTypeKind_MAT2
	AccessorTypeKind_MAT3
	AccessorTypeKind_MAT4
)

type AccessorComponentType interface {
	int8 | uint8 | int16 | uint16 | uint32 | float32
}

type rawBufferData[T AccessorComponentType] struct {
	data []T
}

type BufferData struct {
	rawdata           any // NOTE: rawBufferData[AccessorComponentType]
	componentTypeKind AccessorComponentTypeKind
	typeKind          AccessorTypeKind
	count             int
}

func getAccessorTypeSize(typeKind AccessorTypeKind) int {
	switch typeKind {
	case AccessorTypeKind_SCALAR:
		return 1
	case AccessorTypeKind_VEC2:
		return 2
	case AccessorTypeKind_VEC3:
		return 3
	case AccessorTypeKind_VEC4:
		return 4
	case AccessorTypeKind_MAT2:
		return 2 * 2
	case AccessorTypeKind_MAT3:
		return 3 * 3
	case AccessorTypeKind_MAT4:
		return 4 * 4
	}
	return 0
}

func convertAccessorTypeKindFromString(kind string) AccessorTypeKind {
	switch kind {
	case "SCALAR":
		return AccessorTypeKind_SCALAR
	case "VEC2":
		return AccessorTypeKind_VEC2
	case "VEC3":
		return AccessorTypeKind_VEC3
	case "VEC4":
		return AccessorTypeKind_VEC4
	case "MAT2":
		return AccessorTypeKind_MAT2
	case "MAT3":
		return AccessorTypeKind_MAT3
	case "MAT4":
		return AccessorTypeKind_MAT4
	}

	return AccessorTypeKind_UNDEFINED
}

func newRawBufferData[T AccessorComponentType](dataReader io.Reader, typesize, count int) (rawBufferData[T], error) {
	data := make([]T, typesize*count)
	err := binary.Read(dataReader, binary.LittleEndian, data)
	if err != nil {
		return rawBufferData[T]{}, fmt.Errorf("failed binary.Read: %w", err)
	}

	return rawBufferData[T]{data: data}, nil
}

func NewBufferData(buffers []Buffers, bufferViews []BufferViews, accessor Accessors) (*BufferData, error) {
	var rawBuffer any
	dataCount := accessor.Count
	componentTypeKind := AccessorComponentTypeKind_UNDEFINED
	typeKind := convertAccessorTypeKindFromString(accessor.Type)

	bufferView := bufferViews[accessor.BufferView]

	dataReader := bytes.NewReader(buffers[bufferView.Buffer].Data[bufferView.ByteOffset+accessor.ByteOffset:])
	switch accessor.ComponentType {
	case 5120: // int8
		componentTypeKind = AccessorComponentTypeKind_BYTE
		buffer, err := newRawBufferData[int8](dataReader, getAccessorTypeSize(typeKind), dataCount)
		if err != nil {
			return nil, fmt.Errorf("failed newRawBufferData: %w", err)
		}

		rawBuffer = buffer
	case 5121: // uint8
		componentTypeKind = AccessorComponentTypeKind_UNSIGNED_BYTE
		buffer, err := newRawBufferData[uint8](dataReader, getAccessorTypeSize(typeKind), dataCount)
		if err != nil {
			return nil, fmt.Errorf("failed newRawBufferData: %w", err)
		}

		rawBuffer = buffer
	case 5122: // int16
		componentTypeKind = AccessorComponentTypeKind_SHORT
		buffer, err := newRawBufferData[int16](dataReader, getAccessorTypeSize(typeKind), dataCount)
		if err != nil {
			return nil, fmt.Errorf("failed newRawBufferData: %w", err)
		}

		rawBuffer = buffer
	case 5123: // uint16
		componentTypeKind = AccessorComponentTypeKind_UNSIGNED_SHORT
		buffer, err := newRawBufferData[uint16](dataReader, getAccessorTypeSize(typeKind), dataCount)
		if err != nil {
			return nil, fmt.Errorf("failed newRawBufferData: %w", err)
		}

		rawBuffer = buffer
	case 5125: // uint32
		componentTypeKind = AccessorComponentTypeKind_UNSIGNED_INT
		buffer, err := newRawBufferData[uint32](dataReader, getAccessorTypeSize(typeKind), dataCount)
		if err != nil {
			return nil, fmt.Errorf("failed newRawBufferData: %w", err)
		}

		rawBuffer = buffer
	case 5126: // float32
		componentTypeKind = AccessorComponentTypeKind_FLOAT
		buffer, err := newRawBufferData[float32](dataReader, getAccessorTypeSize(typeKind), dataCount)
		if err != nil {
			return nil, fmt.Errorf("failed newRawBufferData: %w", err)
		}

		rawBuffer = buffer
	}

	return &BufferData{
		rawdata:           rawBuffer,
		count:             dataCount,
		componentTypeKind: componentTypeKind,
		typeKind:          typeKind,
	}, nil
}
