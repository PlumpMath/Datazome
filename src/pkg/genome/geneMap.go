package genome

import (
	"io"
	"gob"
	"os"
	"fmt"
)

type GeneMap map[string]GeneAddress // used to asscaite genes with their hashes //

//probaly will need to change this to return a []byte
func (geneMap GeneMap) GetClosestAddr(addr int64) (key string, valueOff int64) {
	var closest int64 //need to find a way to assign this to a real value so that it does not accidently get to be the closestValuePair
	closest = 1000000
	for curKey, vp := range (geneMap) {
		if vp.Addr-addr == addr {
			key = curKey
			valueOff = 0
			continue
		} else {
			if vp.Addr-addr < closest {
				key = curKey
				valueOff = vp.Addr - addr
			}
		}
	}
	return
}

func (geneMap GeneMap) ContainsKey(key string) (added bool) {
	added = false
	for key, _ := range (geneMap) {
		if key == key {
			added = true
			return
		}
	}
	return
}

func LoadGeneMap(file io.ReadWriteSeeker) (geneMap GeneMap, err os.Error) {
	//file.Seek(0, 0)
	dec := gob.NewDecoder(file)
	var length lens
	err = dec.Decode(&length)
	geneMap = make(GeneMap, length.i)
	fmt.Println("keys being decoded:" ,length.i)
	keyval := new(keyVal)
	for x := 0; x < length.i; x++ {
		err = dec.Decode(keyval)
 		if err == os.EOF { break }
		if err!= nil {panic(err.String())}
		fmt.Println("data", keyval.k, keyval.i)
		geneMap[keyval.k] = keyval.i
	}
	return
}

type lens struct {
	i int
}
type keyVal struct{
	k string
	i GeneAddress
}

func (geneMap GeneMap) Save(file io.ReadWriteSeeker) (err os.Error) {
	//go to the begging of the file
	file.Seek(0, 0)
	enc := gob.NewEncoder(file)
	err = enc.Encode(lens{len(geneMap)})
	for key, item := range (geneMap) {
		enc.Encode(keyVal{ key,item })
		if err != nil {
			panic(err.String())
		}
	}
	return

}
