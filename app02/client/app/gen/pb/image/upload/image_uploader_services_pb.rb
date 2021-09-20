# Generated by the protocol buffer compiler.  DO NOT EDIT!
# Source: image_uploader.proto for package 'image.uploader'

require 'grpc'
require 'image_uploader_pb'

module Image
  module Uploader
    module ImageUploadService
      class Service

        include ::GRPC::GenericService

        self.marshal_class_method = :encode
        self.unmarshal_class_method = :decode
        self.service_name = 'image.uploader.ImageUploadService'

        # リクエストがストリーム　レスポンスは単一
        rpc :Upload, stream(::Image::Uploader::ImageUploadRequest), ::Image::Uploader::ImageUploadResponse
      end

      Stub = Service.rpc_stub_class
    end
  end
end
