require_relative 'lib/fbcrawl_colly/version'

Gem::Specification.new do |spec|
  spec.name          = "fbcrawl-colly"
  spec.version       = FbcrawlColly::VERSION
  spec.authors       = ["Duy Le"]
  spec.email         = ["duyleekun@gmail.com"]

  spec.summary       = %q{Crawl mbasic.facebook.com using GO Colly}
  spec.description   = %q{Crawl mbasic.facebook.com using GO Colly}
  spec.homepage      = "http://github.com/duyleekun/fbcrawl-colly"
  spec.license       = "MIT"
  spec.required_ruby_version = Gem::Requirement.new(">= 2.3.0")

  # spec.metadata["allowed_push_host"] = "TODO: Set to 'http://mygemserver.com'"

  spec.metadata["homepage_uri"] = spec.homepage
  spec.metadata["source_code_uri"] = "http://github.com/duyleekun/fbcrawl-colly"
  spec.metadata["changelog_uri"] = "http://github.com/duyleekun/fbcrawl-colly"

  # Specify which files should be added to the gem when it is released.
  # The `git ls-files -z` loads the files in the RubyGem that have been added into git.
  spec.files         = Dir.chdir(File.expand_path('..', __FILE__)) do
    `git ls-files -z`.split("\x0").reject { |f| f.match(%r{^(test|spec|features)/}) }
  end
  # spec.bindir        = "exe"
  # spec.executables   = spec.files.grep(%r{^exe/}) { |f| File.basename(f) }
  spec.require_paths = %w[lib lib/pb]

  spec.add_runtime_dependency 'google-protobuf'
  spec.add_runtime_dependency 'grpc'
end
