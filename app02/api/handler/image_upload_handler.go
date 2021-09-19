package handler

import (
	"bytes"
	"sync"

	"image/upload/gen/pb"
)

// sync.Mutexの仕組みについて
type ImageUploadHandler struct {
	sync.Mutex
	files map[string][]byte
}

func NewImageUploadHandler() *ImageUploadHandler {
	return &ImageUploadHandler{
		files: make(map[string][]byte),
	}
}

// 画像アップロード処理 メソッド
func (h *ImageUploadHandler) Upload(stream pb.ImageUploadService_UploadServer) error {
	// 最初のリクエストを受け取る
	req, err := stream.Recv()
	if err != nil {
		return err
	}

	// 初回は必ずメタデータが送られる
	meta := req.GetFileMeta()
	filename := meta.Filename
	// UUIDの生成
	u, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	uuid := u.String()

	// 画像データ格納用バッファ
	buf := &bytes.Buffer{}

	// 塊ごとにアップロードされたバイナリをループしながら全て受け取る
	for {
		// Recvは一度リクエストを取得すると、新しいリクエストをクライアントから受け取るまで処理をブロックする
		r, err := stream.Recv()
		// 全てのリクエストを受け取ると、Recvはio.EOFエラーを返す
		if err == io.EOF {
			break
		} if err != nil {
			return err
		}

		// 分割された画像のバイナリを取得
		chunk := r.GetData()
		// バイナリをバッファに追加
		_, err = buf.Write(chunk)
		if err != nil {
			return err
		}
	}

	data := buf.Bytes()
	mimeType := http.DetectContentType(data)

	h.uploaded.files[filename] = data

	err = stream.SendAndClose(&pb.ImageUploadResponse{
		Uuuid: uuid,
		Size: int32(len(data)),
		Filename: filename,
		ContentType: mimeType,
	})
	return err
}
