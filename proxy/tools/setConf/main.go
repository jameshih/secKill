package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/jameshih/secKill/proxy/service"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"
)

const (
	EtcdKey = "/seckill/product"
)

var (
	cli *clientv3.Client
	err error
)

func init() {
	cfg := clientv3.Config{
		Endpoints:   []string{"http://127.0.0.1:2379"},
		DialTimeout: 2 * time.Second,
	}
	cli, err = clientv3.New(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func setLogConfigToEtcd() {
	fmt.Println("connect succ")
	var secInfoArr []service.ProductInfoConf
	secInfoArr = append(secInfoArr, service.ProductInfoConf{
		ProductID: 100,
		StartTime: 1589231756,
		EndTime:   1589404556,
		Status:    0,
		Total:     1000,
		Left:      1000,
	})

	secInfoArr = append(secInfoArr, service.ProductInfoConf{
		ProductID: 102,
		StartTime: 1589231756,
		EndTime:   1589404556,
		Status:    0,
		Total:     2000,
		Left:      2000,
	})

	data, err := json.Marshal(secInfoArr)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	_, err = cli.Put(ctx, EtcdKey, string(data))
	cancel()
	if err != nil {
		switch err {
		case context.Canceled:
			log.Fatalf("ctx is canceled by another routine: %v", err)
		case context.DeadlineExceeded:
			log.Fatalf("ctx is attached with a deadline is exceeded: %v", err)
		case rpctypes.ErrEmptyKey:
			log.Fatalf("client-side error: %v", err)
		default:
			log.Fatalf("bad cluster endpoints, which are not etcd servers: %v", err)
		}
	}
	fmt.Println("setting value to etcd...")

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, EtcdKey)
	cancel()
	if err != nil {
		switch err {
		case context.Canceled:
			log.Fatalf("ctx is canceled by another routine: %v", err)
		case context.DeadlineExceeded:
			log.Fatalf("ctx is attached with a deadline is exceeded: %v", err)
		case rpctypes.ErrEmptyKey:
			log.Fatalf("client-side error: %v", err)
		default:
			log.Fatalf("bad cluster endpoints, which are not etcd servers: %v", err)
		}
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s: %s\n", ev.Key, ev.Value)
	}

}

func deleteFromEtcd() {
	// delete log testing
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	cli.Delete(ctx, EtcdKey)
	fmt.Printf("deleted %s", EtcdKey)
	cancel()
	return
}

func main() {
	flag.Parse()
	arg := flag.Arg(0)
	switch arg {
	case "add":
		setLogConfigToEtcd()
	case "del":
		deleteFromEtcd()
	}
	defer cli.Close()
}
