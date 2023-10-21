/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	confluentKafka "github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/Lucas-Mondini/code-pix/application/kafka"
	"github.com/Lucas-Mondini/code-pix/infrastructure/db"
	"github.com/spf13/cobra"
)

// kafkaCmd represents the kafka command
var kafkaCmd = &cobra.Command{
	Use:   "kafka",
	Short: "Start consuming transaction with apache Kafka",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("kafka called")

		deliveryChan := make(chan confluentKafka.Event)
		producer := kafka.NewKafkaProducer()
		database := db.ConnectDB(os.Getenv("env"))

		//kafka.Publish("olá consumer", "test", producer, deliveryChan)
		go kafka.DeliveryReport(deliveryChan)

		kafkaProcessor := kafka.NewKafkaProcessor(database, producer, deliveryChan)
		kafkaProcessor.Consume()
	},
}

func init() {
	rootCmd.AddCommand(kafkaCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// kafkaCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// kafkaCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
