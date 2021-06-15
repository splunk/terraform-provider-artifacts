package client

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Client represents an HTTP connection and credentials.
type Client struct {
	URL      string
	Context  context.Context
	Username string
	Password string
}

// setBasicAuth adds basic auth username/password when Client has a Username set.
func (c Client) setBasicAuth(request *http.Request) error {
	if c.Username != "" {
		if c.Password == "" {
			return fmt.Errorf("username set in Client, but password unset")
		}

		request.SetBasicAuth(c.Username, c.Password)
	}

	return nil
}

// Do performs a request against the service, returning an http.Response on success. As with http.Client.Do, it is the
// responsibility of the calling function to close the resulting http.Response.body.
func (c Client) Do(request *http.Request) (response *http.Response, err error) {
	client := &http.Client{}
	response, err = client.Do(request)

	return
}

// Checksums returns the Checksums object from a remote path's file info endpoint.
func (c Client) Checksums(path string) (checksums Checksums, err error) {
	url := fmt.Sprintf("%s/api/storage/%s", c.URL, path)

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return checksums, fmt.Errorf("unable to create GET request for url %s", url)
	}

	c.setBasicAuth(request)

	response, err := c.Do(request)
	if err != nil {
		return checksums, fmt.Errorf("unable to read file info at %s: %s", url, err)
	}
	defer response.Body.Close()

	// "not found" isn't an error, it's just an empty checksum
	if response.StatusCode == 404 {
		return checksums, nil
	}

	// anything else other than a 200OK returns an error
	if response.StatusCode != 200 {
		return checksums, fmt.Errorf("response from GET %s: %s", url, response.Status)
	}

	dec := json.NewDecoder(response.Body)
	info := fileInfo{}

	if err := dec.Decode(&info); err != nil {
		return checksums, fmt.Errorf("unable to deserialize JSON properties: %s", err)
	}

	return info.Checksums, nil
}

// SHA1 returns the SHA1 checksum of a file at filename.
func (c Client) SHA1(filename string) (sha1String string, err error) {
	data, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("unable to read file %s: %s", filename, err)
	}
	defer data.Close()

	hash := sha1.New()
	if _, err := io.Copy(hash, data); err != nil {
		return "", fmt.Errorf("unable to compute sha1 for %s: %s", filename, err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// Upload performs a PUT of a file's contents to a path relative to the client's URL.
func (c Client) Upload(path string, filename string) error {
	data, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("unable to read file %s: %s", filename, err)
	}
	defer data.Close()

	url := fmt.Sprintf("%s/%s", c.URL, path)

	request, err := http.NewRequest(http.MethodPut, url, data)
	if err != nil {
		return fmt.Errorf("unable to create PUT request for %s to %s: %s", filename, url, err)
	}

	c.setBasicAuth(request)

	sha1, err := c.SHA1(filename)
	if err != nil {
		return fmt.Errorf("unable to get sha1 for %s: %s", filename, err)
	}

	request.Header.Set("X-Checksum-Sha1", sha1)

	response, err := c.Do(request)
	if err != nil {
		return fmt.Errorf("unable to perform PUT request for %s to %s: %s", filename, url, err)
	}
	response.Body.Close()

	if response.StatusCode != 201 {
		return fmt.Errorf("response from PUT %s: %s", url, response.Status)
	}

	return nil
}

// Delete performs a DELETE of a path relative to the client's URL.
func (c Client) Delete(path string) error {
	url := fmt.Sprintf("%s/%s", c.URL, path)

	request, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("unable to create DELETE request for %s: %s", url, err)
	}

	c.setBasicAuth(request)

	response, err := c.Do(request)
	if err != nil {
		return fmt.Errorf("unable to perform DELETE request for %s: %s", url, err)
	}
	response.Body.Close()

	if response.StatusCode != 204 {
		return fmt.Errorf("response from DELETE %s: %s", url, response.Status)
	}

	return nil
}
