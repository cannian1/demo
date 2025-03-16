package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	esV8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/update"
	"github.com/elastic/go-elasticsearch/v8/typedapi/some"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

func main() {
	// ES 配置
	cfg := esV8.Config{
		Addresses: []string{
			"http://localhost:9200", // 本机的 es 地址
		},
	}

	// 创建客户端连接
	// esV8.NewClient() 低级操作的客户端
	// esV8.NewDefaultClient() 默认连接 本地:9200
	client, err := esV8.NewTypedClient(cfg)
	if err != nil {
		fmt.Printf("elasticsearch.NewTypedClient failed, err:%v\n", err)
		return
	}

	// 创建索引
	//createIndex(add_client)

	// 创建文档
	//indexDocument(add_client)
	//indexDocument2(add_client)

	// 查询文档
	//getDocument(add_client, strconv.Itoa(2))

	// 搜索所有文档
	searchDocument(client)

	// 聚合搜索
	//aggregationDemo(add_client)

	// 用结构体更新文档
	//updateDocument(add_client)
	// 用 json 更新文档
	//updateDocument2(add_client)

	// 删除文档
	//deleteDocument(add_client)

	// 删除索引
	//deleteIndex(add_client)
}

// createIndex 创建索引
func createIndex(client *esV8.TypedClient) {
	resp, err := client.Indices.
		Create("my-review-1").
		Do(context.Background())
	if err != nil {
		fmt.Printf("create index failed, err:%v\n", err)
		return
	}
	fmt.Printf("index:%#v\n", resp.Index)
}

// Review 评价数据
type Review struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"userID"`
	Score       uint8     `json:"score"`
	Content     string    `json:"content"`
	Tags        []Tag     `json:"tags"`
	Status      int       `json:"status"`
	PublishTime time.Time `json:"publishDate"`
}

// Tag 评价标签
type Tag struct {
	Code  int    `json:"code"`
	Title string `json:"title"`
}

// indexDocument 索引文档
func indexDocument(client *esV8.TypedClient) {
	// 定义 document 结构体对象
	d1 := Review{
		ID:      1,
		UserID:  147982601,
		Score:   5,
		Content: "这是一个好评！",
		Tags: []Tag{
			{1000, "好评"},
			{1100, "物超所值"},
			{9000, "有图"},
		},
		Status:      2,
		PublishTime: time.Now(),
	}

	// 添加文档
	resp, err := client.Index("my-review-1").
		Id(strconv.FormatInt(d1.ID, 10)).
		Document(d1).
		Do(context.Background())
	if err != nil {
		fmt.Printf("indexing document failed, err:%v\n", err)
		return
	}
	fmt.Printf("result:%#v\n", resp.Result)
}

func indexDocument2(client *esV8.TypedClient) {
	// 定义 document 结构体对象
	d1 := Review{
		ID:      2,
		UserID:  147982602,
		Score:   1,
		Content: "这是一个差评！",
		Tags: []Tag{
			{2000, "差评"},
		},
		Status:      2,
		PublishTime: time.Now(),
	}

	// 添加文档
	resp, err := client.Index("my-review-1").
		Id(strconv.FormatInt(d1.ID, 10)).
		Document(d1).
		Do(context.Background())
	if err != nil {
		fmt.Printf("indexing document failed, err:%v\n", err)
		return
	}
	fmt.Printf("result:%#v\n", resp.Result)
}

// getDocument 获取文档
func getDocument(client *esV8.TypedClient, id string) {
	resp, err := client.Get("my-review-1", id).
		Do(context.Background())
	if err != nil {
		fmt.Printf("get document by id failed, err:%v\n", err)
		return
	}
	fmt.Printf("fileds:%s\n", resp.Source_)
}

// searchDocument 搜索所有文档
func searchDocument(client *esV8.TypedClient) {
	// 搜索文档
	resp, err := client.Search().
		Index("my-review-1").
		Query(&types.Query{
			MatchAll: &types.MatchAllQuery{},
		}).
		Do(context.Background())
	if err != nil {
		fmt.Printf("search document failed, err:%v\n", err)
		return
	}
	fmt.Printf("total: %d\n", resp.Hits.Total.Value)
	// 遍历所有结果
	for _, hit := range resp.Hits.Hits {
		fmt.Printf("%s\n", hit.Source_)
	}
}

// aggregationDemo 聚合
func aggregationDemo(client *esV8.TypedClient) {
	avgScoreAgg, err := client.Search().
		Index("my-review-1").
		Request(
			&search.Request{
				Size: some.Int(0),
				Aggregations: map[string]types.Aggregations{
					"avg_score": { // 将所有文档的 score 的平均值聚合为 avg_score
						Avg: &types.AverageAggregation{
							Field: some.String("score"),
						},
					},
				},
			},
		).Do(context.Background())
	if err != nil {
		fmt.Printf("aggregation failed, err:%v\n", err)
		return
	}

	aggr := avgScoreAgg.Aggregations["avg_score"]
	fmt.Printf("avgScore:%#v\n", aggr)
	// avgScore:&types.AvgAggregate{Meta:types.Metadata(nil), Value:(*types.Float64)(0xc00028c658), ValueAsString:(*string)(nil)}

	if avgAgg, ok := aggr.(*types.AvgAggregate); ok {
		if avgAgg.Value != nil {
			fmt.Printf("Avg score value: %f\n", *avgAgg.Value)
		} else {
			fmt.Println("Avg score value is nil")
		}
	} else {
		fmt.Println("Failed to assert aggr to *types.AvgAggregate")
	}
}

// updateDocument 更新文档
func updateDocument(client *esV8.TypedClient) {
	// 修改后的结构体变量
	d1 := Review{
		ID:      1,
		UserID:  147982601,
		Score:   5,
		Content: "这是一个修改后的好评！", // 有修改
		Tags: []Tag{ // 有修改
			{1000, "好评"},
			{9000, "有图"},
		},
		Status:      2,
		PublishTime: time.Now(),
	}

	resp, err := client.Update("my-review-1", "1").
		Doc(d1). // 使用结构体变量更新
		Do(context.Background())
	if err != nil {
		fmt.Printf("update document failed, err:%v\n", err)
		return
	}
	fmt.Printf("result:%v\n", resp.Result)
}

// updateDocument2 更新文档
func updateDocument2(client *esV8.TypedClient) {
	// 修改后的JSON字符串
	str := `{
					"id":1,
					"userID":147982601,
					"score":5,
					"content":"这是一个二次修改后的好评！",
					"tags":[
						{
							"code":1000,
							"title":"好评"
						},
						{
							"code":9000,
							"title":"有图"
						}
					],
					"status":2,
					"publishDate":"2023-12-10T15:27:18.219385+08:00"
				}`
	// 直接使用JSON字符串更新
	resp, err := client.Update("my-review-1", "1").
		Request(&update.Request{
			Doc: json.RawMessage(str),
		}).
		Do(context.Background())
	if err != nil {
		fmt.Printf("update document failed, err:%v\n", err)
		return
	}
	fmt.Printf("result:%v\n", resp.Result)
}

// deleteDocument 删除文档
func deleteDocument(client *esV8.TypedClient) {
	resp, err := client.Delete("my-review-1", "1").
		Do(context.Background())
	if err != nil {
		fmt.Printf("delete document failed, err:%v\n", err)
		return
	}
	fmt.Printf("result:%v\n", resp.Result)
}

// deleteIndex 删除 index
func deleteIndex(client *esV8.TypedClient) {
	resp, err := client.Indices.
		Delete("my-review-1").
		Do(context.Background())
	if err != nil {
		fmt.Printf("delete document failed, err:%v\n", err)
		return
	}
	fmt.Printf("Acknowledged:%v\n", resp.Acknowledged)
}
