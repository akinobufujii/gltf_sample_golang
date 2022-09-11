package tiny_gltf

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Buffer[T BufferDataType] struct {
	bufferData *BufferData[T]
	bufferKind BufferDataKind
}

type Instance struct {
	RawData *rawData
}

func NewInstanceFromFile(filename string) (*Instance, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed os.Open %s %w", filename, err)
	}
	defer file.Close()

	instance, err := NewInstance(file)
	if err != nil {
		return nil, fmt.Errorf("failed NewInstance %s %w", filename, err)
	}

	return instance, nil
}

func NewInstance(reader io.Reader) (*Instance, error) {
	decoder := json.NewDecoder(reader)
	rawData := new(rawData)
	if err := decoder.Decode(rawData); err != nil {
		return nil, fmt.Errorf("failed decoder.Decode %w", err)
	}

	instance := &Instance{
		RawData: rawData,
	}

	if err := instance.initData(); err != nil {
		return nil, fmt.Errorf("failed instance.initData() %w", err)
	}

	return instance, nil
}

func (instance *Instance) initData() error {
	for _, node := range instance.RawData.Scenes[instance.RawData.Scene].Nodes {
		mesh := instance.RawData.Nodes[node].Mesh
		for _, prim := range instance.RawData.Meshes[mesh].Primitives {
			// TODO: refactor
			accessor := instance.RawData.Accessors[prim.Indices]
			dataCount := accessor.Count
			var buffer any
			dataReader := bytes.NewReader(instance.RawData.Buffers[accessor.BufferView].Data)
			switch accessor.ComponentType {
			case 5120: // int8
				b := Buffer[int8]{
					bufferData: NewBufferData(make([]int8, dataCount)),
					bufferKind: BufferDataKind_BYTE,
				}
				err := binary.Read(dataReader, binary.LittleEndian, b.bufferData.data)
				if err != nil {
					return fmt.Errorf("failed binary.Read: %w", err)
				}

				buffer = b
			case 5121: // uint8
				b := Buffer[uint8]{
					bufferData: NewBufferData(make([]uint8, dataCount)),
					bufferKind: BufferDataKind_UNSIGNED_BYTE,
				}

				err := binary.Read(dataReader, binary.LittleEndian, b.bufferData.data)
				if err != nil {
					return fmt.Errorf("failed binary.Read: %w", err)
				}

				buffer = b
			case 5122: // int16
				b := Buffer[int16]{
					bufferData: NewBufferData(make([]int16, dataCount)),
					bufferKind: BufferDataKind_SHORT,
				}

				err := binary.Read(dataReader, binary.LittleEndian, b.bufferData.data)
				if err != nil {
					return fmt.Errorf("failed binary.Read: %w", err)
				}

				buffer = b
			case 5123: // uint16
				b := Buffer[uint16]{
					bufferData: NewBufferData(make([]uint16, dataCount)),
					bufferKind: BufferDataKind_UNSIGNED_SHORT,
				}

				err := binary.Read(dataReader, binary.LittleEndian, b.bufferData.data)
				if err != nil {
					return fmt.Errorf("failed binary.Read: %w", err)
				}

				buffer = b
			case 5125: // uint32
				b := Buffer[uint32]{
					bufferData: NewBufferData(make([]uint32, dataCount)),
					bufferKind: BufferDataKind_UNSIGNED_INT,
				}

				err := binary.Read(dataReader, binary.LittleEndian, b.bufferData.data)
				if err != nil {
					return fmt.Errorf("failed binary.Read: %w", err)
				}

				buffer = b
			case 5126: // float32
				b := Buffer[float32]{
					bufferData: NewBufferData(make([]float32, dataCount)),
					bufferKind: BufferDataKind_FLOAT,
				}

				err := binary.Read(dataReader, binary.LittleEndian, b.bufferData.data)
				if err != nil {
					return fmt.Errorf("failed binary.Read: %w", err)
				}

				buffer = b
			}
			fmt.Printf("%+v\n", buffer)
		}
	}

	return nil
}
