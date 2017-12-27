package blockchain

type MerkleTree struct {
	Left  *MerkleTree
	Data  []byte
	Right *MerkleTree
}

func New(data string) *MerkleTree {
	return &MerkleTree{nil,[]byte(data),nil}
}