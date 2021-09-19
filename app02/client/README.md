# README

This README would normally document whatever steps are necessary to get the
application up and running.

Things you may want to cover:

* Ruby version 2.6.0
```
ruby -v
ruby 2.6.3p62 (2019-04-16 revision 67580) [universal.arm64e-darwin20]

➜  client git:(master) ✗ which ruby
/usr/bin/ruby

```

# how to setup
## Initializing rails
4.1.1を参照
```
gem install bundler

# Gem関連ファイル初期化
bundle init

# railsのコメントアウトを外し、railsをインポートするようにする
vim Gemfile

# Gemfileに書かれているパッケージをインストールする
bundle install

# railsアプリの初期化
bundle exec rails new .
```

## .protoファイルからコードを自動生成
- Gemfileに下記のgrpc用ツールが指定されているので、bundle installをする
```bash
# gRPC用gem
gem 'grpc'
gem 'grpc-tools'
```

.protoファイルから、ruby用のclientコードを自動生成します。
生成先は`./app/gen/api/pancake/maker` です
make コマンドで生成用スクリプトを実行します

```
➜  client git:(master) ✗ make gen_proto
bundle exec grpc_tools_ruby_protoc \
                -I ../proto \
                --ruby_out=app/gen/api/pancake/maker \
                --grpc_out=app/gen/api/pancake/maker \
                ../proto/pancake.proto
```

その後、config/application.rbに追記する
```
    # gRPC用に生成したコードのパス指定
    config.paths.add Rails.root.join(
      'app',
      'gen',
      'api',
      'pancake',
      'maker'
    ).to_s, eager_load: true
```
