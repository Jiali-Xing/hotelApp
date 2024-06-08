package config

import (
	"flag"
)

var (
	RunLocally      bool
	SearchAddr      string
	RateAddr        string
	ProfileAddr     string
	UserAddr        string
	ReservationAddr string
	Debug           bool
)

func init() {
	flag.BoolVar(&RunLocally, "local", false, "Run the services locally")
	//flag.BoolVar(&UseSeparatePorts, "useSeparatePorts", false, "Use separate ports for each downstream service")
	flag.BoolVar(&Debug, "debug", false, "Enable debug logging")
	flag.Parse()

	if RunLocally {
		//UseSeparatePorts = true
		SearchAddr = "localhost:50054"
		RateAddr = "localhost:50056"
		ProfileAddr = "localhost:50057"
		UserAddr = "localhost:50053"
		ReservationAddr = "localhost:50055"
	} else {
		SearchAddr = "search:50051"
		RateAddr = "rate:50051"
		ProfileAddr = "profile:50051"
		UserAddr = "user:50051"
		ReservationAddr = "reservation:50051"
	}

	//if UseSeparatePorts {
	//	// Define different ports for each service
	//	UserPort = ":50053"
	//	SearchPort = ":50054"
	//	ReservationPort = ":50055"
	//	RatePort = ":50056"
	//	ProfilePort = ":50057"
	//} else {
	//	// Use the same port but different URLs (Kubernetes-like environment)
	//	UserPort = ":50051"
	//	SearchPort = ":50051"
	//	ReservationPort = ":50051"
	//	RatePort = ":50051"
	//	ProfilePort = ":50051"
	//}

}
