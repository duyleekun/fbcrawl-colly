require "bundler/gem_tasks"
require "rake/testtask"

Rake::TestTask.new(:test) do |t|
  t.libs << "test"
  t.libs << "lib"
  t.test_files = FileList["test/**/*_test.rb"]
end

task :fbcrawl_colly do
  Dir.chdir("./ext/fbcrawl_colly/") do
    require './extconf'
    `make`
  end
end

task :compile => [:fbcrawl_colly]
task :test => :compile
task :default => :test
