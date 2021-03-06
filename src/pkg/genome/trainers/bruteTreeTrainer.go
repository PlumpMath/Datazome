package trainers

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sort"
	"bytes"
	//	"gob"

)

func init() { fmt.Print("Test") }

type GenePart struct {
	DataValue byte
	HitNumber uint16
	Next      *GeneNode
}
type GeneData struct {
	DataValue []byte
	HitNumber uint16
}
type GeneNode struct {
	DataValue byte
	Depth     uint8
	HitNumber uint32
	GeneParts [256]GenePart
}

type BruteTreeTrainer struct {
	root       *GeneNode
	curHighest []GeneData
}

func NewBruteTreeTrainer() *BruteTreeTrainer {
	bt := new(BruteTreeTrainer)
	bt.curHighest = make([]GeneData, 0)
	bt.root = newTree()
	return bt
}

func newTree() *GeneNode {
	root := new(GeneNode)
	root.Depth = 0
	root.DataValue = 0
	root.HitNumber = 0
	for index, part := range (root.GeneParts) {
		part.DataValue = uint8(index)
		part.HitNumber = 1
	}
	return root
}

const MinLin int = 6
const NextNodeBound int = 2
const MaxDepth int = 10

//creates a new geneNode wit b as it counter
func contains(i int, slice []int) bool {
	for _, cur := range (slice) {
		if cur == i {
			return false
		}
	}
	return true
}

func (gn *GeneNode) numLeafs() (num int) {
	tally := 0
	//for all of the decendent nodes
	for _, cur := range (gn.GeneParts) {
		//ask them how mnay leafs they have
		if cur.Next != nil {
			num += cur.Next.numLeafs()
			tally++
		}
	}
	//if this is a leaf return 1
	if tally == 0 {
		num = 1
	}
	return
}

func (gn *GeneNode) numChildern() (num int) {
	num = 0
	for _, cur := range (gn.GeneParts) {
		if cur.Next != nil {
			num++
		}
	}
	return
}

type GeneDataSlice []GeneData

//joins two GeneData Slices
func join(right, left []GeneData) []GeneData {
	//!fmt.Print("j")
	newGeneData := make([]GeneData, len(right)+len(left))
	for index, item := range (right) {
		newGeneData[index] = item
	}
	for index, item := range (left) {
		newGeneData[index+len(right)] = item
	}
	right = newGeneData
	return right
}

//prunes based on len
func shrink(slice []GeneData,minLen int) []GeneData{
	newCount := 0
	for _,item := range(slice){
		if len(item.DataValue) >= minLen{
			newCount++
		}
	}
	newSlice := make([]GeneData,newCount)
	curIndex:=0
	for _,item := range(slice){
		if len(item.DataValue) >= minLen{
			newSlice[curIndex] = item
			curIndex++
		}
	}
	return newSlice
}
	
	

func merge(right, left []GeneData) []GeneData {

	//find the size of the new array by finding the number of matchs
	fmt.Print("m")
	defer fmt.Print("M")
	dups := 0
	for _, item := range (left) {

		for _, itemr := range (right) {
			if bytes.Equal(item.DataValue, itemr.DataValue) {
				dups++
				break
			}
		}
	}
	newGeneData := make([]GeneData, len(right)+len(left)-dups)
	for index, item := range (right) {
		newGeneData[index] = item
	}
	tally := 0
	for _, item := range (left) {
		matched := false
		for indexr, itemr := range (right) {
			if bytes.Equal(item.DataValue, itemr.DataValue) {
				newGeneData[indexr].HitNumber += item.HitNumber
				matched = true
				break
			}
		}
		if !matched {
			newGeneData[len(right)+tally] = item
			tally++
		}
	}

	return newGeneData
}

func (gn *GeneNode) collapseLeafs() []GeneData {
	geneData := make([]GeneData, 0)
	//!fmt.Println("depth:",gn.Depth, "numChildern:", gn.numChildern(),"numLeafs:", gn.numLeafs())
	//if this is a leaf
	//make a slice of the depth that this leaf is at
	//put your data value at the end
	//the the nodes data value at the next from the end
	//return

	if gn.numChildern() == 0 {
		geneData = make([]GeneData, 1)
		geneData[0].HitNumber = uint16(gn.HitNumber)
		geneData[0].DataValue = make([]byte, gn.Depth)
		geneData[0].DataValue[gn.Depth-1] = gn.DataValue
		//!fmt.Println("    geneData:",geneData,"\n")
		return geneData
	}

	//if this is a brach
	//join all the leafs together
	//add this data value to depths place in all of the slice
	for _, child := range (gn.GeneParts) {
		//recurese this function down to its leaf
		if child.Next != nil {
			geneData = join(geneData, child.Next.collapseLeafs())
			//!	fmt.Println("cumlitive geneData:", geneData)
		}
	}
	if gn.Depth > 0 {
		for _, curGeneData := range (geneData) {
			curGeneData.DataValue[gn.Depth-1] = gn.DataValue
		}
	}
	///!fmt.Println("Final geneData:", geneData)
	return geneData
}

func (gn *GeneNode) deleteLeafs(root *GeneNode) {

	for index, child := range (gn.GeneParts) {
		//recurese this function down to its leaf
		if child.Next != nil {
			child.Next.deleteLeafs(root)
			//!	fmt.Println("cumlitive geneData:", geneData)
			gn.GeneParts[index].Next = nil
		}
	}
	return
}

//find the number of leafs on a node

func (gn *GeneNode) addNode(b byte) {
	gn.GeneParts[b].Next = new(GeneNode)
	gn.GeneParts[b].Next.Depth = gn.Depth + 1
	gn.GeneParts[b].Next.DataValue = uint8(b)
	for index, part := range (gn.GeneParts[b].Next.GeneParts) {
		part.DataValue = uint8(index)
		part.HitNumber = 0
	}
	return
}

//will take a byte to add to the tree then if it wants anouther byte for the next node it returns true. if it returns false then it
//hit a leaf
//will take a byte to add to the tree then if it wants anouther byte for the next node it returns true. if it returns false then it
//hit a leaf
func (gn *GeneNode) AddByte(b []byte) {
	//find the next GenePart and add one to its hit number
	//!fmt.Println("at:", b[0], " old hit number:", gn.GeneParts[b[0]].HitNumber, "Depth:",gn.Depth)
	curGN := gn
	for {
		curGN.GeneParts[b[0]].HitNumber++
		curGN.HitNumber++

		switch {
		//if the hit number of a genePart is enugh to create a new node
		case curGN.GeneParts[b[0]].HitNumber == uint16(NextNodeBound):
			//!fmt.Println("    adding node")
			//add a node
			curGN.addNode(b[0])
			//if it is greater tehn then Next node
		case curGN.GeneParts[b[0]].HitNumber > uint16(NextNodeBound):
			//!fmt.Println("    calling Down")
			//then call AddByte with the next byte on the next node if there are bytes left
			if cap(b) > 1 && curGN.Depth < uint8(MaxDepth) {
				curGN = curGN.GeneParts[b[0]].Next
				b = b[1:2]
				continue
			}
		}
		return
	}
}
//inverted becuase it need to be sorted greatest to least
func (gd GeneDataSlice) Less(i, j int) bool {
	if gd[i].HitNumber >= gd[j].HitNumber {
		return true
	} else {
		return false
	}
	return false
}
func (gd GeneDataSlice) Swap(i, j int) {
	tempd := gd[i]
	gd[i] = gd[j]
	gd[j] = tempd
	return
}
func (gd GeneDataSlice) Len() int { return len(gd) }
func (gd GeneDataSlice) Sort() {
	sortable := sort.Interface(gd)
	fmt.Print("s")
	sort.Sort(sortable)
	return
}
func (gd GeneDataSlice) Stats(){
	var avarageLen float= 0
	var maxLen int= 0
	var minLen int= 50000
	var maxO uint16 = 0
	var minO uint16= 50000
	var avarageO float =0 
	
	for _, item := range(gd){
		if len(item.DataValue) > maxLen{
			maxLen = len(item.DataValue)
		}
		if len(item.DataValue) < minLen{
			minLen = len(item.DataValue)
		}
		avarageLen += float(len(item.DataValue))
		
		if item.HitNumber > maxO {
			maxO = item.HitNumber
		}
		if item.HitNumber < minO {
			minO = item.HitNumber
		}
		avarageO += float(uint(item.HitNumber))
	}
	avarageLen /= float(len(gd))
	avarageO /= float(len(gd))
	fmt.Println()
	fmt.Println("min Length:", minLen,
		    "max Length:", maxLen,
		    "a Length:", avarageLen,
		    "min Hits:", minO,
		    "max Hits:", maxO,
		    "a Hits:", avarageO,
		    "total:", len(gd))
	return
			
}

//should consider changeing this to a io.Reader instead of a []byte to deal with huge data sets (a movie for example)
func (bt *BruteTreeTrainer) Train(trainingData []byte) {
	fmt.Println(len(trainingData))
	//from the minimum match size to the max ma
	for index := 0; index < len(trainingData); index++ {
		if bt.root.HitNumber > 70000 {
			go GeneDataSlice(bt.curHighest).Stats()
			bt.curHighest = merge(bt.curHighest, shrink(bt.nodeSet(),MinLin))
			bt.root.deleteLeafs(bt.root)
			GeneDataSlice(bt.curHighest).Stats()
			bt.root = nil
			bt.root = newTree()
		}
		bt.root.AddByte(trainingData[index : index+1])
	}

}

//Returns the nthHighestNodes this is decied from the top down
//so if there are three nodes that are valid at the top level of the tree
//and 2 lower down and the user asks for 5 nodes then it will return
//first the top one then the 2 lower ones
func (bt *BruteTreeTrainer) nthHighest(n int) []GeneData {
	slice := bt.root.collapseLeafs()
	//fmt.Println(slice)
	gds := GeneDataSlice(slice)
	gds.Sort()
	if n > len(gds) {
		n = len(gds) - 1
	}
	return gds[0:n]
}

func (bt *BruteTreeTrainer) nodeSet() []GeneData {
	slice := bt.root.collapseLeafs()
	gds := GeneDataSlice(slice)
	return gds
}

func (trainer *BruteTreeTrainer) TrainOnDir(dir string) (err os.Error) {
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
			file.Close()
		}
	}
	return
}
//func compareArray(bytes []byte, index, i int) (isThereGene bool) { }
func (trainer *BruteTreeTrainer) TrainOnDirPart(dir string, howMany int) (err os.Error) {
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
			file.Close()
			//trainer.Output()
		}
	}
	trainer.curHighest = merge(trainer.curHighest, trainer.nodeSet())
	gds := GeneDataSlice(trainer.curHighest)
	gds.Sort()
	fmt.Println("final Highest:", gds[0:30])
	return
}
// needs to be implented
func (bt *BruteTreeTrainer) GetGenes() (genes [][]byte) {
	bt.curHighest = merge(bt.curHighest, bt.nodeSet())
	gds := GeneDataSlice(bt.curHighest)
	gds.Sort()
	genes = make([][]byte, len(gds))
	for index, curI := range (gds) {
		genes[index] = curI.DataValue
	}
	return
}
/*
func OpenTrainer(data io.Reader) (err os.Error) {
	return
}

func Save(data io.Writer) (err os.Error) {}
*/
func (bt *BruteTreeTrainer) Output() {
	//fmt.Println("current Highest occurance", bt.nthHighest(5))
	return
}
