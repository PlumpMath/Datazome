package trainers

import (
	"genome"
	"fmt"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type PossibleGene struct {
	dataValue       []byte
	occuranceNumber int
}

type BruteTrainer struct {
	possibleGenes                []PossibleGene
	minChunkBound, maxChunkBound int
}

func NewBruteTrainer() *BruteTrainer {
	bt := new(BruteTrainer)
	tempArray := new([900000000]PossibleGene)
	bt.possibleGenes = tempArray[0:1]
	bt.minChunkBound = 5
	bt.maxChunkBound = 9
	return bt
}
//should consider changeing this to a io.Reader instead of a []byte to deal with huge data sets (a movie for example)
func (trainer *BruteTrainer) Train(trainingData []byte) {
	fmt.Println(len(trainingData))
	//from the minimum match size to the max match size
	for i := trainer.minChunkBound; i < trainer.maxChunkBound && i < len(trainingData); i++ {
		fmt.Print(".")
		//Start at the begging of the slice
		for index := 0; index < len(trainingData); index++ {
			//initilzie the is the var this will be used to determine if it is necessary to make a new entery in the possible gene list
			isThereGene := false
			//check each known possible gene
			//fmt.Println(len(trainer.possibleGenes))
			if (index % 1000) == 0 {
				fmt.Print(",")
			}
			for pindex := 0; pindex < len(trainer.possibleGenes); pindex += i {

				//fmt.Println(trainingData[index : index+i])
				//if there is a match to this known gene
				if len(trainer.possibleGenes[pindex].dataValue) == i {
					isThereGene = bytes.Equal(trainingData[index:index+i], trainer.possibleGenes[pindex].dataValue)
					if isThereGene {
						trainer.possibleGenes[pindex].occuranceNumber++

						//fmt.Print(trainer.possibleGenes[pindex])
						break
					}
				}

			}
			if isThereGene == false {
				//fmt.Println("in outside loop", trainingData[index:index+i])
				trainer.possibleGenes = trainer.possibleGenes[0 : len(trainer.possibleGenes)+1]
				trainer.possibleGenes[len(trainer.possibleGenes)-1].dataValue = trainingData[index : index+i]
				trainer.possibleGenes[len(trainer.possibleGenes)-1].occuranceNumber++
			}
		}
	}
	fmt.Println("current possible Gene len", len(trainer.possibleGenes))
	fmt.Println("current Highest occurance", trainer.getNth(2))
}
/* already declated in the outher brute trainer
func contains(i int, slice []int)(bool){
	for _,item := range(slice){
		if i == item { return false }
	}
	return true
}
*/
//
func (trainer *BruteTrainer) getNth(n int) PossibleGene {
	greatest := 0
	index := 0
	last := make([]int, n, 0)
	for i := 0; i < n; i++ {
		for curIndex, gene := range (trainer.possibleGenes) {
			if gene.occuranceNumber > greatest || contains(curIndex, last) {
				index = curIndex
				greatest = gene.occuranceNumber
			}
		}
		last = last[0 : i+1]
		last[i] = index
	}
	return trainer.possibleGenes[index]

}
func (trainer *BruteTrainer) TrainOnDir(dir string) (err os.Error) {
	//open the directory
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, fileDesc := range (files) {
		if fileDesc.IsRegular() {
			file, err := os.Open(strings.Join([]string{dir, fileDesc.Name}, "/"), os.O_RDONLY, 0666)
			if err != nil {
				fmt.Println("oops!!", err)
			}
			fReader := io.Reader(file)
			data, err := ioutil.ReadAll(fReader)
			trainer.Train(data)
		}
	}
	return
}
//func compareArray(bytes []byte, index, i int) (isThereGene bool) { }
func (trainer *BruteTrainer) TrainOnDirPart(dir string, howMany int) (err os.Error) {
	//open the directory
	files, err := ioutil.ReadDir(dir)
	fmt.Println("trainging on dir ! = finished file . = finished pass , = finished pass part")
	if err != nil {
		return err
	}
	for x := 0; x < howMany; x++ {
		fileDesc := files[x]
		if fileDesc.IsRegular() {
			file, err := os.Open(strings.Join([]string{dir, fileDesc.Name}, "/"), os.O_RDONLY, 0666)
			if err != nil {
				fmt.Println("oops!!", err)
			}
			fReader := io.Reader(file)
			data, err := ioutil.ReadAll(fReader)
			trainer.Train(data)
			fmt.Print("!")
		}
	}
	return
}
// needs to be implented
func (trainer *BruteTrainer) GetGenes() (genes []genome.Gene, err os.Error) {
	genes = nil
	err = nil
	return
}
/*
func OpenTrainer(data io.Reader) (err os.Error) {
	return
}

func Save(data io.Writer) (err os.Error) {}
*/
func (trainer *BruteTrainer) Output() {
	fmt.Print(len(trainer.possibleGenes))
	for x := 0; x < len(trainer.possibleGenes); x++ {
		fmt.Print(trainer.possibleGenes[x], "a")
	}
	return
}
