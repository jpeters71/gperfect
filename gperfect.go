package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	defer trace("gperfect")()
	perPtr := flag.Uint64("number", 0, "Perfect number to find")
	coreCountPtr := flag.Uint("numCores", 8, "Number of cores to use when calculating random perfect numbers")
	numRandPtr := flag.Uint("numRandom", 100, "Number of random perfect numbers to look for.  Omit for continuous calculations")
	flag.Parse()

	if *perPtr > 0 {
		calcPerfect(*perPtr, *coreCountPtr)
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
			if uNumTest%i == 0 {
				uCalc += i
			}
		}
	} else {
		//		chans := make(chan uint64, uNumCores)
		var uChan uint
		var uLeftOver = uHalf % uint64(uNumCores)
		var uChunkStart uint64 = 1
		var chans []chan uint64

		uInterval := uHalf / uint64(uNumCores)
		uChunkEnd := uLeftOver

		for uChan = 0; uChan < uNumCores; uChan++ {
			chans = append(chans, make(chan uint64))

			uChunkEnd += uInterval
			go calcProc(uChan, uNumTest, uChunkStart, uChunkEnd, chans[uChan])
			uChunkStart = uChunkEnd + 1
		}
		for uChan = 0; uChan < uNumCores; uChan++ {
			fmt.Printf("uChan=%d, uCalc=%d\n", uChan, uCalc)
			uCalc += <-chans[uChan]
		}
	}
	if uCalc == uNumTest {
		fmt.Printf("%d,PERFECT\n", uNumTest)
	} else {
		fmt.Printf("%d,not perfect\n", uNumTest)
	}
}
func calcProc(uTmpChan uint, uTmpNumTest uint64, uChunkStart uint64, uChunkEnd uint64, chans chan (uint64)) {
	fmt.Printf("%d: PROC STARTED -start=%d, end=%d\n", uTmpChan, uChunkStart, uChunkEnd)
	var uTmp uint64
	for i := (uChunkStart + 1); i <= uChunkEnd; i++ {
		//fmt.Printf("chan=%d, idx=%d\n", uTmpChan, i)
		if uTmpNumTest%i == 0 {
			uTmp += i
		}
	}
	fmt.Printf("DONE - %d; val=%d\n", uTmpChan, uTmp)
	chans <- uTmp
}

func calcRandPerfects(coureCountPtr *uint, numRandPtr *uint) {
	var i uint
	uMax := *numRandPtr
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i = 0; i < uMax; i++ {
		calcPerfect(uint64(r.Int31()), *coureCountPtr)
	}
}

func trace(msg string) func() {
	start := time.Now()
	fmt.Printf("%s started at: %s\n", msg, start.Local())
	return func() {
		end := time.Now()
		fmt.Printf("%s finished at: %s; elapsed time: %s\n", msg, end.Local(), end.Sub(start))
	}
}
