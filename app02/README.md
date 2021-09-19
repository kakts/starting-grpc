# app02 image streaming uploader
chapter06

- server: Go
- client: Ruby


# 仕様
- clientからモデルのメソッド引数に画像ファイルパスを渡してアップロードする
- クライアントは第4章(app01)と同様にRailsコンソールから実行する
- 画像は100KBごとに分割してAPIサーバに送る(chunked upload)
- 初回リクエストでメタデータ送信。2回目から最後のリクエストまでは画像のチャンクアップロードを行う
- APIサーバ(Go)は全てのリクエストを受け取ったら、画像ファイルのバイナリをサーバのメモリ内に保存
- APIサーバはレスポンスとしてクライアントに画像ファイルのメタデータを返す

# how to develop server
- protoファイルでapi定義実装
- protoファイルからGoのソース生成
```
make protoc
```
これによりapi/gen/pb配下にimage_uploader.pb.goが生成