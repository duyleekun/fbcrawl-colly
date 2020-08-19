# Generated by the protocol buffer compiler.  DO NOT EDIT!
# Source: fbcrawl.proto for package 'fbcrawl_colly'

require 'grpc'
require 'fbcrawl_pb'

module FbcrawlColly
  module Grpc
    class Service

      include GRPC::GenericService

      self.marshal_class_method = :encode
      self.unmarshal_class_method = :decode
      self.service_name = 'fbcrawl_colly.Grpc'

      # Sends a greeting
      rpc :Login, FbcrawlColly::LoginRequest, FbcrawlColly::LoginResponse
      rpc :FetchMyGroups, FbcrawlColly::FetchMyGroupsRequest, FbcrawlColly::FacebookGroupList
      rpc :FetchGroupInfo, FbcrawlColly::FetchGroupInfoRequest, FbcrawlColly::FacebookGroup
      rpc :FetchUserInfo, FbcrawlColly::FetchUserInfoRequest, FbcrawlColly::FacebookUser
      rpc :FetchGroupFeed, FbcrawlColly::FetchGroupFeedRequest, FbcrawlColly::FacebookPostList
      rpc :FetchPost, FbcrawlColly::FetchPostRequest, FbcrawlColly::FacebookPost
      rpc :FetchContentImages, FbcrawlColly::FetchContentImagesRequest, FbcrawlColly::FacebookImageList
      rpc :FetchImageUrl, FbcrawlColly::FetchImageUrlRequest, FbcrawlColly::FacebookImage
    end

    Stub = Service.rpc_stub_class
  end
end
