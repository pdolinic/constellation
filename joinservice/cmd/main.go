package main

import (
	"flag"
	"path/filepath"
	"strconv"

	"github.com/edgelesssys/constellation/internal/atls"
	"github.com/edgelesssys/constellation/internal/constants"
	"github.com/edgelesssys/constellation/internal/file"
	"github.com/edgelesssys/constellation/internal/grpc/atlscredentials"
	"github.com/edgelesssys/constellation/internal/logger"
	"github.com/edgelesssys/constellation/internal/watcher"
	"github.com/edgelesssys/constellation/joinservice/kms"
	"github.com/edgelesssys/constellation/joinservice/kubeadm"
	"github.com/edgelesssys/constellation/joinservice/kubernetesca"
	"github.com/edgelesssys/constellation/joinservice/server"
	"github.com/spf13/afero"
	"go.uber.org/zap"
)

func main() {
	provider := flag.String("cloud-provider", "", "cloud service provider this binary is running on")
	kmsEndpoint := flag.String("kms-endpoint", "", "endpoint of Constellations key management service")
	verbosity := flag.Int("v", 0, logger.CmdLineVerbosityDescription)

	flag.Parse()
	log := logger.New(logger.JSONLog, logger.VerbosityFromInt(*verbosity))

	log.With(zap.String("version", constants.VersionInfo), zap.String("cloudProvider", *provider)).
		Infof("Constellation Node Activation Service")

	handler := file.NewHandler(afero.NewOsFs())

	validator, err := watcher.NewValidator(log.Named("validator"), *provider, handler)
	if err != nil {
		flag.Usage()
		log.With(zap.Error(err)).Fatalf("Failed to create validator")
	}

	creds := atlscredentials.New(nil, []atls.Validator{validator})

	kubeadm, err := kubeadm.New(log.Named("kubeadm"))
	if err != nil {
		log.With(zap.Error(err)).Fatalf("Failed to create kubeadm")
	}
	kms := kms.New(log.Named("kms"), *kmsEndpoint)

	server := server.New(
		log.Named("server"),
		handler,
		kubernetesca.New(log.Named("certificateAuthority"), handler),
		kubeadm,
		kms,
	)

	watcher, err := watcher.New(log.Named("fileWatcher"), validator)
	if err != nil {
		log.With(zap.Error(err)).Fatalf("Failed to create watcher for measurements updates")
	}
	defer watcher.Close()

	go func() {
		log.Infof("starting file watcher for measurements file %s", filepath.Join(constants.ServiceBasePath, constants.MeasurementsFilename))
		if err := watcher.Watch(filepath.Join(constants.ServiceBasePath, constants.MeasurementsFilename)); err != nil {
			log.With(zap.Error(err)).Fatalf("Failed to watch measurements file")
		}
	}()

	if err := server.Run(creds, strconv.Itoa(constants.ActivationServicePort)); err != nil {
		log.With(zap.Error(err)).Fatalf("Failed to run server")
	}
}