require 'rake/file_list'
Gem::Specification.new do |s|
  s.name = %q{fbcrawl-colly}
  s.version = "0.0.1"
  s.author = "duyleekun"
  s.date = %q{2020-08-04}
  s.summary = %q{fbcrawl_colly}
  s.files = Rake::FileList['ext/**/*','fbcolly/**/*','fbcrawl.proto','go.mod','go.sum','lib/**/*','main.go'].map(&:to_s)
  s.extensions = [
      'ext/fbcrawl_colly/extconf.rb'
  ]
  s.require_paths = ["lib"]
  s.add_runtime_dependency 'ffi'
  s.add_runtime_dependency 'google-protobuf'
end
