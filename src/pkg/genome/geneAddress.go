//gene
package genome

type GeneAddress struct {
	Addr   int64
	Length int64
	Hash   string
}

func NewGeneAddress(addr, length int64) (ga GeneAddress) {
	ga.Addr = addr
	ga.Length = length
	return
}

type IndexValuePair struct {
	Index Index
	Value GeneAddress
}
type Index []byte //this is the hashed index of a gene in the geneMap has to be a string to be compatable with maps
