package bin

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	mst "github.com/flywave/go-mst"
	"github.com/flywave/go3d/mat4"
	"github.com/flywave/go3d/quaternion"
	"github.com/flywave/go3d/vec2"
	"github.com/flywave/go3d/vec3"
	"github.com/flywave/go3d/vec4"
	"golang.org/x/image/bmp"
)

func ThreejsBin2Mst(fpath string) (*mst.Mesh, error) {
	f, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}

	jsobj, err := ThreeJSObjFromJson(f)
	if err != nil {
		fmt.Println(err.Error())
	}

	mesh := mst.NewMesh()
	nd := &mst.MeshNode{}

	binpath, _ := filepath.Split(fpath)
	binpath = filepath.Join(binpath, jsobj.BinBuffer)
	bf, _ := os.Open(binpath)
	binobj, _ := Decode(bf)

	var sc float32 = 1
	rot := jsobj.Topology.Rotation
	var quat quaternion.T
	if len(rot) != 0 {
		quat = quaternion.FromVec4(&vec4.T{float32(rot[0]), float32(rot[1]), float32(rot[2]), float32(rot[3])})
	}

	if jsobj.Topology.Scale != 0 {
		sc = float32(jsobj.Topology.Scale)
	}
	off := vec3.T{}
	if len(jsobj.Topology.Offset) != 0 {
		of := jsobj.Topology.Offset
		off = vec3.T{float32(of[0]), float32(of[1]), float32(of[2])}
	}
	mat := mat4.Compose(&off, &quat, &vec3.T{sc, sc, sc})
	for i := range binobj.Vectilers {
		v := (*vec3.T)(&binobj.Vectilers[i])
		v1 := mat.MulVec3(v)
		nd.Vertices = append(nd.Vertices, v1)
	}

	if len(binobj.Normals) > 0 {
		for i := range binobj.Normals {
			nl := &vec3.T{float32(binobj.Normals[i][0] / 127), float32(binobj.Normals[i][1] / 127), float32(binobj.Normals[i][2] / 127)}
			if nl.IsZero() {
				nl[2] = 1
			}
			nd.Normals = append(nd.Normals, *nl)
		}
	}

	if len(binobj.UVs) > 0 {
		for i := range binobj.UVs {
			uv := (*vec2.T)(&binobj.UVs[i])
			nd.TexCoords = append(nd.TexCoords, *uv)
		}
	}

	mtlcount := len(jsobj.Materials)
	nd.FaceGroup = make([]*mst.MeshTriangle, mtlcount)
	for i := range nd.FaceGroup {
		g := &mst.MeshTriangle{}
		g.Batchid = int32(i)
		nd.FaceGroup[i] = g
	}
	if binobj.Header.TriFlatCount > 0 {
		mtls := binobj.FlatTriangle.Material
		for i, id := range mtls {
			g := nd.FaceGroup[int(id)]
			f := &mst.Face{
				Vertex: binobj.FlatTriangle.Vertices[i],
			}
			g.Faces = append(g.Faces, f)
		}
	}
	if binobj.Header.TriFlatUVCount > 0 {
		mtls := binobj.FlatUVTriangle.Material
		for i, id := range mtls {
			g := nd.FaceGroup[int(id)]
			f := &mst.Face{
				Vertex: binobj.FlatUVTriangle.Vertices[i],
				Uv:     &binobj.FlatUVTriangle.Uvs[i],
			}
			g.Faces = append(g.Faces, f)
		}
	}

	if binobj.Header.TriSmoothCount > 0 {
		mtls := binobj.SmoothTriangle.Material
		for i, id := range mtls {
			g := nd.FaceGroup[int(id)]
			f := &mst.Face{
				Vertex: binobj.SmoothTriangle.Vertices[i],
				Normal: &binobj.SmoothTriangle.Normals[i],
			}
			g.Faces = append(g.Faces, f)
		}
	}

	if binobj.Header.TriSmoothUVCount > 0 {
		mtls := binobj.SmoothUVTriangle.Material
		for i, id := range mtls {
			g := nd.FaceGroup[int(id)]
			f := &mst.Face{
				Vertex: binobj.SmoothUVTriangle.Vertices[i],
				Uv:     &binobj.SmoothUVTriangle.Uvs[i],
				Normal: &binobj.SmoothUVTriangle.Normals[i],
			}
			g.Faces = append(g.Faces, f)
		}
	}

	if binobj.Header.QuadFlatCount > 0 {
		mtls := binobj.FlatQuad.Material
		for i, id := range mtls {
			vt := binobj.FlatQuad.Vertices[i]

			g := nd.FaceGroup[int(id)]
			f := &mst.Face{
				Vertex: [3]uint32{vt[0], vt[1], vt[2]},
			}
			g.Faces = append(g.Faces, f)
			f = &mst.Face{
				Vertex: [3]uint32{vt[2], vt[3], vt[0]},
			}
			g.Faces = append(g.Faces, f)
		}
	}

	if binobj.Header.QuadFlatUVCount > 0 {
		mtls := binobj.FlatUVQuad.Material
		for i, id := range mtls {
			vt := binobj.FlatUVQuad.Vertices[i]
			uv := binobj.FlatUVQuad.Uvs[i]

			g := nd.FaceGroup[int(id)]
			f := &mst.Face{
				Vertex: [3]uint32{vt[0], vt[1], vt[2]},
				Uv:     &[3]uint32{uv[0], uv[1], uv[2]},
			}
			g.Faces = append(g.Faces, f)
			f = &mst.Face{
				Vertex: [3]uint32{vt[2], vt[3], vt[0]},
				Uv:     &[3]uint32{uv[2], uv[3], uv[0]},
			}
			g.Faces = append(g.Faces, f)
		}
	}

	if binobj.Header.QuadSmoothCount > 0 {
		mtls := binobj.SmoothQuad.Material
		for i, id := range mtls {
			vt := binobj.SmoothQuad.Vertices[i]
			nl := binobj.SmoothQuad.Normals[i]

			g := nd.FaceGroup[int(id)]
			f := &mst.Face{
				Vertex: [3]uint32{vt[0], vt[1], vt[2]},
				Normal: &[3]uint32{nl[0], nl[1], nl[2]},
			}
			g.Faces = append(g.Faces, f)
			f = &mst.Face{
				Vertex: [3]uint32{vt[2], vt[3], vt[0]},
				Normal: &[3]uint32{nl[2], nl[3], nl[0]},
			}
			g.Faces = append(g.Faces, f)
		}
	}

	if binobj.Header.QuadSmoothUVCount > 0 {
		mtls := binobj.SmoothUVQuad.Material
		for i, id := range mtls {
			vt := binobj.SmoothUVQuad.Vertices[i]
			nl := binobj.SmoothUVQuad.Normals[i]
			uv := binobj.SmoothUVQuad.Uvs[i]

			g := nd.FaceGroup[int(id)]
			f := &mst.Face{
				Vertex: [3]uint32{vt[0], vt[1], vt[2]},
				Normal: &[3]uint32{nl[0], nl[1], nl[2]},
				Uv:     &[3]uint32{uv[0], uv[1], uv[2]},
			}
			g.Faces = append(g.Faces, f)
			f = &mst.Face{
				Vertex: [3]uint32{vt[2], vt[3], vt[0]},
				Normal: &[3]uint32{nl[2], nl[3], nl[0]},
				Uv:     &[3]uint32{uv[2], uv[3], uv[0]},
			}
			g.Faces = append(g.Faces, f)
		}
	}

	for id, mtl := range jsobj.Materials {
		ml := &mst.PbrMaterial{Roughness: 1, Metallic: 0}
		if len(mtl.ColorDiffuse) != 0 {
			ml.Color[0] = byte(mtl.ColorDiffuse[0] * 255.0)
			ml.Color[1] = byte(mtl.ColorDiffuse[1] * 255.0)
			ml.Color[2] = byte(mtl.ColorDiffuse[2] * 255.0)
		}
		ml.Transparency = 1 - float32(mtl.Opacity)
		if ml.Transparency == 1 {
			ml.Transparency = 0
		}

		if mtl.MapDiffuse != "" {
			dir, _ := filepath.Split(fpath)
			var ap *string
			if mtl.MapAlpha != "" {
				ph := filepath.Join(dir, mtl.MapAlpha)
				ap = &ph
			}
			tex, err := convertTex(filepath.Join(dir, mtl.MapDiffuse), ap, id)
			if err == nil {
				ml.Texture = tex
			}
		}

		mesh.Materials = append(mesh.Materials, ml)
	}

	nd.ResortVtVn(mesh)
	// nd.ReComputeNormal()
	mesh.Nodes = append(mesh.Nodes, nd)
	return mesh, nil
}

func readDir(root, path string, ext_filter []string) ([]string, error) {
	root = filepath.Clean(root)
	path = filepath.Clean(path)
	res := []string{}
	fs, er := os.ReadDir(path)
	if er != nil {
		return res, er
	}
	for _, info := range fs {
		ph := filepath.Join(path, info.Name())
		if info.IsDir() {
			ls, err := readDir(root, ph, ext_filter)
			if err != nil {
				return nil, err
			}
			res = append(res, ls...)
		} else {
			if len(ext_filter) > 0 {
				for _, ext := range ext_filter {
					et := filepath.Ext(ph)
					if et == ext {
						res = append(res, strings.Replace(ph, root, "", 1))
						break
					}
				}
			} else {
				res = append(res, strings.Replace(ph, root, "", 1))
			}
		}
	}
	return res, nil
}

func convertTex(path string, alphPh *string, texId int) (*mst.Texture, error) {
	img1, err := readImageByPath(path)
	if err != nil {
		return nil, err
	}
	var img2 image.Image
	if alphPh != nil {
		img2, err = readImageByPath(*alphPh)
	}
	if err != nil {
		return nil, err
	}
	bd := img1.Bounds()
	buf := []byte{}
	for y := 0; y < bd.Dy(); y++ {
		for x := 0; x < bd.Dx(); x++ {
			cl := img1.At(x, y)
			r, g, b, a := color.RGBAModel.Convert(cl).RGBA()
			sc := float32(1)
			if img2 != nil {
				cl2 := img2.At(x, y)
				r, _, _, _ := color.RGBAModel.Convert(cl2).RGBA()
				sc = float32(r&0xff) / 255
			}
			buf = append(buf, byte(r&0xff), byte(g&0xff), byte(b&0xff), byte(float32(a&0xff)*sc))
		}
	}
	_, name := filepath.Split(path)
	name = strings.Replace(name, ".jpg", ".png", 1)
	name = strings.Replace(name, ".jpeg", ".png", 1)

	t := &mst.Texture{}
	t.Id = int32(texId)
	t.Name = name
	t.Format = mst.TEXTURE_FORMAT_RGBA
	t.Size = [2]uint64{uint64(bd.Dx()), uint64(bd.Dy())}
	t.Compressed = mst.TEXTURE_COMPRESSED_ZLIB
	t.Data = mst.CompressImage(buf)
	return t, nil
}

func readImageByPath(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	_, ft, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	f.Seek(0, 0)
	return readImage(f, ft)
}

func readImage(rd io.Reader, ft string) (image.Image, error) {
	switch ft {
	case "jpeg", "jpg":
		return jpeg.Decode(rd)
	case "png":
		return png.Decode(rd)
	case "gif":
		return gif.Decode(rd)
	case "bmp":
		return bmp.Decode(rd)
	default:
		return nil, errors.New("unknow format")
	}
}
