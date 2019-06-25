# simplistic-mystrom-mqtt-gateway

## Steps

0. install golang

1. create thing, cert & policy in AWS IoT and download certificate and private key

2. create config file (conf.json)

    ```
    {
        "host": "YOUR-AWS-IOT-ENDPOINT-HERE-ats.iot.eu-west-1.amazonaws.com",
        "port": 8883,
        "caCert": "AmazonRootCA1.pem",
        "clientCert": "YOUR-CERT-HERE-certificate.pem.crt",
        "privateKey": "YOUR-PRIVATE-KEY-HERE-e7c08dd12d-private.pem.key"
    }
    ```

3. build

    ```
    go build
    ```

4. run (you might want to create test subscription to "go-mqtt/sample" in the AWS IoT web console)

    ```
    ./simplistic-mystrom-mqtt-gateway -conf conf.json -client-id testclient -mystrom-switch-url https://httpecho-167219.appspot.com/mystrom/switch/sampleresponse1.json
    ```

## Sources partially "borrowed" from here

- https://github.com/chtzuehlke/helloworld-go-mqtt-awsiot
