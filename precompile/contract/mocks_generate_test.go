package contract

//go:generate go run go.uber.org/mock/mockgen@v0.4.0 -package=$GOPACKAGE -destination=mocks.go . BlockContext,AccessibleState,StateDB
