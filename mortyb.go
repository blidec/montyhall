package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)
var (
	mode = flag.String("mode", "stay", "playing mode [stay|switch|random]")
	plays = flag.Int("plays", 10, "# of times to play")
	numBalls = flag.Int("balls", 3, "number of balls in bag")
)

type ball struct {
	hasPrize bool
}

type bag struct {
	ballMap map[int]ball
}

func main(){
	flag.Parse()

	if *mode != "stay" && *mode != "switch" && *mode != "random" {
		fmt.Println("Please specify one of the supported mode are [stay|switch|random]")
		os.Exit(1)
	}
	// bag of balls with one has prize
	prizeBallIndex := genRandomNum(0, *numBalls)

	ballMap := make(map[int]ball)
	for i := 0; i < prizeBallIndex; i++ {
		ballMap[i] = ball{false}
	}
	ballMap[prizeBallIndex] = ball{true}
	for i := prizeBallIndex+1; i < *numBalls; i++ {
		ballMap[i] = ball{false}
	}
	bag := bag{ballMap: ballMap}
	
	var winCount int
	for i := 0; i < *plays; i++ {
		// simulate player picks a ball from the bag
		pick := genRandomNum(0, *numBalls)

		// play with mode
		rc := playGame(bag.ballMap, prizeBallIndex, pick, *mode)
		//fmt.Printf("prize ball: %d, initial pick: %d, strategy: %s, win: %t\n", prizeBallIndex, pick, *mode, rc)

		if rc {
			winCount++
		}
	}

	fmt.Printf("strategy: %s, wins: %d out of %d\n", *mode, winCount, *plays)
}

// simulate playing game with strategy. Return win or not
func playGame(balls map[int]ball, prizeIndex, pickIndex int, mode string) bool {
	// build bag of remaining balls
	remainingBalls := make(map[int]ball)
	for key, value := range balls {
		if key != pickIndex {
			remainingBalls[key] = value
		}
	}

	pickBall,_ := balls[pickIndex]
	return playStrategy(remainingBalls, prizeIndex, pickIndex, pickBall, mode)
}

// playStrategy simulate removing one non-prize ball, and player picks one from remaining pool until
// the last one, then we check if that one is win or not
// balls - original bag of balls
// remainingBalls - bag of balls after one non-prize ball is removed
func playStrategy(remainingBalls map[int]ball, prizeIndex, pickIndex int, pickBall ball, mode string) bool {
	if len(remainingBalls) == 1 {
		for _,v := range remainingBalls {
			return !v.hasPrize
		}
	}

	// remove one non-prize from the pool
	nextRemainingPool := make(map[int]ball, 0)
	removeIndex := pickFromPool(remainingBalls, true, prizeIndex)
	for k,v := range remainingBalls {
		if k != removeIndex {
			nextRemainingPool[k] = v
		}
	}

	changePick := false
	if mode == "switch" {
		changePick = true
	} else if mode == "random" {
		if genRandomNum(0, 2) == 1 {
			changePick = true
		}
	}
	
	newPickIndex := pickIndex
	newPickBall := pickBall
	if changePick {
		// if  change pick, we'll pick one from and put back the original back to nextRemainingPool
		newPickIndex = pickFromPool(nextRemainingPool, false, 0)
		//fmt.Printf("New pick %d\n", newPickIndex)
		newPickBall, _ =  nextRemainingPool[newPickIndex]
		delete(nextRemainingPool, newPickIndex)
		nextRemainingPool[pickIndex] = pickBall
	}

	return playStrategy(nextRemainingPool, prizeIndex, newPickIndex, newPickBall, mode)
}

func genRandomNum(min, max int) int {
	if min == max {
		return min
	}
	rand.Seed(time.Now().UnixNano())
    return rand.Intn(max - min) + min
}

// pick ball from pool. excludeIndex is index to be excluded if exclude is set
func pickFromPool(balls map[int]ball, exclude bool, excludeIndex int) int {
	keys := make([]int, 0)
	for k := range balls {
		keys = append(keys, k)
	}
	
	for {
		index := genRandomNum(0, len(keys))
		if !exclude {
			return keys[index]
		} else {
			if keys[index] != excludeIndex {
				return keys[index]
			}
		}
	}
}