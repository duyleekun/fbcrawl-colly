module FbcrawlColly
  class Client

    def initialize(host_and_port)
      @host_and_port = host_and_port
      @client = new_grpc_client
      @context = nil
    end

    def login(email, password, totp_secret = "")
      cookies = @client.login(FbcrawlColly::LoginRequest.new(email: email, password: password, totp_secret: totp_secret)).cookies
      @context = FbcrawlColly::Context.new(cookies: cookies)
      cookies
    end

    def login_with_cookies(cookies)
      @context = FbcrawlColly::Context.new(cookies: cookies)
    end

    def fetch_user_info(username)
      s = @client.fetch_user_info(FbcrawlColly::FetchUserInfoRequest.new(context: @context, username: username))
    end

    def fetch_my_groups
      s = @client.fetch_my_groups(FbcrawlColly::FetchMyGroupsRequest.new(context: @context))
    end

    def fetch_group_info(group_id_or_username)
      s = @client.fetch_group_info(FbcrawlColly::FetchGroupInfoRequest.new(context: @context, group_username: group_id_or_username))
    end

    def fetch_group_feed(group_id, next_cursor = nil)
      s = @client.fetch_group_feed(FbcrawlColly::FetchGroupFeedRequest.new(context: @context, group_id: group_id, next_cursor: next_cursor))
    end

    def fetch_post(group_id, post_id, comment_next_cursor = nil)
      s = @client.fetch_post(FbcrawlColly::FetchPostRequest.new(context: @context, group_id: group_id, post_id: post_id, comment_next_cursor: comment_next_cursor))
    end

    def fetch_content_images(post_id, next_cursor = nil)
      s = @client.fetch_content_images(FbcrawlColly::FetchContentImagesRequest.new(context: @context, post_id: post_id, next_cursor: next_cursor))
    end

    def fetch_image_url(image_id)
      s = @client.fetch_image_url(FbcrawlColly::FetchImageUrlRequest.new(context: @context, image_id: image_id))
    end

    private

    def new_grpc_client
      FbcrawlColly::Grpc::Stub.new(@host_and_port, :this_channel_is_insecure)
    end
  end
end
