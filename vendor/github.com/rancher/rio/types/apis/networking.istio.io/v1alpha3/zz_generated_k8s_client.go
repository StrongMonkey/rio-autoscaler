package v1alpha3

import (
	"context"
	"sync"

	"github.com/rancher/norman/controller"
	"github.com/rancher/norman/objectclient"
	"github.com/rancher/norman/objectclient/dynamic"
	"github.com/rancher/norman/restwatch"
	"k8s.io/client-go/rest"
)

type (
	contextKeyType        struct{}
	contextClientsKeyType struct{}
)

type Interface interface {
	RESTClient() rest.Interface
	controller.Starter

	GatewaysGetter
	VirtualServicesGetter
	DestinationRulesGetter
	ServiceEntriesGetter
}

type Clients struct {
	Interface Interface

	Gateway         GatewayClient
	VirtualService  VirtualServiceClient
	DestinationRule DestinationRuleClient
	ServiceEntry    ServiceEntryClient
}

type Client struct {
	sync.Mutex
	restClient rest.Interface
	starters   []controller.Starter

	gatewayControllers         map[string]GatewayController
	virtualServiceControllers  map[string]VirtualServiceController
	destinationRuleControllers map[string]DestinationRuleController
	serviceEntryControllers    map[string]ServiceEntryController
}

func Factory(ctx context.Context, config rest.Config) (context.Context, controller.Starter, error) {
	c, err := NewForConfig(config)
	if err != nil {
		return ctx, nil, err
	}

	cs := NewClientsFromInterface(c)

	ctx = context.WithValue(ctx, contextKeyType{}, c)
	ctx = context.WithValue(ctx, contextClientsKeyType{}, cs)
	return ctx, c, nil
}

func ClientsFrom(ctx context.Context) *Clients {
	return ctx.Value(contextClientsKeyType{}).(*Clients)
}

func From(ctx context.Context) Interface {
	return ctx.Value(contextKeyType{}).(Interface)
}

func NewClients(config rest.Config) (*Clients, error) {
	iface, err := NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return NewClientsFromInterface(iface), nil
}

func NewClientsFromInterface(iface Interface) *Clients {
	return &Clients{
		Interface: iface,

		Gateway: &gatewayClient2{
			iface: iface.Gateways(""),
		},
		VirtualService: &virtualServiceClient2{
			iface: iface.VirtualServices(""),
		},
		DestinationRule: &destinationRuleClient2{
			iface: iface.DestinationRules(""),
		},
		ServiceEntry: &serviceEntryClient2{
			iface: iface.ServiceEntries(""),
		},
	}
}

func NewForConfig(config rest.Config) (Interface, error) {
	if config.NegotiatedSerializer == nil {
		config.NegotiatedSerializer = dynamic.NegotiatedSerializer
	}

	restClient, err := restwatch.UnversionedRESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &Client{
		restClient: restClient,

		gatewayControllers:         map[string]GatewayController{},
		virtualServiceControllers:  map[string]VirtualServiceController{},
		destinationRuleControllers: map[string]DestinationRuleController{},
		serviceEntryControllers:    map[string]ServiceEntryController{},
	}, nil
}

func (c *Client) RESTClient() rest.Interface {
	return c.restClient
}

func (c *Client) Sync(ctx context.Context) error {
	return controller.Sync(ctx, c.starters...)
}

func (c *Client) Start(ctx context.Context, threadiness int) error {
	return controller.Start(ctx, threadiness, c.starters...)
}

type GatewaysGetter interface {
	Gateways(namespace string) GatewayInterface
}

func (c *Client) Gateways(namespace string) GatewayInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &GatewayResource, GatewayGroupVersionKind, gatewayFactory{})
	return &gatewayClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}

type VirtualServicesGetter interface {
	VirtualServices(namespace string) VirtualServiceInterface
}

func (c *Client) VirtualServices(namespace string) VirtualServiceInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &VirtualServiceResource, VirtualServiceGroupVersionKind, virtualServiceFactory{})
	return &virtualServiceClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}

type DestinationRulesGetter interface {
	DestinationRules(namespace string) DestinationRuleInterface
}

func (c *Client) DestinationRules(namespace string) DestinationRuleInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &DestinationRuleResource, DestinationRuleGroupVersionKind, destinationRuleFactory{})
	return &destinationRuleClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}

type ServiceEntriesGetter interface {
	ServiceEntries(namespace string) ServiceEntryInterface
}

func (c *Client) ServiceEntries(namespace string) ServiceEntryInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &ServiceEntryResource, ServiceEntryGroupVersionKind, serviceEntryFactory{})
	return &serviceEntryClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}