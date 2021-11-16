package bin

import (
	"encoding/json"
	"strconv"
)

type Anchor struct {
	Normal []float64 `json:"normal"`
	Center []float64 `json:"center"`
	Unit   float64   `json:"unit"`
	Name   string    `json:"name"`
}

func (a *Anchor) UnmarshalJSON(bt []byte) error {
	var mp map[string]interface{}
	err := json.Unmarshal(bt, &mp)
	if err != nil {
		return err
	}

	if rt, ok := mp["normal"]; ok {
		nl := rt.([]interface{})
		for _, n := range nl {
			if n == nil {
				a.Normal = append(a.Normal, 0)

			} else {
				a.Normal = append(a.Normal, n.(float64))
			}
		}
	}

	if rt, ok := mp["center"]; ok {
		ct := rt.([]interface{})
		for _, n := range ct {
			if n == nil {
				a.Center = append(a.Center, 0)

			} else {
				a.Center = append(a.Center, n.(float64))
			}
		}
	}

	if v, ok := mp["name"]; ok {
		a.Name = v.(string)
	}

	if v, ok := mp["unit"]; ok {
		switch unit := v.(type) {
		case float64:
			a.Unit = unit
		case string:
			var err error
			a.Unit, err = strconv.ParseFloat(unit, 64)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type Topology struct {
	Scale       float64   `json:"scale"`
	Rotation    []float64 `json:"rotation"`
	AnchorCount int       `json:"anchorcount"`
	Anchors     []Anchor  `json:"anchors"`
	Offset      []float64 `json:"offset"`
}

func (a *Topology) UnmarshalJSON(bt []byte) error {
	var mp map[string]interface{}
	err := json.Unmarshal(bt, &mp)
	if err != nil {
		return err
	}
	if v, ok := mp["rotation"]; ok {
		nl := v.([]interface{})
		for _, n := range nl {
			if n == nil {
				a.Rotation = append(a.Rotation, 0)

			} else {
				a.Rotation = append(a.Rotation, n.(float64))
			}
		}
	}

	if v, ok := mp["offset"]; ok {
		off := v.([]interface{})
		for _, n := range off {
			if n == nil {
				a.Offset = append(a.Offset, 0)

			} else {
				a.Offset = append(a.Offset, n.(float64))
			}
		}
	}

	if v, ok := mp["anchorcount"]; ok {
		a.AnchorCount = int(v.(float64))
	}
	if v, ok := mp["anchors"]; ok {
		bts, _ := json.Marshal(v)
		json.Unmarshal(bts, &a.Anchors)
	}
	if v, ok := mp["scale"]; ok {
		switch sc := v.(type) {
		case float64:
			a.Scale = sc
		case string:
			var err error
			a.Scale, err = strconv.ParseFloat(sc, 64)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
