package config

import (
	"flag"
	"os"
)

var (
	RunLocally       bool
	SearchAddr       string
	RateAddr         string
	ProfileAddr      string
	UserAddr         string
	ReservationAddr  string
	FrontendAddr     string
	UserTimelineAddr string
	HomeTimelineAddr string
	PostStorageAddr  string
	SocialGraphAddr  string
	ComposePostAddr  string
	Debug            bool
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
		UserTimelineAddr = "localhost:50058"
		HomeTimelineAddr = "localhost:50059"
		PostStorageAddr = "localhost:50060"
		SocialGraphAddr = "localhost:50061"
		ComposePostAddr = "localhost:50062"
	} else {
		// Get the frontend service addresses from the environment variables
		FrontendAddr = os.Getenv("SERVICE_A_URL")
		//if not set, ask from cmd std input
		if FrontendAddr == "" {
			FrontendAddr = "frontend:50052"
		}

		SearchAddr = "search:50051"
		RateAddr = "rate:50051"
		ProfileAddr = "profile:50051"
		UserAddr = "user:50051"
		ReservationAddr = "reservation:50051"

		UserTimelineAddr = "usertimeline:50051"
		HomeTimelineAddr = "hometimeline:50051"
		PostStorageAddr = "poststorage:50051"
		SocialGraphAddr = "socialgraph:50051"
		ComposePostAddr = "composepost:50051"
	}
}
