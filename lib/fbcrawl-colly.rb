require 'ffi'
require 'fbcrawl_pb'
module FbcrawlColly
  extend FFI::Library

  ffi_lib File.expand_path("../ext/fbcrawl_colly/fbcolly.so", File.dirname(__FILE__))
  attach_function :free, [ :pointer ], :void

  attach_function :Init, [], :pointer
  attach_function :Login, [:pointer, :string, :string], :void
  attach_function :FetchGroupFeed, [:pointer, :string], :string
  attach_function :FetchPost, [:pointer, :string, :string], :string
  # attach_function :FetchGroup, [:pointer, :string], :pointer
end
