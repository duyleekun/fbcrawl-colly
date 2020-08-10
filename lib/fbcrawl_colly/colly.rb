require 'ffi'
require_relative '../fbcrawl_pb'
require_relative './ffi'

module FbcrawlColly
  class Colly
    def initialize
      super
      @colly = ::FFI::AutoPointer.new(FbcrawlColly::FFI::Init(), FbcrawlColly::FFI.method(:FreeColly))
    end

    def login(email, password)
      FbcrawlColly::FFI.Login(@colly, email, password)
    end

    def fetch_group_feed(group_id)
      s, ptr = FbcrawlColly::FFI.FetchGroupFeed(@colly, group_id)
      list = FacebookPostList.decode(s)
      FbcrawlColly::FFI.free(ptr)
      list
    end

    def fetch_post(group_id, post_id)
      FbcrawlColly::FFI.Login(@colly, email, password)
      s, ptr = FbcrawlColly::FFI.FetchPost(@colly, group_id, post_id)
      post = FacebookPost.decode(s)
      FbcrawlColly::FFI.free(ptr)
      post
    end
  end
end
