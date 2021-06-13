package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"
	"encoding/binary"
	"log"
)


const TargetBits = 24
const MaxNonce = math.MaxInt64

type Block struct{
	//time of creating
	Timestamp int64
	Data []byte
	//previous hash block
	PrevBlockHash []byte
	Hash []byte
	Nonce int
}

type Blockchain struct {
	blocks []*Block
}

type ProofOfWork struct {
	block *Block
	targer *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256 - TargetBits))
	proof := &ProofOfWork{b,target}
	return proof
}
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
func (proof *ProofOfWork) PrepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			proof.block.PrevBlockHash,
			proof.block.Data,
			IntToHex(proof.block.Timestamp),
			IntToHex(int64(TargetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)
	return data
}

func (proof *ProofOfWork) Run() (int,[]byte) {
	var (
		hashInt big.Int
		hash [32]byte
	)
	nonce := 0
	fmt.Printf("Mining the block containing \"%s\"\n", proof.block.Data)
	for nonce < MaxNonce {
		data := proof.PrepareData(nonce)
		hash := sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(proof.targer) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Println()
	return nonce, hash[:]
}


// func (b *Block) SetHash() {
// 	timestamp := []byte(strconv.FormatInt(int64(b.Timestamp),10))
// 	headers := bytes.Join([][]byte{b.PrevBlockHash,b.Data,timestamp},[]byte{})
// 	hash :=sha256.Sum256(headers)
// 	b.Hash = hash[:]
// }

func NewBlock(data string,PrevBlockHash []byte) *Block {
	block :=&Block{time.Now().Unix(),[]byte(data),PrevBlockHash,[]byte{},0}
	proof := NewProofOfWork(block)
	nonce , hash := proof.Run()
	block.Hash = hash
	block.Nonce = nonce
	return block
}

func (proof *ProofOfWork) Valid() bool {
	var hashInt big.Int

	data := proof.PrepareData(proof.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	isvalid := 	hashInt.Cmp(proof.targer) == -1
	return isvalid 
}

func (bc *Blockchain) AddBlock(data string) {
	prevb := bc.blocks[len(bc.blocks) - 1]
	newblock := NewBlock(data,prevb.Hash)
	bc.blocks = append(bc.blocks, newblock)
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis block", []byte{})
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

func main() {
	bc := NewBlockchain()
	bc.AddBlock("One ETH from Anton")
	bc.AddBlock("Two ETH from Anton")
	for _, v := range bc.blocks {
		fmt.Printf("Previous hash is %x\n",v.PrevBlockHash)
		fmt.Printf("Data is %s\n",v.Data)
		fmt.Printf("Hush is %x\n",v.Hash)
		pow := NewProofOfWork(v)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Valid()))
		fmt.Println()
	}
}