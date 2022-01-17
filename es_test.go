package elasticsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	es7 "github.com/elastic/go-elasticsearch/v7"
	iLog "github.com/illidaris/logger"
	"testing"
)

func TestNewClient(t *testing.T) {
	// https://www.jianshu.com/p/075c0ed51053
	iLog.OnlyConsole()
	cfg := es7.Config{
		Addresses: []string{"http://192.168.21.205:9200/"},
		Logger:    NewLogger(false, false),
	}
	es, err := es7.NewClient(cfg)
	if err != nil {
		t.Errorf("Error creating the client: %s", err)
	}
	res, err := es.Info()
	if err != nil {
		t.Errorf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	fmt.Println(res.String())
	var buf bytes.Buffer
	doc := map[string]interface{}{
		"title":   "你看到外面的世界是什么样的？",
		"content": "外面的世界真的很精彩",
	}
	if err := json.NewEncoder(&buf).Encode(doc); err != nil {
		t.Errorf("Error encoding doc.%s", err)
	}
	res, err = es.Create("demo2", "x", &buf, es.Create.WithDocumentType("doc"))
	if err != nil {
		t.Errorf("Error Index response.%s", err)
	}
	defer res.Body.Close()
	fmt.Println(res.String())
}
