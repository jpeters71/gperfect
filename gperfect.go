package main

import (
    "flag"
    "fmt"
    "math/rand"
    "time"    
)

func main() {
    perPtr := flag.Uint64("number", 0, "Perfect number to find")
    coreCountPtr := flag.Uint("numCores", 8, "Number of cores to use when calculating random perfect numbers")
    numRandPtr := flag.Uint("numRandom", 100, "Number of random perfect numbers to look for.  Omit for continuous calculations")
    flag.Parse()
    
    if *perPtr > 0 {
        calcPerfect(*perPtr);
    } else {
        calcRandPerfects(coreCountPtr, numRandPtr)
    }
}


func calcPerfect(uNumTest uint64, uNumCores uint) {
    var uCalc uint64
    var i uint64
    uHalf := uNumTest / 2
    
    // If the domain of possible numbers is > 100,000,000 then split it up on multiple threads
    if uNumTest < 100000000 {
        for i = 1; i <= uHalf; i++ {
            if uNumTest % i == 0 {
                uCalc += i
            } 
        }
    } else {
        chans := make(chan uint64, uNumCores)
        uLeftOver := uHalf % uNumCores 
        uInterval := uHalf / uNumCores
        uChunk := uLeftOver 
        uChunkStart := 0
        
        for uChan := 0; uChan < uNumCores; uChan++ {
            uChunk += uInterval
            
            go func(uTmpChunkStart uint64, uTmpChunk uint64, chans chan(uint64)) {
                var uTmp uint64
                for i = uTmpChunkStart; i <= uTmpChunk; i++ {
                    if uNumTest % i == 0 {
                        uTmp += i
                    } 
                }
                chans <- uTmp
            }
            uChunkStart += uChunk
        } 
    }
    
    
    for i = 1; i <= uHalf; i++ {
        if uNumTest % i == 0 {
            uCalc += i
        } 
    }
    if uCalc == uNumTest {
        fmt.Printf("%d,PERFECT\n", uNumTest)
    } else {
        fmt.Printf("%d,not perfect\n", uNumTest)
    }
}

func calcRandPerfects(coureCountPtr *uint, numRandPtr *uint) {
    var i uint
    uMax := *numRandPtr
    r := rand.New(rand.NewSource(time.Now().UnixNano()))    

    for i=0; i < uMax; i++ {
        calcPerfect(uint64(r.Int31()))
    }
}
