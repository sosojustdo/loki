package stats

import (
	"context"
	"testing"
	"time"

	"github.com/cortexproject/cortex/pkg/util"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
)

func TestSnapshot(t *testing.T) {
	ctx := NewContext(context.Background())

	GetChunkData(ctx).HeadChunkBytes += 10
	GetChunkData(ctx).HeadChunkLines += 20
	GetChunkData(ctx).DecompressedBytes += 40
	GetChunkData(ctx).DecompressedLines += 20
	GetChunkData(ctx).CompressedBytes += 30
	GetChunkData(ctx).TotalDuplicates += 10

	GetStoreData(ctx).TotalChunksRef += 50
	GetStoreData(ctx).TotalChunksDownloaded += 60
	GetStoreData(ctx).ChunksDownloadTime += time.Second

	fakeIngesterQuery(ctx)
	fakeIngesterQuery(ctx)

	res := Snapshot(ctx, 2*time.Second)
	res.Log(util.Logger)
	expected := Result{
		Ingester: Ingester{
			TotalChunksMatched: 200,
			TotalBatches:       50,
			TotalLinesSent:     60,
			HeadChunkBytes:     10,
			HeadChunkLines:     20,
			DecompressedBytes:  24,
			DecompressedLines:  40,
			CompressedBytes:    60,
			TotalDuplicates:    2,
			TotalReached:       2,
		},
		Store: Store{
			TotalChunksRef:        50,
			TotalChunksDownloaded: 60,
			ChunksDownloadTime:    time.Second.Seconds(),
			HeadChunkBytes:        10,
			HeadChunkLines:        20,
			DecompressedBytes:     40,
			DecompressedLines:     20,
			CompressedBytes:       30,
			TotalDuplicates:       10,
		},
		Summary: Summary{
			ExecTime:                 2 * time.Second.Seconds(),
			BytesProcessedPerSeconds: int64(42),
			LinesProcessedPerSeconds: int64(50),
			TotalBytesProcessed:      int64(84),
			TotalLinesProcessed:      int64(100),
		},
	}
	require.Equal(t, expected, res)
}

func fakeIngesterQuery(ctx context.Context) {
	d, _ := ctx.Value(trailersKey).(*trailerCollector)
	meta := d.addTrailer()

	c, _ := jsoniter.MarshalToString(ChunkData{
		HeadChunkBytes:    5,
		HeadChunkLines:    10,
		DecompressedBytes: 12,
		DecompressedLines: 20,
		CompressedBytes:   30,
		TotalDuplicates:   1,
	})
	meta.Set(chunkDataKey, c)
	i, _ := jsoniter.MarshalToString(IngesterData{
		TotalChunksMatched: 100,
		TotalBatches:       25,
		TotalLinesSent:     30,
	})
	meta.Set(ingesterDataKey, i)
}

func TestResult_Merge(t *testing.T) {
	var res Result

	toMerge := Result{
		Ingester: Ingester{
			TotalChunksMatched: 200,
			TotalBatches:       50,
			TotalLinesSent:     60,
			HeadChunkBytes:     10,
			HeadChunkLines:     20,
			DecompressedBytes:  24,
			DecompressedLines:  40,
			CompressedBytes:    60,
			TotalDuplicates:    2,
			TotalReached:       2,
		},
		Store: Store{
			TotalChunksRef:        50,
			TotalChunksDownloaded: 60,
			ChunksDownloadTime:    time.Second.Seconds(),
			HeadChunkBytes:        10,
			HeadChunkLines:        20,
			DecompressedBytes:     40,
			DecompressedLines:     20,
			CompressedBytes:       30,
			TotalDuplicates:       10,
		},
		Summary: Summary{
			ExecTime:                 2 * time.Second.Seconds(),
			BytesProcessedPerSeconds: int64(42),
			LinesProcessedPerSeconds: int64(50),
			TotalBytesProcessed:      int64(84),
			TotalLinesProcessed:      int64(100),
		},
	}

	res.Merge(toMerge)
	require.Equal(t, toMerge, res)
}
