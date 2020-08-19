require "test_helper"

class FbcrawlCollyTest < Minitest::Test
  DEFAULT_GROUP_ID = 658075901719147
  EMAIL = ENV["FACEBOOK_EMAIL"]
  PASSWORD = ENV["FACEBOOK_PASSWORD"]
  MOTP_SECRET = ENV["FACEBOOK_MOTP_SECRET"]
  HOST_AND_PORT = ENV["FBCRAWL_HOST_PORT"]

  def setup
    super
    puts login_cookies
  end

  def test_init_should_return_pointer
    assert new_colly != nil, "Colly should not be nil"
  end

  def test_login_ok
    assert login_cookies.size > 0
  end

  def test_my_groups
    p = new_logged_in_colly.fetch_my_groups
    p.groups.each do |g|
      assert g.id > 0
      assert g.name
    end
  end

  def test_group_info
    p = new_logged_in_colly.fetch_group_info "fbcolly"
    assert p.name.size > 0
    assert p.member_count > 0
    assert p.id > 0
    puts p
  end

  def test_group_feed
    p = new_logged_in_colly.fetch_group_feed DEFAULT_GROUP_ID
    assert p.posts.size > 0
  end

  def test_fetch_user_info
    p = new_logged_in_colly.fetch_user_info "duyleekun"
    puts p
  end

  def test_text_only_post
    p = new_logged_in_colly.fetch_post DEFAULT_GROUP_ID, 658076021719135
    assert p
    assert p.user.id
    assert p.id
    assert_equal p.content, 'test_text_only_post'
    assert p.created_at > 0
    assert p.reaction_count > 0
    first_comment = p.comments.comments.first
    assert first_comment
    assert first_comment.id
    assert first_comment.user.username
    assert first_comment.content.size > 0
    assert first_comment.created_at > 0
    puts Time.at first_comment.created_at
  end


  def test_text_post_with_background
    p = new_logged_in_colly.fetch_post DEFAULT_GROUP_ID, 659998108193593
    puts p
    assert p
    assert p.user.id
    assert p.id
    assert_equal p.content, 'test_text_post_with_background'
  end

  def test_link_post
    p = new_logged_in_colly.fetch_post DEFAULT_GROUP_ID, 660007198192684
    puts p
    assert p
    assert p.user.id
    assert p.id
    assert_equal p.content, 'test_link_post'
    assert_includes p.content_link, 'vnexpress.net'
  end

  def test_one_photo_post
    p = new_logged_in_colly.fetch_post DEFAULT_GROUP_ID, 660012564858814
    puts p
    assert p
    assert p.user.id
    assert p.id
    assert_equal p.content, 'test_one_photo_post'
    assert p.content_image
    pic = new_logged_in_colly.fetch_image_url(p.content_image.id)
    assert pic.url
  end

  def test_three_photo_post
    p = new_logged_in_colly.fetch_post DEFAULT_GROUP_ID, 660012668192137
    puts p
    assert p
    assert p.user.id
    assert p.id
    assert_equal p.content, 'test_three_photo_post'
    assert p.content_image
    assert_equal 3, p.content_images.size
  end

  def test_five_photo_post
    p = new_logged_in_colly.fetch_post DEFAULT_GROUP_ID, 660012811525456
    puts p
    assert p
    assert p.user.id
    assert p.id
    assert_equal p.content, 'test_five_photo_post'
    assert p.content_image
    assert_equal 5, p.content_images.size
  end

  def test_ten_photo_post
    p = new_logged_in_colly.fetch_post DEFAULT_GROUP_ID, 660017394858331
    puts p
    assert p
    assert p.user.id
    assert p.id
    assert_equal p.content, 'test_ten_photo_post'
    assert p.content_image
    assert_equal 5, p.content_images.size
    images = new_logged_in_colly.fetch_content_images 660017394858331
    assert_equal images.images.size, 10
  end

  def test_fifty_photo_post
    p = new_logged_in_colly.fetch_post DEFAULT_GROUP_ID, 660082011518536
    puts p
    assert p
    assert p.user.id
    assert p.id
    assert_equal p.content, 'test_fifty_photo_post'
    assert p.content_image
    assert_equal 5, p.content_images.size
    images = []
    next_cursor = nil
    loop do
      r = new_logged_in_colly.fetch_content_images 660082011518536, next_cursor
      next_cursor = r.next_cursor
      images += r.images
      break if next_cursor.empty?
    end

    assert_equal 50, images.size
  end

  def test_ten_comments
    p = new_logged_in_colly.fetch_post DEFAULT_GROUP_ID, 660137804846290
    puts p
    assert p
    assert p.user.id
    assert p.id
    assert_equal p.content, 'test_ten_comments'
    assert_equal p.comments.comments.size, 10
  end

  def test_twenty_comments
    p = new_logged_in_colly.fetch_post DEFAULT_GROUP_ID, 660138831512854
    puts p
    assert p
    assert p.user.id
    assert p.id
    assert_equal p.content, 'test_twenty_comments'

    comments = []
    next_cursor = nil
    loop do
      r = new_logged_in_colly.fetch_post DEFAULT_GROUP_ID, 660138831512854, next_cursor
      next_cursor = r.comments.next_cursor
      comments += r.comments.comments
      break if next_cursor.empty?
    end

    assert_equal 20, comments.size
  end

  def test_thirty_comments
    p = new_logged_in_colly.fetch_post DEFAULT_GROUP_ID, 660143521512385
    puts p
    assert p
    assert p.user.id
    assert p.id
    assert_equal p.content, 'test_thirty_comments'

    comments = []
    next_cursor = nil
    loop do
      r = new_logged_in_colly.fetch_post DEFAULT_GROUP_ID, 660143521512385, next_cursor
      next_cursor = r.comments.next_cursor
      comments += r.comments.comments
      break if next_cursor.empty?
    end

    assert_equal 30, comments.size
  end


  def test_one_video_post

  end

  private

  # @return [FbcrawlColly::Client]
  def new_colly
    FbcrawlColly::Client.new HOST_AND_PORT
  end

  # @return [FbcrawlColly::Client]
  def new_logged_in_colly
    colly = new_colly
    colly.login_with_cookies(login_cookies)
    colly
  end

  def login_cookies
    colly = FbcrawlColly::Client.new HOST_AND_PORT
    @@login_cookies ||= colly.login EMAIL, PASSWORD, MOTP_SECRET
  end

end
