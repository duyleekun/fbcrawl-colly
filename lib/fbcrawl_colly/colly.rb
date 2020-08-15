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
      s, ptr = FbcrawlColly::FFI.Login(@colly, email, password)
      FbcrawlColly::FFI.free(ptr)
      s
    end

    def login_with_cookies(cookies)
      FbcrawlColly::FFI.LoginWithCookies(@colly, cookies)
    end

    def fetch_group_info(group_id)
      s, ptr = FbcrawlColly::FFI.FetchGroupInfo(@colly, group_id)
      list = FbcrawlColly::FacebookGroup.decode(s)
      FbcrawlColly::FFI.free(ptr)
      list
    end

    def fetch_group_feed(group_id)
      s, ptr = FbcrawlColly::FFI.FetchGroupFeed(@colly, group_id)
      list = FbcrawlColly::FacebookPostList.decode(s)
      FbcrawlColly::FFI.free(ptr)
      list
    end

    def fetch_post(group_id, post_id)
      s, ptr = FbcrawlColly::FFI.FetchPost(@colly, group_id, post_id)
      post = FbcrawlColly::FacebookPost.decode(s)
      FbcrawlColly::FFI.free(ptr)
      post
    end

    def fetch_content_images(post_id)
      s, ptr = FbcrawlColly::FFI.FetchContentImages(@colly, post_id)
      imageList = FbcrawlColly::FacebookImageList.decode(s)
      FbcrawlColly::FFI.free(ptr)
      imageList
    end

    def fetch_image_url(image_id)
      s, ptr = FbcrawlColly::FFI.FetchImageUrl(@colly, image_id)
      image = FbcrawlColly::FacebookImage.decode(s)
      FbcrawlColly::FFI.free(ptr)
      image
    end
  end
end
