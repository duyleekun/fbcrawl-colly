all:
	cd ../../ && protoc fbcrawl.proto --go_out ./ --ruby_out ./lib/ && go build -o ./ext/fbcrawl_colly/fbcolly.so -buildmode=c-shared #-v -ldflags=-v .
clean:
	rm -rf fbcolly.h fbcolly.so ../../lib/fbcrawl_pb.rb
install:
	true
