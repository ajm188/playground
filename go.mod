module github.com/ajm188/playground

go 1.15

replace playground => ./

require (
	github.com/golang/protobuf v1.4.3
	github.com/gorilla/mux v1.8.0
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	google.golang.org/grpc v1.33.2
	playground v0.0.0-00010101000000-000000000000
)
