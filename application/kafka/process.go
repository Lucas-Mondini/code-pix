package kafka

import (
	"fmt"
	"os"

	confluentKafka "github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/Lucas-Mondini/code-pix/application/factory"
	appmodel "github.com/Lucas-Mondini/code-pix/application/model"
	"github.com/Lucas-Mondini/code-pix/application/usecase"
	"github.com/Lucas-Mondini/code-pix/domain/model"

	"github.com/jinzhu/gorm"
)

type KafkaProcessor struct {
	Database     *gorm.DB
	Producer     *confluentKafka.Producer
	DeliveryChan chan confluentKafka.Event
}

func NewKafkaProcessor(database *gorm.DB,
	producer *confluentKafka.Producer,
	deliveryChan chan confluentKafka.Event) *KafkaProcessor {
	return &KafkaProcessor{
		Database:     database,
		Producer:     producer,
		DeliveryChan: deliveryChan,
	}
}

func (k *KafkaProcessor) Consume() {
	configMap := &confluentKafka.ConfigMap{
		"bootstrap.servers": os.Getenv("kafkaBootstrapServers"),
		"group.id":          os.Getenv("kafkaConsumerGroupId"),
		"auto.offset.reset": "earliest",
	}

	consumer, err := confluentKafka.NewConsumer(configMap)

	if err != nil {
		panic(err)
	}

	topics := []string{os.Getenv("kafkaTransactionTopic"), os.Getenv("kafkaTransactionConfirmationTopic")}

	consumer.SubscribeTopics(topics, nil)

	fmt.Println("Kafka consumer has been started")

	for {
		msg, err := consumer.ReadMessage(-1)

		if err == nil {
			k.processMessage(msg)
		}
	}
}

func (k *KafkaProcessor) processMessage(msg *confluentKafka.Message) {
	transactionsTopic := "transactions"
	transactionsConfirmationTopic := "transaction_confirmation"

	switch topic := *msg.TopicPartition.Topic; topic {
	case transactionsTopic:
		k.processTransaction(msg)
	case transactionsConfirmationTopic:
		k.processConfirmedTransaction(msg)
	default:
		fmt.Println("message in a no valid topic: ", string(msg.Value))
	}
}

func (k *KafkaProcessor) processTransaction(msg *confluentKafka.Message) error {
	t := appmodel.NewTransaction()
	err := t.ParseJson(msg.Value)
	if err != nil {
		return err
	}

	TransactionUseCase := factory.TransactionUseCaseFactory(k.Database)

	createdTransaction, err := TransactionUseCase.Register(
		t.AccountID,
		t.Value,
		t.PixKeyTo,
		t.PixKeyKindTo,
		t.Description,
	)

	if err != nil {
		fmt.Println("error registering transaction", err)
		return err
	}

	topic := "bank" + createdTransaction.PixKeyTo.Account.Bank.Code
	t.ID = createdTransaction.ID
	t.Status = model.TransactionPending
	transactionJson, err := t.ToJson()

	if err != nil {
		return err
	}

	err = Publish(string(transactionJson), topic, k.Producer, k.DeliveryChan)

	return err
}

func (k *KafkaProcessor) processConfirmedTransaction(msg *confluentKafka.Message) error {
	t := appmodel.NewTransaction()
	err := t.ParseJson(msg.Value)
	if err != nil {
		return err
	}

	TransactionUseCase := factory.TransactionUseCaseFactory(k.Database)

	if t.Status == model.TransactionConfirmed {
		err = k.confirmTransaction(t, TransactionUseCase)
	} else if t.Status == model.TransactionCompleted {
		_, err := TransactionUseCase.Complete(t.ID)
		return err
	}
	return err
}

func (k *KafkaProcessor) confirmTransaction(transaction *appmodel.Transaction, transactionUseCase usecase.TransactionUseCase) error {
	confirmedTransaction, err := transactionUseCase.Confirm(transaction.ID)
	if err != nil {
		return err
	}

	topic := "bank" + confirmedTransaction.PixKeyTo.Account.Bank.Code
	//transaction.Status = model.TransactionConfirmed
	transactionJson, err := transaction.ToJson()
	return Publish(string(transactionJson), topic, k.Producer, k.DeliveryChan)
}
