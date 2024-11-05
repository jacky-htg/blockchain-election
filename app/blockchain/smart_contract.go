package blockchain

import (
	"fmt"
)

type Election struct {
	Candidates []string
	Votes      map[string]int
	Voters     map[string]bool
}

func NewElection(candidates []string) *Election {
	return &Election{
		Candidates: candidates,
		Votes:      make(map[string]int),
		Voters:     make(map[string]bool),
	}
}

func (e *Election) AddCandidate(candidate string) error {
	if e.isValidCandidate(candidate) {
		return fmt.Errorf("candidate %s already exists", candidate)
	}
	e.Candidates = append(e.Candidates, candidate)
	return nil
}

func (e *Election) Vote(voterID, candidate string) error {
	// Validasi kandidat
	if !e.isValidCandidate(candidate) {
		return fmt.Errorf("candidate %s is not valid", candidate)
	}

	// Cek apakah pemilih sudah memberikan suara
	if e.Voters[voterID] {
		return fmt.Errorf("voter %s has already voted", voterID)
	}

	// Tambahkan suara
	e.Votes[candidate]++
	e.Voters[voterID] = true
	return nil
}

func (e *Election) isValidCandidate(candidate string) bool {
	for _, c := range e.Candidates {
		if c == candidate {
			return true
		}
	}
	return false
}

func (e *Election) GetResults() map[string]int {
	return e.Votes
}

func (e *Election) DisplayResults() {
	results := e.GetResults()
	for candidate, votes := range results {
		fmt.Printf("Candidate: %s, Votes: %d\n", candidate, votes)
	}
}
