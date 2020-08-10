require "test_helper"
require_relative '../lib/fbcrawl_colly/colly'

class FbcrawlCollyTest < Minitest::Test
  def setup
    super
    @colly = FbcrawlColly::Colly.new
  end

  def test_init_should_return_pointer
    puts @colly
    assert @colly != nil, "Colly should not be nil"
  end

  def test_login_ok
    # @colly.login "xxx@gmail.com", "xxx"
  end

  def test_group_feed
    @colly.fetch_group_feed "fbcolly"
    puts "DM"
  end
end
