package tiny_gltf

type rawData struct {
	Scene       int           `json:"scene"`
	Scenes      []Scenes      `json:"scenes"`
	Nodes       []Nodes       `json:"nodes"`
	Meshes      []Meshes      `json:"meshes"`
	Buffers     []Buffers     `json:"buffers"`
	BufferViews []BufferViews `json:"bufferViews"`
	Accessors   []Accessors   `json:"accessors"`
	Asset       Asset         `json:"asset"`
}

type Scenes struct {
	Nodes []int `json:"nodes"`
}

type Nodes struct {
	Mesh int `json:"mesh"`
}

type Attributes struct {
	Position int `json:"POSITION"`
}

type Primitives struct {
	Attributes Attributes `json:"attributes"`
	Indices    int        `json:"indices"`
}

type Meshes struct {
	Primitives []Primitives `json:"primitives"`
}

type Buffers struct {
	URI        string `json:"uri"`
	ByteLength int    `json:"byteLength"`
}

type BufferViews struct {
	Buffer     int `json:"buffer"`
	ByteOffset int `json:"byteOffset"`
	ByteLength int `json:"byteLength"`
	Target     int `json:"target"`
}

type Accessors struct {
	BufferView    int       `json:"bufferView"`
	ByteOffset    int       `json:"byteOffset"`
	ComponentType int       `json:"componentType"`
	Count         int       `json:"count"`
	Type          string    `json:"type"`
	Max           []float32 `json:"max"`
	Min           []float32 `json:"min"`
}

type Asset struct {
	Version string `json:"version"`
}
