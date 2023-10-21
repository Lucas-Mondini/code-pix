/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/Lucas-Mondini/code-pix/application/grpc"
	"github.com/Lucas-Mondini/code-pix/application/kafka"
	"github.com/Lucas-Mondini/code-pix/infrastructure/db"
	confluentKafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/spf13/cobra"
)

var gRpcPortNumber int

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Run all systems",
	Run: func(cmd *cobra.Command, args []string) {
		database := db.ConnectDB(os.Getenv("env"))
		deliveryChan := make(chan confluentKafka.Event)
		producer := kafka.NewKafkaProducer()

		go grpc.StartGrpcServer(database, gRpcPortNumber)
		go kafka.DeliveryReport(deliveryChan)

		kafkaProcessor := kafka.NewKafkaProcessor(database, producer, deliveryChan)
		kafkaProcessor.Consume()
	},
}

func init() {
	rootCmd.AddCommand(allCmd)
	allCmd.Flags().IntVarP(&gRpcPortNumber, "grpc-port", "g", 50051, "gRPC port")
}
