// ゲーム盤面のロジック
package game

import "fmt"

// リバーシ盤面の構造体
// Colorの2次元配列
type Board struct {
	Cells [][]Color
}

// TODO game.goなどでも使っているので共通化
const CELL_SIZE int = 10

// 盤面の作成
func NewBoard() *Board {
	// 8*8のセルに4方の壁を加えて、10*10の盤面を二次元配列で作る
	b := &Board{
		Cells: make([][]Color, CELL_SIZE),
	}

	// 盤面の各行の初期化。各行それぞれ10列ある
	for i := 0; i < CELL_SIZE; i++ {
		b.Cells[i] = make([]Color, 10)
	}

	// 盤面の端に壁を設置する
	for i := 0; i < CELL_SIZE; i++ {
		b.Cells[0][i] = Wall // Wallオブジェクトを配置
	}
	// 各行の左・右端に壁を設置
	for i := 1; i < CELL_SIZE-1; i++ {
		b.Cells[i][0] = Wall
		b.Cells[i][CELL_SIZE-1] = Wall
	}

	// 一番最後の行に壁を設置する
	// TODO 終了条件 9の意味は？
	for i := 0; i < 9; i++ {
		b.Cells[9][i]
	}

	// 初期石を置く
	b.Cells[4][4] = White
	b.Cells[5][5] = White
	b.Cells[5][4] = Black
	b.Cells[4][5] = Black

	return b
}

// 石を置く
func (b *Board) PutStone(x int32, y int32, c Color) error {
	// 指定した座標のセルに石を置けるかのチェック
	if !b.CanPutStone(x, y, c) {
		return fmt.Errorf("Can not put stone x=%v, y=%v color=%v", x, y, ColorToStr(c))
	}

	// セルに石を置く
	b.Cells[x][y] = c

	// おいた石の縦・横・斜め方向でひっくり返すことのできる石を全部ひっくり返す
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {

			// 置いた石と同じ座標はチェックしない
			if dx == 0 && dy == 0 {
				continue
			}

			b.TurnStonesByDirection(x, y, c, int32(dx), int32(dy))
		}
	}

	return nil
}

// セルに石を置けるか判定する
func (b *Board) CanPutStone(x, y int32, c Color) bool {
	// すでに石がある場合は石を置けない
	if b.Cells[x][y] != Empty {
		fmt.Printf("[Board.CanPutStone] can not put stone. x=%v, y=%v", x, y)
		return false
	}

	// 縦・横・斜めの各方向をチェック
	// おいた石の縦・横・斜め方向でひっくり返すことのできる石を全部ひっくり返す
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {

			// 置いた石と同じ座標はチェックしない
			if dx == 0 && dy == 0 {
				continue
			}

			//　ひっくり返すことのできる石が1つでもあれば石をおける
			if b.CountTurnableStonesByDirection(x, y, c, int32(dx), int32(dy)) > 0 {
				return true
			}
		}
	}

	// 1つもひっくり返せる石がなければ、石を置けない
	return false
}

// あるセルに石を置いた場合、ある方向にひっくり返すことのできる石がいくつあるか教える
func (b *Board) CountTurnableStonesByDirection(x, y int32, c Color, dx, dy int32) int {
	cnt := 0

	nx := x + dx
	ny := y + dy
	for {
		nc := b.Cells[nx][ny]

		// 壁か自分の石であればループ終了
		if nc != OpponentColor(c) {
			fmt.Println("[Board.CountTurnableStonesByDirection] stop loop")
			break
		}
		// 相手の石なので数え上げ
		cnt++
		nx += dx
		ny += dy
	}

	// その方向にある相手の石がゼロより大きく、かつその先に自分の石がある場合は数を返す
	if cnt > 0 && b.Cells[nx][ny] == c {
		return cnt
	}

	// それ以外の場合はゼロ
	return 0
}

// ある方向の石をひっくり返す。ひっくり返してもいい場合のみ呼ぶ
func (b *Board) TurnStonesByDirection(x, y int32, c Color, dx, dy int32) {
	nx := x + dx
	ny := y + dy

	for {
		nc := b.Cells[nx][ny]

		// 自分の石の場合はループ終了
		if nc != OpponentColor(c) {
			break
		}

		b.Cells[nx][ny] = c

		nx += dx
		ny += dy
	}
}

// 盤面内である色の石をおけるセルの数を数える
func (b *Board) AvailableCellCount(c Color) int {
	cnt := 0

	for i := 1; i < CELL_SIZE-1; i++ {
		for j := 1; j < CELL_SIZE-1; j++ {
			if b.CanPutStone(int32(i), int32(j), c) {
				cnt++
			}
		}
	}

	return cnt
}

// 盤面内に置かれている石の数を数える
func (b *Board) Score(c Color) int {
	cnt := 0

	for i := 1; i < CELL_SIZE-1; i++ {
		for j := 1; j < CELL_SIZE-1; j++ {
			// 自分の石でない場合はスキップ
			if b.Cells[i][j] != c {
				continue
			}
			cnt++
		}
	}
	return cnt
}

// 盤面内で石が置かれていないセルの数を数える
func (b *Board) Rest() int {
	cnt := 0

	for i := 1; i < CELL_SIZE-1; i++ {
		for j := 1; j < CELL_SIZE-1; j++ {
			if b.Cells[i][j] == Empty {
				cnt++
			}
		}
	}
	return cnt
}
