# Generate protobuf classes for Go and Python.
echo "Compiling protobufs..."
protoc --python_out=. proto/*.proto
