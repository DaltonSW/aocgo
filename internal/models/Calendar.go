package models

type Year struct {
	numStars    int
	days        []*Day
	leaderboard Leaderboard
}

type Day struct {
	puzzle      Puzzle
	leaderboard Leaderboard
}
