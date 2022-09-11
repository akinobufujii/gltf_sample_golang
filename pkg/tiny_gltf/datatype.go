package tiny_gltf

type BufferDataKind int

const (
	BufferDataKind_UNDEFINED BufferDataKind = iota
	BufferDataKind_BYTE
	BufferDataKind_UNSIGNED_BYTE
	BufferDataKind_SHORT
	BufferDataKind_UNSIGNED_SHORT
	BufferDataKind_UNSIGNED_INT
	BufferDataKind_FLOAT
)

type BufferDataType interface {
	int8 | uint8 | int16 | uint16 | uint32 | float32
}

type BufferData[T BufferDataType] struct {
	data []T
}

func NewBufferData[T BufferDataType](d []T) *BufferData[T] {
	return &BufferData[T]{data: d}
}
