package config

import (
	"flag"
)

var (
	UseSeparatePorts bool
	SearchPort       string
	RatePort         string
	ProfilePort      string
	UserPort         string
	ReservationPort  string
)

func init() {
	flag.BoolVar(&UseSeparatePorts, "useSeparatePorts", false, "Use separate ports for each downstream service")
	flag.Parse()

	if UseSeparatePorts {
		// Define different ports for each service
		UserPort = ":50053"
		SearchPort = ":50054"
		ReservationPort = ":50055"
		RatePort = ":50056"
		ProfilePort = ":50057"
	} else {
		// Use the same port but different URLs (Kubernetes-like environment)
		UserPort = ":50051"
		SearchPort = ":50051"
		ReservationPort = ":50051"
		RatePort = ":50051"
		ProfilePort = ":50051"
	}

}
