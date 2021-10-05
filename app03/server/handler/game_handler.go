package handler

import (
	"fmt"
	"sync"

	"reversi/build"
	"reversi/game"
	"reversi/gen/pb"
)

type GameHandler struct {
	sync.RWMutex
	games  map[int32]*game.Game                  // ゲーム情報　盤面などを格納する
	client map[int32][]pb.GameService_PlayServer // 状態変更時にクライアントにストリーミングを返すために格納する
}

func NewGameHandler() *GameHandler {
	return &GameHandler{
		games:  make(map[int32]*game.Game),
		client: make(map[int32][]pb.GameService_PlayServer),
	}
}

// ゲーム開始アクション
func (h *GameHandler) Play(stream pb.GameService_PlayServer) error {
	for {
		// クライアントからリクエスト受信でreqにリクエスト代入される
		req, err := stream.Recv()
		if err != nil {
			return err
		}

		roomID := req.GetRoomId()
		player := build.Player(req.GetPlayer())

		// oneofで複数の型のリクエストがくるのでswitch文で処理
		switch req.GetAction().(type) {
		case *pb.PlayRequest_Start:
			// ゲーム開始リクエスト
			err := h.start(stream, roomID, player)
			if err != nil {
				return err
			}
		case *pb.PlayRequest_Move:
			// 石を置いた時のリクエスト
			action := req.GetMove()
			x := action.GetMove().GetX()
			y := action.GetMove().GetY()
			fmt.Printf("[game_handler] trying to put stone. x=%d, y=%d\n", x, y)
			err := h.move(roomID, x, y, player)
			if err != nil {
				return err
			}
		}

	}
}

// ゲーム開始アクション
func (h *GameHandler) start(stream pb.GameService_PlayServer, roomID int32, me *game.Player) error {
	h.Lock()
	defer h.Unlock()
	// ゲーム情報がなければ作成する
	g := h.games[roomID]
	if g == nil {
		g = game.NewGame(game.None)
		h.games[roomID] = g
		h.client[roomID] = make([]pb.GameService_PlayServer, 0, 2)
	}

	// 自分のクライアントを格納する
	h.client[roomID] = append(h.client[roomID], stream)

	if len(h.client[roomID]) == 2 {
		// 二人揃ったので開始する
		for _, s := range h.client[roomID] {
			// クライアントにゲーム開始を通知する
			// サーバはプレイヤーにPlayREsponse＿REadyEventを返す
			err := s.Send(&pb.PlayResponse{
				Event: &pb.PlayResponse_Ready{
					Ready: &pb.PlayResponse_ReadyEvent{},
				},
			})
			if err != nil {
				return err
			}
		}
		fmt.Printf("game has started room_id=%v\n", roomID)
	} else {
		// まだ揃ってないので待機中であることをクライアントに通知
		err := stream.Send(&pb.PlayResponse{
			Event: &pb.PlayResponse_Waiting{
				Waiting: &pb.PlayResponse_WaitingEvent{},
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// 石を置くアクション
func (h *GameHandler) move(roomID, x, y int32, p *game.Player) error {
	h.Lock()
	defer h.Unlock()

	g := h.games[roomID]
	finished, err := g.Move(x, y, p.Color)
	if err != nil {
		return err
	}
	// エラーでなければ、石を置けている
	fmt.Printf("[move] can put stone x=%d, y=%d", x, y)
	for _, s := range h.client[roomID] {
		// 手が打たれたことをクライアントに通知する
		err := s.Send(&pb.PlayResponse{
			Event: &pb.PlayResponse_Move{
				Move: &pb.PlayResponse_MoveEvent{
					Player: build.PBPlayer(p),
					Move: &pb.Move{
						X: x,
						Y: y,
					},
					Board: build.PBBoard(g.Board),
				},
			},
		})
		if err != nil {
			return err
		}
		if finished {
			// ゲーム終了通知
			err := s.Send(
				&pb.PlayResponse{
					Event: &pb.PlayResponse_Finished{
						Finished: &pb.PlayResponse_FinishedEvent{
							Winner: build.PBColor(g.Winner()),
							Board:  build.PBBoard(g.Board),
						},
					},
				},
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
