package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/edgelesssys/constellation/coordinator/cloudprovider/cloudtypes"
	"github.com/edgelesssys/constellation/coordinator/role"
	"github.com/edgelesssys/constellation/hack/qemu-metadata-api/virtwrapper"
	"github.com/edgelesssys/constellation/internal/file"
	"github.com/edgelesssys/constellation/internal/logger"
	"go.uber.org/zap"
)

const exportedPCRsDir = "/pcrs/"

type Server struct {
	log  *logger.Logger
	virt virConnect
	file file.Handler
}

func New(log *logger.Logger, conn virConnect, file file.Handler) *Server {
	return &Server{
		log:  log,
		virt: conn,
		file: file,
	}
}

func (s *Server) ListenAndServe(port string) error {
	mux := http.NewServeMux()
	mux.Handle("/self", http.HandlerFunc(s.listSelf))
	mux.Handle("/peers", http.HandlerFunc(s.listPeers))
	mux.Handle("/log", http.HandlerFunc(s.postLog))
	mux.Handle("/pcrs", http.HandlerFunc(s.exportPCRs))

	server := http.Server{
		Handler: mux,
	}

	lis, err := net.Listen("tcp", net.JoinHostPort("", port))
	if err != nil {
		return err
	}

	s.log.Infof("Starting QEMU metadata API on %s", lis.Addr())
	return server.Serve(lis)
}

// listSelf returns peer information about the instance issuing the request.
func (s *Server) listSelf(w http.ResponseWriter, r *http.Request) {
	log := s.log.With(zap.String("peer", r.RemoteAddr))
	log.Infof("Serving GET request for /self")

	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.With(zap.Error(err)).Errorf("Failed to parse remote address")
		http.Error(w, fmt.Sprintf("Failed to parse remote address: %s\n", err), http.StatusInternalServerError)
		return
	}

	peers, err := s.listAll()
	if err != nil {
		log.With(zap.Error(err)).Errorf("Failed to list peer metadata")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, peer := range peers {
		for _, ip := range peer.PublicIPs {
			if ip == remoteIP {
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(peer); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				log.Infof("Request successful")
				return
			}
		}
	}

	log.Errorf("Failed to find peer in active leases")
	http.Error(w, "No matching peer found", http.StatusNotFound)
}

// listPeers returns a list of all active peers.
func (s *Server) listPeers(w http.ResponseWriter, r *http.Request) {
	log := s.log.With(zap.String("peer", r.RemoteAddr))
	log.Infof("Serving GET request for /peers")

	peers, err := s.listAll()
	if err != nil {
		log.With(zap.Error(err)).Errorf("Failed to list peer metadata")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(peers); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Infof("Request successful")
}

// postLog writes implements cloud-logging for QEMU instances.
func (s *Server) postLog(w http.ResponseWriter, r *http.Request) {
	log := s.log.With(zap.String("peer", r.RemoteAddr))
	if r.Method != http.MethodPost {
		log.With(zap.String("method", r.Method)).Errorf("Invalid method for /log")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Infof("Serving POST request for /log")

	if r.Body == nil {
		log.Errorf("Request body is empty")
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		return
	}

	msg, err := io.ReadAll(r.Body)
	if err != nil {
		log.With(zap.Error(err)).Errorf("Failed to read request body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.With(zap.String("message", string(msg))).Infof("Cloud-logging entry")
}

// exportPCRs allows QEMU instances to export their TPM state during boot.
// This can be used to check expected PCRs for GCP/Azure cloud images locally.
func (s *Server) exportPCRs(w http.ResponseWriter, r *http.Request) {
	log := s.log.With(zap.String("peer", r.RemoteAddr))
	if r.Method != http.MethodPost {
		log.With(zap.String("method", r.Method)).Errorf("Invalid method for /log")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Infof("Serving POST request for /pcrs")

	if r.Body == nil {
		log.Errorf("Request body is empty")
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		return
	}

	// unmarshal the request body into a map of PCRs
	var pcrs map[uint32][]byte
	if err := json.NewDecoder(r.Body).Decode(&pcrs); err != nil {
		log.With(zap.Error(err)).Errorf("Failed to read request body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get name of the node sending the export request
	var nodeName string
	peers, err := s.listAll()
	if err != nil {
		log.With(zap.Error(err)).Errorf("Failed to list peer metadata")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.With(zap.Error(err)).Errorf("Failed to parse remote address")
		http.Error(w, fmt.Sprintf("Failed to parse remote address: %s\n", err), http.StatusInternalServerError)
		return
	}
	for _, peer := range peers {
		if peer.PublicIPs[0] == remoteIP {
			nodeName = peer.Name
		}
	}

	// write PCRs as JSON and YAML to disk
	if err := s.file.WriteJSON(exportedPCRsDir+nodeName+"_pcrs.json", pcrs, file.OptOverwrite); err != nil {
		log.With(zap.Error(err)).Errorf("Failed to write pcrs to JSON")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// convert []byte to base64 encoded strings for YAML encoding
	pcrsYAML := make(map[uint32]string)
	for k, v := range pcrs {
		pcrsYAML[k] = base64.StdEncoding.EncodeToString(v)
	}
	if err := s.file.WriteYAML(exportedPCRsDir+nodeName+"_pcrs.yaml", pcrsYAML, file.OptOverwrite); err != nil {
		log.With(zap.Error(err)).Errorf("Failed to write pcrs to YAML")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// listAll returns a list of all active peers.
func (s *Server) listAll() ([]cloudtypes.Instance, error) {
	net, err := s.virt.LookupNetworkByName("constellation")
	if err != nil {
		return nil, err
	}
	defer net.Free()

	leases, err := net.GetDHCPLeases()
	if err != nil {
		return nil, err
	}
	var peers []cloudtypes.Instance

	for _, lease := range leases {
		instanceRole := role.Node
		if strings.HasPrefix(lease.Hostname, "control-plane") {
			instanceRole = role.Coordinator
		}

		peers = append(peers, cloudtypes.Instance{
			Name:       lease.Hostname,
			Role:       instanceRole,
			PrivateIPs: []string{lease.IPaddr},
			PublicIPs:  []string{lease.IPaddr},
			ProviderID: "qemu:///hostname/" + lease.Hostname,
		})
	}

	return peers, nil
}

type virConnect interface {
	LookupNetworkByName(name string) (*virtwrapper.Network, error)
}