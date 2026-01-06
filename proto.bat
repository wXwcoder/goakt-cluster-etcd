::protoc --go_out=internal/opspb --go_opt=paths=source_relative --connect-go_out=internal/opspb --connect-go_opt=paths=source_relative protos/ops/ops.proto

::protoc --go_out=. --go_opt=module=goakt-actors-cluster protos/ops/ops.proto
::protoc --connect-go_out=. --connect-go_opt=module=goakt-actors-cluster protos/ops/ops.proto

protoc --proto_path=protos/ops --go_out=internal/opspb --go_opt=paths=source_relative --connect-go_out=internal/opspb --connect-go_opt=paths=source_relative ops.proto
protoc --proto_path=protos/sample --go_out=internal/samplepb --go_opt=paths=source_relative --connect-go_out=internal/samplepb --connect-go_opt=paths=source_relative sample.proto service.proto