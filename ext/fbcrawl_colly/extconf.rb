require 'mkmf'
MakeMakefile::find_executable 'go'
MakeMakefile::find_executable 'protoc'
MakeMakefile::find_executable 'protoc-gen-go'
$makefile_created = true
