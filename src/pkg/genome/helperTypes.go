//Helper types
package genome

import (
	"io"
	"os"
)
//////////////////////////////////////////////////////////////////////////
type Adder interface {
	Add(genome io.ReadWriteSeeker, gene Gene, geneMap GeneMap) (pos int64, err os.Error)
}

type AdderFunc func(io.ReadWriteSeeker, Gene, GeneMap) (pos int64, err os.Error)

func (a AdderFunc) Add(genome io.ReadWriteSeeker, gene Gene, geneMap GeneMap) (pos int64, err os.Error) {
	return a(genome, gene, geneMap)
}
///////////////////////////////////////////////////////////////////////
type Trainer interface {
	Train(trainingData []byte)
	GetGenes() (genes [][]byte)
}
//////////////////////////////////////////////////////////////////////
type Indexer interface {
	Index(gene Gene) (index []byte)
}

type IndexerFunc func(Gene) []byte

func (i IndexerFunc) Index(gene Gene) (index []byte) {
	return i(gene)
}

////////////////////////////////////////////////////////////////////
type Finder interface {
	Find(data io.ReadSeeker, genomeRW io.ReadSeeker,theGenome Genome ,geneMap GeneMap, indexer Indexer) (geneAddress []GeneAddress, unknownGenes []byte, ordering []bool, size int64)
}

type FinderFunc func(io.ReadSeeker, io.ReadSeeker,Genome, GeneMap, Indexer) ([]GeneAddress, []byte, []bool, int64)

func (f FinderFunc) Find(data io.ReadSeeker, genomeRW io.ReadSeeker,theGenome Genome ,geneMap GeneMap, indexer Indexer) (geneAddress []GeneAddress, unknownGenes []byte, ordering []bool, size int64) {
	return f(data, genomeRW, theGenome ,geneMap, indexer)
}
