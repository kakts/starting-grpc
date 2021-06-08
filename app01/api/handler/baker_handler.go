package handler

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"pancake.maker/gen/api"
)

// パンケーキの仕上がりに影響するseedを初期化する
func init() {
	rand.Seed(time.Now().UnixNano())
}

// BakerHandlerはパンケーキを焼く
type BakerHandler struct {
	report *report
}

type report struct {
	sync.Mutex // 複数人が同時に焼いても大丈夫にしておく
	data       map[api.Pancake_Menu]int
}

// NewBakerHandler はBakerHandlerを初期化する
func NewBakerHandler() *BakerHandler {
	// データの保存に使うmapも初期化する
	return &BakerHandler{
		report: &report{
			data: make(map[api.Pancake_Menu]int),
		},
	}
}

// Bake は指定されたメニューのパンケーキを焼いて、焼けたパンをResponseとして返す
func (h *BakerHandler) Bake(ctx context.Context, req *api.BakeRequest) (*api.BakeResponse, error) {
	// バリデーション
	if req.Menu == api.Pancake_UNKNOWN || req.Menu > api.Pancake_SPICY_CURRY {
		// エラーコードはHTTPのステータスコードとは異なる
		/**
		 * 0: 成功
		 * 1: Canceled
		 * ...
		 */
		return nil, status.Errorf(codes.InvalidArgument, "パンケーキを選んでください!")
	}

	// パンケーキを焼いて数を記録する
	now := time.Now()
	/**
	 * 複数リクエストが行われても大丈夫なように,Mutexを使って排他的にアクセスできるようにする
	 * DBを使うのであれば不要
	 */
	h.report.Lock()
	h.report.data[req.Menu] = h.report.data[req.Menu] + 1
	h.report.Unlock()

	// レスポンスを作って返す
	return &api.BakeResponse{
		Pancake: &api.Pancake{
			Menu:           req.Menu,
			ChefName:       "gami",
			TechnicalScore: rand.Float32(),
			CreateTime: &timestamp.Timestamp{
				Seconds: now.Unix(),
				Nanos:   int32(now.Nanosecond()),
			},
		},
	}, nil
}

// Report は焼けたパンケーキの数を報告する
func (h *BakerHandler) Report(ctx context.Context, req *api.ReportRequest) (*api.ReportResponse, error) {

	counts := make([]*api.Report_BakeCount, 0)

	// レポートを作る
	h.report.Lock()
	for k, v := range h.report.data {
		counts = append(counts,
			&api.Report_BakeCount{
				Menu:  k,
				Count: int32(v),
			},
		)
	}

	h.report.Unlock()

	// レスポンスを作って返す
	return &api.ReportResponse{
		Report: &api.Report{
			BakeCounts: counts,
		},
	}, nil
}
