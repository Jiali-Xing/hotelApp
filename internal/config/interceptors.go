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
	"github.com/Jiali-Xing/topdown-grpc"
	"github.com/tgiannoukos/charon"
	"google.golang.org/grpc"
)

var (
	serviceName  = getEnv("SERVICE_NAME", "Client")
	entryService bool
	Intercept    string
	serviceData  ServiceData
	PriceTable   *charon.PriceTable
	Breakwater   *bw.Breakwater
	Dg           *dagor.Dagor
	Topdown      *topdown.TopDownRL
	Breakwaterd  map[string]*bw.Breakwater
	// logLevel     string
	// serverConfig []Config
	yamlFile = getEnv("MSGRAPH_YAML", "msgraph.yaml")

	// nodes []Node
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
	entryService = false
	// if serviceName == "frontend", "nginx" then set entryService to true
	if serviceName == "frontend" || serviceName == "nginx" {
		entryService = true
		DebugLog("Service %s is the entry service", serviceName)
	}
	// nodes = GetNodes()

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
			Intercept = config.Value
			DebugLog("Reading config of interceptor %s", Intercept)
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

	switch Intercept {
	case "charon":
		charonOptions := map[string]interface{}{
			"initprice":          int64(0),
			"rateLimiting":       rateLimiting,
			"loadShedding":       loadShedding,
			"pinpointQueuing":    true,
			"pinpointThroughput": false,
			"pinpointLatency":    false,
			"debug":              Debug,
			"lazyResponse":       lazyUpdate,
			"priceUpdateRate":    priceUpdateRate,
			"guidePrice":         int64(-1),
			"priceStrategy":      priceStrategy,
			"latencyThreshold":   latencyThreshold,
			"priceStep":          priceStep,
			"priceAggregation":   "maximal",
			"recordPrice":        charonTrackPrice,
		}

		DebugLog("Initializing Charon with options: %v", charonOptions)
		PriceTable = charon.NewCharon(
			serviceName,
			serviceData.CallGraph,
			charonOptions,
		)

		DebugLog("Charon call graph: %v", serviceData.CallGraph)
		DebugLog("Charon Config: %v", charonOptions)

	case "breakwater":
		bwConfig = bw.BWParameters{
			Verbose:          Debug,
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

		DebugLog("Initializing Breakwater with config: %v", bwConfig)
		Breakwater = bw.InitBreakwater(bwConfig)
		// log.Printf("Breakwater Config: %v", bwConfig)

	case "breakwaterd":
		if entryService {
			// apply the same configuration as the breakwater.
			bwConfig = bw.BWParameters{
				Verbose:          Debug,
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

			DebugLog("Initializing Breakwater frontend of BreakwaterD with config: %v", bwConfig)
			Breakwater = bw.InitBreakwater(bwConfig)
			// log.Printf("Breakwater Config: %v", bwConfig)
		} else {
			bwConfig = bw.BWParameters{
				Verbose:          Debug,
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

			DebugLog("Initializing BreakwaterD Backend with config: %v", bwConfig)
			Breakwater = bw.InitBreakwater(bwConfig)
			// log.Printf("BreakwaterD Backend Config: %v", bwConfig)
		}
	case "dagor":
		dagorParams := dagor.DagorParam{
			NodeName: serviceName,
			BusinessMap: map[string]int{
				"search-hotel":  1,
				"store-hotel":   2,
				"reserve-hotel": 3,
				"compose":       1,
				"home-timeline": 2,
				"user-timeline": 3,
			},
			EntryService:                 entryService,
			IsEnduser:                    false,
			QueuingThresh:                dagorQueuingThresh,
			AdmissionLevelUpdateInterval: dagorAdmissionLevelUpdateInterval,
			Alpha:                        dagorAlpha,
			Beta:                         dagorBeta,
			Umax:                         dagorUmax,
			Bmax:                         4,
			Debug:                        Debug,
			UseSyncMap:                   false,
		}

		DebugLog("Initializing Dagor with config: %v", dagorParams)
		Dg = dagor.NewDagorNode(dagorParams)
		// log.Printf("Dagor Config: %v", dagorParams)
	case "topdown":
		sloMap := make(map[string]time.Duration)
		// manually set SLOs for each service
		sloMap["search-hotel"] = 60 * time.Millisecond
		sloMap["reserve-hotel"] = 60 * time.Millisecond

		Topdown = topdown.NewTopDownRL(10000, 1000, sloMap)
	case "plain":
		// No special initialization required for the plain interceptor
	default:
		// No interceptors or unknown interceptor type
	}

	if Intercept == "breakwaterd" {
		// Initialize the Breakwater instances for each downstream service
		InitializeBreakwaterd(bwConfig)
		DebugLog("[BreakwaterD] Initialized multiple instances of Breakwaters for downstream services of %s: %v", serviceName, Breakwaterd)
	}

}

// InitializeBreakwaterd initializes Breakwater instances for each downstream service
func InitializeBreakwaterd(bwConfig bw.BWParameters) {
	Breakwaterd = make(map[string]*bw.Breakwater)
	for _, downstream := range serviceData.Downstreams {
		// Customize the Breakwater config per downstream if needed
		// For example, you might have different SLOs or other parameters per downstream service
		downstreamConfig := bwConfig
		bwConfig.ServerSide = false
		addr := getURL(downstream)
		Breakwaterd[addr] = bw.InitBreakwater(downstreamConfig)
		DebugLog("[BreakwaterD] Initializing Breakwater for downstream service %s", downstream)
	}
}

func CreateGRPCConn(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure(), grpc.WithBlock())

	// Append interceptor if it exists in the map and serverSideInterceptOnly is false
	if !serverSideInterceptOnly {
		// Apply the selected interceptor
		DebugLog("[As a Client/Sender] Creating Interceptor %s for service %s", Intercept, serviceName)
		switch Intercept {
		case "charon":
			opts = append(opts, grpc.WithUnaryInterceptor(PriceTable.UnaryInterceptorClient))
		case "breakwater":
			// this case should not be reached
			DebugLog("Breakwater interceptor should not be used for client-side on service %s", serviceName)
			// 	opts = append(opts, grpc.WithUnaryInterceptor(Breakwater.UnaryInterceptorClient))
		case "breakwaterd":
			opts = append(opts, grpc.WithUnaryInterceptor(Breakwaterd[addr].UnaryInterceptorClient))
		case "dagor":
			opts = append(opts, grpc.WithUnaryInterceptor(Dg.UnaryInterceptorClient))
		// case "topdown":
		// opts = append(opts)
		case "plain":
			opts = append(opts, grpc.WithUnaryInterceptor(plain.UnaryInterceptorClient))
		}
	}

	conn, err := grpc.DialContext(ctx, addr, opts...)
	DebugLog("Creating gRPC connection to %s with intercept %s", addr, Intercept)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return conn, nil
}
