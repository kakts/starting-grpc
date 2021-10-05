package game

import "fmt"

// リバーシゲーム自体の構造体
type Game struct {
	Board    *Board
	started  bool
	finished bool
	me       Color
}

// ゲームの初期化
func NewGame(me Color) *Game {
	return &Game{
		Board: NewBoard(),
		me:    me,
	}
}

// 手を打ち、その後盤面を出力する
// 返り値として、ゲームが終了したかどうかを返す
func (g *Game) Move(x, y int32, c Color) (bool, error) {
	fmt.Println("[Move] Check if game had finished")
	if g.finished {
		return true, nil
	}

	fmt.Printf("[Move] Before PutStone x=%d, y=%d", x, y)
	err := g.Board.PutStone(x, y, c)
	if err != nil {
		return false, err
	}
	g.Display(g.me)
	if g.IsGameOver() {
		fmt.Println("Game finished!")
		g.finished = true
		return true, nil
	}

	return false, nil
}

// ゲームが終了したかの判定
// 黒と白双方における場所がなければ終了
func (g *Game) IsGameOver() bool {
	if g.Board.AvailableCellCount(Black) > 0 {
		return false
	}

	if g.Board.AvailableCellCount(White) > 0 {
		return false
	}

	return true
}

// 勝者の色を返す 引き分けの場合はNoneを返す
func (g *Game) Winner() Color {
	black := g.Board.Score(Black)
	white := g.Board.Score(White)

	if black == white {
		return None
	} else if black > white {
		return Black
	}
	return White
}

// 盤面の出力
func (g *Game) Display(me Color) {
	fmt.Println("[Display] displaying game board.")
	fmt.Println("")
	if me != None {
		fmt.Printf("You: %v\n", ColorToStr(me))
	}

	fmt.Print("  | ")

	// ?
	rs := []rune("ABCDEFGH")
	for i, r := range rs {
		fmt.Printf("%v", string(r))
		if i < len(rs)-1 {
			fmt.Print(" | ")
		}
	}

	// 行末
	fmt.Println("\n")
	fmt.Println("--------------")

	for j := 1; j < cell_size-1; j++ {
		fmt.Printf("%d", j)
		fmt.Print(" | ")
		for i := 1; i < cell_size-1; i++ {
			fmt.Print(ColorToStr(g.Board.Cells[i][j]))
			fmt.Print(" | ")
		}
		fmt.Print("\n")
	}
	fmt.Println("--------------")

	fmt.Printf("Score: BLACK=%d, WHITE=%d REST=%d\n", g.Board.Score(Black), g.Board.Score(White), g.Board.Rest())
	fmt.Print("\n")
}
