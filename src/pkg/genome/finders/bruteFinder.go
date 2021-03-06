package finders

import (
	"io"
	"genome"
	"os"
	"bytes"
	"fmt"
)

func HashFinder(data io.ReadSeeker, genomeRW io.ReadSeeker,theGenome genome.Genome, geneMap genome.GeneMap, indexer genome.Indexer) (geneAddressSlice []genome.GeneAddress, unfound []byte, ordering []bool, size int64) {
	//find the longest item in geneMap
	maxSize := 0
	for _, item := range (geneMap) {
		if item.Length > int64(maxSize) {
			maxSize = int(item.Length)
		}
	}
	fmt.Println("max Gene size:",maxSize)
	//make an array of byte arrays to hold all of the hashes that are generated that is as long as the longest
	strings := make([]string, maxSize+1)
	fileSize, _ := data.Seek(0, 2)
	fmt.Println("file Size:", fileSize)
	geneAddressSlice = make([]genome.GeneAddress, fileSize/4)
	unfound = make([]byte, fileSize)
	ordering = make([]bool, fileSize)
	//the starting position
	startPos := 0
	foundGenes := 0
	unfoundGenes := 0
	fmt.Println(data)
	//number of times the loop has run through
	i := 0
	//grab n data  to make room for the zero case 
	dataSlice := make([]byte, maxSize)
	for {
		if (int64(i) % (fileSize/2000)) ==0 {
			fmt.Println((float(i)/float(fileSize))*100,"% foundGenes:", foundGenes,
				    "unfoundGenes:", unfoundGenes,
				    "precent Compressed:", ((float(startPos)-float(unfoundGenes))/float(startPos))*100, "%", 
				    "total bytes:",startPos,
				    "total iters:", i)
		}
		//fmt.Println(data)
		//start at the begging of this data segment
		_, err := data.Seek(int64(startPos), 0 )
		if err != nil && err != os.EOF {
			panic(err.String())
		}
			
		//get the data remebering how much of the data is valid
		validData, err := data.Read(dataSlice)
		//fmt.Println("read :", validData)
		if err != nil && err != os.EOF {
			panic(err.String())
		}
		//for each pos in the data array
		//generate hashes from minSize (aka 5) size of an int to n length
		for index := 1; index < validData && index <= maxSize; index++ {
			tempBytes := indexer.Index(genome.Gene(dataSlice[0:index]))
			buf := bytes.NewBuffer(tempBytes)
			strings[index] = buf.String()
		}
		falsePositive := false
		biggestHashSize := 0
		biggestHashIndex := ""
		//for each item in the genome cmp the strings
		for key, address := range (geneMap) {
			if int(address.Length) > biggestHashSize && key == strings[address.Length] {
				testData := make([]byte, int(address.Length))
				genomeRW.Seek(address.Addr,0 )
				genomeRW.Read(testData)
				//if these are not infact the same skip over this one
				if !bytes.Equal(testData, dataSlice[0:address.Length]) {
					fmt.Println(testData)
					fmt.Println(dataSlice[0:address.Length])
					println("false Positive")
					falsePositive = true
					continue
				}
				//fmt.Println(address.Length)
				falsePositive = false
				biggestHashSize = int(address.Length)
				biggestHashIndex = key

			}
		}
		//now after the item has been found (or not) put the new infromation in the right list
		if biggestHashSize != 0 {
			//if something has been found assign it to the found genes
			geneAddressSlice[foundGenes] = genome.GeneAddress{0, 0, biggestHashIndex}
			//geneAddressSlice[foundGenes].Hash = biggestHashIndex
			ordering[i] = true
			foundGenes++
			startPos += biggestHashSize
		} else {
			if !falsePositive {
				//theGenome.AddGene(dataSlice)
				//foundGenes++
				//ordering[i]
			}//else{
			ordering[i] = false
			unfound[unfoundGenes] = dataSlice[0]
			unfoundGenes++
			startPos += 1
			//}
		}
		i++
		if 1 == validData {
			println("end of file")
			break
		}
	}
	//now finish up the ordering slice
	//println("ordering size", i)
	newOrdering := make([]bool, i)
	newAddressSlice := make([]genome.GeneAddress, foundGenes)
	newUnfound := make([]byte, unfoundGenes)
	//then copy both the unfound and geneAddressSlice over to well sized ones
	for place := 0; place < len(newOrdering); place++ {
		newOrdering[place] = ordering[place]
	}
	for j := 0; j < len(newAddressSlice); j++ {
		newAddressSlice[j] = geneAddressSlice[j]
	}
	for j := 0; j < len(newUnfound); j++ {
		newUnfound[j] = unfound[j]
	}
	ordering = newOrdering
	geneAddressSlice = newAddressSlice
	unfound = newUnfound

	size = fileSize
	fmt.Print("\n")
	fmt.Println("found", foundGenes)
	fmt.Println("unfound", unfoundGenes)
	fmt.Println("precent unFound: ",  (float(unfoundGenes)/float(size))*100)
	return

}


//TODO: Find a way to return the found parts of the  of the data
func BruteFinder(data io.ReadSeeker, genomeRW io.ReadSeeker, geneMap genome.GeneMap) (geneAddressSlice []genome.GeneAddress, unknownGenes []genome.Gene, ordering []int, size int64) {
	//yes i konw 9 is a magic number... the reason it is picked is becasue it is 4 more then 8 (64/4) thous ensuring it is more effcent to transfer that then the data
	//max read size for now is 1000 bytes
	///run this after checking all of the genes that have been added
	maxCache := 200000
	genCache := make([]byte, maxCache)
	dataCache := make([]byte, maxCache)
	maxChunk := 100000
	var err os.Error
	dataSize := 0
	dataBound := 30
	offset := 0
	dataOffset := 0
	i := 0
	err = nil
	//this should be changed depending on the length of the data...
	// geneAddress := make([]genome.GeneAddress,50000)
	//for as long as there is data in the file
	for err != os.EOF {
		_, err = data.Read(dataCache)
		for dataChunkBegs := 0; dataChunkBegs < (maxCache/4)*3; dataChunkBegs += (maxCache / 4) {
			dataChunk := dataCache[dataChunkBegs : dataChunkBegs+maxChunk]
			//for the whole genome
			var innerErr os.Error
			for os.EOF != innerErr {
				_, innerErr = genomeRW.Read(genCache)
				// grab the nth chunk of the currently read genmo
				for chunkBegs := 0; chunkBegs < (maxCache/4)*3; chunkBegs += (maxCache / 4) {
					offset += (maxCache / 4)

					cacheChunk := genCache[chunkBegs : chunkBegs+maxChunk]
					//pick a begging place in the gene and in the data file

					//for the whole data set
					sizeFound := 0
					for dataBeg := 0; dataBeg < maxCache/4; dataBeg++ {
						//only up to the end of the current chunck
						//keep track of the number genes have been found in this section of the genome for this pice of data
						//this could be used lator to find the beset possible gene
						numFound := 0

						//compair it to every position genome and save the geneAddress
						for genBeg := 0; genBeg < maxCache/4; genBeg++ {
							//see how long the match is
							var geneSize int
							for geneSize = 0; cacheChunk[genBeg+dataSize] == dataChunk[dataBeg+dataSize]; geneSize++ {
							}
							// if it is smaller then a previous match to the same thing throw it out
							if dataSize > dataBound && (numFound == 0 || int64(dataSize) > geneAddressSlice[i-1].Length) {
								//create a gene
								var geneAddress genome.GeneAddress
								geneAddress.Length = int64(geneSize)
								geneAddress.Addr = int64(dataOffset + genBeg)
								if numFound != 0 {
									geneAddressSlice[i-1] = geneAddress
								} else {
									geneAddressSlice[i] = geneAddress
									i++
								}

								sizeFound++
								//move the data reader forward to compensate for the data used in the gene
								dataBeg += geneSize
								if dataBeg > maxCache/4 {
									dataOffset = int(dataBeg + (maxCache / 4))
								} else {
									dataOffset = 0
								}
							}
						}
					}
				}
				//move back to make the cacheing work right
				genomeRW.Seek(int64(-(maxCache / 4)), 1)
			}
		}
		//move back to make the caching work right
		data.Seek(int64(-(maxCache / 4)), 1)
	}
	return
}
