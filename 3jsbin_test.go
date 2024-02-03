package bin

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	mst "github.com/flywave/go-mst"
)

func TestJsLoadFiles(t *testing.T) {
	f, err := os.Open("../../tests/binary/palm.js")
	defer f.Close()

	if err != nil {
		t.Error("error")
	}
	obj, err := ThreeJSObjFromJson(f)
	if err != nil {
		t.Error("error")
	}
	if obj != nil {
		t.Error("error")
	}
}

func TestBinLoadFiles(t *testing.T) {
	f, err := os.Open("../../tests/binary/palm.bin")
	defer f.Close()

	if err != nil {
		t.Error("error")
	}

	obj, err := Decode(f)
	if err != nil {
		t.Error("error")
	}
	if obj != nil {
		t.Error("error")
	}
}

const absPath = "/home/hj/workspace/GISCore/build/public/Resources/"

func TestToMst(t *testing.T) {
	dirs := []string{"model/thsk/thsk_sw_xj"} //"anchormodel"
	for _, d := range dirs {
		dr := absPath + d
		fs, _ := readDir(dr, dr, []string{".json"})
		for _, f := range fs {
			fpath := dr + f
			mstPh := strings.Replace(fpath, ".json", ".mst", 1)
			glbPh := strings.Replace(mstPh, ".mst", ".glb", 1)
			// if info, err := os.Stat(mstPh); err == nil {
			// 	t := time.Date(2021, time.November, 16, 12, 40, 0, 0, time.UTC)
			// 	if info.ModTime().After(t) {
			// 		continue
			// 	}
			// }
			mh, _ := ThreejsBin2Mst(fpath)
			for _, nd := range mh.Nodes {
				nd.ReComputeNormal()
			}
			doc := mst.CreateDoc()
			mst.BuildGltf(doc, mh, false)
			bt, _ := mst.GetGltfBinary(doc, 8)
			ioutil.WriteFile(glbPh, bt, os.ModePerm)

			wt := bytes.NewBuffer([]byte{})
			mst.MeshMarshal(wt, mh)
			ioutil.WriteFile(mstPh, wt.Bytes(), os.ModePerm)
		}
	}
}

func TestGltf(t *testing.T) {
	mesh, _ := ThreejsBin2Mst("/home/hj/workspace/GISCore/build/public/Resources/model/public/HHRQQiTiWoLunLiuLiangJi/HHRQQiTiWoLunLiuLiangJi.json")
	doc := mst.CreateDoc()
	mst.BuildGltf(doc, mesh, false)
	bt, _ := mst.GetGltfBinary(doc, 8)
	ioutil.WriteFile("/home/hj/workspace/GISCore/build/public/Resources/model/public/HHRQQiTiWoLunLiuLiangJi/HHRQQiTiWoLunLiuLiangJi.glb", bt, os.ModePerm)
}

func TestGltf2(t *testing.T) {
	mh, _ := ThreejsBin2Mst("/home/hj/workspace/GISCore/build/public/Resources/model/zbrl/ZBRL_BY/ZBRL_BY_1.json")

	doc := mst.CreateDoc()
	mst.BuildGltf(doc, mh, false)
	bt, _ := mst.GetGltfBinary(doc, 8)
	ioutil.WriteFile("/home/hj/workspace/GISCore/build/public/Resources/model/zbrl/ZBRL_BY/ZBRL_BY_1.glb", bt, os.ModePerm)
}

func TestBin2(t *testing.T) {
	mh, _ := ThreejsBin2Mst("/home/hj/workspace/GISCore/build/public/Resources/anchormodel/public/psqitong/psqitong.json")
	doc := mst.CreateDoc()
	mst.BuildGltf(doc, mh, false)
	bt, _ := mst.GetGltfBinary(doc, 8)
	ioutil.WriteFile("/home/hj/workspace/GISCore/build/public/Resources/anchormodel/public/psqitong/psqitong.glb", bt, os.ModePerm)

}
