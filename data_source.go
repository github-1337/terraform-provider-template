package main

import (
	"crypto/sha1"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIDPFingerprintRead(d *schema.ResourceData, m interface{}) error {
	idpURL := d.Get("idp_url").(string)

	fingerprint, domain, err := fetchFingerprintFromIDP(idpURL)
	if err != nil {
		return err
	}

	d.SetId("IDPFingerprint")
	d.Set("fingerprint", fingerprint)
	d.SetId("IDPDomain")
	d.Set("domain", domain)
	return nil
}

func fetchFingerprintFromIDP(idpURL string) (string, string, error) {
	resp, err := http.Get(idpURL)
	if err != nil {
		return "", "", fmt.Errorf("failed to retrieve OIDC configuration: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("received non-OK HTTP status from IdP: %s", resp.Status)
	}

	var config struct {
		JWKSURI string `json:"jwks_uri"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return "", "", fmt.Errorf("failed to decode OIDC configuration JSON: %v", err)
	}

	jwksDomain, err := url.Parse(config.JWKSURI)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse domain name from url: %v", err)
	}
	conn, err := tls.Dial("tcp", jwksDomain.Hostname()+":443", &tls.Config{ServerName: jwksDomain.Hostname()})
	if err != nil {
		return "", "", fmt.Errorf("failed to connect: %v", err)
	}
	defer conn.Close()

	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return "", "", fmt.Errorf("no certificates found")
	}
	topIntermediateCA := certs[len(certs)-1]

	// Compute the SHA-1 fingerprint of the certificate
	fingerprint := sha1.Sum(topIntermediateCA.Raw)
	return fmt.Sprintf("%x", fingerprint), jwksDomain.Hostname(), nil
}
