package commDomain

type JacocoResult struct {
	Sessioninfo JacocoSessioninfo `xml:"sessioninfo"`
	Package     JacocoPackage     `xml:"package"`
	Counter     []JacocoCounter   `xml:"counter"`
	Name        string            `xml:"name,attr"`
}

type JacocoCounter struct {
	Type    Type   `xml:"type,attr"`
	Missed  string `xml:"missed,attr"`
	Covered string `xml:"covered,attr"`
}

type JacocoPackage struct {
	Class      JacocoClassClass `xml:"class"`
	Sourcefile JacocoSourcefile `xml:"sourcefile"`
	Counter    []JacocoCounter  `xml:"counter"`
	Name       string           `xml:"name,attr"`
}

type JacocoClassClass struct {
	Method         []JacocoMethodElement `xml:"method"`
	Counter        []JacocoCounter       `xml:"counter"`
	Name           string                `xml:"name,attr"`
	Sourcefilename string                `xml:"sourcefilename,attr"`
}

type JacocoMethodElement struct {
	Counter []JacocoCounter `xml:"counter"`
	Name    string          `xml:"name,attr"`
	Desc    string          `xml:"desc,attr"`
	Line    string          `xml:"line,attr"`
}

type JacocoSourcefile struct {
	Line    []JacocoLineElement `xml:"line"`
	Counter []JacocoCounter     `xml:"counter"`
	Name    string              `xml:"name,attr"`
}

type JacocoLineElement struct {
	Nr string `xml:"nr,attr"`
	Mi string `xml:"mi,attr"`
	Ci string `xml:"ci,attr"`
	MB string `xml:"mb,attr"`
	Cb string `xml:"cb,attr"`
}

type JacocoSessioninfo struct {
	ID    string `xml:"id,attr"`
	Start string `xml:"start,attr"`
	Dump  string `xml:"dump,attr"`
}

type Type string

const (
	Branch      Type = "BRANCH"
	Class       Type = "CLASS"
	Complexity  Type = "COMPLEXITY"
	Instruction Type = "INSTRUCTION"
	Line        Type = "LINE"
	Method      Type = "METHOD"
)
