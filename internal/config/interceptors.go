package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	bw "github.com/Jiali-Xing/breakwater-grpc/breakwater"
	dagor "github.com/Jiali-Xing/dagor-grpc/dagor"
	"github.com/Jiali-Xing/plain"
	"github.com/tgiannoukos/charon"
	"google.golang.org/grpc"
)

var (
	serviceName  = getEnv("SERVICE_NAME", "hotelApp")
	intercept    string
	serviceData  ServiceData
	priceTable   *charon.PriceTable
	breakwater   *bw.Breakwater
	dg           *dagor.Dagor
	breakwaterd  map[string]*bw.Breakwater
	logLevel     string
	serverConfig []Config
	yamlFile     = getEnv("MSGRAPH_YAML", "msgraph.yaml")

	nodes []Node

	priceUpdateRate  time.Duration
	latencyThreshold time.Duration
	priceStep        int64
	priceStrategy    string
	lazyUpdate       bool
	rateLimiting     bool
	loadShedding     bool
	charonTrackPrice bool

	breakwaterSLO           time.Duration
	breakwaterClientTimeout time.Duration
	breakwaterInitialCredit int64
	breakwaterA             float64
	breakwaterB             float64
	breakwaterLoadShedding  bool
	breakwaterRTT           time.Duration
	breakwaterTrackCredit   bool

	// add a separate set of parameters for breakwaterd
	breakwaterdSLO           time.Duration
	breakwaterdClientTimeout time.Duration
	breakwaterdInitialCredit int64
	breakwaterdA             float64
	breakwaterdB             float64
	breakwaterdLoadShedding  bool
	breakwaterdRTT           time.Duration

	dagorQueuingThresh                time.Duration
	dagorAlpha                        float64
	dagorBeta                         float64
	dagorUmax                         int
	dagorAdmissionLevelUpdateInterval time.Duration

	// serverSideInterceptOnly is false by default
	serverSideInterceptOnly = false
)

type ServiceData struct {
	CallGraph    map[string][]string
	Downstreams  []string
	ServerConfig []Config
	CharonConfig []Config
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
func init() {

	nodeList := getNodeList()
	// raise error if the serviceName is not in the nodeList (slice of strings)
	if !contains(nodeList, serviceName) {
		log.Fatalf("Service %s is not in the nodeList", serviceName)
	}

	nodes = GetNodes()

	callGraph := GetCallGraph()
	fmt.Println("Call Graph of the Service:")
	fmt.Println(callGraph[serviceName])

	downstreams := GetDownstreamNames()
	fmt.Println("Downstreams of the Service:")
	fmt.Println(downstreams[serviceName])

	serverConfigs := GetServerConfigs()
	fmt.Println("Server Configurations:")
	fmt.Println(serverConfigs[serviceName])

	charonConfigs := GetCharonConfigs()
	fmt.Println("Charon Configurations:")
	fmt.Println(charonConfigs[serviceName])

	// Initialize the global serviceData struct
	serviceData = ServiceData{
		CallGraph:   callGraph[serviceName],
		Downstreams: downstreams[serviceName],
		// DownstreamURLs: downstreamURLs[serviceName],
		ServerConfig: serverConfigs[serviceName],
		CharonConfig: charonConfigs[serviceName],
	}

	for _, config := range serviceData.CharonConfig {
		switch config.Name {
		case "INTERCEPT":
			intercept = config.Value
		// charon parameters
		case "PRICE_UPDATE_RATE":
			priceUpdateRate, _ = time.ParseDuration(config.Value)
		case "LATENCY_THRESHOLD":
			latencyThreshold, _ = time.ParseDuration(config.Value)
		case "PRICE_STEP":
			priceStep, _ = strconv.ParseInt(config.Value, 10, 64)
		case "PRICE_STRATEGY":
			priceStrategy = config.Value
		case "LAZY_UPDATE":
			lazyUpdate, _ = strconv.ParseBool(config.Value)
		case "RATE_LIMITING":
			rateLimiting, _ = strconv.ParseBool(config.Value)
		case "LOAD_SHEDDING":
			loadShedding, _ = strconv.ParseBool(config.Value)
		case "CHARON_TRACK_PRICE":
			charonTrackPrice, _ = strconv.ParseBool(config.Value)
		// breakwater parameters
		case "BREAKWATER_SLO":
			breakwaterSLO, _ = time.ParseDuration(config.Value)
		case "BREAKWATER_CLIENT_EXPIRATION":
			breakwaterClientTimeout, _ = time.ParseDuration(config.Value)
		case "BREAKWATER_INITIAL_CREDIT":
			breakwaterInitialCredit, _ = strconv.ParseInt(config.Value, 10, 64)
		case "BREAKWATER_A":
			breakwaterA, _ = strconv.ParseFloat(config.Value, 64)
		case "BREAKWATER_B":
			breakwaterB, _ = strconv.ParseFloat(config.Value, 64)
		case "BREAKWATER_LOAD_SHEDDING":
			breakwaterLoadShedding, _ = strconv.ParseBool(config.Value)
		case "BREAKWATER_RTT":
			breakwaterRTT, _ = time.ParseDuration(config.Value)
		case "BREAKWATER_TRACK_CREDIT":
			breakwaterTrackCredit, _ = strconv.ParseBool(config.Value)
		// and one optional field: 'SIDE': 'server_only'
		case "SIDE":
			// if the side is server_only, then set the serverSideInterceptOnly to true
			if config.Value == "server_only" {
				serverSideInterceptOnly = true
			}
		case "BREAKWATERD_SLO":
			breakwaterdSLO, _ = time.ParseDuration(config.Value)
		case "BREAKWATERD_CLIENT_EXPIRATION":
			breakwaterdClientTimeout, _ = time.ParseDuration(config.Value)
		case "BREAKWATERD_INITIAL_CREDIT":
			breakwaterdInitialCredit, _ = strconv.ParseInt(config.Value, 10, 64)
		case "BREAKWATERD_A":
			breakwaterdA, _ = strconv.ParseFloat(config.Value, 64)
		case "BREAKWATERD_B":
			breakwaterdB, _ = strconv.ParseFloat(config.Value, 64)
		case "BREAKWATERD_LOAD_SHEDDING":
			breakwaterdLoadShedding, _ = strconv.ParseBool(config.Value)
		case "BREAKWATERD_RTT":
			breakwaterdRTT, _ = time.ParseDuration(config.Value)
		case "DAGOR_QUEUING_THRESHOLD":
			dagorQueuingThresh, _ = time.ParseDuration(config.Value)
		case "DAGOR_ALPHA":
			dagorAlpha, _ = strconv.ParseFloat(config.Value, 64)
		case "DAGOR_BETA":
			dagorBeta, _ = strconv.ParseFloat(config.Value, 64)
		case "DAGOR_ADMISSION_LEVEL_UPDATE_INTERVAL":
			dagorAdmissionLevelUpdateInterval, _ = time.ParseDuration(config.Value)
		case "DAGOR_UMAX":
			dagorUmax, _ = strconv.Atoi(config.Value)
		}
	}
	bwConfig := bw.BWParametersDefault

	switch intercept {
	case "charon":
		charonOptions := map[string]interface{}{
			"initprice":          int64(0),
			"rateLimiting":       rateLimiting,
			"loadShedding":       loadShedding,
			"pinpointQueuing":    true,
			"pinpointThroughput": false,
			"pinpointLatency":    false,
			"debug":              logLevel == "debug",
			"lazyResponse":       lazyUpdate,
			"priceUpdateRate":    priceUpdateRate,
			"guidePrice":         int64(-1),
			"priceStrategy":      priceStrategy,
			"latencyThreshold":   latencyThreshold,
			"priceStep":          priceStep,
			"priceAggregation":   "maximal",
			"recordPrice":        charonTrackPrice,
		}

		priceTable = charon.NewCharon(
			serviceName,
			nil, // Assuming no call graph is provided
			charonOptions,
		)
		log.Printf("Charon Config: %v", charonOptions)

	case "breakwater":
		bwConfig = bw.BWParameters{
			Verbose:          logLevel == "debug",
			SLO:              breakwaterSLO.Microseconds(),
			ClientExpiration: breakwaterClientTimeout.Microseconds(),
			InitialCredits:   breakwaterInitialCredit,
			LoadShedding:     breakwaterLoadShedding,
			ServerSide:       true,
			AFactor:          breakwaterA,
			BFactor:          breakwaterB,
			RTT_MICROSECOND:  breakwaterRTT.Microseconds(),
			TrackCredits:     breakwaterTrackCredit,
		}

		breakwater = bw.InitBreakwater(bwConfig)
		log.Printf("Breakwater Config: %v", bwConfig)

	case "breakwaterd":
		if serviceName == "frontend" {
			// apply the same configuration as the breakwater.
			bwConfig = bw.BWParameters{
				Verbose:          logLevel == "debug",
				SLO:              breakwaterSLO.Microseconds(),
				ClientExpiration: breakwaterClientTimeout.Microseconds(),
				InitialCredits:   breakwaterInitialCredit,
				LoadShedding:     breakwaterLoadShedding,
				ServerSide:       true,
				AFactor:          breakwaterA,
				BFactor:          breakwaterB,
				RTT_MICROSECOND:  breakwaterRTT.Microseconds(),
				TrackCredits:     breakwaterTrackCredit,
			}

			breakwater = bw.InitBreakwater(bwConfig)
			log.Printf("Breakwater Config: %v", bwConfig)
		} else {
			bwConfig := bw.BWParameters{
				Verbose:          logLevel == "debug",
				SLO:              breakwaterdSLO.Microseconds(),
				ClientExpiration: breakwaterdClientTimeout.Microseconds(),
				InitialCredits:   breakwaterdInitialCredit,
				LoadShedding:     true,
				ServerSide:       true,
				AFactor:          breakwaterdA,
				BFactor:          breakwaterdB,
				RTT_MICROSECOND:  breakwaterdRTT.Microseconds(),
				TrackCredits:     breakwaterTrackCredit,
			}
			breakwater = bw.InitBreakwater(bwConfig)
			log.Printf("BreakwaterD Backend Config: %v", bwConfig)
		}
	case "dagor":
		dagorParams := dagor.DagorParam{
			NodeName: serviceName,
			BusinessMap: map[string]int{
				"SearchHotels":        1,
				"StoreHotel":          2,
				"FrontendReservation": 3,
			},
			EntryService:                 true,
			IsEnduser:                    false,
			QueuingThresh:                dagorQueuingThresh,
			AdmissionLevelUpdateInterval: dagorAdmissionLevelUpdateInterval,
			Alpha:                        dagorAlpha,
			Beta:                         dagorBeta,
			Umax:                         dagorUmax,
			Bmax:                         4,
			Debug:                        logLevel == "debug",
			UseSyncMap:                   false,
		}

		dg = dagor.NewDagorNode(dagorParams)
		log.Printf("Dagor Config: %v", dagorParams)

	case "plain":
		// No special initialization required for the plain interceptor
	default:
		// No interceptors or unknown interceptor type
	}

	if intercept == "breakwaterd" {
		// Initialize the Breakwater instances for each downstream service
		InitializeBreakwaterd(bwConfig)
		log.Printf("Initialized multiple instances of Breakwaters for downstream services of %s: %v", serviceName, breakwaterd)
	}

}

// InitializeBreakwaterd initializes Breakwater instances for each downstream service
func InitializeBreakwaterd(bwConfig bw.BWParameters) {
	breakwaterd = make(map[string]*bw.Breakwater)
	for _, downstream := range serviceData.Downstreams {
		// Customize the Breakwater config per downstream if needed
		// For example, you might have different SLOs or other parameters per downstream service
		downstreamConfig := bwConfig
		bwConfig.ServerSide = false
		addr := getURL(downstream)
		breakwaterd[addr] = bw.InitBreakwater(downstreamConfig)
	}
}

func createGRPCConn(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure(), grpc.WithBlock())

	// Append interceptor if it exists in the map and serverSideInterceptOnly is false
	if !serverSideInterceptOnly {
		// Apply the selected interceptor
		switch intercept {
		case "charon":
			opts = append(opts, grpc.WithUnaryInterceptor(priceTable.UnaryInterceptorClient))
		case "breakwater":
			opts = append(opts, grpc.WithUnaryInterceptor(breakwater.UnaryInterceptorClient))
		case "breakwaterd":
			opts = append(opts, grpc.WithUnaryInterceptor(breakwaterd[addr].UnaryInterceptorClient))
		case "dagor":
			opts = append(opts, grpc.WithUnaryInterceptor(dg.UnaryInterceptorClient))
		case "plain":
			opts = append(opts, grpc.WithUnaryInterceptor(plain.UnaryInterceptorClient))
		}
	}

	conn, err := grpc.DialContext(ctx, addr, opts...)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return conn, nil
}
