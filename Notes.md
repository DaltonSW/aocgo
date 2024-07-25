# AoC Util Notes

## Open Ended Ponderings

## Specific Callouts

### Getting AoC Session Token

- Open browser and go to adventofcode.com
- Open Dev Tools (F12) and go to the Network tab
- Navigate to any puzzle input page and wait for it to load
- Inspect your request's header and find the SESSION cookie
- Copy this to either AOC_SESSION environment variable, or ~/.config/aocutil/session

## Pseudocode

### User / Session loading
- Try to load in User
    - If able, check if SESSION key is stored
    - If unable, create new User
- If no SESSION key, check file and then env var
    - If found, store to User and continue
    - If not found, error and print out help text

## Data Structure (WIP)

- User
    - username : string
    - SESSION token : string
    - numStars : int

- <Something High Level>
    - Years : map[int]*Year

- Year
    - Puzzles : []*Puzzle, where its index in the array is the puzzle's day

- Puzzle
    - Year : int
    - Day : int
    - PartA : *PuzzlePart
    - PartB : *PuzzlePart

- PuzzlePart
    - isPartB : bool (if false, is Part A)
    - submissions : []Submissions
    - bestSolve : datetime
    - correctAnswer : int
