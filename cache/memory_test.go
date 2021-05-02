package cache

import (
	"bytes"
	"fmt"
	"reflect"
	"sync"
	"testing"
)

func TestMemoryCache_Get(t *testing.T) {
	type fields struct {
		m map[string][]byte
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "should be able to get cached media",
			fields: fields{m: map[string][]byte{
				"/test.jpeg": []byte("result"),
			}},
			args:    args{key: "/test.jpeg"},
			want:    []byte("result"),
			wantErr: false,
		},
		{
			name:    "should return error if not defined in cache",
			fields:  fields{m: map[string][]byte{}},
			args:    args{key: "/test.jpeg"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &MemoryCache{
				mtx: sync.RWMutex{},
				m:   tt.fields.m,
			}
			got, err := c.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemoryCache.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemoryCache.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryCache_Set(t *testing.T) {
	type fields struct {
		m map[string][]byte
	}
	type args struct {
		key string
		v   []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "should be able to set cached media",
			fields:  fields{m: map[string][]byte{}},
			args:    args{key: "/test.jpeg", v: []byte("result")},
			wantErr: false,
		},
		{
			name: "should be able to replace cached media",
			fields: fields{m: map[string][]byte{
				"/test.jpeg": []byte("result"),
			}},
			args:    args{key: "/test.jpeg", v: []byte("result2")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &MemoryCache{
				mtx: sync.RWMutex{},
				m:   tt.fields.m,
			}
			if err := c.Set(tt.args.key, tt.args.v, 0); (err != nil) != tt.wantErr {
				t.Errorf("MemoryCache.Set() error = %v, wantErr %v", err, tt.wantErr)
			}

			v, err := c.Get(tt.args.key)
			if err != nil {
				panic(err)
			}

			if !bytes.Equal(v, tt.args.v) {
				t.Errorf("result differ")
			}
		})
	}
}

func BenchmarkMemoryCache_Get(b *testing.B) {
	testCacheKey := "test-key"

	c := &MemoryCache{
		mtx: sync.RWMutex{},
		m:   map[string][]byte{},
	}

	_ = c.Set(testCacheKey, []byte("a good media value"), 0)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = c.Get(testCacheKey)
	}
}

func BenchmarkMemoryCache_SetSameKey(b *testing.B) {
	testCacheKey := "test-key"

	c := &MemoryCache{
		mtx: sync.RWMutex{},
		m:   map[string][]byte{},
	}

	_ = c.Set(testCacheKey, []byte("a good media value"), 0)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = c.Set(testCacheKey, []byte("a good media value"), 0)
	}
}

func BenchmarkMemoryCache_SetNewKey(b *testing.B) {
	testCacheKey := "test-key"

	c := &MemoryCache{
		mtx: sync.RWMutex{},
		m:   map[string][]byte{},
	}

	_ = c.Set(testCacheKey, []byte("a good media value"), 0)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = c.Set(fmt.Sprintf("%s-%d", testCacheKey, i), []byte("a good media value"), 0)
	}
}
