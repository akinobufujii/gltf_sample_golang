package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	runtime.LockOSThread()
}

func main() {

	err := glfw.Init()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 0)

	window, err := glfw.CreateWindow(640, 480, "gltf_samples", nil, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	glfw.SwapInterval(1)

	// instance, err := tiny_gltf.NewInstanceFromFile("pkg/tiny_gltf/testdata/minimal.gltf")
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, err)
	// 	os.Exit(1)
	// }

	var vertexbuffer uint32
	g_vertex_buffer_data := []float32{
		-1.0, -1.0, 0.0,
		1.0, -1.0, 0.0,
		0.0, 1.0, 0.0,
	}
	gl.GenBuffers(1, &vertexbuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexbuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(g_vertex_buffer_data)*4, gl.Ptr(g_vertex_buffer_data), gl.STATIC_DRAW)

	for !window.ShouldClose() {
		gl.EnableVertexAttribArray(0)
		gl.BindBuffer(gl.ARRAY_BUFFER, vertexbuffer)
		gl.VertexAttribPointer(
			0,
			3,
			gl.FLOAT,
			false,
			0,
			gl.Ptr(nil),
		)

		gl.DrawArrays(gl.TRIANGLES, 0, 3)

		gl.DisableVertexAttribArray(0)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
