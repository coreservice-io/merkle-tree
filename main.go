package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/coreservice-io/merkleTree/merkle"
	"github.com/ethereum/go-ethereum/common/math"
	"golang.org/x/crypto/sha3"
)

// testcontent is indeed leaf
type TestContent struct {
	index   string
	address string
	amount  string
}

func encodePacked(input ...[]byte) []byte {
	return bytes.Join(input, nil)
}

//CalculateHash hashes the values of a TestContent for leaf
func (t TestContent) CalculateHash() ([]byte, error) {
	address, _ := hex.DecodeString(strings.Trim(t.address, "0x"))
	pre_result := encodePacked(
		encodeUint256(t.index),
		address,
		encodeUint256(t.amount),
	)

	hash := sha3.NewLegacyKeccak256()
	hash.Write(pre_result)
	result := hash.Sum(nil)
	return result, nil
}

func encodeUint256(v string) []byte {
	bn := new(big.Int)
	bn.SetString(v, 10)
	return math.U256Bytes(bn)
}

//Equals tests for equality of two Contents
func (t TestContent) Equals(other merkle.Content) (bool, error) {
	if (t.index == other.(TestContent).index) &&
		(t.address == other.(TestContent).address) &&
		(t.amount == other.(TestContent).amount) {
		return true, nil
	} else {
		return false, nil
	}
}

func main() {

	var list []merkle.Content

	list = append(list, TestContent{
		index:   "0",
		address: "0x37e9e835171e40ceb35cdb0a05346f9c451c6156",
		amount:  "0",
	})

	list = append(list, TestContent{
		index:   "1",
		address: "0x37e9e835171e40ceb35cdb0a05346f9c451c6156",
		amount:  "1",
	})

	list = append(list, TestContent{
		index:   "2",
		address: "0x37e9e835171e40ceb35cdb0a05346f9c451c6156",
		amount:  "2",
	})

	list = append(list, TestContent{
		index:   "3",
		address: "0x37e9e835171e40ceb35cdb0a05346f9c451c6156",
		amount:  "3",
	})

	t, err := merkle.NewTreeWithHashStrategy(list,
		func(left []byte, right []byte) []byte {
			// fmt.Println("////////////////")
			// fmt.Printf("left:%x\n", left)
			// fmt.Printf("right:%x\n", right)
			// fmt.Println("////////////////")
			hash := sha3.NewLegacyKeccak256()
			if string(left) == string(right) {
				return left
			} else if string(left) < string(right) {
				hash.Write(append(left, right...))
				return hash.Sum(nil)
			}
			hash.Write(append(right, left...))
			return hash.Sum(nil)
		})

	if err != nil {
		fmt.Println(err)
	}

	//Get the Merkle Root of the tree
	mr := t.MerkleRoot()
	//fmt.Println(mr)
	//fmt.Println(hex.EncodeToString(mr))
	fmt.Printf("%x\n", mr)

}
