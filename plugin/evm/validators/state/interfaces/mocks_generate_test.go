package interfaces

//go:generate go run go.uber.org/mock/mockgen@v0.4.0 -package=$GOPACKAGE -copyright_file=../../../../../license_header -destination=mock_listener.go . StateCallbackListener
