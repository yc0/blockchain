# BlockChain 
Building blockchain in Golang

refer to [jeiwan blockchain building tutorial]( https://jeiwan.cc/posts/building-blockchain-in-go-part-1/)
## Part-1 Branch
build the simple block with sha256 hash; however, I ignore this part, and move toward part-2 directly

## Part-2 Branch
build the proof of work hash for block.
The different part is that I utilize the goroutine, sync.once and limit goroutine to expedite the pow calculation.

This modification gave me runable and reasonable time-process for PoW with high compleixty on my laptop (new macbook)

Here are my constanst parameters

```
targetBits = 26 
maxNonce = 1<<28 - 1
maxConcurrencies = 32
```

### Workable Parameter 
While trying to test and implement your blockchain, you might need reasonable parameters, saying below 1.0s

Here are the suitable configurations.
```
targetBits = 16
maxNonce = 1<<24 - 1
maxConcurrencies = 32
```