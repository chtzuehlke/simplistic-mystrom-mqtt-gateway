# simplistic-mystrom-mqtt-gateway

This is a non-production-ready PoC to peridoically gather power information from mystrom switches and to publish the 
values to an AWS IoT MQTT topic.

## Steps

1. To compile the gateway: install golang

2. To conveniently detect mystrom devices (see simplistic-mystrom-mqtt-gateway.sh): install nmap and ruby and ensure ifconfig and curl are available (tested with macOS)

3. Required AWS IoT Core setup: create a thing with a certificate and an appropriate policy to publish MQTT messages. Also download the generated certificate and the private key. Figure out your AWS IoT endpoint.

4. Create a configuration file (conf.json)

    ```
    {
        "host": "YOUR-AWS-IOT-ENDPOINT-HERE-ats.iot.eu-west-1.amazonaws.com",
        "port": 8883,
        "caCert": "AmazonRootCA1.pem",
        "clientCert": "YOUR-CERT-HERE-certificate.pem.crt",
        "privateKey": "YOUR-PRIVATE-KEY-HERE-e7c08dd12d-private.pem.key"
    }
    ```

5. Build the gateway

    ```
    go build
    ```

6. Run the gatweay

    ```
    ./simplistic-mystrom-mqtt-gateway.sh
    ```

6. Subscribe to topic "mystrom/power/testclient" in the MQTT client in the AWS IoT web console to see the messages produced by the gateway
