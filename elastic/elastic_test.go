package elastic

import (
	"fmt"
	"testing"

	"github.com/olivere/elastic/v7"
)

func TestName(t *testing.T) {
	conf := Config{
		Address: "http://127.0.0.1:9288",
		Auth: AuthConfig{
			Enable:   true,
			Username: "elastic",
			Password: "phamthanhha1896",
		},
	}
	es, err := NewService(conf)
	if err == nil {
		agg := elastic.NewTermsAggregation().Field("symbol").Size(1000)
		agg.SubAggregation("top_hit", elastic.NewTopHitsAggregation().Size(1).Sort("time", false))
		var result map[string]interface{}
		err := es.Aggregate("binance-spot", "", map[string]interface{}{"match_all": map[string]interface{}{}}, "agg", agg, &result)
		if err == nil {
			fmt.Println("Done")
		}
	}
}
