all:../../../lib/genome.a
	
genome.8: gene.go geneAddress.go geneMap.go genome.go helperTypes.go tGene.go
	8g -I ../../../lib -o genome.8 gene.go geneAddress.go geneMap.go genome.go helperTypes.go tGene.go

../../../lib/genome.a: genome.8
	gopack gru ../../../lib/genome.a genome.8

test: ../../../lib/genome.a genomeTests.go
	8g -I ../../../lib genomeTests.go
	8l -L ../../../lib -o genomeTester genomeTests.8
	

clean:
	-rm *.8
	-rm *.a
	-rm ../../../lib/genome.a