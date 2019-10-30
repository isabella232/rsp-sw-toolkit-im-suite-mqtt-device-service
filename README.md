
# RSP MQTT Device Service
Based on the Edgex Go MQTT Device Service, the RSP MQTT Device Service is a specific connector for the Intel® RSP Controller Application to EdgeX. 

RSP MQTT Device Service:
*   Registers the Intel® RSP Controller Application device with the EdgeX platform
*   Sends commands from EdgeX's [Command](https://docs.edgexfoundry.org/Ch-Command.html) service to the Intel® RSP Controller Application
*   Sends responses from the Intel® RSP Controller Application to the EdgeX [Command](https://docs.edgexfoundry.org/Ch-Command.html) service
*   Sends RFID reads from an Intel® RSP Sensor to EdgeX Core Data service

To accomplish this, modifications were made to:
*   Add multiple topics support
*   Consume RSP Controller Application messages 
*   Send commands to RSP controller Application and receive responses

## Contents
  * [Make Targets](#make-targets)
  * [Building and Launching the MQTT Device Service with EdgeX](#building-and-launching-the-mqtt-device-service-with-edgeX)
    + [Prerequisites](#prerequisites)
    + [Getting the source code](#getting-the-source-code)
    + [Building and creating the docker image](#building-and-creating-the-docker-image)
    + [Adding to EdgeX](#adding-to-edgeX)
    + [Starting the services](#starting-the-services)
  * [Sending Commands to RSP Controller Application](#sending-commands-to-rsp-controller-application)
  * [Retrieving raw sensor data from EdgeX Core Data](#retrieving-raw-sensor-data-from-edgeX-core-data)

## Make Targets
The included [Makefile](Makefile) has some useful targets for building and 
testing the service. Here's a description of these targets:

- `$(SERVICE_NAME)` (default is `mqtt-device-service`): builds the service 
- `build`: alias for `$(SERVICE_NAME)` 
- `test`: runs the test suite with coverage 
- `clean`: deletes the service executable
- `image`: builds and tags a Docker image
- `clean-img` deletes the Docker image

## Building and Launching the MQTT Device Service with EdgeX

### Prerequisites
#### Make
```bash
sudo apt -y install make
```

#### Golang (1.12+)
*   [Install Instructions](https://golang.org/doc/install)

#### EdgeX - Edinburgh Release
*   Download the latest EdgeX Edinburgh docker-compose file [here](https://raw.githubusercontent.com/edgexfoundry/developer-scripts/master/releases/edinburgh/compose-files/docker-compose-edinburgh-no-secty-1.0.1.yml) and save this as docker-compose.yml in your local directory. This file contains everything you need to deploy EdgeX with docker.

#### Intel® RSP Controller Application
*   Must have the Intel® RSP Controller Application [*Getting Started with Intel® RFID Sensor Platform (RSP) on Linux*](https://software.intel.com/en-us/getting-started-with-intel-rfid-sensor-platform-on-linux) installed and running.  This will allow for the RSP MQTT Device service to register the RSP Controller Application and the list of commands that are made available.

*   :heavy_check_mark: If you installed the DOCKER version of the Intel® RSP Controller Application, go straight to the [Getting the source code](#getting-the-source-code) section.

*   :warning: If you installed the NATIVE version of the Intel® RSP Controller Application you will need the following:

##### Curl 
```bash
sudo apt -y install curl
```

##### Docker
```bash
sudo apt -y install docker.io
```

##### Docker Compose
```bash
sudo curl -L "https://github.com/docker/compose/releases/download/1.24.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose && sudo chmod a+x /usr/local/bin/docker-compose
```

### Getting the source code
```bash
git clone https://github.impcloud.net/RSP-Inventory-Suite/mqtt-device-service.git
```

### Building and creating the docker image
```bash
cd mqtt-device-service
```

```bash
sudo make build image 
```

### Adding to EdgeX
1. To use this service with Docker go to the directory with the EdgeX `docker-compose.yml` file you downloaded in the [EdgeX prerequisites section](#EdgeX). 
2. Add the following code snippet to the DEVICE SERVICES section of the EdgeX `docker-compose.yml`.  This snippet also gives it network access to the EdgeX services and the MQTT broker. If the EdgeX services are reachable on a network named `edgex-network` (this is the default name in the EdgeX Edinburgh docker-compose.yml) and the MQTT broker is reachable via `172.17.0.1`. 

Section to add to the `docker-compose.yml` (remember spacing and alignment is important!):

```yaml
  mqtt-device-service:
    image: mqtt-device-service:latest
    networks:
        - edgex-network 
    extra_hosts:
      - "mosquitto-server:172.17.0.1"
    depends_on:
      - logging
```

### Starting the services
```bash
sudo docker-compose up -d
```


## Sending Commands to RSP Controller Application
To demonstrate sending commands from Edgex to RSP Controller Application use a web browser or tool similar to [Postman](https://www.getpostman.com/) .
 
Execute the following apis:

- Replace `localhost` in the below api with your respective server IP address if not running on localhost. This api is
used to find all the executable commands for a particular device (rsp-controller is the default name of the RSP Controller)
```
GET to http://localhost:48082/api/v1/device/name/rsp-controller
```
- If the GET request is successful a json response is received from which all the executable commands can be found

![GET device](docs/Command_list.png)

- The commands can be be sent by modifying the above api. For e.g. the below api is used to send a command known as
`behavior_get_all` 
```
GET to http://localhost:48082/api/v1/device/name/rsp-controller/command/behavior_get_all
```

- If the above request is successful a json response is received from which the RSP Controller response can be found in the
`value` field.

![GET command](docs/Response.png)


## Retrieving raw sensor data from EdgeX Core Data
To demonstrate retrieving raw RSP sensor data, the below api can be executed. If successful, a response is sent back similar to:

```json
[
    {
        "id": "ff74476a-c741-48a5-8533-22f946f29ff8",
        "created": 1572475398900,
        "origin": 1572475398882,
        "modified": 1572475398900,
        "device": "rsp-controller",
        "name": "inventory_data",
        "value": "{\"jsonrpc\":\"2.0\",\"method\":\"inventory_data\",\"params\":{\"sent_on\":1572475398919,\"period\":500,\"device_id\":\"RSP-1508b2\",\"location\":{\"latitude\":0.0,\"longitude\":0.0,\"altitude\":0.0},\"facility_id\":\"DEFAULT_FACILITY\",\"motion_detected\":false,\"data\":[{\"epc\":\"300C0000000000000000006B\",\"tid\":null,\"antenna_id\":0,\"last_read_on\":1572475398409,\"rssi\":-591,\"phase\":20,\"frequency\":911250},{\"epc\":\"300C0000000000000000006B\",\"tid\":null,\"antenna_id\":0,\"last_read_on\":1572475398484,\"rssi\":-608,\"phase\":-43,\"frequency\":911250},{\"epc\":\"300C0000000000000000006B\",\"tid\":null,\"antenna_id\":0,\"last_read_on\":1572475398602,\"rssi\":-636,\"phase\":20,\"frequency\":911250},{\"epc\":\"300C0000000000000000006B\",\"tid\":null,\"antenna_id\":0,\"last_read_on\":1572475398678,\"rssi\":-618,\"phase\":17,\"frequency\":911750},{\"epc\":\"300C0000000000000000006B\",\"tid\":null,\"antenna_id\":0,\"last_read_on\":1572475398723,\"rssi\":-618,\"phase\":-53,\"frequency\":911750},{\"epc\":\"300C0000000000000000006B\",\"tid\":null,\"antenna_id\":0,\"last_read_on\":1572475398821,\"rssi\":-618,\"phase\":15,\"frequency\":911750},{\"epc\":\"300C0000000000000000006B\",\"tid\":null,\"antenna_id\":0,\"last_read_on\":1572475398897,\"rssi\":-591,\"phase\":-43,\"frequency\":911750}]}}"
    }
]
```

  
