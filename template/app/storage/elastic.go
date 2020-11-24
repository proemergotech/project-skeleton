//%: {{- if .Elastic }}
package storage

import (
	"context"

	"github.com/olivere/elastic"
	"github.com/proemergotech/uuid"
)

//%: {{- if .Examples }}
const dummyIndex = "dummy_index_v1"

// todo: remove
//  Example for elastic sort script
const sortScript = `
if ('offline' == doc['state'].value ) {
  if (doc['last_broadcast'].value == null) {
    return 0;
  }
  return doc['last_broadcast'].value.getMillis();
}

return (doc['_id'].value + params.sort_hash).hashCode();`

//%: {{- end }}

type Elastic struct {
	ElasticClient *elastic.Client
}

func NewElastic(elasticClient *elastic.Client) *Elastic {
	return &Elastic{ElasticClient: elasticClient}
}

//%: {{- if .Examples }}
// todo: remove
//  Implementation example for search in elastic
func (e *Elastic) DummySearch(ctx context.Context, siteGroupCode, queryStr string, limit, offset int) ([]uuid.UUID, int, uuid.UUID, error) {
	query := elastic.NewBoolQuery().Filter(
		elastic.NewTermQuery("site_group_code", siteGroupCode),
	)

	if len(queryStr) > 0 {
		query = query.Must(
			elastic.NewMultiMatchQuery(queryStr).
				Type("most_fields").
				FieldWithBoost("display_name", 10).
				Field("display_name.ngram"),
		)
	}

	sortUUID := uuid.Nil

	sortBy := make([]elastic.Sorter, 0, 2)
	sortBy = append(sortBy, elastic.NewScriptSort(elastic.NewScript(sortScript), "number").Desc())
	if len(queryStr) > 0 {
		sortBy = append(sortBy, elastic.NewScoreSort())
	}

	result, err := e.ElasticClient.
		Search(dummyIndex).
		Type("_doc").
		FetchSource(true).
		Query(query).
		SortBy(sortBy...).
		Size(limit).
		From(offset).
		Do(ctx)
	if err != nil {
		return nil, 0, uuid.Nil, elasticError{Err: err}.E()
	}

	totalCount := int(result.Hits.TotalHits)
	res, err := elasticHitsToUUIDs(result.Hits.Hits)
	if err != nil {
		return nil, 0, uuid.Nil, err
	}

	return res, totalCount, sortUUID, nil
}

func elasticHitsToUUIDs(hits []*elastic.SearchHit) ([]uuid.UUID, error) {
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

//%: {{- end }}

//%: {{- end }}
