require 'ffi'
module FbcrawlColly::FFI
  extend FFI::Library

  ffi_lib File.expand_path("../../ext/fbcrawl_colly/fbcolly.so", File.dirname(__FILE__))
  attach_function :free, [ :pointer ], :void

  attach_function :Init, [], :pointer
  attach_function :FreeColly, [:pointer], :pointer
  attach_function :Login, [:pointer, :string, :string], :strptr
  attach_function :LoginWithCookies, [:pointer, :string], :void
  attach_function :FetchGroupInfo, [:pointer, :string], :strptr
  attach_function :FetchGroupFeed, [:pointer, :int64], :strptr
  attach_function :FetchPost, [:pointer, :int64, :int64], :strptr
  attach_function :FetchContentImages, [:pointer, :int64], :strptr
  attach_function :FetchImageUrl, [:pointer, :int64], :strptr
  # attach_function :FetchGroup, [:pointer, :string], :pointer
end
