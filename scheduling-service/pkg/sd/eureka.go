package sd

import (
	"errors"
	"fmt"
	"github.com/go-kit/kit/sd/eureka"
	"github.com/go-kit/log"
	"github.com/hudl/fargo"
	"net"
	"os"
	"strconv"
	"strings"
)

// BuildFargoInstance build a Fargo Instance
func BuildFargoInstance() eureka.Registrar {
	eurekaAddr := os.Getenv("EUREKA_SERVER_URL")
	if eurekaAddr == "" {
		fmt.Println("EUREKA_SERVER_URL is not set")
	}

	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestamp)

	var fargoConfig fargo.Config
	fargoConfig.Eureka.ServiceUrls = []string{eurekaAddr}
	fargoConfig.Eureka.PollIntervalSeconds = 1

	fargoConnection := fargo.NewConnFromConfig(fargoConfig)
	fInstance := buildFargoInstanceBody("scheduling-service", "UP")
	return *eureka.NewRegistrar(&fargoConnection, fInstance, log.With(logger, "component", "registrar"))
}

func buildFargoInstanceBody(appName, status string) *fargo.Instance {
	ipAddr, err := externalIP()
	if err != nil {
		fmt.Println(err)
	}
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		fmt.Println(err)
		port = 9000
	}
	return &fargo.Instance{
		InstanceId:        ipAddr + ":" + appName + ":" + os.Getenv("PORT"),
		HostName:          "localhost",
		App:               strings.ToUpper(appName),
		IPAddr:            ipAddr,
		VipAddress:        appName,
		SecureVipAddress:  appName,
		Status:            fargo.StatusType(status),
		Overriddenstatus:  "UNKNOWN",
		Port:              port,
		PortEnabled:       true,
		SecurePort:        8443,
		SecurePortEnabled: false,
		HomePageUrl:       "http://localhost:" + os.Getenv("PORT") + os.Getenv("BASE_PATH"),
		StatusPageUrl:     "http://localhost:" + os.Getenv("PORT") + "/status",
		HealthCheckUrl:    "http://localhost:" + os.Getenv("PORT") + "/health",
		CountryId:         0,
		DataCenterInfo: fargo.DataCenterInfo{
			Name: "MyOwn", Class: "com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo",
		},
		LeaseInfo: fargo.LeaseInfo{},
		Metadata:  fargo.InstanceMetadata{},
		UniqueID:  nil,
	}
}

// aux func to get external ip from
func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("not connected to any network")
}
