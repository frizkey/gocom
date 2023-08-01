clear
go env -w GOOS=darwin
go env -w GOARCH=arm64
go build
cd test
../distobjgen -src Test
cd ..