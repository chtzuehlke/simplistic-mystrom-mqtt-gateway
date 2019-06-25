# simplistic-mystrom-mqtt-gateway

This is a non-production-ready PoC to peridoically gather power information from mystrom switches and to publish the 
values to an AWS IoT MQTT topic.

## Steps

1. To compile the gateway: install golang

2. Required AWS IoT Core setup: create a thing with a certificate and an appropriate policy to publish MQTT messages. Also download the generated certificate and the private key. Figure out your AWS IoT endpoint.

3. Create a configuration file (conf.json)

    ```
    {
        "host": "YOUR-AWS-IOT-ENDPOINT-HERE-ats.iot.eu-west-1.amazonaws.com",
        "port": 8883,
        "caCert": "AmazonRootCA1.pem",
        "clientCert": "YOUR-CERT-HERE-certificate.pem.crt",
        "privateKey": "YOUR-PRIVATE-KEY-HERE-e7c08dd12d-private.pem.key"
    }
    ```

4. Build the gateway

    ```
    go build
    ```

5. Detect or configure mystrom switch IPs

macOS (requires nmap, ruby, curl, ifconfig):

    ```
    SWITCH_IPS=$(./macos-experimental-mystrom-switch-detection.sh)
    ```

Manual:

    ```
    export SWITCH_IPS=192.168.178.38,192.168.178.40
    ```

6. Run the gatweay

    ```
    ./simplistic-mystrom-mqtt-gateway.sh $SWITCH_IPS
    ```

7. Subscribe to topic "mystrom/power/testclient" in the MQTT client in the AWS IoT web console to see the messages produced by the gateway
