# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: fbcrawl.proto

require 'google/protobuf'

Google::Protobuf::DescriptorPool.generated_pool.build do
  add_file("fbcrawl.proto", :syntax => :proto3) do
    add_message "fbcrawl_colly.Empty" do
    end
    add_message "fbcrawl_colly.Pointer" do
      optional :address, :int64, 1
    end
    add_message "fbcrawl_colly.LoginRequest" do
      optional :pointer, :message, 1, "fbcrawl_colly.Pointer"
      optional :email, :string, 2
      optional :password, :string, 3
      optional :totp_secret, :string, 4
    end
    add_message "fbcrawl_colly.LoginResponse" do
      optional :cookies, :string, 1
    end
    add_message "fbcrawl_colly.LoginWithCookiesRequest" do
      optional :pointer, :message, 1, "fbcrawl_colly.Pointer"
      optional :cookies, :string, 2
    end
    add_message "fbcrawl_colly.FetchGroupInfoRequest" do
      optional :pointer, :message, 1, "fbcrawl_colly.Pointer"
      optional :group_username, :string, 2
    end
    add_message "fbcrawl_colly.FetchGroupFeedRequest" do
      optional :pointer, :message, 1, "fbcrawl_colly.Pointer"
      optional :group_id, :int64, 2
      optional :next_cursor, :string, 3
    end
    add_message "fbcrawl_colly.FetchPostRequest" do
      optional :pointer, :message, 1, "fbcrawl_colly.Pointer"
      optional :group_id, :int64, 2
      optional :post_id, :int64, 3
      optional :comment_next_cursor, :string, 4
    end
    add_message "fbcrawl_colly.FetchContentImagesRequest" do
      optional :pointer, :message, 1, "fbcrawl_colly.Pointer"
      optional :post_id, :int64, 2
      optional :next_cursor, :string, 3
    end
    add_message "fbcrawl_colly.FetchImageUrlRequest" do
      optional :pointer, :message, 1, "fbcrawl_colly.Pointer"
      optional :image_id, :int64, 2
    end
    add_message "fbcrawl_colly.FacebookGroup" do
      optional :id, :int64, 1
      optional :name, :string, 2
      optional :member_count, :int64, 3
    end
    add_message "fbcrawl_colly.FacebookUser" do
      optional :id, :int64, 1
      optional :name, :string, 2
    end
    add_message "fbcrawl_colly.FacebookPost" do
      optional :id, :int64, 1
      optional :group, :message, 2, "fbcrawl_colly.FacebookGroup"
      optional :user, :message, 3, "fbcrawl_colly.FacebookUser"
      optional :content, :string, 4
      optional :comments, :message, 5, "fbcrawl_colly.CommentList"
      optional :content_link, :string, 6
      repeated :content_images, :message, 7, "fbcrawl_colly.FacebookImage"
      optional :content_image, :message, 8, "fbcrawl_colly.FacebookImage"
      optional :created_at, :int64, 9
      optional :reaction_count, :int64, 10
      optional :comment_count, :int64, 11
    end
    add_message "fbcrawl_colly.CommentList" do
      repeated :comments, :message, 5, "fbcrawl_colly.FacebookComment"
      optional :next_cursor, :string, 12
    end
    add_message "fbcrawl_colly.FacebookImage" do
      optional :id, :int64, 1
      optional :url, :string, 2
    end
    add_message "fbcrawl_colly.FacebookComment" do
      optional :id, :int64, 1
      optional :post, :message, 2, "fbcrawl_colly.FacebookPost"
      optional :user, :message, 3, "fbcrawl_colly.FacebookUser"
      optional :content, :string, 4
      optional :created_at, :int64, 5
    end
    add_message "fbcrawl_colly.FacebookPostList" do
      repeated :posts, :message, 1, "fbcrawl_colly.FacebookPost"
      optional :next_cursor, :string, 2
    end
    add_message "fbcrawl_colly.FacebookImageList" do
      repeated :images, :message, 1, "fbcrawl_colly.FacebookImage"
      optional :next_cursor, :string, 2
    end
  end
end

module FbcrawlColly
  Empty = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("fbcrawl_colly.Empty").msgclass
  Pointer = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("fbcrawl_colly.Pointer").msgclass
  LoginRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("fbcrawl_colly.LoginRequest").msgclass
  LoginResponse = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("fbcrawl_colly.LoginResponse").msgclass
  LoginWithCookiesRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("fbcrawl_colly.LoginWithCookiesRequest").msgclass
  FetchGroupInfoRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("fbcrawl_colly.FetchGroupInfoRequest").msgclass
  FetchGroupFeedRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("fbcrawl_colly.FetchGroupFeedRequest").msgclass
  FetchPostRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("fbcrawl_colly.FetchPostRequest").msgclass
  FetchContentImagesRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("fbcrawl_colly.FetchContentImagesRequest").msgclass
  FetchImageUrlRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("fbcrawl_colly.FetchImageUrlRequest").msgclass
  FacebookGroup = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("fbcrawl_colly.FacebookGroup").msgclass
  FacebookUser = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("fbcrawl_colly.FacebookUser").msgclass
  FacebookPost = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("fbcrawl_colly.FacebookPost").msgclass
  CommentList = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("fbcrawl_colly.CommentList").msgclass
  FacebookImage = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("fbcrawl_colly.FacebookImage").msgclass
  FacebookComment = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("fbcrawl_colly.FacebookComment").msgclass
  FacebookPostList = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("fbcrawl_colly.FacebookPostList").msgclass
  FacebookImageList = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("fbcrawl_colly.FacebookImageList").msgclass
end