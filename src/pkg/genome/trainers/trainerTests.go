//gg:{trainerTester: trainerTests.go, trainers.a: bruteTreeTrainer.go}
package main

import (
	"fmt"
	"./trainers"
)

func main() {
	//fmt.Printf (possibleGene[0].dataValue[0])
	fmt.Printf("Running trainers\n")
	treeTrainer()
}
/* func bruteTrainer(){
	train := trainers.NewBruteTrainer()
	train.TrainOnDirPart("/usr/bin", 100)
	train.GetGenes()

}*/
func treeTrainer() {
	train := trainers.NewBruteTreeTrainer()
	train.TrainOnDirPart("/usr/bin", 100)
	train.Output()
}
