package storage

import (
	"github.com/olivere/elastic"

	uuid "gitlab.com/proemergotech/uuid-go"
)

type Elastic struct {
	ElasticClient *elastic.Client
}

func NewElastic(elasticClient *elastic.Client) *Elastic {
	return &Elastic{ElasticClient: elasticClient}
}

func (e *Elastic) elasticHitsToUUIDs(hits []*elastic.SearchHit) ([]uuid.UUID, error) {
	res := make([]uuid.UUID, 0, len(hits))
	for _, hit := range hits {
		var uid uuid.UUID
		uid, err := uuid.FromString(hit.Id)
		if err != nil {
			return nil, elasticError{Err: err}.E()
		}
		res = append(res, uid)
	}
	return res, nil
}
