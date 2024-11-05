package blockchain

type VoteData struct {
	VoterID     string `json:"voter_id"`     // ID pemilih (misalnya NIK atau hash unik)
	CandidateID string `json:"candidate_id"` // Kandidat yang dipilih
	Timestamp   int64  `json:"timestamp"`    // Waktu saat pemilih memberikan suara
}
