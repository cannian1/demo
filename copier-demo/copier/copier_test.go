package copier

import (
	"encoding/json"
	"github.com/jinzhu/copier"
	"testing"
	"time"
)

type SourceStruct struct {
	ID        int
	Name      string
	CreatedAt time.Time
	Tags      []string
	Data      map[string]interface{}
}

type DestStruct struct {
	ID        int
	Name      string
	CreatedAt int64 `json:"-"`
	Tags      []string
	Data      map[string]interface{}
}

func BenchmarkDirectAssignment(b *testing.B) {
	src := SourceStruct{
		ID:        1,
		Name:      "test",
		CreatedAt: time.Now(),
		Tags:      []string{"tag1", "tag2"},
		Data:      map[string]interface{}{"key": "value"},
	}

	var dst DestStruct

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Direct field-by-field assignment
		dst.ID = src.ID
		dst.Name = src.Name
		dst.CreatedAt = src.CreatedAt.Unix()
		dst.Tags = make([]string, len(src.Tags))
		copy(dst.Tags, src.Tags)
		dst.Data = make(map[string]interface{})
		for k, v := range src.Data {
			dst.Data[k] = v
		}
	}
}

func BenchmarkCopierCopy(b *testing.B) {
	src := SourceStruct{
		ID:        1,
		Name:      "test",
		CreatedAt: time.Now(),
		Tags:      []string{"tag1", "tag2"},
		Data:      map[string]interface{}{"key": "value"},
	}

	var dst DestStruct

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Copy(&dst, &src)
		dst.CreatedAt = src.CreatedAt.Unix()
	}
}

func BenchmarkCopierCopyWithOption(b *testing.B) {
	src := SourceStruct{
		ID:        1,
		Name:      "test",
		CreatedAt: time.Now(),
		Tags:      []string{"tag1", "tag2"},
		Data:      map[string]interface{}{"key": "value"},
	}

	var dst DestStruct

	// Test with some options
	ignoreEmpty := func(o *copier.Option) {
		o.IgnoreEmpty = true
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CopyWithOption(&dst, &src, ignoreEmpty, TimeToInt64())
	}
}

func BenchmarkJSONMarshalUnmarshalCopy(b *testing.B) {
	src := SourceStruct{
		ID:        1,
		Name:      "test",
		CreatedAt: time.Now(),
		Tags:      []string{"tag1", "tag2"},
		Data:      map[string]interface{}{"key": "value"},
	}

	var dst DestStruct

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Marshal to JSON
		data, err := json.Marshal(src)
		if err != nil {
			b.Fatal(err)
		}

		// Unmarshal to destination
		err = json.Unmarshal(data, &dst)
		if err != nil {
			b.Fatal(err)
		}
		dst.CreatedAt = src.CreatedAt.Unix()
	}
}
