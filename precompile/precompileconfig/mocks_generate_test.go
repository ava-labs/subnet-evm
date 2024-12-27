package precompileconfig

//go:generate mockgen -package=$GOPACKAGE -copyright_file=../../license_header -destination=mocks.go . Predicater,Config,ChainConfig,Accepter
