package game

// マッチングした部屋を表す構造体
type Room struct {
	ID    int32
	Host  *Player
	Guest *Player
}
