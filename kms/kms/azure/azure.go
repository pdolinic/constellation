package azure

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"github.com/edgelesssys/constellation/kms/config"
	"github.com/edgelesssys/constellation/kms/kms"
	"github.com/edgelesssys/constellation/kms/kms/util"
	"github.com/edgelesssys/constellation/kms/storage"
)

const (
	vaultPrefix = "https://"
	// DefaultCloud is the suffix for the default Vault URL.
	DefaultCloud VaultSuffix = ".vault.azure.net/"
	// ChinaCloud is the suffix for Vaults in Azure China Cloud.
	ChinaCloud VaultSuffix = ".vault.azure.cn/"
	// USGovCloud is the suffix for Vaults in Azure US Government Cloud.
	USGovCloud VaultSuffix = ".vault.usgovcloudapi.net/"
	// GermanCloud is the suffix for Vaults in Azure German Cloud.
	GermanCloud VaultSuffix = ".vault.microsoftazure.de/"
)

// VaultSuffix is the suffix added to a Vault name to create a valid Vault URL.
type VaultSuffix string

// KMSClient implements the CloudKMS interface for Azure Key Vault.
type KMSClient struct {
	client  *azsecrets.Client
	storage kms.Storage
	opts    *Opts
}

// Opts are optional settings for AKV clients.
type Opts struct {
	credentials *azidentity.DefaultAzureCredentialOptions
	client      *azsecrets.ClientOptions
}

// New initializes a KMS client for Azure Key Vault.
func New(ctx context.Context, vaultName string, vaultType VaultSuffix, store kms.Storage, opts *Opts) (*KMSClient, error) {
	if opts == nil {
		opts = &Opts{}
	}
	cred, err := azidentity.NewDefaultAzureCredential(opts.credentials)
	if err != nil {
		return nil, fmt.Errorf("loading credentials: %w", err)
	}
	client, err := azsecrets.NewClient(vaultPrefix+vaultName+string(vaultType), cred, opts.client)
	if err != nil {
		return nil, fmt.Errorf("creating vault client: %w", err)
	}

	// `azsecrets.NewClient()` does not error if the vault is non existent
	// Test here if we can reach the vault, and error otherwise
	pager := client.ListSecrets(nil)
	pager.NextPage(ctx)
	if pager.Err() != nil {
		return nil, fmt.Errorf("AKV not reachable: %w", pager.Err())
	}

	if store == nil {
		store = storage.NewMemMapStorage()
	}
	return &KMSClient{client: client, storage: store, opts: opts}, nil
}

// CreateKEK saves a new Key Encryption Key using Azure Key Vault.
//
// Keys are saved as software protected secrets.
// If no key material is provided, a new random 32 Byte key is generated and imported to the Vault.
func (c *KMSClient) CreateKEK(ctx context.Context, keyID string, key []byte) error {
	if len(key) == 0 {
		var err error
		key, err = util.GetRandomKey(config.SymmetricKeyLength)
		if err != nil {
			return fmt.Errorf("could not generate key: %w", err)
		}
	}

	// Saving symmetric keys in Azure Key Vault requires encoding them to base64
	_, err := c.client.SetSecret(ctx, keyID, base64.StdEncoding.EncodeToString(key), &azsecrets.SetSecretOptions{
		ContentType: to.StringPtr("KeyEncryptionKey"),
		Tags:        config.KmsTags,
	})
	if err != nil {
		return fmt.Errorf("importing KEK to Azure Key Vault: %w", err)
	}

	return nil
}

// GetDEK decrypts a DEK from storage.
func (c *KMSClient) GetDEK(ctx context.Context, kekID, keyID string, dekSize int) ([]byte, error) {
	kek, err := c.getKEK(ctx, kekID)
	if err != nil {
		return nil, fmt.Errorf("loading KEK from key vault: %w", err)
	}

	encryptedDEK, err := c.storage.Get(ctx, keyID)
	if err != nil {
		if !errors.Is(err, storage.ErrDEKUnset) {
			return nil, fmt.Errorf("loading encrypted DEK from storage: %w", err)
		}

		// If the DEK does not exist we generate a new random DEK and save it to storage
		newDEK, err := util.GetRandomKey(dekSize)
		if err != nil {
			return nil, fmt.Errorf("could not generate key: %w", err)
		}
		return newDEK, c.putDEK(ctx, keyID, kek, newDEK)
	}

	// Azure Key Vault does not support crypto operations with secrets, so we do the unwrapping ourselves
	return util.UnwrapAES(encryptedDEK, kek)
}

// putDEK encrypts a DEK and saves it to storage.
func (c *KMSClient) putDEK(ctx context.Context, keyID string, kek, plainDEK []byte) error {
	// Azure Key Vault does not support crypto operations with secrets, so we do the wrapping ourselves
	encryptedDEK, err := util.WrapAES(plainDEK, kek)
	if err != nil {
		return fmt.Errorf("encrypting DEK: %w", err)
	}
	return c.storage.Put(ctx, keyID, encryptedDEK)
}

// getKEK loads a Key Encryption Key from Azure Key Vault.
func (c *KMSClient) getKEK(ctx context.Context, kekID string) ([]byte, error) {
	res, err := c.client.GetSecret(ctx, kekID, nil)
	if err != nil {
		if strings.Contains(err.Error(), "SecretNotFound") {
			return nil, kms.ErrKEKUnknown
		}
		return nil, err
	}

	// Keys are saved in base64, decode them
	return base64.StdEncoding.DecodeString(*res.Value)
}