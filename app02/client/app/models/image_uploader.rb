# frozen_string_literal:true


require Rails.root.join('app', 'gen', 'pb', 'image', 'upload', 'image_uploader_pb')
require Rails.root.join('app', 'gen', 'pb', 'image', 'upload', 'image_uploader_services_pb')

class ImageUploader
    include ActiveModel::Model

    # 画像アップロード
    def self.chunked_upload(file_path)
        # puts "calling chunked_upload #{file_path}"
        # ストリーミングで送るリクエストを作成
        # それぞれのリクエスト生成はサーバにリクエストを送信するタイミングで遅延して実行される
        reqs = Enumerator.new do |yielder|
            # 最初のリクエスト
            filename = File.basename(file_path)
            puts "filename #{filename}"
            file_meta = Image::Uploader::ImageUploadRequest::FileMeta.new(
                filename: name
            )
            puts "sent name=#{filename}"
            yielder << Image::Uploader::ImageUploadRequest.new(
                file_meta: file_meta
            )

            # チャンクドアップロードリクエスト
            File.open(file_path, 'r') do |f|
                # 100KBごとに分割された画像のチャンクを送る
                while (chunk = f.read(100.kilobytes))
                    puts "sent #{chunk.size}"
                    yielder << Image::Uploader::ImageUploadRequest.new(data: chunk)
                end
            end
        end

        puts "upload start #{file_path}"
        # APIへリクエストを送り、レスポンスを受け取る
        res = stub.upload(reqs)

        # レスポンスをHashにして返す
        {
            uuid: res.uuid,
            size: res.size,
            content_type: res.content_type,
            filename: res.filename
        }
    end
    
    def self.config_dns
        '127.0.0.1:50051'
    end
    
    def self.stub
        Image::Uploader::ImageUploadService::Stub.new(
            config_dns,
            :this_channel_is_insecure,
            timeout: 1,
        )
    end
end
