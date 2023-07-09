package main

import (
	"fmt"
	"hash"
	"time"

	"github.com/google/uuid"
	"github.com/spaolacci/murmur3"
)

//use uint of 8 bits or byte stream to reduce the space, since we can store 8 values at each index unlike bool slice, where we can store only 1 value at each index(either true or false)
type BloomFilter struct{
	Filter []uint8
	Size int32
}

var mHasher hash.Hash32

func init(){
	mHasher = murmur3.New32WithSeed(uint32(time.Now().Unix()))
}

func MurmurHash(key string,size int32) int32{
	mHasher.Write([]byte(key))
	result := mHasher.Sum32() % uint32(size)
	mHasher.Reset()
	return int32(result)
}

func NewBloomFilter(size int32) *BloomFilter{
	return &BloomFilter{
		make([]uint8, size),
		size,
	}
}

func (b *BloomFilter) Add(key string){
	//multiply by 8, since we can store 8 values at each index(each bit acts as a bucket)
	index := MurmurHash(key,b.Size*8)
	aIdx := index/8
	aBit := 1<<(index%8)
	b.Filter[aIdx] |= uint8(aBit)
}

func (b *BloomFilter) Exists(key string) bool{
	index := MurmurHash(key,b.Size*8) 
	aIdx := index/8
	return b.Filter[aIdx] >> (uint8(index)%8) & 1 == 1
}

func main(){
	datasetExists := make(map[string] bool)
	datasetNotExists := make(map[string] bool)
	
	for i:=0;i<500;i++{
		u := uuid.New()
		datasetExists[u.String()]=true
	}

	for i:=0;i<500;i++{
		u := uuid.New()
		datasetNotExists[u.String()]=true
	}

	bloom := NewBloomFilter(1000)
	
	for key := range datasetExists{
		bloom.Add(key)
	}

	var falsePositive int = 0
	for key := range datasetNotExists {	if bloom.Exists(key){
			falsePositive++
		}
	}
	fmt.Println("False positive rate:",float64(falsePositive)/float64(len(datasetNotExists)))
}