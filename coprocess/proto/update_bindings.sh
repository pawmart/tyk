#!/bin/sh
echo "Generating bindings for Go."
protoc -I. --go_out=plugins=grpc:../ *.proto
cd .. ; ls -l *.pb.go ; cd -

echo "Generating bindings for Python."
mkdir -p ../bindings/python
protoc -I. --python_out=../bindings/python *.proto
python codegen.py
cd ../bindings/python ; ls -l *.py ; cd -

echo "Generating bindings for Ruby."
mkdir -p ../bindings/ruby
protoc -I. --ruby_out=plugins=grpc:../bindings/ruby *.proto
cd ../bindings/ruby ; ls -l *.rb ; cd -

echo "Generating bindings for C++."
mkdir -p ../bindings/cpp
protoc -I. --cpp_out=../bindings/cpp *.proto
cd ../bindings/cpp ; ls -l ; cd -

echo
echo "Done"
