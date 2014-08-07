package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/citadel/citadel"
)

type (
	Manager struct {
		baseUrl string
	}
)

func NewManager(baseUrl string) *Manager {
	m := &Manager{
		baseUrl: baseUrl,
	}
	return m
}

func (m *Manager) buildUrl(path string) string {
	return fmt.Sprintf("%s%s", m.baseUrl, path)
}

func (m *Manager) GetContainers() ([]*citadel.Container, error) {
	containers := []*citadel.Container{}
	url := m.buildUrl("/containers")
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(r.Body).Decode(&containers); err != nil {
		return nil, err
	}
	return containers, nil
}

func (m *Manager) Run(image *citadel.Image) (*citadel.Container, error) {
	b, err := json.Marshal(image)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(b)
	url := m.buildUrl("/run")
	resp, err := http.Post(url, "application/json", buf)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 201 {
		c, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(string(c))
	}
	var container citadel.Container
	if err := json.NewDecoder(resp.Body).Decode(&container); err != nil {
		return nil, err
	}
	return &container, nil
}

func (m *Manager) Destroy(container *citadel.Container) error {
	b, err := json.Marshal(container)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(b)
	url := m.buildUrl("/destroy")
	req, err := http.NewRequest("DELETE", url, buf)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 204 {
		c, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(string(c))
	}
	return nil
}
