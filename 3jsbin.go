package bin

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
)

const (
	JS3_BIN_SIGN = "Three.js 002"
)

var (
	littleEndian = binary.LittleEndian
)

type FlatTriangle struct {
	Vertices [][3]uint32
	Material []int16
}

func (obj *FlatTriangle) GetVertices() []uint32 {
	ret := make([]uint32, len(obj.Vertices)*3)
	for i := range obj.Vertices {
		ret[i*3] = obj.Vertices[i][0]
		ret[i*3+1] = obj.Vertices[i][1]
		ret[i*3+2] = obj.Vertices[i][2]
	}
	return ret
}

func (obj *FlatTriangle) GetMaterials() []int16 {
	return obj.Material
}

func (obj *FlatTriangle) encode(wr io.Writer) error {
	if len(obj.Vertices) != len(obj.Material) {
		return errors.New("Vertices size must eq")
	}

	err := binary.Write(wr, littleEndian, obj.GetVertices())
	if err != nil {
		return err
	}

	err = binary.Write(wr, littleEndian, obj.GetMaterials())
	if err != nil {
		return err
	}

	return nil
}

func (obj *FlatTriangle) decode(rd io.ReadSeeker, size uint32) error {
	verbuf := make([]uint32, size*3)
	err := binary.Read(rd, littleEndian, verbuf)
	if err != nil {
		return err
	}
	err = obj.SetVertices(verbuf)
	if err != nil {
		return err
	}

	mtlbuf := make([]int16, size)
	err = binary.Read(rd, littleEndian, mtlbuf)
	if err != nil {
		return err
	}
	err = obj.SetMaterial(mtlbuf)
	if err != nil {
		return err
	}
	return nil
}

func (obj *FlatTriangle) SetVertices(vers interface{}) error {
	switch t := vers.(type) {
	case []uint32:
		if len(t)%3 != 0 {
			return errors.New("[]float32 must 3^")
		}
		obj.Vertices = make([][3]uint32, len(t)/3)
		for i := range obj.Vertices {
			obj.Vertices[i][0] = t[i*3]
			obj.Vertices[i][1] = t[i*3+1]
			obj.Vertices[i][2] = t[i*3+2]
		}
	case [][]uint32:
		obj.Vertices = make([][3]uint32, len(t))
		for i := range obj.Vertices {
			obj.Vertices[i][0] = t[i][0]
			obj.Vertices[i][1] = t[i][1]
			obj.Vertices[i][2] = t[i][2]
		}
	case [][3]uint32:
		obj.Vertices = t
	}
	return nil
}

func (obj *FlatTriangle) SetMaterial(vers interface{}) error {
	switch t := vers.(type) {
	case []int16:
		obj.Material = t
	default:
		return errors.New("Material must []int16")
	}
	return nil
}

type SmoothTriangle struct {
	FlatTriangle
	Normals [][3]uint32
}

func (obj *SmoothTriangle) GetNormals() []uint32 {
	ret := make([]uint32, len(obj.Normals)*3)
	for i := range obj.Normals {
		ret[i*3] = obj.Normals[i][0]
		ret[i*3+1] = obj.Normals[i][1]
		ret[i*3+2] = obj.Normals[i][2]
	}
	return ret
}

func (obj *SmoothTriangle) encode(wr io.Writer) error {
	err := obj.FlatTriangle.encode(wr)
	if err != nil {
		return err
	}
	err = binary.Write(wr, littleEndian, obj.GetNormals())
	if err != nil {
		return err
	}

	return nil
}

func (obj *SmoothTriangle) decode(rd io.ReadSeeker, size uint32) error {
	err := obj.FlatTriangle.decode(rd, size)
	if err != nil {
		return err
	}
	verbuf := make([]uint32, size*3)
	err = binary.Read(rd, littleEndian, verbuf)
	if err != nil {
		return err
	}
	err = obj.SetNormals(verbuf)
	if err != nil {
		return err
	}
	return nil
}

func (obj *SmoothTriangle) SetNormals(vers interface{}) error {
	switch t := vers.(type) {
	case []uint32:
		if len(t)%3 != 0 {
			return errors.New("[]float32 must 3^")
		}
		obj.Normals = make([][3]uint32, len(t)/3)
		for i := range obj.Normals {
			obj.Normals[i][0] = t[i*3]
			obj.Normals[i][1] = t[i*3+1]
			obj.Normals[i][2] = t[i*3+2]
		}
	case [][]uint32:
		obj.Normals = make([][3]uint32, len(t))
		for i := range obj.Normals {
			obj.Normals[i][0] = t[i][0]
			obj.Normals[i][1] = t[i][1]
			obj.Normals[i][2] = t[i][2]
		}
	case [][3]uint32:
		obj.Normals = t
	}
	return nil
}

type FlatUVTriangle struct {
	FlatTriangle
	Uvs [][3]uint32
}

func (obj *FlatUVTriangle) GetUvs() []uint32 {
	ret := make([]uint32, len(obj.Uvs)*3)
	for i := range obj.Uvs {
		ret[i*3] = obj.Uvs[i][0]
		ret[i*3+1] = obj.Uvs[i][1]
		ret[i*3+2] = obj.Uvs[i][2]
	}
	return ret
}

func (obj *FlatUVTriangle) encode(wr io.Writer) error {
	err := obj.FlatTriangle.encode(wr)
	if err != nil {
		return err
	}
	err = binary.Write(wr, littleEndian, obj.GetUvs())
	if err != nil {
		return err
	}

	return nil
}

func (obj *FlatUVTriangle) decode(rd io.ReadSeeker, size uint32) error {
	err := obj.FlatTriangle.decode(rd, size)
	if err != nil {
		return err
	}
	uvbuf := make([]uint32, size*3)
	err = binary.Read(rd, littleEndian, uvbuf)
	if err != nil {
		return err
	}
	err = obj.SetUVs(uvbuf)
	if err != nil {
		return err
	}
	return nil
}

func (obj *FlatUVTriangle) SetUVs(vers interface{}) error {
	switch t := vers.(type) {
	case []uint32:
		if len(t)%3 != 0 {
			return errors.New("[]float32 must 3^")
		}
		obj.Uvs = make([][3]uint32, len(t)/3)
		for i := range obj.Uvs {
			obj.Uvs[i][0] = t[i*3]
			obj.Uvs[i][1] = t[i*3+1]
			obj.Uvs[i][2] = t[i*3+2]
		}
	case [][]uint32:
		obj.Uvs = make([][3]uint32, len(t))
		for i := range obj.Uvs {
			obj.Uvs[i][0] = t[i][0]
			obj.Uvs[i][1] = t[i][1]
			obj.Uvs[i][2] = t[i][2]
		}
	case [][3]uint32:
		obj.Uvs = t
	}
	return nil
}

type SmoothUVTriangle struct {
	FlatTriangle
	Normals [][3]uint32
	Uvs     [][3]uint32
}

func (obj *SmoothUVTriangle) GetNormals() []uint32 {
	ret := make([]uint32, len(obj.Normals)*3)
	for i := range obj.Normals {
		ret[i*3] = obj.Normals[i][0]
		ret[i*3+1] = obj.Normals[i][1]
		ret[i*3+2] = obj.Normals[i][2]
	}
	return ret
}

func (obj *SmoothUVTriangle) GetUvs() []uint32 {
	ret := make([]uint32, len(obj.Uvs)*3)
	for i := range obj.Uvs {
		ret[i*3] = obj.Uvs[i][0]
		ret[i*3+1] = obj.Uvs[i][1]
		ret[i*3+2] = obj.Uvs[i][2]
	}
	return ret
}

func (obj *SmoothUVTriangle) encode(wr io.Writer) error {
	err := obj.FlatTriangle.encode(wr)
	if err != nil {
		return err
	}
	err = binary.Write(wr, littleEndian, obj.GetNormals())
	if err != nil {
		return err
	}

	err = binary.Write(wr, littleEndian, obj.GetUvs())
	if err != nil {
		return err
	}

	return nil
}

func (obj *SmoothUVTriangle) decode(rd io.ReadSeeker, size uint32) error {
	err := obj.FlatTriangle.decode(rd, size)
	if err != nil {
		return err
	}
	verbuf := make([]uint32, size*3)
	err = binary.Read(rd, littleEndian, verbuf)
	if err != nil {
		return err
	}
	err = obj.SetNormals(verbuf)
	if err != nil {
		return err
	}

	uvbuf := make([]uint32, size*3)
	err = binary.Read(rd, littleEndian, uvbuf)
	if err != nil {
		return err
	}
	err = obj.SetUVs(uvbuf)
	if err != nil {
		return err
	}
	return nil
}

func (obj *SmoothUVTriangle) SetNormals(vers interface{}) error {
	switch t := vers.(type) {
	case []uint32:
		if len(t)%3 != 0 {
			return errors.New("[]float32 must 3^")
		}
		obj.Normals = make([][3]uint32, len(t)/3)
		for i := range obj.Normals {
			obj.Normals[i][0] = t[i*3]
			obj.Normals[i][1] = t[i*3+1]
			obj.Normals[i][2] = t[i*3+2]
		}
	case [][]uint32:
		obj.Normals = make([][3]uint32, len(t))
		for i := range obj.Normals {
			obj.Normals[i][0] = t[i][0]
			obj.Normals[i][1] = t[i][1]
			obj.Normals[i][2] = t[i][2]
		}
	case [][3]uint32:
		obj.Normals = t
	}
	return nil
}

func (obj *SmoothUVTriangle) SetUVs(vers interface{}) error {
	switch t := vers.(type) {
	case []uint32:
		if len(t)%3 != 0 {
			return errors.New("[]float32 must 3^")
		}
		obj.Uvs = make([][3]uint32, len(t)/3)
		for i := range obj.Uvs {
			obj.Uvs[i][0] = t[i*3]
			obj.Uvs[i][1] = t[i*3+1]
			obj.Uvs[i][2] = t[i*3+2]
		}
	case [][]uint32:
		obj.Uvs = make([][3]uint32, len(t))
		for i := range obj.Uvs {
			obj.Uvs[i][0] = t[i][0]
			obj.Uvs[i][1] = t[i][1]
			obj.Uvs[i][2] = t[i][2]
		}
	case [][3]uint32:
		obj.Uvs = t
	}
	return nil
}

type FlatQuad struct {
	Vertices [][4]uint32
	Material []int16
}

func (obj *FlatQuad) GetVertices() []uint32 {
	ret := make([]uint32, len(obj.Vertices)*4)
	for i := range obj.Vertices {
		ret[i*3] = obj.Vertices[i][0]
		ret[i*3+1] = obj.Vertices[i][1]
		ret[i*3+2] = obj.Vertices[i][2]
		ret[i*3+3] = obj.Vertices[i][3]
	}
	return ret
}

func (obj *FlatQuad) GetMaterials() []int16 {
	return obj.Material
}

func (obj *FlatQuad) encode(wr io.Writer) error {
	if len(obj.Vertices) != len(obj.Material) {
		return errors.New("Vertices size must eq")
	}

	err := binary.Write(wr, littleEndian, obj.GetVertices())
	if err != nil {
		return err
	}

	err = binary.Write(wr, littleEndian, obj.GetMaterials())
	if err != nil {
		return err
	}
	return nil
}

func (obj *FlatQuad) decode(rd io.ReadSeeker, size uint32) error {
	verbuf := make([]uint32, size*4)
	err := binary.Read(rd, littleEndian, verbuf)
	if err != nil {
		return err
	}
	err = obj.SetVertices(verbuf)
	if err != nil {
		return err
	}

	mtlbuf := make([]int16, size)
	err = binary.Read(rd, littleEndian, mtlbuf)
	if err != nil {
		return err
	}
	err = obj.SetMaterial(mtlbuf)
	if err != nil {
		return err
	}
	return nil
}

func (obj *FlatQuad) SetVertices(vers interface{}) error {
	switch t := vers.(type) {
	case []uint32:
		if len(t)%4 != 0 {
			return errors.New("[]uint32 must 4^")
		}
		obj.Vertices = make([][4]uint32, len(t)/4)
		for i := range obj.Vertices {
			obj.Vertices[i][0] = t[i*3]
			obj.Vertices[i][1] = t[i*3+1]
			obj.Vertices[i][2] = t[i*3+2]
			obj.Vertices[i][3] = t[i*3+3]
		}
	case [][]uint32:
		obj.Vertices = make([][4]uint32, len(t))
		for i := range obj.Vertices {
			obj.Vertices[i][0] = t[i][0]
			obj.Vertices[i][1] = t[i][1]
			obj.Vertices[i][2] = t[i][2]
			obj.Vertices[i][3] = t[i][3]
		}
	case [][4]uint32:
		obj.Vertices = t
	}
	return nil
}

func (obj *FlatQuad) SetMaterial(vers interface{}) error {
	switch t := vers.(type) {
	case []int16:
		obj.Material = t
	default:
		return errors.New("Material must []int16")
	}
	return nil
}

type SmoothQuad struct {
	FlatQuad
	Normals [][4]uint32
}

func (obj *SmoothQuad) GetNormals() []uint32 {
	ret := make([]uint32, len(obj.Normals)*4)
	for i := range obj.Normals {
		ret[i*3] = obj.Normals[i][0]
		ret[i*3+1] = obj.Normals[i][1]
		ret[i*3+2] = obj.Normals[i][2]
		ret[i*3+3] = obj.Normals[i][3]
	}
	return ret
}

func (obj *SmoothQuad) encode(wr io.Writer) error {
	err := obj.FlatQuad.encode(wr)
	if err != nil {
		return err
	}

	err = binary.Write(wr, littleEndian, obj.GetNormals())
	if err != nil {
		return err
	}
	return nil
}

func (obj *SmoothQuad) decode(rd io.ReadSeeker, size uint32) error {
	err := obj.FlatQuad.decode(rd, size)
	if err != nil {
		return err
	}
	verbuf := make([]uint32, size*4)
	err = binary.Read(rd, littleEndian, verbuf)
	if err != nil {
		return err
	}
	err = obj.SetNormals(verbuf)
	if err != nil {
		return err
	}
	return nil
}

func (obj *SmoothQuad) SetNormals(vers interface{}) error {
	switch t := vers.(type) {
	case []uint32:
		if len(t)%4 != 0 {
			return errors.New("[]uint32 must 4^")
		}
		obj.Normals = make([][4]uint32, len(t)/4)
		for i := range obj.Normals {
			obj.Normals[i][0] = t[i*3]
			obj.Normals[i][1] = t[i*3+1]
			obj.Normals[i][2] = t[i*3+2]
			obj.Normals[i][3] = t[i*3+3]
		}
	case [][]uint32:
		obj.Normals = make([][4]uint32, len(t))
		for i := range obj.Normals {
			obj.Normals[i][0] = t[i][0]
			obj.Normals[i][1] = t[i][1]
			obj.Normals[i][2] = t[i][2]
			obj.Normals[i][3] = t[i][3]
		}
	case [][4]uint32:
		obj.Normals = t
	}
	return nil
}

type FlatUVQuad struct {
	FlatQuad
	Uvs [][4]uint32
}

func (obj *FlatUVQuad) GetUvs() []uint32 {
	ret := make([]uint32, len(obj.Uvs)*4)
	for i := range obj.Uvs {
		ret[i*3] = obj.Uvs[i][0]
		ret[i*3+1] = obj.Uvs[i][1]
		ret[i*3+2] = obj.Uvs[i][2]
		ret[i*3+3] = obj.Uvs[i][3]
	}
	return ret
}

func (obj *FlatUVQuad) encode(wr io.Writer) error {
	err := obj.FlatQuad.encode(wr)
	if err != nil {
		return err
	}

	err = binary.Write(wr, littleEndian, obj.GetUvs())
	if err != nil {
		return err
	}
	return nil
}

func (obj *FlatUVQuad) decode(rd io.ReadSeeker, size uint32) error {
	err := obj.FlatQuad.decode(rd, size)
	if err != nil {
		return err
	}
	uvbuf := make([]uint32, size*4)
	err = binary.Read(rd, littleEndian, uvbuf)
	if err != nil {
		return err
	}
	err = obj.SetUVs(uvbuf)
	if err != nil {
		return err
	}
	return nil
}

func (obj *FlatUVQuad) SetUVs(vers interface{}) error {
	switch t := vers.(type) {
	case []uint32:
		if len(t)%4 != 0 {
			return errors.New("[]uint32 must 4^")
		}
		obj.Uvs = make([][4]uint32, len(t)/4)
		for i := range obj.Uvs {
			obj.Uvs[i][0] = t[i*3]
			obj.Uvs[i][1] = t[i*3+1]
			obj.Uvs[i][2] = t[i*3+2]
			obj.Uvs[i][3] = t[i*3+3]
		}
	case [][]uint32:
		obj.Uvs = make([][4]uint32, len(t))
		for i := range obj.Uvs {
			obj.Uvs[i][0] = t[i][0]
			obj.Uvs[i][1] = t[i][1]
			obj.Uvs[i][2] = t[i][2]
			obj.Uvs[i][3] = t[i][3]
		}
	case [][4]uint32:
		obj.Uvs = t
	}
	return nil
}

type SmoothUVQuad struct {
	FlatQuad
	Normals [][4]uint32
	Uvs     [][4]uint32
}

func (obj *SmoothUVQuad) GetNormals() []uint32 {
	ret := make([]uint32, len(obj.Normals)*4)
	for i := range obj.Normals {
		ret[i*3] = obj.Normals[i][0]
		ret[i*3+1] = obj.Normals[i][1]
		ret[i*3+2] = obj.Normals[i][2]
		ret[i*3+3] = obj.Normals[i][3]
	}
	return ret
}

func (obj *SmoothUVQuad) GetUvs() []uint32 {
	ret := make([]uint32, len(obj.Uvs)*4)
	for i := range obj.Uvs {
		ret[i*3] = obj.Uvs[i][0]
		ret[i*3+1] = obj.Uvs[i][1]
		ret[i*3+2] = obj.Uvs[i][2]
		ret[i*3+3] = obj.Uvs[i][3]
	}
	return ret
}

func (obj *SmoothUVQuad) encode(wr io.Writer) error {
	err := obj.FlatQuad.encode(wr)
	if err != nil {
		return err
	}

	err = binary.Write(wr, littleEndian, obj.GetNormals())
	if err != nil {
		return err
	}

	err = binary.Write(wr, littleEndian, obj.GetUvs())
	if err != nil {
		return err
	}
	return nil
}

func (obj *SmoothUVQuad) decode(rd io.ReadSeeker, size uint32) error {
	err := obj.FlatQuad.decode(rd, size)
	if err != nil {
		return err
	}
	verbuf := make([]uint32, size*4)
	err = binary.Read(rd, littleEndian, verbuf)
	if err != nil {
		return err
	}
	err = obj.SetNormals(verbuf)
	if err != nil {
		return err
	}

	uvbuf := make([]uint32, size*4)
	err = binary.Read(rd, littleEndian, uvbuf)
	if err != nil {
		return err
	}
	err = obj.SetUVs(uvbuf)
	if err != nil {
		return err
	}
	return nil
}

func (obj *SmoothUVQuad) SetNormals(vers interface{}) error {
	switch t := vers.(type) {
	case []uint32:
		if len(t)%4 != 0 {
			return errors.New("[]uint32 must 4^")
		}
		obj.Normals = make([][4]uint32, len(t)/4)
		for i := range obj.Normals {
			obj.Normals[i][0] = t[i*3]
			obj.Normals[i][1] = t[i*3+1]
			obj.Normals[i][2] = t[i*3+2]
			obj.Normals[i][3] = t[i*3+3]
		}
	case [][]uint32:
		obj.Normals = make([][4]uint32, len(t))
		for i := range obj.Normals {
			obj.Normals[i][0] = t[i][0]
			obj.Normals[i][1] = t[i][1]
			obj.Normals[i][2] = t[i][2]
			obj.Normals[i][3] = t[i][3]
		}
	case [][4]uint32:
		obj.Normals = t
	}
	return nil
}

func (obj *SmoothUVQuad) SetUVs(vers interface{}) error {
	switch t := vers.(type) {
	case []uint32:
		if len(t)%4 != 0 {
			return errors.New("[]uint32 must 4^")
		}
		obj.Uvs = make([][4]uint32, len(t)/4)
		for i := range obj.Uvs {
			obj.Uvs[i][0] = t[i*3]
			obj.Uvs[i][1] = t[i*3+1]
			obj.Uvs[i][2] = t[i*3+2]
			obj.Uvs[i][3] = t[i*3+3]
		}
	case [][]uint32:
		obj.Uvs = make([][4]uint32, len(t))
		for i := range obj.Uvs {
			obj.Uvs[i][0] = t[i][0]
			obj.Uvs[i][1] = t[i][1]
			obj.Uvs[i][2] = t[i][2]
			obj.Uvs[i][3] = t[i][3]
		}
	case [][4]uint32:
		obj.Uvs = t
	}
	return nil
}

type Header struct {
	Signature   [12]uint8
	HeaderBytes uint8

	VertexCoordinateBytes uint8
	NormalCoordinateBytes uint8
	UVCoordinateBytes     uint8

	VertexIndexBytes   uint8
	NormalIndexBytes   uint8
	UVIndexBytes       uint8
	MaterialIndexBytes uint8

	VerticeCount uint32
	NormalCount  uint32
	UVCount      uint32

	TriFlatCount     uint32
	TriSmoothCount   uint32
	TriFlatUVCount   uint32
	TriSmoothUVCount uint32

	QuadFlatCount     uint32
	QuadSmoothCount   uint32
	QuadFlatUVCount   uint32
	QuadSmoothUVCount uint32
}

func (h *Header) SetDefault() {
	for i := range h.Signature {
		h.Signature[i] = JS3_BIN_SIGN[i]
	}
	h.HeaderBytes = 64
	h.VertexCoordinateBytes = 4
	h.NormalCoordinateBytes = 1
	h.UVCoordinateBytes = 4

	h.VertexIndexBytes = 4
	h.NormalIndexBytes = 4
	h.UVIndexBytes = 4
	h.MaterialIndexBytes = 2

	h.VerticeCount = 0
	h.NormalCount = 0
	h.UVCount = 0

	h.TriFlatCount = 0
	h.TriSmoothCount = 0
	h.TriFlatUVCount = 0
	h.TriSmoothUVCount = 0

	h.QuadFlatCount = 0
	h.QuadSmoothCount = 0
	h.QuadFlatUVCount = 0
	h.QuadSmoothUVCount = 0
}

func (h *Header) setup(obj *Binobj) {
	h.SetDefault()

	h.VerticeCount = uint32(len(obj.Vectilers))
	h.NormalCount = uint32(len(obj.Normals))
	h.UVCount = uint32(len(obj.UVs))

	h.TriFlatCount = uint32(len(obj.FlatTriangle.Material))
	h.TriSmoothCount = uint32(len(obj.SmoothTriangle.Material))
	h.TriFlatUVCount = uint32(len(obj.FlatUVTriangle.Material))
	h.TriSmoothUVCount = uint32(len(obj.SmoothUVTriangle.Material))

	h.QuadFlatCount = uint32(len(obj.FlatQuad.Material))
	h.QuadSmoothCount = uint32(len(obj.SmoothQuad.Material))
	h.QuadFlatUVCount = uint32(len(obj.FlatUVQuad.Material))
	h.QuadSmoothUVCount = uint32(len(obj.SmoothUVQuad.Material))
}

type Binobj struct {
	Header           Header
	Vectilers        [][3]float32
	Normals          [][3]int8
	UVs              [][2]float32
	FlatTriangle     FlatTriangle
	SmoothTriangle   SmoothTriangle
	FlatUVTriangle   FlatUVTriangle
	SmoothUVTriangle SmoothUVTriangle
	FlatQuad         FlatQuad
	SmoothQuad       SmoothQuad
	FlatUVQuad       FlatUVQuad
	SmoothUVQuad     SmoothUVQuad
}

func (obj *Binobj) Setup() {
	obj.Header.setup(obj)
}

func (obj *Binobj) GetVectilers() []float32 {
	ret := make([]float32, len(obj.Vectilers)*3)
	for i := range obj.Vectilers {
		ret[i*3] = obj.Vectilers[i][0]
		ret[i*3+1] = obj.Vectilers[i][1]
		ret[i*3+2] = obj.Vectilers[i][2]
	}
	return ret
}

func (obj *Binobj) GetNormals() []int8 {
	ret := make([]int8, len(obj.Normals)*3)
	for i := range obj.Normals {
		ret[i*3] = obj.Normals[i][0]
		ret[i*3+1] = obj.Normals[i][1]
		ret[i*3+2] = obj.Normals[i][2]
	}
	return ret
}

func (obj *Binobj) GetUVs() []float32 {
	ret := make([]float32, len(obj.UVs)*2)
	for i := range obj.UVs {
		ret[i*2] = obj.UVs[i][0]
		ret[i*2+1] = obj.UVs[i][1]
	}
	return ret
}

func (obj *Binobj) SetVectilers(vers interface{}) error {
	switch t := vers.(type) {
	case []float32:
		if len(t)%3 != 0 {
			return errors.New("[]float32 must 3^")
		}
		obj.Vectilers = make([][3]float32, len(t)/3)
		for i := range obj.Vectilers {
			obj.Vectilers[i][0] = t[i*3]
			obj.Vectilers[i][1] = t[i*3+1]
			obj.Vectilers[i][2] = t[i*3+2]
		}
	case [][]float32:
		obj.Vectilers = make([][3]float32, len(t))
		for i := range obj.Vectilers {
			obj.Vectilers[i][0] = t[i][0]
			obj.Vectilers[i][1] = t[i][1]
			obj.Vectilers[i][2] = t[i][2]
		}
	case [][3]float32:
		obj.Vectilers = t
	}
	return nil
}

func (obj *Binobj) SetNormals(norm interface{}) error {
	switch t := norm.(type) {
	case []int8:
		if len(t)%3 != 0 {
			return errors.New("[]int8 must 3^")
		}
		obj.Normals = make([][3]int8, len(t)/3)
		for i := range obj.Normals {
			obj.Normals[i][0] = int8(t[i*3])
			obj.Normals[i][1] = int8(t[i*3+1])
			obj.Normals[i][2] = int8(t[i*3+2])
		}
	case [][]int8:
		obj.Normals = make([][3]int8, len(t))
		for i := range obj.Normals {
			obj.Normals[i][0] = int8(t[i][0])
			obj.Normals[i][1] = int8(t[i][1])
			obj.Normals[i][2] = int8(t[i][2])
		}
	case [][3]int8:
		obj.Normals = t
	}
	return nil
}

func (obj *Binobj) SetUVs(vers interface{}) error {
	switch t := vers.(type) {
	case []float32:
		if len(t)%2 != 0 {
			return errors.New("[]float32 must 2^")
		}
		obj.UVs = make([][2]float32, len(t)/2)
		for i := range obj.UVs {
			obj.UVs[i][0] = t[i*2]
			obj.UVs[i][1] = t[i*2+1]
		}
	case [][]float32:
		obj.UVs = make([][2]float32, len(t))
		for i := range obj.UVs {
			obj.UVs[i][0] = t[i][0]
			obj.UVs[i][1] = t[i][1]
		}
	case [][2]float32:
		obj.UVs = t
	}
	return nil
}

func handlePadding(n uint32) uint32 {
	if n%4 > 0 {
		return 4 - n%4
	}
	return 0
}

func Decode(rd io.ReadSeeker) (*Binobj, error) {
	obj := new(Binobj)
	err := binary.Read(rd, littleEndian, &obj.Header)
	if err != nil {
		return nil, err
	}
	if string(obj.Header.Signature[:5]) != JS3_BIN_SIGN[:5] {
		return nil, errors.New("file not Three.js bin")
	}

	verbuf := make([]float32, obj.Header.VerticeCount*3)
	err = binary.Read(rd, littleEndian, verbuf)
	if err != nil {
		return nil, err
	}
	err = obj.SetVectilers(verbuf)
	if err != nil {
		return nil, err
	}

	norbuf := make([]int8, obj.Header.NormalCount*3)
	err = binary.Read(rd, littleEndian, norbuf)
	if err != nil {
		return nil, err
	}
	err = obj.SetNormals(norbuf)
	if err != nil {
		return nil, err
	}

	if obj.Header.NormalCount > 0 {
		pading := handlePadding(obj.Header.NormalCount * 3)
		if pading > 0 {
			rd.Seek(int64(pading), io.SeekCurrent)
		}
	}

	uvbuf := make([]float32, obj.Header.UVCount*2)
	err = binary.Read(rd, littleEndian, uvbuf)
	if err != nil {
		return nil, err
	}
	err = obj.SetUVs(uvbuf)
	if err != nil {
		return nil, err
	}

	if obj.Header.TriFlatCount > 0 {
		err = obj.FlatTriangle.decode(rd, obj.Header.TriFlatCount)
		if err != nil {
			return nil, err
		}

		pading := handlePadding(obj.Header.TriFlatCount * 2)
		if pading > 0 {
			rd.Seek(int64(pading), io.SeekCurrent)
		}
	}

	if obj.Header.TriSmoothCount > 0 {
		err = obj.SmoothTriangle.decode(rd, obj.Header.TriSmoothCount)
		if err != nil {
			return nil, err
		}

		pading := handlePadding(obj.Header.TriSmoothCount * 2)
		if pading > 0 {
			rd.Seek(int64(pading), io.SeekCurrent)
		}
	}

	if obj.Header.TriFlatUVCount > 0 {
		err = obj.FlatUVTriangle.decode(rd, obj.Header.TriFlatUVCount)
		if err != nil {
			return nil, err
		}

		pading := handlePadding(obj.Header.TriFlatUVCount * 2)
		if pading > 0 {
			rd.Seek(int64(pading), io.SeekCurrent)
		}
	}

	if obj.Header.TriSmoothUVCount > 0 {
		err = obj.SmoothUVTriangle.decode(rd, obj.Header.TriSmoothUVCount)
		if err != nil {
			return nil, err
		}

		pading := handlePadding(obj.Header.TriSmoothUVCount * 2)
		if pading > 0 {
			rd.Seek(int64(pading), io.SeekCurrent)
		}
	}

	if obj.Header.QuadFlatCount > 0 {
		err = obj.FlatQuad.decode(rd, obj.Header.QuadFlatCount)
		if err != nil {
			return nil, err
		}

		pading := handlePadding(obj.Header.QuadFlatCount * 2)
		if pading > 0 {
			rd.Seek(int64(pading), io.SeekCurrent)
		}
	}

	if obj.Header.QuadSmoothCount > 0 {
		err = obj.SmoothQuad.decode(rd, obj.Header.QuadSmoothCount)
		if err != nil {
			return nil, err
		}

		pading := handlePadding(obj.Header.QuadSmoothCount * 2)
		if pading > 0 {
			rd.Seek(int64(pading), io.SeekCurrent)
		}
	}

	if obj.Header.QuadFlatUVCount > 0 {
		err = obj.FlatUVQuad.decode(rd, obj.Header.QuadFlatUVCount)
		if err != nil {
			return nil, err
		}

		pading := handlePadding(obj.Header.QuadFlatUVCount * 2)
		if pading > 0 {
			rd.Seek(int64(pading), io.SeekCurrent)
		}
	}

	if obj.Header.QuadSmoothUVCount > 0 {
		err = obj.SmoothUVQuad.decode(rd, obj.Header.QuadSmoothUVCount)
		if err != nil {
			return nil, err
		}
	}
	return obj, nil
}

func writePading(wr io.Writer, padding uint32) error {
	bytes := make([]byte, padding)
	for i := uint32(0); i < padding; i++ {
		bytes = append(bytes, '0')
	}
	_, err := wr.Write(bytes)
	return err
}

func Encode(wr io.Writer, obj *Binobj) error {
	if string(obj.Header.Signature[:]) != JS3_BIN_SIGN {
		obj.Setup()
	}
	err := binary.Write(wr, littleEndian, obj.Header)
	if err != nil {
		return err
	}

	err = binary.Write(wr, littleEndian, obj.GetVectilers())
	if err != nil {
		return err
	}

	err = binary.Write(wr, littleEndian, obj.GetNormals())
	if err != nil {
		return err
	}

	if obj.Header.NormalCount > 0 {
		pading := handlePadding(obj.Header.NormalCount * 3)
		if pading > 0 {
			err = writePading(wr, pading)
			if err != nil {
				return err
			}
		}
	}

	err = binary.Write(wr, littleEndian, obj.GetUVs())
	if err != nil {
		return err
	}

	if obj.Header.TriFlatCount > 0 {
		err = obj.FlatTriangle.encode(wr)
		if err != nil {
			return err
		}

		pading := handlePadding(obj.Header.TriFlatCount * 2)
		if pading > 0 {
			err = writePading(wr, pading)
			if err != nil {
				return err
			}
		}
	}

	if obj.Header.TriSmoothCount > 0 {
		err = obj.SmoothTriangle.encode(wr)
		if err != nil {
			return err
		}

		pading := handlePadding(obj.Header.TriSmoothCount * 2)
		if pading > 0 {
			err = writePading(wr, pading)
			if err != nil {
				return err
			}
		}
	}

	if obj.Header.TriFlatUVCount > 0 {

		err = obj.FlatUVTriangle.encode(wr)
		if err != nil {
			return err
		}

		pading := handlePadding(obj.Header.TriFlatUVCount * 2)
		if pading > 0 {
			err = writePading(wr, pading)
			if err != nil {
				return err
			}
		}
	}

	if obj.Header.TriSmoothUVCount > 0 {

		err = obj.SmoothUVTriangle.encode(wr)
		if err != nil {
			return err
		}

		pading := handlePadding(obj.Header.TriSmoothUVCount * 2)
		if pading > 0 {
			err = writePading(wr, pading)
			if err != nil {
				return err
			}
		}
	}

	if obj.Header.QuadFlatCount > 0 {

		err = obj.FlatQuad.encode(wr)
		if err != nil {
			return err
		}

		pading := handlePadding(obj.Header.QuadFlatCount * 2)
		if pading > 0 {
			err = writePading(wr, pading)
			if err != nil {
				return err
			}
		}
	}

	if obj.Header.QuadSmoothCount > 0 {

		err = obj.SmoothQuad.encode(wr)
		if err != nil {
			return err
		}

		pading := handlePadding(obj.Header.QuadSmoothCount * 2)
		if pading > 0 {
			err = writePading(wr, pading)
			if err != nil {
				return err
			}
		}
	}

	if obj.Header.QuadFlatUVCount > 0 {

		err = obj.FlatUVQuad.encode(wr)
		if err != nil {
			return err
		}

		pading := handlePadding(obj.Header.QuadFlatUVCount * 2)
		if pading > 0 {
			err = writePading(wr, pading)
			if err != nil {
				return err
			}
		}
	}

	if obj.Header.QuadSmoothUVCount > 0 {

		err = obj.SmoothUVQuad.encode(wr)
		if err != nil {
			return err
		}

		pading := handlePadding(obj.Header.QuadSmoothUVCount * 2)
		if pading > 0 {
			err = writePading(wr, pading)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type Metadata struct {
	Version      float64 `json:"formatVersion"`
	Source       string  `json:"sourceFile"`
	GeneratedBy  string  `json:"generatedBy"`
	VerticeCount uint32  `json:"vertices"`
	FaceCount    uint32  `json:"faces"`
	NormalCount  uint32  `json:"normals"`
	ColorsCount  uint32  `json:"colors"`
	UVCount      uint32  `json:"uvs"`
	Materials    uint32  `json:"materials"`
}

type ThreeJSObj struct {
	Metadata  Metadata   `json:"metadata"`
	Materials []Material `json:"materials"`
	BinBuffer string     `json:"buffers"`
	Topology  Topology   `json:"topology,omitempty"`
}

func (ts *ThreeJSObj) ToJson() string {
	b, _ := json.Marshal(ts)
	return string(b)
}

func ThreeJSObjFromJson(data io.Reader) (*ThreeJSObj, error) {
	var ts *ThreeJSObj
	err := json.NewDecoder(data).Decode(&ts)
	if err != nil {
		return nil, err
	}
	return ts, nil
}
