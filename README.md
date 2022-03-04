# FbcrawlColly

This project is to crawl mbasic.facebook.com using GO Colly. 

Ruby gem is only client for GRPC colly service

## Installation

Add this line to your application's Gemfile:

```ruby
gem 'fbcrawl-colly'
```

And then execute:

    $ bundle install

Or install it yourself as:

    $ gem install fbcrawl-colly


## Contributing

### Update proto file

```shell
protoc fbcrawl.proto --go_out ./ --ruby_out ./lib/pb
```

### Bump gem ver

```shell

```

## License

The gem is available as open source under the terms of the [MIT License](https://opensource.org/licenses/MIT).

## Code of Conduct

Everyone interacting in the FbcrawlColly project's codebases, issue trackers, chat rooms and mailing lists is expected to follow the [code of conduct](https://github.com/[USERNAME]/fbcrawl-colly/blob/master/CODE_OF_CONDUCT.md).
