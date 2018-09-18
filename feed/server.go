package feed

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"github.com/johnnywale/mdb/dao"
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"time"
)

type feedServer struct {
	db      *dao.FdbDao
	decoder *ProtoDecoder
}

func GenX509KeyPair() (tls.Certificate, error) {
	now := time.Now()
	template := &x509.Certificate{
		SerialNumber: big.NewInt(now.Unix()),
		Subject: pkix.Name{
			CommonName:         "quickserve.example.com",
			Country:            []string{"USA"},
			Organization:       []string{"example.com"},
			OrganizationalUnit: []string{"quickserve"},
		},
		NotBefore:             now,
		NotAfter:              now.AddDate(0, 0, 1), // Valid for one day
		SubjectKeyId:          []byte{113, 117, 105, 99, 107, 115, 101, 114, 118, 101},
		BasicConstraintsValid: true,
		IsCA:        true,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		KeyUsage: x509.KeyUsageKeyEncipherment |
			x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return tls.Certificate{}, err
	}

	cert, err := x509.CreateCertificate(rand.Reader, template, template,
		priv.Public(), priv)
	if err != nil {
		return tls.Certificate{}, err
	}

	var outCert tls.Certificate
	outCert.Certificate = append(outCert.Certificate, cert)
	outCert.PrivateKey = priv

	return outCert, nil
}

func (s *feedServer) Start() {
	cert, _ := GenX509KeyPair()
	lis, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		panic(err)
	}
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(ensureValidToken),

		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
	}
	grpcServer := grpc.NewServer(opts...)
	RegisterProtoServiceServer(grpcServer, s)
	grpcServer.Serve(lis)
}
func (s *feedServer) Register(dao *dao.FdbDao) {
	s.db = dao
}

var res empty.Empty

func (s *feedServer) OneToOne(ctx context.Context, p *Request) (*empty.Empty, error) {
	clientId := ctx.Value(userKey).(*User).ClientId
	path := append([]string{clientId}, p.Path...)
	msg := s.decoder.GetInstance("feed.Transaction")
	s.decoder.Decode(msg, p.GetProtobuf())
	id := msg.GetFieldByName("id")
	fmt.Printf("user %+v save into path %s with key %s", ctx.Value(userKey), path, id)
	s.db.Save(path, id, p.Protobuf)
	return &res, nil
}

func (s *feedServer) OneToMany(r *Request, p ProtoService_OneToManyServer) error {

	return nil
}

func (s *feedServer) ManyToOne(stream ProtoService_ManyToOneServer) error {
	Id := stream.Context().Value(userKey)
	fmt.Printf("failed to find %s", Id)
	clientId := "client_id"
	for {
		p, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&res)
		}
		if err != nil {
			return err
		}
		path := append([]string{clientId}, p.Path...)
		msg := s.decoder.GetInstance("feed.Transaction")
		s.decoder.Decode(msg, p.GetProtobuf())
		id := msg.GetFieldByName("id")
		s.db.Save(path, id, p.Protobuf)
	}
	return nil
}
func (s *feedServer) ManyToMany(ProtoService_ManyToManyServer) error {
	return nil
}

func NewServer() *feedServer {
	s := &feedServer{}	
	decode := &ProtoDecoder{}
	f, err := os.Open("./test.protoset")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	bb, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	decode.Init(bb)
	s.decoder = decode
	return s
}
