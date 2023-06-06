package model

type SearchResponse struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		Hits []struct {
			Score  float64 `json:"_score"`
			Source struct {
				Key    string `json:"key"`
				Title  string `json:"title"`
				Number int    `json:"number"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

type BulkError struct {
	Error struct {
		RootCause []struct {
			Type   string `json:"type"`
			Reason string `json:"reason"`
		} `json:"root_cause"`
		Type   string `json:"type"`
		Reason string `json:"reason"`
	} `json:"error"`
	Status int `json:"status"`
}

type BulkCreateResponse struct {
	Took   int  `json:"took"`
	Errors bool `json:"errors"`
	Items  []struct {
		Create struct {
			Index   string `json:"_index"`
			Id      string `json:"_id"`
			Version int    `json:"_version"`
			Result  string `json:"result"`
			Shards  struct {
				Total      int `json:"total"`
				Successful int `json:"successful"`
				Failed     int `json:"failed"`
			} `json:"_shards"`
			SeqNo       int `json:"_seq_no"`
			PrimaryTerm int `json:"_primary_term"`
			Status      int `json:"status"`
		} `json:"create,omitempty"`
	} `json:"items"`
}

// ---bulk Insert ---
type Action struct {
	Create *CreateDetail `json:"create,omitempty"`
}

type CreateDetail struct {
	Index string `json:"_index"`
}

//若有其他content需求 改這邊
type ContentDetail struct {
	Title    string `json:"title"`
	Director string `json:"director"`
	Year     string `json:"year"`
}

// ---bulk Insert ---
