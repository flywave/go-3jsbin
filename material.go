package bin

import (
	"math/rand"
)

var keyMap = map[string]string{
	"kd":          "colorDiffuse",   // Diffuse color: Kd 1.000 1.000 1.000
	"ka":          "colorAmbient",   // Ambient color: Ka 1.000 1.000 1.000 disabled for THREE.js > r80
	"ke":          "colorEmissive",  // Emissive color: Ke 0.000 0.000 0.000
	"ks":          "colorSpecular",  // Specular color: Ks 1.000 1.000 1.000
	"ns":          "specularCoef",   // Specular coefficient: Ns 154.000
	"tr":          "opacity",        // Transparency: Tr 0.9 or d 0.9
	"d":           "opacity",        // Transparency: Tr 0.9 or d 0.9
	"ni":          "opticalDensity", // Optical density: Ni 1.0
	"map_kd":      "mapDiffuse",     // Diffuse texture: map_Kd texture_diffuse.jpg
	"map_ka":      "mapAmbient",     // Ambient texture: map_Ka texture_ambient.jpg
	"map_ke":      "mapEmissive",    // Emissive texture: map_Ka texture_ambient.jpg
	"map_ks":      "mapSpecular",    // Specular texture: map_Ks texture_specular.jpg
	"map_ns":      "mapSpecular",    // Specular texture: map_Ns texture_specular.jpg
	"map_d":       "mapAlpha",       // Alpha texture: map_d texture_alpha.png
	"map_opacity": "mapAlpha",       // alternate alias for alpha map...
	"map_bump":    "mapBump",        // Bump texture: map_bump texture_bump.jpg
	"bump":        "mapBump",        // or bump texture_bump.jpg
	"illum":       "illumination",   // * Reflection (check footnote)
	"refl":        "illumination",   // second alias
}

// dummy colors
var COLORS = []uint32{0xeeeeee, 0xee0000, 0x00ee00, 0x0000ee, 0xeeee00, 0x00eeee, 0xee00ee}

type Material struct {
	DbgName        string    `json:"DbgName"`
	DbgIndex       uint32    `json:"DbgIndex"`
	DbgColor       uint32    `json:"DbgColor"`
	ColorDiffuse   []float64 `json:"colorDiffuse,omitempty"`
	ColorAmbient   []float64 `json:"colorAmbient,omitempty"`
	ColorEmissive  []float64 `json:"colorEmissive,omitempty"`
	ColorSpecular  []float64 `json:"colorSpecular,omitempty"`
	SpecularCoef   float64   `json:"specularCoef"`
	Opacity        float64   `json:"opacity"`
	OpticalDensity float64   `json:"opticalDensity"`
	MapDiffuse     string    `json:"mapDiffuse,omitempty"`
	MapAmbient     string    `json:"mapAmbient,omitempty"`
	MapEmissive    string    `json:"mapEmissive,omitempty"`
	MapSpecular    string    `json:"mapSpecular,omitempty"`
	MapAlpha       string    `json:"mapAlpha,omitempty"`
	MapBump        string    `json:"mapBump,omitempty"`
	Illumination   uint32    `json:"illumination"`
}

func GenerateColor(i int) uint32 {
	if i < len(COLORS) {
		return COLORS[i]
	}
	return uint32(0xffffff * rand.Float32())
}
