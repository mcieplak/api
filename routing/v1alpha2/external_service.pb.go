// Code generated by protoc-gen-go. DO NOT EDIT.
// source: routing/v1alpha2/external_service.proto

package v1alpha2

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Different ways of discovering the IP addresses associated with the
// service.
type ExternalService_Discovery int32

const (
	// If set to "none", the proxy will assume that incoming connections
	// have already been resolved (to a specific destination IP
	// address). Such connections are typically routed via the proxy using
	// mechanisms such as IP table REDIRECT/ eBPF. After performing any
	// routing related transformations, the proxy will forward the
	// connection to the IP address to which the connection was bound.
	ExternalService_NONE ExternalService_Discovery = 0
	// If set to "static", the proxy will use the IP addresses specified in
	// endpoints (See below) as the backing nodes associated with the
	// external service.
	ExternalService_STATIC ExternalService_Discovery = 1
	// If set to "dns", the proxy will attempt to resolve the DNS address
	// during request processing. If no endpoints are specified, the proxy
	// will resolve the DNS address specified in the hosts field, if
	// wildcards are not used. If endpoints are specified, the DNS
	// addresses specified in the endpoints will be resolved to determine
	// the destination IP address.
	ExternalService_DNS ExternalService_Discovery = 2
)

var ExternalService_Discovery_name = map[int32]string{
	0: "NONE",
	1: "STATIC",
	2: "DNS",
}
var ExternalService_Discovery_value = map[string]int32{
	"NONE":   0,
	"STATIC": 1,
	"DNS":    2,
}

func (x ExternalService_Discovery) String() string {
	return proto.EnumName(ExternalService_Discovery_name, int32(x))
}
func (ExternalService_Discovery) EnumDescriptor() ([]byte, []int) { return fileDescriptor1, []int{0, 0} }

// External service describes the endpoints, ports and protocols of a
// white-listed set of mesh-external domains and IP blocks that services in
// the mesh are allowed to access.
//
// For example, the following external service configuration describes the
// set of services at https://example.com to be accessed internally over
// plaintext http (i.e. http://example.com:443), with the sidecar originating
// TLS.
//
//     apiVersion: config.istio.io/v1alpha2
//     kind: ExternalService
//     metadata:
//       name: external-svc-example
//     spec:
//       hosts:
//       - example.com
//       ports:
//       - number: 443
//         name: example-http
//         protocol: http # not HTTPS.
//       discovery: dns
//
// and a destination rule to initiate TLS connections to the external service.
//
//     apiVersion: config.istio.io/v1alpha2
//     kind: DestinationRule
//     metadata:
//       name: tls-example
//     spec:
//       destination:
//         name: example.com
//       tls:
//         mode: simple # initiates HTTPS when talking to example.com
//
// The following specification specifies a static set of backend nodes for
// a MongoDB cluster behind a set of virtual IPs, and sets up a destination
// rule to initiate mTLS connections upstream.
//
//     apiVersion: config.istio.io/v1alpha2
//     kind: ExternalService
//     metadata:
//       name: external-svc-mongocluster
//     spec:
//       hosts:
//       - 192.192.192.192/24
//       ports:
//       - number: 27018
//         name: mongodb
//         protocol: mongo
//       discovery: static
//       endpoints:
//       - address: 2.2.2.2
//       - address: 3.3.3.3
//
// and the associated destination rule
//
//     apiVersion: config.istio.io/v1alpha2
//     kind: DestinationRule
//     metadata:
//       name: mtls-mongocluster
//     spec:
//       destination:
//         name: 192.192.192.192/24
//       tls:
//         mode: mutual
//         clientCertificate: /etc/certs/myclientcert.pem
//         privateKey: /etc/certs/client_private_key.pem
//         caCertificates: /etc/certs/rootcacerts.pem
//
// The following example demonstrates the use of wildcards in the hosts. If
// the connection has to be routed to the IP address requested by the
// application (i.e. application resolves DNS and attempts to connect to a
// specific IP), the discovery mode must be set to "none".
//
//     apiVersion: config.istio.io/v1alpha2
//     kind: ExternalService
//     metadata:
//       name: external-svc-wildcard-example
//     spec:
//       hosts:
//       - *.bar.com
//       ports:
//       - number: 80
//         name: http
//         protocol: http
//       discovery: none
//
// For HTTP based services, it is possible to create a virtual service
// backed by multiple DNS addressible endpoints. In such a scenario, the
// application can use the HTTP_PROXY environment variable to transparently
// reroute API calls for the virtual service to a chosen backend. For
// example, the following configuration creates a non-existent service
// called foo.bar.com backed by three domains: us.foo.bar.com:8443,
// uk.foo.bar.com:9443, and in.foo.bar.com:7443
//
//     apiVersion: config.istio.io/v1alpha2
//     kind: ExternalService
//     metadata:
//       name: external-svc-dns
//     spec:
//       hosts:
//       - foo.bar.com
//       ports:
//       - number: 443
//         name: https
//         protocol: http
//       discovery: dns
//       endpoints:
//       - address: us.foo.bar.com
//         ports:
//         - https: 8443
//       - address: uk.foo.bar.com
//         ports:
//         - https: 9443
//       - address: in.foo.bar.com
//         ports:
//         - https: 7443
//
// and a destination rule to initiate TLS connections to the external service.
//
//     apiVersion: config.istio.io/v1alpha2
//     kind: DestinationRule
//     metadata:
//       name: tls-foobar
//     spec:
//       destination:
//         name: foo.bar.com
//       tls:
//         mode: simple # initiates HTTPS
//
// With HTTP_PROXY=http://localhost:443, calls from the application to
// http://foo.bar.com will be upgraded to HTTPS and load balanced across
// the three domains specified above. In other words, a call to
// http://foo.bar.com/baz would be translated to
// https://uk.foo.bar.com/baz.
//
// NOTE: In the scenario above, the value of the HTTP Authority/host header
// associated with the outbound HTTP requests will be based on the
// endpoint's DNS name, i.e. ":authority: uk.foo.bar.com". Refer to Envoy's
// auto_host_rewrite for further details. The automatic rewrite can be
// overridden using a host rewrite route rule.
//
type ExternalService struct {
	// REQUIRED. The hosts associated with the external service. Could be a
	// DNS name with wildcard prefix or a CIDR prefix. Note that the hosts
	// field applies to all protocols. DNS names in hosts will be ignored if
	// the application accesses the service over non-HTTP protocols such as
	// mongo/opaque TCP/even HTTPS. In such scenarios, the port on which the
	// external service is being accessed must not be shared by any other
	// service in the mesh. In other words, the sidecar will behave as a
	// simple TCP proxy, forwarding incoming traffic on a specified port to
	// the specified destination endpoint IP/host.
	Hosts []string `protobuf:"bytes,1,rep,name=hosts" json:"hosts,omitempty"`
	// REQUIRED. The ports associated with the external service.
	Ports []*Port `protobuf:"bytes,2,rep,name=ports" json:"ports,omitempty"`
	// Service discovery mode for the hosts. If not set, Istio will attempt
	// to infer the discovery mode based on the value of hosts and endpoints.
	Discovery ExternalService_Discovery `protobuf:"varint,3,opt,name=discovery,enum=istio.routing.v1alpha2.ExternalService_Discovery" json:"discovery,omitempty"`
	// One or more endpoints associated with the service. Endpoints must be
	// accessible over the set of outPorts defined at the service level.
	Endpoints []*ExternalService_Endpoint `protobuf:"bytes,4,rep,name=endpoints" json:"endpoints,omitempty"`
}

func (m *ExternalService) Reset()                    { *m = ExternalService{} }
func (m *ExternalService) String() string            { return proto.CompactTextString(m) }
func (*ExternalService) ProtoMessage()               {}
func (*ExternalService) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *ExternalService) GetHosts() []string {
	if m != nil {
		return m.Hosts
	}
	return nil
}

func (m *ExternalService) GetPorts() []*Port {
	if m != nil {
		return m.Ports
	}
	return nil
}

func (m *ExternalService) GetDiscovery() ExternalService_Discovery {
	if m != nil {
		return m.Discovery
	}
	return ExternalService_NONE
}

func (m *ExternalService) GetEndpoints() []*ExternalService_Endpoint {
	if m != nil {
		return m.Endpoints
	}
	return nil
}

// Endpoint defines a network address (IP or hostname) associated with
// the external service.
type ExternalService_Endpoint struct {
	// REQUIRED: Address associated with the network endpoint without the
	// port ( IP or fully qualified domain name without wildcards).
	Address string `protobuf:"bytes,1,opt,name=address" json:"address,omitempty"`
	// Set of ports associated with the endpoint. The ports must be
	// associated with a port name that was declared as part of the
	// service.
	Ports map[string]uint32 `protobuf:"bytes,2,rep,name=ports" json:"ports,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"varint,2,opt,name=value"`
	// One or more labels associated with the endpoint.
	Labels map[string]string `protobuf:"bytes,3,rep,name=labels" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *ExternalService_Endpoint) Reset()                    { *m = ExternalService_Endpoint{} }
func (m *ExternalService_Endpoint) String() string            { return proto.CompactTextString(m) }
func (*ExternalService_Endpoint) ProtoMessage()               {}
func (*ExternalService_Endpoint) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0, 0} }

func (m *ExternalService_Endpoint) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *ExternalService_Endpoint) GetPorts() map[string]uint32 {
	if m != nil {
		return m.Ports
	}
	return nil
}

func (m *ExternalService_Endpoint) GetLabels() map[string]string {
	if m != nil {
		return m.Labels
	}
	return nil
}

func init() {
	proto.RegisterType((*ExternalService)(nil), "istio.routing.v1alpha2.ExternalService")
	proto.RegisterType((*ExternalService_Endpoint)(nil), "istio.routing.v1alpha2.ExternalService.Endpoint")
	proto.RegisterEnum("istio.routing.v1alpha2.ExternalService_Discovery", ExternalService_Discovery_name, ExternalService_Discovery_value)
}

func init() { proto.RegisterFile("routing/v1alpha2/external_service.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 374 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x92, 0x4b, 0x6f, 0x9b, 0x40,
	0x10, 0xc7, 0x0b, 0xf8, 0xc5, 0x58, 0x6d, 0xd1, 0xaa, 0xaa, 0x10, 0xea, 0x03, 0xf9, 0x52, 0xd4,
	0x03, 0xd4, 0xf4, 0xe2, 0x3e, 0x2e, 0x6d, 0xcc, 0x21, 0x52, 0x84, 0x13, 0xec, 0x53, 0x2e, 0xd1,
	0xda, 0xac, 0xec, 0x55, 0x10, 0x8b, 0x76, 0xd7, 0x24, 0x7c, 0xbd, 0x7c, 0x9f, 0x7c, 0x87, 0xc8,
	0x3c, 0x62, 0xc7, 0x49, 0x94, 0xf8, 0xb6, 0xb3, 0x9a, 0xff, 0x6f, 0xe6, 0x3f, 0x33, 0xf0, 0x8d,
	0xb3, 0xb5, 0xa4, 0xe9, 0xd2, 0xcb, 0x87, 0x38, 0xc9, 0x56, 0xd8, 0xf7, 0xc8, 0xb5, 0x24, 0x3c,
	0xc5, 0xc9, 0x85, 0x20, 0x3c, 0xa7, 0x0b, 0xe2, 0x66, 0x9c, 0x49, 0x86, 0x3e, 0x52, 0x21, 0x29,
	0x73, 0xeb, 0x74, 0xb7, 0x49, 0xb7, 0xbe, 0x3c, 0x02, 0x2c, 0xb1, 0x24, 0x57, 0xb8, 0xa8, 0x74,
	0x83, 0xdb, 0x16, 0xbc, 0x0f, 0x6a, 0xe4, 0xb4, 0x22, 0xa2, 0x0f, 0xd0, 0x5e, 0x31, 0x21, 0x85,
	0xa9, 0xd8, 0x9a, 0xa3, 0x47, 0x55, 0x80, 0x7c, 0x68, 0x67, 0x8c, 0x4b, 0x61, 0xaa, 0xb6, 0xe6,
	0xf4, 0xfd, 0x4f, 0xee, 0xd3, 0x15, 0xdd, 0x53, 0xc6, 0x65, 0x54, 0xa5, 0xa2, 0x09, 0xe8, 0x31,
	0x15, 0x0b, 0x96, 0x13, 0x5e, 0x98, 0x9a, 0xad, 0x38, 0xef, 0xfc, 0xe1, 0x73, 0xba, 0xbd, 0x2e,
	0xdc, 0x71, 0x23, 0x8c, 0xb6, 0x0c, 0x14, 0x82, 0x4e, 0xd2, 0x38, 0x63, 0x34, 0x95, 0xc2, 0x6c,
	0x95, 0x8d, 0xfc, 0x78, 0x2d, 0x30, 0xa8, 0x85, 0xd1, 0x16, 0x61, 0xdd, 0xa8, 0xd0, 0x6b, 0xfe,
	0x91, 0x09, 0x5d, 0x1c, 0xc7, 0x9c, 0x88, 0x8d, 0x73, 0xc5, 0xd1, 0xa3, 0x26, 0x44, 0x67, 0x0f,
	0xbd, 0xff, 0x39, 0xb4, 0x64, 0x39, 0x14, 0x11, 0xa4, 0x92, 0x17, 0xcd, 0x68, 0x66, 0xd0, 0x49,
	0xf0, 0x9c, 0x24, 0xc2, 0xd4, 0x4a, 0xe6, 0xdf, 0x83, 0x99, 0x27, 0xa5, 0xbc, 0x82, 0xd6, 0x2c,
	0x6b, 0x04, 0xb0, 0x2d, 0x85, 0x0c, 0xd0, 0x2e, 0x49, 0x51, 0x9b, 0xd9, 0x3c, 0x37, 0xab, 0xcd,
	0x71, 0xb2, 0x26, 0xa6, 0x6a, 0x2b, 0xce, 0xdb, 0xa8, 0x0a, 0x7e, 0xab, 0x23, 0xc5, 0xfa, 0x05,
	0xfd, 0x1d, 0xe0, 0x4b, 0x52, 0x7d, 0x47, 0x3a, 0xf8, 0x0e, 0xfa, 0xfd, 0xb2, 0x50, 0x0f, 0x5a,
	0xe1, 0x24, 0x0c, 0x8c, 0x37, 0x08, 0xa0, 0x33, 0x9d, 0xfd, 0x9b, 0x1d, 0x1f, 0x19, 0x0a, 0xea,
	0x82, 0x36, 0x0e, 0xa7, 0x86, 0xfa, 0xff, 0xeb, 0xf9, 0xe7, 0xca, 0x27, 0x65, 0x1e, 0xce, 0xa8,
	0xb7, 0x7f, 0x9e, 0xf3, 0x4e, 0x79, 0x97, 0x3f, 0xef, 0x02, 0x00, 0x00, 0xff, 0xff, 0xfe, 0xc4,
	0x3a, 0x6d, 0xfa, 0x02, 0x00, 0x00,
}
