//gg:{genome.a:gene.go geneMap.go geneAddress.go tGene.go helperTypes.go genome.go }

package genome
///////////////////
///////TODO: Need to split this file up
/////////////////
import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"os"
	"bytes"
)

//The Genome class
//This will be the publicly accessable interface to the datazome project

type Genome struct {
	geneMap   GeneMap //this is a map into the file of where all the genes are needs to be saved every time the genome is closed
	genomeRW  io.ReadWriteSeeker
	geneMapRW io.ReadWriteSeeker
	indexer   Indexer
	finder    Finder
	trainer   Trainer
	adder     Adder
}

//
//Need to make a simple Genome interface so that you dont *HAVE* to call a function that takes 6 params
//
func init() {}
//////GENOME IO///////////////////////////////////////////////////////////////
//create a genome that is saved to the writer
//this should be the default func for a new user to call
//The indexing service to be used , the trainer to be used, and the geneAdder to be used when genes are being added to the system
func NewGenome(genomeData, geneMapData io.ReadWriteSeeker, indexer Indexer, trainer Trainer, finder Finder, adder Adder) (genome *Genome, err os.Error) {
	genome = new(Genome)
	genome.geneMap = make(GeneMap, 10000000)
	genome.geneMapRW = geneMapData
	genome.genomeRW = genomeData
	genome.indexer = indexer
	genome.trainer = trainer
	genome.finder = finder
	genome.adder = adder
	err = nil
	return
	//may need to make a new map for the Genome struct
}

//loads a genome from a file
//this should be the function that a returning user uses
func LoadGenome(genomeData, geneMapData io.ReadWriteSeeker, indexer Indexer, trainer Trainer, finder Finder, adder Adder) (genome *Genome, err os.Error) {
	genome, _ = NewGenome(genomeData, geneMapData, indexer, trainer, finder, adder)
	genome.geneMap, _ = LoadGeneMap(geneMapData)
	//err = decoder.Decode(genome.geneMap)
	return
}

//should be called at the end of use
//saves all changes to the genome file
func (genome *Genome) SaveGenome() (err os.Error) {
	_, err = genome.geneMapRW.Seek(0, 0)
	if err != nil {
		return
	}

	genome.geneMap.Save(genome.geneMapRW)
	if err != nil {
		return
	}
	return
}
/////END GENOMEIO///////////////////////////////////////////////////////////////////////

//////GENOME TRAINING/////////////////////////////////////////////////////////////////
//Trains a genome according to the contents of a directory
func (genome *Genome) TrainOnDir(dir string) (err os.Error) {
	//open the directory
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	fmt.Println("analyzing", len(files), "files")
	for fileNum, fileDesc := range (files) {
		if fileDesc.IsRegular() {
			file, err := os.Open(strings.Join([]string{dir, fileDesc.Name}, "/"), os.O_RDONLY, 0666)
			if err != nil {
				fmt.Println("oops!!", err)
			}
			fReader := io.Reader(file)
			data, err := ioutil.ReadAll(fReader)
			genome.trainer.Train(data)
			fmt.Println("on file", fileNum)
			file.Close()
		}
	}
	genes := genome.trainer.GetGenes()
	for _, gene := range (genes) {
		genome.addGene(Gene(gene))
	}
	return
}

//Trains a Genome to a data
// This functoin will need to be changed to use a stream to handle larger data piceces
func (genome *Genome) TrainOnData(data []byte) (err os.Error) {
	genome.trainer.Train(data)
	//genome.addGene(genome.trainer.GetGenes()[0]) //this is odd i am not really sure why this is the way it is
	err = nil
	return
}

func (genome *Genome) UpdateGenome() (err os.Error) {
	for _, gene := range (genome.trainer.GetGenes()) {
		genome.addGene(gene)
	}
	err = nil
	return
}
//adds a gene using
func (genome *Genome) addGene(gene Gene) (address GeneAddress) {
	//-fmt.Print("a")
	index := genome.indexer.Index(gene) //rember the index is a hash of some sort
	buf := bytes.NewBuffer(index)
	str := buf.String()
	//try to git rid of redunency
	if len(gene) < 5 {
		return GeneAddress{}
	}
	_, present := genome.geneMap[str]
	if present {
		return genome.geneMap[str]
	}
	//call the adder helper to handle the file operations
	loc, _ := genome.adder.Add(genome.genomeRW, gene, genome.geneMap)
	genome.geneMap[str] = NewGeneAddress(loc, int64(len(gene)))
	return NewGeneAddress(loc, int64(len(gene)))

}

func (genome *Genome) AddGene(gene Gene) (address GeneAddress) {
	return genome.addGene(gene)
}

func (genome *Genome) getGene(address GeneAddress) (gene Gene, err os.Error) {
	//should add in some code to use the Addr as an offset
	if address.Hash != "" {
		temAddr, exsits := genome.geneMap[address.Hash]
		if exsits {
			address = temAddr
		}
	}
	genome.genomeRW.Seek(address.Addr, 0)
	data := make([]byte, address.Length)
	genome.genomeRW.Read(data)
	if err != nil {
		panic("ohh dear", err.String())
	}
	gene = data
	return
}

/* to do ... Get it right the first time damn it
func (genome *Genome) removeGene(address GeneAddress)(err os.Error){

}
*/

///////////////////////////////////////////////////////////////////////////////////
func (genome *Genome) NewTGene(data io.ReadSeeker, name string) (tGene TGene, err os.Error) {
	geneAddressSet, unGeneSet, ordering, size := genome.finder.Find(data, genome.genomeRW,*genome ,genome.geneMap,genome.indexer)
	tGene = NewTGene(geneAddressSet, unGeneSet, ordering, size)
	err = nil
	return
}

func (genome *Genome) RecreateData(tGene TGene) (data []byte, err os.Error) {
	//for every index in the tGene set
	genes := make([][]byte, len(tGene.Ordering))
	fmt.Println("ordering:", tGene.Ordering)
	fmt.Println("konwn Genes:", tGene.KnownGenes)
	unknownIndex := 0
	knownIndex := 0
	for index, dataIndex  := range (tGene.Ordering) {
		//unknonwn genes are stacked on to
		if dataIndex == false{
			genes[index] = Gene(tGene.UnknownGenes[unknownIndex:unknownIndex+1])
			unknownIndex++
		} else {
			genes[index], _ = genome.getGene(tGene.KnownGenes[knownIndex])
			knownIndex++
		}
	}
	data = bytes.Join(genes, []byte{})
	err = nil
	return
}
