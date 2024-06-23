package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"unsafe"

	"github.com/cilium/ebpf"
	"github.com/rs/cors"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go xdp_lb ../loadbalancer/xdp/xdp_lb.c

// to export data use PascleCase

// request data handling for /api/create_net
type subnet struct {
	SubnetCIDR string `json:"subnet"`
}

// request data handling for /api/launch_client
type ip struct {
	IP string `json:"ip"`
}

// mac len
const (
	ETH_ALEN = 6
)

// mac address
type MacAddr struct {
	Addr [ETH_ALEN]byte
}

type IpMac struct {
	Ip  uint32
	Mac MacAddr
}

// global_data
type global_data struct {
	Directory   string
	CIDR        string
	IfaceName   string
	ClientIp    string
	ClientMac   string
	LbIp        string
	LbMac       string
	Xdp_ProgObj *xdp_lbObjects
}

var GlobalData global_data

// parseMAC parses a MAC address string and returns an EthAddr.
func parseMAC(macStr string) MacAddr {
	// Split the MAC address string by ":"
	parts := strings.Split(macStr, ":")

	// Initialize a byte slice of length 6
	var macAddress MacAddr

	// Convert each part from hexadecimal string to byte
	for i := 0; i < 6; i++ {
		b, err := hex.DecodeString(parts[i])
		if err != nil {
			fmt.Printf("[error] error decoding MAC address: %v\n", err)

		}
		macAddress.Addr[i] = b[0]
	}

	return macAddress
}

// parseIP parses an IP address string and returns it as a uint32.
func parseIP(ipStr string) uint32 {
	// Split IP address into octets
	parts := strings.Split(ipStr, ".")

	// Convert octets to integers
	var ip uint32
	for i := 0; i < 4; i++ {
		octet, _ := strconv.Atoi(parts[i])
		if octet < 0 || octet > 255 {
			fmt.Printf("Invalid octet value: %s\n", parts[i])
		}
		ip |= uint32(octet) << (uint(i) * 8)
	}
	return ip
}

// create docker network and get the bridge interface within that subnet
func handleSubnetCreation(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var data subnet
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Printf("[info] received CIDR: %s\n", data.SubnetCIDR)
		GlobalData.CIDR = data.SubnetCIDR
		// adding the network setup
		cmd := exec.Command("make", "create_net", fmt.Sprintf("CIDR=%s", data.SubnetCIDR))
		cmd.Dir = GlobalData.Directory
		Output, err := cmd.Output()
		if err != nil {
			fmt.Printf("[error] failed to create network -> [\n %s ]\n", err)
			return
		}
		fmt.Printf("[success] network creation successfull -> [\n %s ]\n", Output)
		// get the iface of the setup network
		cmd = exec.Command("make", "get_iface")
		cmd.Dir = GlobalData.Directory
		Output, err = cmd.Output()
		if err != nil {
			fmt.Printf("[error] failed to get network interface -> [\n %s ]\n", err)
			return
		}
		fmt.Printf("[info] selected network interface -> [\n %s ]\n", Output)
		GlobalData.IfaceName = strings.TrimSpace(string(Output))
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{"message": "network creation successful with name lb-net"}
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// launch the client and get the client mac address
func handleClientLaunch(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var data ip
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Printf("[info] received IP: %s\n", data.IP)
		GlobalData.ClientIp = strings.TrimSpace(data.IP)
		// adding the network setup
		cmd := exec.Command("make", "run_client", fmt.Sprintf("CLIENT_IP=%s", data.IP))
		cmd.Dir = GlobalData.Directory
		Output, err := cmd.Output()
		if err != nil {
			fmt.Printf("[error] failed to launch client -> [\n %s ]\n", err)
			return
		}
		fmt.Printf("[success] client launched successfully -> [\n %s ]\n", Output)

		// getting client mac address
		cmd = exec.Command("make", "get_client_mac")
		cmd.Dir = GlobalData.Directory
		Output, err = cmd.Output()
		if err != nil {
			fmt.Printf("[error] failed to get mac address of client-server -> [\n %s ]\n", err)
			return
		}
		fmt.Printf("[info] client-server mac address -> [\n %s ]\n", Output)
		GlobalData.ClientMac = strings.TrimSpace(string(Output))

		parsedIp := parseIP(GlobalData.ClientIp)
		parsedMac := parseMAC(GlobalData.ClientMac)
		ipMac := IpMac{Ip: parsedIp, Mac: parsedMac}
		fmt.Printf("parsedIp : %d\n", parsedIp)
		fmt.Printf("parsedMac: %v\n", parsedMac)

		clientMacMap := GlobalData.Xdp_ProgObj.xdp_lbMaps.ClientMacMap
		key := uint32(0)
		if err := clientMacMap.Update(unsafe.Pointer(&key), unsafe.Pointer(&ipMac), ebpf.UpdateAny); err != nil {
			fmt.Printf("[error] client_mac_map update failed -> [\n %s ]\n", err)
			return
		}
		fmt.Printf("[success] client_mac_map update successfull\n")
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{"message": "client launch successful with name client-server"}
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// launch the load-balancer server and the load-balancer server mac address
func handleLbLaunch(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var data ip
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Printf("[info] received IP: %s\n", data.IP)
		GlobalData.LbIp = strings.TrimSpace(data.IP)
		// adding the network setup
		cmd := exec.Command("make", "run_lb", fmt.Sprintf("LB_IP=%s", data.IP))
		cmd.Dir = GlobalData.Directory
		Output, err := cmd.Output()
		if err != nil {
			fmt.Printf("[error] failed to launch load-balancer -> [\n %s ]\n", err)
			return
		}
		fmt.Printf("[success] load-balancer launched successfully -> [\n %s ]\n", Output)

		// getting client mac address
		cmd = exec.Command("make", "get_lb_mac")
		cmd.Dir = GlobalData.Directory
		Output, err = cmd.Output()
		if err != nil {
			fmt.Printf("[error] failed to get mac address of lb-server -> [\n %s ]\n", err)
			return
		}
		fmt.Printf("[info] lb-server mac address -> [\n %s ]\n", Output)
		GlobalData.LbMac = strings.TrimSpace(string(Output))

		parsedIp := parseIP(GlobalData.LbIp)
		parsedMac := parseMAC(GlobalData.LbMac)
		ipMac := IpMac{Ip: parsedIp, Mac: parsedMac}
		fmt.Printf("parsedIp : %d\n", parsedIp)
		fmt.Printf("parsedMac: %v\n", parsedMac)

		lbMacMap := GlobalData.Xdp_ProgObj.xdp_lbMaps.LbMacMap
		key := uint32(0)
		if err := lbMacMap.Update(unsafe.Pointer(&key), unsafe.Pointer(&ipMac), ebpf.UpdateAny); err != nil {
			fmt.Printf("[error] lb_mac_map update failed -> [\n %s ]\n", err)
			return
		}
		fmt.Printf("[success] lb_mac_map update successfull\n")
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{"message": "loadbalancer launch successful with name lb-server"}
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// launch the server and get the server mac address
func handleServerLaunch(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var data ip
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Printf("[info] received IP: %s\n", data.IP)
		// adding the network setup
		cmd := exec.Command("make", "run_server", fmt.Sprintf("SERVER_IP=%s", data.IP))
		cmd.Dir = GlobalData.Directory
		Output, err := cmd.Output()
		if err != nil {
			fmt.Printf("[error] failed to launch server -> [\n %s ]\n", err)
			return
		}
		fmt.Printf("[success] server launched successfully -> [\n %s ]\n", Output)

		// getting client mac address
		cmd = exec.Command("make", "get_server_mac", fmt.Sprintf("SERVER_IP=%s", data.IP))
		cmd.Dir = GlobalData.Directory
		Output, err = cmd.Output()
		if err != nil {
			fmt.Printf("[error] failed to get mac address of server-backend -> [\n %s ]\n", err)
			return
		}
		fmt.Printf("[info] server-backend mac address -> [\n %s ]\n", Output)

		parsedIp := parseIP(strings.TrimSpace(data.IP))
		parsedMac := parseMAC(strings.TrimSpace(string(Output)))
		ipMac := IpMac{Ip: parsedIp, Mac: parsedMac}
		fmt.Printf("parsedIp : %d\n", parsedIp)
		fmt.Printf("parsedMac: %v\n", parsedMac)

		// get the current size
		size_map := GlobalData.Xdp_ProgObj.xdp_lbMaps.BsMapSizeMap
		var size uint32
		key := uint32(0)
		if err := size_map.Lookup(unsafe.Pointer(&key), unsafe.Pointer(&size)); err != nil {
			fmt.Printf("[error] failed to read the bs_map_size_map -> [\n %s ]\n", err)
			return
		}
		backendServerMacMap := GlobalData.Xdp_ProgObj.xdp_lbMaps.BackendServerMap
		if err := backendServerMacMap.Update(unsafe.Pointer(&size), unsafe.Pointer(&ipMac), ebpf.UpdateAny); err != nil {
			fmt.Printf("[error] backend_server_map update failed -> [\n %s ]\n", err)
			return
		}
		// update the size
		size += uint32(1)

		if err := size_map.Update(unsafe.Pointer(&key), unsafe.Pointer(&size), ebpf.UpdateAny); err != nil {
			fmt.Printf("[error] bbs_map_size_map update failed -> [\n %s ]\n", err)
			return
		}
		fmt.Printf("[success] backend_server_map update successfull\n")
		fmt.Printf("[info] total server-backend launched %d\n", size)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{"message": "server launch successful with name server-backend"}
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func handleAttachXdp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// adding the network setup
		cmd := exec.Command("make", "build_obj")
		cmd.Dir = GlobalData.Directory
		_, err := cmd.Output()
		if err != nil {
			fmt.Printf("[error] failed to build xdp_lb.c -> [\n %s ]\n", err)
			return
		}
		fmt.Printf("[success] successfully build xdp_lb.c and generate object file xdp_lb.o\n")
		cmd = exec.Command("make", "attach_xdp_lb", fmt.Sprintf("IFACE=%s", GlobalData.IfaceName))
		cmd.Dir = GlobalData.Directory
		_, err = cmd.Output()
		if err != nil {
			fmt.Printf("[error] failed to attach xdp_lb.o -> [\n %s ]\n", err)
			return
		}
		fmt.Printf("[success] successfully attached xdp_lb to interface %s\n", GlobalData.IfaceName)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{"message": "xdp_lb attached successfully"}
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
func main() {
	// -port <port_no>  8080 is set to default
	port := flag.Int("port", 8080, "port to listen on")

	// Parse the command-line flags
	flag.Parse()

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("[error] getting the current_working_directory -> ", err)
		return
	}
	// Split path by "/"
	pathComponents := strings.Split(cwd, "/")
	// Change the last component to "ip_info"
	if len(pathComponents) > 0 {
		pathComponents[len(pathComponents)-1] = ""
	}
	// Join components back into a path
	appRoot := strings.Join(pathComponents, "/")

	fmt.Printf("[info]  app_root_path %s\n", appRoot)
	GlobalData.Directory = appRoot

	// loading the ebpf and set the pin path
	var objs xdp_lbObjects
	if err := loadXdp_lbObjects(&objs, &ebpf.CollectionOptions{Maps: ebpf.MapOptions{PinPath: "/sys/fs/bpf/tc/globals"}}); err != nil {
		fmt.Printf("[error] failed loading eBPF objects: %v", err)
		return
	}
	defer objs.Close()

	// storing it to use in another functions
	GlobalData.Xdp_ProgObj = &objs

	bsMapSizeMap := GlobalData.Xdp_ProgObj.xdp_lbMaps.BsMapSizeMap
	key := uint32(0)
	size := uint32(0)
	if err := bsMapSizeMap.Update(unsafe.Pointer(&key), unsafe.Pointer(&size), ebpf.UpdateAny); err != nil {
		fmt.Printf("[error] failed to initialize size for the backend_server_map -> [\n %s ]\n", err)
		return
	}
	fmt.Printf("[info] initialize size for backend_server_map\n")
	// Declare the routes
	mux := http.NewServeMux()
	mux.HandleFunc("/api/create_subnet", handleSubnetCreation)
	mux.HandleFunc("/api/launch_client", handleClientLaunch)
	mux.HandleFunc("/api/launch_lb", handleLbLaunch)
	mux.HandleFunc("/api/launch_server", handleServerLaunch)
	mux.HandleFunc("/api/attach_xdp", handleAttachXdp)
	// Setup the CORS origin
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(mux)

	// Use the provided port
	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("[info] server starting on port %d\n", *port)
	if err := http.ListenAndServe(addr, handler); err != nil {
		fmt.Fprintf(os.Stderr, "[error] failed to start server: %v", err)
		os.Exit(1)
	}
}
