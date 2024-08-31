package resources

type Year struct {
	numStars    int
	days        []*Day
	leaderboard YearLB
}

type Day struct {
	puzzle      Puzzle
	leaderboard DayLB
}
