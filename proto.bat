::protoc --go_out=internal/opspb --go_opt=paths=source_relative --connect-go_out=internal/opspb --connect-go_opt=paths=source_relative protos/ops/ops.proto

protoc --go_out=. --go_opt=module=goakt-actors-cluster protos/ops/ops.proto
protoc --connect-go_out=. --connect-go_opt=module=goakt-actors-cluster protos/ops/ops.proto
