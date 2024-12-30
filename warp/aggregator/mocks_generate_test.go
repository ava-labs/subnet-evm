package aggregator

//go:generate go run go.uber.org/mock/mockgen@v0.4 -package=$GOPACKAGE -source=signature_getter.go -destination=mock_signature_getter.go -exclude_interfaces=NetworkClient
