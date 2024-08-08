package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
)

type Client struct {
	client    *api.Client
	serviceID string
}

func newClient(serviceID string) *Client {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Printf("fialed to create consul client %v", err)
		return nil
	}
	return &Client{
		client:    client,
		serviceID: serviceID,
	}
}

func (c *Client) getService(serviceName string) (string, error) {
	//获取服务信息
	services, err := c.client.Agent().Services()
	if err != nil {
		log.Printf("failed to get services %v", err)
		return "", err
	}

	var serviceAddr string
	for _, service := range services {
		if service.Service == serviceName {
			serviceAddr = fmt.Sprintf("%s:%d", service.Address, service.Port)
			break
		}
	}

	if serviceAddr == "" {
		log.Printf("service %v not found in consul", serviceName)
		return "", err
	}

	return serviceAddr, nil
}

func (c *Client) registerService(serviceName string, ip string, port int) error {

	//注册服务信息
	registration := &api.AgentServiceRegistration{
		ID:      c.serviceID,
		Name:    serviceName,
		Address: ip,
		Port:    port,
	}

	err := c.client.Agent().ServiceRegister(registration)
	if err != nil {
		return err
	}
	return nil

}

func (c *Client) Close() error {
	return c.client.Agent().ServiceDeregister(c.serviceID)
}
