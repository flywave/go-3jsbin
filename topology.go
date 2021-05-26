package bin

type Anchor struct {
	Normal []float64 `json:"normal"`
	Center []float64 `json:"center"`
	Unit   float64   `json:"unit"`
	Name   string    `json:"name"`
}

type Topology struct {
	Scale       float64   `json:"scale"`
	Rotation    []float64 `json:"rotation"`
	AnchorCount int       `json:"anchorcount"`
	Anchors     []Anchor  `json:"anchors"`
	Offset      []float64 `json:"offset"`
}
