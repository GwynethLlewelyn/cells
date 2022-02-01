/*
 * Copyright (c) 2019-2022. Abstrium SAS <team (at) pydio.com>
 * This file is part of Pydio Cells.
 *
 * Pydio Cells is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Pydio Cells is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Pydio Cells.  If not, see <http://www.gnu.org/licenses/>.
 *
 * The latest code can be found at <https://pydio.com>.
 */

package docstore

import (
	"context"
	"fmt"
	"github.com/blevesearch/bleve/v2"
	query2 "github.com/blevesearch/bleve/v2/search/query"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/pydio/cells/v4/common/dao/mongodb"
	"github.com/pydio/cells/v4/common/proto/docstore"
	"github.com/pydio/cells/v4/common/utils/configx"
	json "github.com/pydio/cells/v4/common/utils/jsonx"
)

const (
	collDocuments = "documents"
)

var (
	upsert = true
	model  = mongodb.Model{
		Collections: []mongodb.Collection{
			{
				Name: collDocuments,
				Indexes: []map[string]int{
					{"store_id": 1},
					{"store_id": 1, "doc_id": 1},
				},
			},
		},
	}
)

type mDoc struct {
	StoreId string      `bson:"store_id"`
	DocId   string      `bson:"doc_id"`
	Owner   string      `bson:"owner"`
	Data    interface{} `bson:"data"`
	Raw     string      `bson:"raw"`
}

type mongoImpl struct {
	mongodb.DAO
}

func (m *mongoImpl) Init(conf configx.Values) error {
	if e := model.Init(context.Background(), m.DB()); e != nil {
		return e
	}
	return m.DAO.Init(conf)
}

func (m *mongoImpl) PutDocument(storeID string, doc *docstore.Document) error {
	ctx := context.Background()
	mdoc := m.toMdoc(storeID, doc)
	filter := bson.D{{"store_id", storeID}, {"doc_id", doc.ID}}
	res, e := m.DB().Collection(collDocuments).ReplaceOne(ctx, filter, mdoc, &options.ReplaceOptions{Upsert: &upsert})
	if e != nil {
		return e
	}
	fmt.Println("Docstore upsert", res.UpsertedCount, "modified", res.ModifiedCount)
	return nil
}

func (m *mongoImpl) GetDocument(storeID string, docId string) (*docstore.Document, error) {
	ctx := context.Background()
	filter := bson.D{{"store_id", storeID}, {"doc_id", docId}}
	res := m.DB().Collection(collDocuments).FindOne(ctx, filter)
	if res.Err() != nil {
		return nil, res.Err()
	}
	mdoc := &mDoc{}
	if er := res.Decode(mdoc); er == nil {
		return m.toDocument(mdoc), nil
	} else {
		return nil, er
	}
}

func (m *mongoImpl) QueryDocuments(storeID string, query *docstore.DocumentQuery) (chan *docstore.Document, error) {
	ctx := context.Background()
	filter, err := m.buildFilters(storeID, query)
	if err != nil {
		return nil, err
	}
	cursor, e := m.DB().Collection(collDocuments).Find(ctx, filter)
	if e != nil {
		return nil, e
	}
	res := make(chan *docstore.Document)
	go func() {
		bg := context.Background()
		for cursor.Next(bg) {
			mdoc := &mDoc{}
			if cursor.Decode(mdoc) == nil {
				res <- m.toDocument(mdoc)
			}
		}
		close(res)
	}()
	return res, nil
}

func (m *mongoImpl) ListStores() (ss []string, e error) {
	ctx := context.Background()
	ii, er := m.DB().Collection(collDocuments).Distinct(ctx, "store_id", nil)
	if er != nil {
		e = er
		return
	}
	for _, i := range ii {
		ss = append(ss, i.(string))
	}
	return
}

func (m *mongoImpl) DeleteDocument(storeID string, docID string) error {
	ctx := context.Background()
	filter := bson.D{{"store_id", storeID}, {"doc_id", docID}}
	res, e := m.DB().Collection(collDocuments).DeleteOne(ctx, filter)
	if e != nil {
		return e
	}
	fmt.Println("docstore deleted doc", res.DeletedCount)
	return nil
}

func (m *mongoImpl) DeleteDocuments(storeID string, query *docstore.DocumentQuery) (int, error) {
	ctx := context.Background()
	filter, err := m.buildFilters(storeID, query)
	if err != nil {
		return 0, err
	}
	r, e := m.DB().Collection(collDocuments).DeleteMany(ctx, filter)
	if e != nil {
		return 0, e
	}
	return int(r.DeletedCount), nil
}

func (m *mongoImpl) CountDocuments(storeID string, query *docstore.DocumentQuery) (int, error) {
	ctx := context.Background()
	filter, err := m.buildFilters(storeID, query)
	if err != nil {
		return 0, err
	}
	r, e := m.DB().Collection(collDocuments).CountDocuments(ctx, filter)
	if e != nil {
		return 0, e
	}
	return int(r), nil
}

func (m *mongoImpl) Reset() error {
	return nil
}

func (m *mongoImpl) CloseDAO() error {
	return m.CloseConn()
}

func (m *mongoImpl) toMdoc(storeId string, document *docstore.Document) *mDoc {
	var data interface{}
	d := make(map[string]interface{})
	if er := json.Unmarshal([]byte(document.Data), &d); er == nil {
		data = d
	} else {
		data = document.Data
	}
	return &mDoc{
		StoreId: storeId,
		DocId:   document.ID,
		Owner:   document.Owner,
		Data:    data,
		Raw:     document.Data,
	}
}

func (m *mongoImpl) toDocument(doc *mDoc) *docstore.Document {
	return &docstore.Document{
		ID:    doc.DocId,
		Type:  docstore.DocumentType_JSON,
		Owner: doc.Owner,
		Data:  doc.Raw,
	}
}

func (m *mongoImpl) buildFilters(storeID string, query *docstore.DocumentQuery) (interface{}, error) {

	filter := bson.D{
		{"store_id", storeID},
	}
	if query.ID != "" {
		filter = append(filter, bson.E{Key: "doc_id", Value: query.ID})
	}
	if query.Owner != "" {
		filter = append(filter, bson.E{Key: "owner", Value: query.Owner})
	}
	if query.MetaQuery != "" {
		q, e := bleve.NewQueryStringQuery(query.MetaQuery).Parse()
		if e != nil {
			return nil, e
		}
		if bQ, o := q.(*query2.BooleanQuery); o {
			if cj, o2 := bQ.Must.(*query2.ConjunctionQuery); o2 {
				for _, m := range cj.Conjuncts {
					switch v := m.(type) {
					case *query2.WildcardQuery:
						wc := v.Wildcard
						regexp := ""
						if !strings.HasPrefix(wc, "*") {
							regexp += "^"
						}
						regexp += strings.Trim(wc, "*")
						if !strings.HasSuffix(wc, "*") {
							regexp += "$"
						}
						filter = append(filter, bson.E{Key: "data." + v.Field(), Value: bson.M{"$regex": regexp}})
					case *query2.MatchQuery:
						filter = append(filter, bson.E{Key: "data." + v.Field(), Value: v.Match})
					case *query2.MatchPhraseQuery:
						filter = append(filter, bson.E{Key: "data." + v.Field(), Value: v.MatchPhrase})
					case *query2.NumericRangeQuery:
						if v.Min != nil {
							ref := "$gt"
							if v.InclusiveMin != nil && *v.InclusiveMin {
								ref = "$gte"
							}
							filter = append(filter, bson.E{Key: "data." + v.Field(), Value: bson.M{ref: v.Min}})
						}
						if v.Max != nil {
							ref := "$lt"
							if v.InclusiveMax != nil && *v.InclusiveMax {
								ref = "$lte"
							}
							filter = append(filter, bson.E{Key: "data." + v.Field(), Value: bson.M{ref: v.Max}})
						}
					}
				}
			}
		}
	}

	return filter, nil
}
