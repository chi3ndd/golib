package elastic

import (
	"encoding/json"
	"errors"
	"reflect"

	es "github.com/olivere/elastic/v7"
)

type service struct {
	model *ES
}

func NewService(conf Config) (Database, error) {
	opts := make([]es.ClientOptionFunc, 0)
	opts = append(opts, es.SetURL(conf.Address))
	opts = append(opts, es.SetSniff(false))
	if conf.Auth.Enable {
		opts = append(opts, es.SetBasicAuth(conf.Auth.Username, conf.Auth.Password))
	}
	con, err := es.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	// Success
	return &service{model: &ES{model: con}}, nil
}

func (con *service) Get(database, _, id string, result interface{}) error {
	res, err := con.model.Get(database, id)
	if err != nil {
		if es.IsNotFound(err) {
			return errors.New(NotFoundError)
		}
		return err
	}
	err = json.Unmarshal(res.Source, result)
	if err != nil {
		return err
	}
	// Success
	return nil
}

func (con *service) Exists(database, _, id string) (bool, error) {
	// Success
	return con.model.Exists(database, id)
}

func (con *service) Count(database, _ string, query Query) (int64, error) {
	// Success
	return con.model.Count(database, query)
}

func (con *service) FindOne(database, _ string, query Query, sorts []string, result interface{}) error {
	res, err := con.model.SearchOffset(database, query, sorts, 0, 1)
	if err != nil || res.Hits == nil {
		return err
	}
	if res.Hits.TotalHits.Value == 0 {
		return errors.New(NotFoundError)
	}
	err = json.Unmarshal(res.Hits.Hits[0].Source, result)
	if err != nil {
		return err
	}
	// Success
	return nil
}

func (con *service) FindPaging(database, _ string, query Query, sorts []string, page, size int, results interface{}) (int64, error) {
	res, err := con.model.SearchPaging(database, query, sorts, page, size)
	if err != nil || res.Hits == nil {
		return 0, err
	}
	if res.Hits.TotalHits.Value == 0 {
		return 0, errors.New(NotFoundError)
	}
	resultType := reflect.TypeOf(results)
	resultValue := reflect.ValueOf(results)
	resultElemType := resultType.Elem().Elem()
	if resultType.Kind() != reflect.Ptr {
		return 0, errors.New(ResultNotAPointer)
	}
	count := res.Hits.TotalHits.Value
	for _, hit := range res.Hits.Hits {
		itemValue := reflect.New(resultElemType)
		err = json.Unmarshal(hit.Source, itemValue.Interface())
		if err != nil {
			return count, err
		}
		resultValue.Elem().Set(reflect.Append(resultValue.Elem(), itemValue.Elem()))
	}
	// Success
	return count, nil
}

func (con *service) FindOffset(database, _ string, query Query, sorts []string, offset, size int, results interface{}) (int64, error) {
	res, err := con.model.SearchOffset(database, query, sorts, offset, size)
	if err != nil || res.Hits == nil {
		return 0, err
	}
	if res.Hits.TotalHits.Value == 0 {
		return 0, errors.New(NotFoundError)
	}
	resultType := reflect.TypeOf(results)
	resultValue := reflect.ValueOf(results)
	resultElemType := resultType.Elem().Elem()
	if resultType.Kind() != reflect.Ptr {
		return 0, errors.New(ResultNotAPointer)
	}
	count := res.Hits.TotalHits.Value
	for _, hit := range res.Hits.Hits {
		itemValue := reflect.New(resultElemType)
		err = json.Unmarshal(hit.Source, itemValue.Interface())
		if err != nil {
			return count, err
		}
		resultValue.Elem().Set(reflect.Append(resultValue.Elem(), itemValue.Elem()))
	}
	// Success
	return count, nil
}

func (con *service) FindScroll(database, _ string, query Query, sorts []string, size int, scrollID, keepAlive string, results interface{}) (string, int64, error) {
	res, err := con.model.SearchScroll(database, query, sorts, size, scrollID, keepAlive)
	if err != nil || res.Hits == nil {
		return "", 0, err
	}
	if res.Hits.TotalHits.Value == 0 {
		return "", 0, errors.New(NotFoundError)
	}
	resultType := reflect.TypeOf(results)
	resultValue := reflect.ValueOf(results)
	resultElemType := resultType.Elem().Elem()
	if resultType.Kind() != reflect.Ptr {
		return "", 0, errors.New(ResultNotAPointer)
	}
	count := res.Hits.TotalHits.Value
	for _, hit := range res.Hits.Hits {
		itemValue := reflect.New(resultElemType)
		err = json.Unmarshal(hit.Source, itemValue.Interface())
		if err != nil {
			return "", count, err
		}
		resultValue.Elem().Set(reflect.Append(resultValue.Elem(), itemValue.Elem()))
	}
	// Success
	return res.ScrollId, count, nil
}

func (con *service) InsertOne(database, _ string, doc Document) error {
	_, err := con.model.Index(database, doc)
	if err != nil {
		return err
	}
	// Success
	return nil
}

func (con *service) InsertMany(database, _ string, docs []Document) error {
	bulk := con.model.Bulk()
	for idx := range docs {
		bulk.Index(database, docs[idx].GetID(), docs[idx])
	}
	err := bulk.Do()
	if err != nil {
		return err
	}
	// Success
	return nil
}

func (con *service) UpdateByID(database, _, id string, update interface{}, upsert bool) error {
	_, err := con.model.Update(database, id, update, upsert)
	if err != nil {
		return err
	}
	// Success
	return nil
}

func (con *service) UpdateOne(database, _ string, query Query, update interface{}, upsert bool) error {
	// TODO
	return nil
}

func (con *service) UpdateMany(database, _ string, query Query, update interface{}, upsert bool) error {
	// TODO
	return nil
}

func (con *service) DeleteByID(database, _, id string) error {
	_, err := con.model.DeleteByID(database, id)
	if err != nil {
		return err
	}
	// Success
	return nil
}

func (con *service) DeleteMany(database, _ string, query Query) error {
	_, err := con.model.DeleteByQuery(database, query)
	if err != nil {
		return err
	}
	// Success
	return nil
}

func (con *service) Aggregate(database, _ string, query Query, name string, aggregation es.Aggregation, result interface{}) error {
	res, err := con.model.Aggregate(database, query, name, aggregation)
	if err != nil {
		return err
	}
	bts, err := res.Aggregations[name].MarshalJSON()
	if err != nil {
		return err
	}
	if err = json.Unmarshal(bts, result); err != nil {
		return err
	}
	// Success
	return nil
}
