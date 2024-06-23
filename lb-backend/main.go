package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/rs/cors"
)

// to export data use PascleCase

// request data handling for /api/create_net
type subnet struct {
	SubnetCIDR string `json:"subnet"`
}

// request data handling for /api/launch_client
type clientIp struct {
	ClientIp string `json:"clientIp"`
}

// global_data
type global_data struct {
	Directory string
	CIDR      string
	IfaceName string
	ClientIp  string
	ClientMac string
}

var GlobalData global_data

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
		var data clientIp
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Printf("[info] received IP: %s\n", data.ClientIp)
		GlobalData.ClientIp = strings.TrimSpace(data.ClientIp)
		// adding the network setup
		cmd := exec.Command("make", "run_client", fmt.Sprintf("CLIENT_IP=%s", data.ClientIp))
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

		fmt.Printf("%s\n", GlobalData.ClientIp)
		fmt.Printf("%s\n", GlobalData.ClientMac)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{"message": "client launch successful with name client-server"}
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
	// Declare the routes
	mux := http.NewServeMux()
	mux.HandleFunc("/api/create_subnet", handleSubnetCreation)
	mux.HandleFunc("/api/launch_client", handleClientLaunch)
	// mux.HandleFunc("/api/delete_entry",deleteMapData)

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
