package miner

// DifficultyInterval sets the desired number of blocks between difficulty updates
var DifficultyInterval int = 144

// DifficultyTargetTimespan sets the desired time between difficulty updates
var DifficultyTargetTimespan int = 86400

// MiningIterationsPerCall sets the limit of hashes to be computed for each call of Mine()
var MiningIterationsPerCall int = 1000000
