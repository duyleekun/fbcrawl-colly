require 'mkmf'
requirement_passed = true
requirement_passed &&= MakeMakefile::find_executable 'go'
requirement_passed &&= MakeMakefile::find_executable 'protoc'
requirement_passed &&= MakeMakefile::find_executable 'protoc-gen-go'
$makefile_created = requirement_passed
