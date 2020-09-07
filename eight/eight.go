package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

//------- Trie implementation

type TrieNode struct {
	cellValue int
	children  []*TrieNode
}

func newTrieNode(value int) *TrieNode {
	node := new(TrieNode)
	node.cellValue = value
	return node
}

func (node *TrieNode) addState(cellData []int) {
	if len(cellData) == 0 {
		return
	}
	var nextNode *TrieNode
	found := false
	for _, childNode := range node.children {
		if childNode.cellValue == cellData[0] {
			nextNode = childNode
			found = true
			break
		}
	}
	if !found {
		nextNode = newTrieNode(cellData[0])
	}
	nextNode.addState(cellData[1:])
	node.children = append(node.children, nextNode)
}

func (node *TrieNode) contains(cellData []int) bool {
	if len(cellData) == 0 {
		return true
	}
	for _, childNode := range node.children {
		if childNode.cellValue == cellData[0] {
			return childNode.contains(cellData[1:])
		}
	}
	return false
}

//---------- End trie implementation

type BoardState [9]int

var endState = BoardState([9]int{0, 1, 2, 3, 4, 5, 6, 7, 8})

type GameState struct {
	boardState BoardState
	parent     *GameState
}

func (gameState GameState) nextStates() []*GameState {
	var nextStates []*GameState
	zeroIndex := -1
	for index, value := range(gameState.boardState) {
		if value == 0 {
			zeroIndex = index
			break
		}
	}
	if zeroIndex == -1 {
		panic("No zero cell in board state")
	}
	if zeroIndex / 3 < 2 {
		// Switch with below
		var newBoardState [9]int = gameState.boardState
		val := newBoardState[zeroIndex+3]
		newBoardState[zeroIndex] = val
		newBoardState[zeroIndex+3] = 0
		state := GameState{newBoardState, &gameState}
		nextStates = append(nextStates, &state)
	}
	if zeroIndex / 3 > 0 {
		// Switch with above
		var newBoardState [9]int = gameState.boardState
		val := newBoardState[zeroIndex-3]
		newBoardState[zeroIndex] = val
		newBoardState[zeroIndex-3] = 0
		state := GameState{newBoardState, &gameState}
		nextStates = append(nextStates, &state)
	}
	if zeroIndex % 3 > 0 {
		// Switch with left
		var newBoardState [9]int = gameState.boardState
		val := newBoardState[zeroIndex-1]
		newBoardState[zeroIndex] = val
		newBoardState[zeroIndex-1] = 0
		state := GameState{newBoardState, &gameState}
		nextStates = append(nextStates, &state)
	}
	if zeroIndex % 3 < 2 {
		// Switch with right
		var newBoardState [9]int = gameState.boardState
		val := newBoardState[zeroIndex+1]
		newBoardState[zeroIndex] = val
		newBoardState[zeroIndex+1] = 0
		state := GameState{newBoardState, &gameState}
		nextStates = append(nextStates, &state)
	}
	return nextStates
}

func (bs *BoardState) fromString(stateData string)  error {
	for index, line := range strings.Split(strings.Trim(stateData, " \n"), "\n") {
		for secondIndex, column := range strings.Split(line, " ") {
			arrayIndex := 3*index + secondIndex
			if column == "X" {
				bs[arrayIndex] = 0
			} else {
				colVal, err := strconv.Atoi(column)
				if err != nil {
					return fmt.Errorf("Error converting value to int: %s", column)
				}
				bs[arrayIndex] = colVal
			}
		}
	}
	return nil
}

func (bs *BoardState) isEndState() bool {
	return *bs == endState
}

func (bs *BoardState) printRepr() string {
	var letters [9]string
	for index, value := range bs {
		if value == 0 {
			letters[index] = "X"
		} else {
			letters[index] = strconv.Itoa(value)
		}
	}
	return strings.Join([]string{strings.Join(letters[:3], " "), strings.Join(letters[3:6], " "), strings.Join(letters[6:9], " ")}, "\n")
}

// A linked-list queue implementation. Since we're not going to traverse the
// queue, this should be good enough
type QueueNode struct {
	state *GameState
	next  *QueueNode
}

type Queue struct {
	head *QueueNode
	tail *QueueNode
}

func newQueue() Queue {
	return Queue{nil, nil}
}

func (queue *Queue) add(state *GameState) {
	newNode := &QueueNode{state, nil}
	if queue.head == nil {
		queue.head = newNode
	}
	if queue.tail == nil {
		queue.tail = newNode
	} else {
		queue.tail.next = newNode
		queue.tail = newNode
	}
}

func (queue *Queue) length() int {
	qlength := 0
	currentNode := queue.head
	for (currentNode != nil) {
		qlength += 1
		currentNode = currentNode.next
	}
	return qlength
}

func (queue *Queue) popLeft() (*GameState, error) {
	if queue.head == nil {
		return nil, fmt.Errorf("Queue empty")
	}
	retval := queue.head.state
	queue.head = queue.head.next
	if queue.head == nil {
		queue.tail = nil
	}
	return retval, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please specify start state")
		os.Exit(1)
	}
	filename := os.Args[1]
	inputData, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Could not open input file, error: ", err)
	}
	startData := string(inputData)
	var startState BoardState
	error := startState.fromString(startData)
	if error != nil {
		fmt.Println(error)
		os.Exit(1)
	}
	trie := newTrieNode(-1)
	trie.addState(startState[:])
	firstMove := GameState{startState, nil}
	queue := newQueue()
	queue.add(&firstMove)
	found := false
	var currentState *GameState
	for !found {
		currentState, err = queue.popLeft()
		if err != nil {
			fmt.Println("Queue empty")
			return
		}
		if currentState.boardState.isEndState() {
			fmt.Println("Solution found")
			found = true
			break
		}
		nextStates := currentState.nextStates()
		for _, move := range nextStates {
			moveData := move.boardState[:]
			if !trie.contains(moveData) {
				queue.add(move)
				trie.addState(moveData)
			}
		}
	}
	if found {
		printState := currentState
		for (printState != nil) {
			fmt.Printf("%s \n\n", printState.boardState.printRepr())
			printState = printState.parent
		}
	}
}

// TODO
// Seen states trie
// getting state children
// push to queue and process
