module FbcrawlColly
  class Client

    def initialize
      super
      @client = new_grpc_client
      @colly = @client.init(Google::Protobuf::Empty.new)

      ObjectSpace.define_finalizer(self) do
        new_grpc_client.free_colly(@colly)
      end
    end

    def login(email, password)
      s = @client.login(FbcrawlColly::LoginRequest.new(pointer: @colly, email: email, password: password)).cookies
    end

    def login_with_cookies(cookies)
      s = @client.login_with_cookies(FbcrawlColly::LoginWithCookiesRequest.new(pointer: @colly, cookies: cookies))
    end

    def fetch_group_info(group_id_or_username)
      s = @client.fetch_group_info(FbcrawlColly::FetchGroupInfoRequest.new(pointer: @colly, group_username: group_id_or_username))
    end

    def fetch_group_feed(group_id, next_cursor = nil)
      s = @client.fetch_group_feed(FbcrawlColly::FetchGroupFeedRequest.new(pointer: @colly, group_id: group_id, next_cursor: next_cursor))
    end

    def fetch_post(group_id, post_id, comment_next_cursor = nil)
      s = @client.fetch_post(FbcrawlColly::FetchPostRequest.new(pointer: @colly, group_id: group_id, post_id: post_id, comment_next_cursor: comment_next_cursor))
    end

    def fetch_content_images(post_id, next_cursor = nil)
      s = @client.fetch_content_images(FbcrawlColly::FetchContentImagesRequest.new(pointer: @colly, post_id: post_id, next_cursor: next_cursor))
    end

    def fetch_image_url(image_id)
      s = @client.fetch_image_url(FbcrawlColly::FetchImageUrlRequest.new(pointer: @colly, image_id: image_id))
    end
    private
    def new_grpc_client
      FbcrawlColly::Grpc::Stub.new('fbcrawl.de3.qmanga.com:81', :this_channel_is_insecure)
    end
  end
end
