package config

import (
	"flag"
	"fmt"
	"os"
)

var (
	RunLocally      bool
	SearchAddr      string
	RateAddr        string
	ProfileAddr     string
	UserAddr        string
	ReservationAddr string
	FrontendAddr    string
	Debug           bool
)

func init() {
	flag.BoolVar(&RunLocally, "local", false, "Run the services locally")
	//flag.BoolVar(&UseSeparatePorts, "useSeparatePorts", false, "Use separate ports for each downstream service")
	flag.BoolVar(&Debug, "debug", false, "Enable debug logging")
	flag.Parse()

	// if config.RunLocally, the frontend service address is localhost:50052, otherwise ask from the environment variable

	if RunLocally {
		FrontendAddr = "localhost:50052"
		SearchAddr = "localhost:50054"
		RateAddr = "localhost:50056"
		ProfileAddr = "localhost:50057"
		UserAddr = "localhost:50053"
		ReservationAddr = "localhost:50055"
	} else {
		// Get the frontend service addresses from the environment variables
		FrontendAddr = os.Getenv("SERVICE_A_URL")
		//if not set, ask from cmd std input
		if FrontendAddr == "" {
			fmt.Print("Enter Frontend service address: ")
			fmt.Scanln(&FrontendAddr)
		}

		SearchAddr = "search:50051"
		RateAddr = "rate:50051"
		ProfileAddr = "profile:50051"
		UserAddr = "user:50051"
		ReservationAddr = "reservation:50051"
	}
}
