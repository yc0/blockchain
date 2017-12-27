package blockchain

type BlockChain struct {
	Blocks []*Block
}

func (bc *BlockChain) AddBlock(data string) {
	prev := bc.Blocks[len(bc.Blocks)-1]
	var b *Block
	bc.Blocks = append(bc.Blocks,b.New(data,prev.Hash))
}

func (bc *BlockChain) New() *BlockChain {
	return &BlockChain{[]*Block{Genesis()}}
}