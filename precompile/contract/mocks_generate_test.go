package contract

//go:generate mockgen -package=$GOPACKAGE -copyright_file=../../license_header -destination=mocks.go . BlockContext,AccessibleState,StateDB
