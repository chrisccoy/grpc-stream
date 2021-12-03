/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a simple gRPC client that demonstrates how to use gRPC-Go libraries
// to perform unary, client streaming, server streaming and full duplex RPCs.
//
// It interacts with the route guide service whose definition can be found in routeguide/route_guide.proto.
package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc/credentials"
	"log"
	"time"

	pb "github.com/chrisccoy/grpc-stream"
	"google.golang.org/grpc"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("server_addr", "localhost:10000", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.example.com", "The server name used to verify the hostname returned by the TLS handshake")
	timeout            = flag.Int("server_timeout", 60, "The server deadline context for client connections")
)

// runRouteChat receives a sequence of route notes, while sending notes for various locations.
func doStream(client pb.StreamerClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	stream, err := client.AckNack(ctx)
	if err != nil {
		log.Fatalf("%v.RouteChat(_) = _, %v", client, err)
	}
	defer stream.CloseSend()

	waitc := make(chan struct{})
	msg:=pb.AckMessage{}
	msg.Ping=0

	go func() {
		for {
			stream.Send(&msg)
			in, err := stream.Recv()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Ping=%d Pong=%d\n",msg.Ping+1, in.Pong)
			msg.Ping++
			time.Sleep(1*time.Second)
		}
	}()
	<-waitc
}

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	if *tls {
		if *caFile == "" {
			log.Fatal("Illegal: caFile options is blank with tls option=true")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	opts = append(opts, grpc.WithBlock())
	log.Printf("Dialing %s\n", *serverAddr)
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewStreamerClient(conn)
	log.Println("Streaming...")
	doStream(client)
}
