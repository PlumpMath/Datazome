all:../../../../lib/genome/trainers.a

trainers.8: bruteTreeTrainer.go bruteTrainer.go
	8g -I ../../../../lib -o trainers.8 bruteTreeTrainer.go bruteTrainer.go

../../../../lib/genome/trainers.a: trainers.8 
	gopack gru ../../../../lib/genome/trainers.a trainers.8

clean:
	-rm *.8
	-rm ../../../../lib/genome/trainers.a

test: clean ../../../../lib/genome/trainers.a
	8g -I ../../../../lib trainerTests.go
	8l -o trainerTester trainerTests.8