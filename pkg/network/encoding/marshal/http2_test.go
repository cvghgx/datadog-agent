// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package marshal

import (
	"runtime"
	"testing"

	model "github.com/DataDog/agent-payload/v5/process"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/DataDog/datadog-agent/pkg/network"
	"github.com/DataDog/datadog-agent/pkg/network/protocols/http"
	"github.com/DataDog/datadog-agent/pkg/process/util"
)

type connTag = uint64

// ConnTag constant must be the same for all platform
const (
	tagGnuTLS  connTag = 0x01 // network.ConnTagGnuTLS
	tagOpenSSL connTag = 0x02 // network.ConnTagOpenSSL
	tagTLS     connTag = 0x10 // network.ConnTagTLS
)

type HTTP2Suite struct {
	suite.Suite
}

func TestHTTP2Stats(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("the feature is only supported on linux.")
	}
	suite.Run(t, &HTTP2Suite{})
}

func (s *HTTP2Suite) TestFormatHTTP2Stats() {
	t := s.T()
	t.Run("status code", func(t *testing.T) {
		testFormatHTTP2Stats(t, true)
	})
	t.Run("status class", func(t *testing.T) {
		testFormatHTTP2Stats(t, false)
	})
}

func testFormatHTTP2Stats(t *testing.T, aggregateByStatusCode bool) {
	var (
		clientPort  = uint16(52800)
		serverPort  = uint16(8080)
		localhost   = util.AddressFromString("127.0.0.1")
		statusCodes = []uint16{101, 202, 307, 404, 503}
	)

	httpKey1 := http.NewKey(
		localhost,
		localhost,
		clientPort,
		serverPort,
		[]byte("/testpath-1"),
		true,
		http.MethodGet,
	)
	http2Stats1 := http.NewRequestStats(aggregateByStatusCode)
	for _, i := range statusCodes {
		http2Stats1.AddRequest(i, 10, 1<<(i/100-1), nil)
	}

	httpKey2 := httpKey1
	httpKey2.Path = http.Path{
		Content:  http.Interner.GetString("/testpath-2"),
		FullPath: true,
	}
	http2Stats2 := http.NewRequestStats(aggregateByStatusCode)
	for _, i := range statusCodes {
		http2Stats2.AddRequest(i, 20, 1<<(i/100-1), nil)
	}

	in := &network.Connections{
		BufferedData: network.BufferedData{
			Conns: []network.ConnectionStats{
				{
					Source: localhost,
					Dest:   localhost,
					SPort:  clientPort,
					DPort:  serverPort,
					HTTP2Stats: []network.USMKeyValue[http.Key, *http.RequestStats]{
						{Key: httpKey1, Value: http2Stats1},
						{Key: httpKey2, Value: http2Stats2},
					},
				},
			},
		},
	}
	out := &model.HTTP2Aggregations{
		EndpointAggregations: []*model.HTTPStats{
			{
				Path:              "/testpath-1",
				Method:            model.HTTPMethod_Get,
				FullPath:          true,
				StatsByStatusCode: make(map[int32]*model.HTTPStats_Data),
			},
			{
				Path:              "/testpath-2",
				Method:            model.HTTPMethod_Get,
				FullPath:          true,
				StatsByStatusCode: make(map[int32]*model.HTTPStats_Data),
			},
		},
	}

	for _, statusCode := range statusCodes {
		code := int32(http2Stats1.NormalizeStatusCode(statusCode))
		out.EndpointAggregations[0].StatsByStatusCode[code] = &model.HTTPStats_Data{Count: 1, FirstLatencySample: 10, Latencies: nil}
		out.EndpointAggregations[1].StatsByStatusCode[code] = &model.HTTPStats_Data{Count: 1, FirstLatencySample: 20, Latencies: nil}
	}

	http2Encoder := newHTTP2Encoder()
	aggregations, tags, _ := getHTTP2Aggregations(t, http2Encoder, in.Conns[0])

	require.NotNil(t, aggregations)
	assert.ElementsMatch(t, out.EndpointAggregations, aggregations.EndpointAggregations)

	assert.Equal(t, uint64((1<<len(statusCodes))-1), tags)
}

func (s *HTTP2Suite) TestFormatHTTP2StatsByPath() {
	t := s.T()
	t.Run("status code", func(t *testing.T) {
		testFormatHTTP2StatsByPath(t, true)
	})
	t.Run("status class", func(t *testing.T) {
		testFormatHTTP2StatsByPath(t, false)
	})
}

func testFormatHTTP2StatsByPath(t *testing.T, aggregateByStatusCode bool) {
	http2ReqStats := http.NewRequestStats(aggregateByStatusCode)

	http2ReqStats.AddRequest(100, 12.5, 0, nil)
	http2ReqStats.AddRequest(100, 12.5, tagGnuTLS, nil)
	http2ReqStats.AddRequest(405, 3.5, tagOpenSSL, nil)
	http2ReqStats.AddRequest(405, 3.5, 0, nil)

	// Verify the latency data is correct prior to serialization

	latencies := http2ReqStats.Data[http2ReqStats.NormalizeStatusCode(100)].Latencies
	assert.Equal(t, 2.0, latencies.GetCount())
	verifyQuantile(t, latencies, 0.5, 12.5)

	latencies = http2ReqStats.Data[http2ReqStats.NormalizeStatusCode(405)].Latencies
	assert.Equal(t, 2.0, latencies.GetCount())
	verifyQuantile(t, latencies, 0.5, 3.5)

	key := http.NewKey(
		util.AddressFromString("10.1.1.1"),
		util.AddressFromString("10.2.2.2"),
		60000,
		80,
		[]byte("/testpath"),
		true,
		http.MethodGet,
	)

	payload := &network.Connections{
		BufferedData: network.BufferedData{
			Conns: []network.ConnectionStats{
				{
					Source: util.AddressFromString("10.1.1.1"),
					Dest:   util.AddressFromString("10.2.2.2"),
					SPort:  60000,
					DPort:  80,
					HTTP2Stats: []network.USMKeyValue[http.Key, *http.RequestStats]{
						{Key: key, Value: http2ReqStats},
					},
				},
			},
		},
	}
	http2Encoder := newHTTP2Encoder()
	http2Aggregations, tags, _ := getHTTP2Aggregations(t, http2Encoder, payload.Conns[0])

	require.NotNil(t, http2Aggregations)
	endpointAggregations := http2Aggregations.EndpointAggregations
	require.Len(t, endpointAggregations, 1)
	assert.Equal(t, "/testpath", endpointAggregations[0].Path)
	assert.Equal(t, model.HTTPMethod_Get, endpointAggregations[0].Method)

	assert.Equal(t, tagGnuTLS|tagOpenSSL, tags)

	// Deserialize the encoded latency information & confirm it is correct
	statsByResponseStatus := endpointAggregations[0].StatsByStatusCode
	assert.Len(t, statsByResponseStatus, 2)

	serializedLatencies := statsByResponseStatus[int32(http2ReqStats.NormalizeStatusCode(100))].Latencies
	sketch := unmarshalSketch(t, serializedLatencies)
	assert.Equal(t, 2.0, sketch.GetCount())
	verifyQuantile(t, sketch, 0.5, 12.5)

	serializedLatencies = statsByResponseStatus[int32(http2ReqStats.NormalizeStatusCode(405))].Latencies
	sketch = unmarshalSketch(t, serializedLatencies)
	assert.Equal(t, 2.0, sketch.GetCount())
	verifyQuantile(t, sketch, 0.5, 3.5)

	_, exists := statsByResponseStatus[200]
	assert.False(t, exists)
}

func (s *HTTP2Suite) TestHTTP2IDCollisionRegression() {
	t := s.T()
	t.Run("status code", func(t *testing.T) {
		testHTTP2IDCollisionRegression(t, true)
	})
	t.Run("status class", func(t *testing.T) {
		testHTTP2IDCollisionRegression(t, false)
	})
}

func testHTTP2IDCollisionRegression(t *testing.T, aggregateByStatusCode bool) {
	http2Stats := http.NewRequestStats(aggregateByStatusCode)
	assert := assert.New(t)
	httpKey := http.NewKey(
		util.AddressFromString("1.1.1.1"),
		util.AddressFromString("2.2.2.2"),
		60000,
		80,
		[]byte("/"),
		true,
		http.MethodGet,
	)
	http2Stats.AddRequest(104, 1.0, 0, nil)

	connections := []network.ConnectionStats{
		{
			Source: util.AddressFromString("1.1.1.1"),
			SPort:  60000,
			Dest:   util.AddressFromString("2.2.2.2"),
			DPort:  80,
			Pid:    1,
			HTTP2Stats: []network.USMKeyValue[http.Key, *http.RequestStats]{
				{Key: httpKey, Value: http2Stats},
			},
		},
		{
			Source: util.AddressFromString("1.1.1.1"),
			SPort:  60000,
			Dest:   util.AddressFromString("2.2.2.2"),
			DPort:  80,
			Pid:    2,
		},
	}

	http2Encoder := newHTTP2Encoder()

	// assert that the first connection matching the HTTP2 data will get
	// back a non-nil result
	aggregations, _, _ := getHTTP2Aggregations(t, http2Encoder, connections[0])
	assert.Equal("/", aggregations.EndpointAggregations[0].Path)
	assert.Equal(uint32(1), aggregations.EndpointAggregations[0].StatsByStatusCode[int32(http2Stats.NormalizeStatusCode(104))].Count)

	// assert that the other connections sharing the same (source,destination)
	// addresses but different PIDs *won't* be associated with the HTTP2 stats
	// object
	streamer := NewProtoTestStreamer[*model.Connection]()
	http2Encoder.WriteHTTP2AggregationsAndTags(connections[1], model.NewConnectionBuilder(streamer))

	var conn model.Connection
	streamer.Unwrap(t, &conn)
	assert.Empty(conn.Http2Aggregations)
}

func getHTTP2Aggregations(t *testing.T, encoder *http2Encoder, c network.ConnectionStats) (*model.HTTP2Aggregations, uint64, map[string]struct{}) {
	streamer := NewProtoTestStreamer[*model.Connection]()
	staticTags, dynamicTags := encoder.WriteHTTP2AggregationsAndTags(c, model.NewConnectionBuilder(streamer))

	var conn model.Connection
	streamer.Unwrap(t, &conn)

	var aggregations model.HTTP2Aggregations
	err := proto.Unmarshal(conn.Http2Aggregations, &aggregations)
	require.NoError(t, err)

	return &aggregations, staticTags, dynamicTags
}
