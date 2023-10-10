package kafka

import "fmt"

func WithPartition(partition int) KafkaReaderOpts {
	return func(k *KafkaReader) error {
		if partition < 0 {
			return fmt.Errorf("invalid partition: %d", partition)
		}

		k.partition = partition
		return nil
	}
}

func WithTlsCert(tlsCertFile string) KafkaReaderOpts {
	return func(k *KafkaReader) error {
		k.tlsCert = tlsCertFile
		return nil
	}
}

func WithTlsKey(tlsKeyFile string) KafkaReaderOpts {
	return func(k *KafkaReader) error {
		k.tlsKey = tlsKeyFile
		return nil
	}
}
