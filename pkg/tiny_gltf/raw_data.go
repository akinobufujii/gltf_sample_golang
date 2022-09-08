package tiny_gltf

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

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
	Data []byte
}

func (buffer *Buffers) UnmarshalJSON(data []byte) error {
	type rawdata struct {
		URI        string `json:"uri"`
		ByteLength int    `json:"byteLength"`
	}

	var raw rawdata
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	bytedata := make([]byte, 0, raw.ByteLength)
	if strings.HasPrefix(raw.URI, "data:") {
		// NOTE:  decode embeded data
		// FIXME: refactor
		raw.URI = raw.URI[len("data:"):]
		assetType := raw.URI[0:strings.Index(raw.URI, ";")]
		raw.URI = raw.URI[len(assetType)+1:]
		switch assetType {
		case "application/octet-stream":
			decodeType := raw.URI[0:strings.Index(raw.URI, ",")]
			switch decodeType {
			case "base64":
				bytedata, err = base64.StdEncoding.DecodeString(raw.URI[len(decodeType)+1:])
				if err != nil {
					return fmt.Errorf("failed base64.RawStdEncoding.DecodeString %w", err)
				}
			}
		}
	} else {
		// NOTE: load a bytedata from file
		file, err := os.Open(raw.URI)
		if err != nil {
			return fmt.Errorf("failed os.Open %s %w", raw.URI, err)
		}

		if _, err := file.Read(bytedata); err != nil {
			return fmt.Errorf("failed file.Read %s %w", raw.URI, err)
		}
	}

	buffer.Data = bytedata

	return nil
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
