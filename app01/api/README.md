# How to start
via make command

1. Run gRPC Server
```bash
$ make run

go run server/server.go
2021/06/09 02:22:32 Start gRPC server port: 50051
```

2. Request to the server via grpc_cli  
別ターミナルを立ち上げ、grpc_cliを使ってgrpcコールをする

まず grpc_cli lsで立ち上げたip:portで利用できるgRPCサービスを確認する
```
➜  api git:(main) ✗ grpc_cli ls localhost:50051
grpc.reflection.v1alpha.ServerReflection
pancake.maker.PancakeBakerService
```

```
filename: pancake.proto
package: pancake.maker;
service PancakeBakerService {
  rpc Bake(pancake.maker.BakeRequest) returns (pancake.maker.BakeResponse) {}
  rpc Report(pancake.maker.ReportRequest) returns (pancake.maker.ReportResponse) {}
}
```

3. pancake.maker.PancakeBakerService.Bakeを呼び出して指定したメニューのパンを処理する

```
➜  api git:(main) ✗ grpc_cli call localhost:50051 pancake.maker.PancakeBakerService.Bake 'menu: 1'
connecting to localhost:50051
pancake {
  chef_name: "gami"
  menu: CLASSIC
  technical_score: 0.425277114
  create_time {
    seconds: 1629724529
    nanos: 326613000
  }
}
Rpc succeeded with OK status
```
menu: 1を指定することで1番目のメニューのパンを焼く

サーバ起動後から焼いたメニューの個数を表示する

```
Rpc succeeded with OK status
➜  api git:(main) ✗ grpc_cli call localhost:50051 pancake.maker.PancakeBakerService.Report ''
connecting to localhost:50051
report {
  bake_counts {
    menu: CLASSIC
    count: 1
  }
}
Rpc succeeded with OK status
```

他のメニューのパンも処理して、カウントされるのを確認する

```
➜  api git:(main) ✗ grpc_cli call localhost:50051 pancake.maker.PancakeBakerService.Bake 'menu: 2'
connecting to localhost:50051
pancake {
  chef_name: "gami"
  menu: BANANA_AND_WHIP
  technical_score: 0.992827237
  create_time {
    seconds: 1629724631
    nanos: 641098000
  }
}
Rpc succeeded with OK status
➜  api git:(main) ✗ grpc_cli call localhost:50051 pancake.maker.PancakeBakerService.Bake 'menu: 2'
connecting to localhost:50051
pancake {
  chef_name: "gami"
  menu: BANANA_AND_WHIP
  technical_score: 0.540432751
  create_time {
    seconds: 1629724633
    nanos: 493803000
  }
}
Rpc succeeded with OK status
➜  api git:(main) ✗ grpc_cli call localhost:50051 pancake.maker.PancakeBakerService.Report ''     
connecting to localhost:50051
report {
  bake_counts {
    menu: CLASSIC
    count: 1
  }
  bake_counts {
    menu: BANANA_AND_WHIP
    count: 2
  }
}
Rpc succeeded with OK status
```