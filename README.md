# RSP MQTT Device Service
Based on the Edgex Go MQTT Device Service, the RSP MQTT Device Service can be 
used to connect the Intel® RSP Controller Application to EdgeX. 

The RSP MQTT Device Service:
- Registers the Intel® RSP Controller Application device with the EdgeX platform
- Accepts commands from EdgeX's [Command](https://docs.edgexfoundry.org/Ch-Command.html) 
    service and translates/forwards them to the Intel® RSP Controller Application
- Translates/forwards responses from the Intel® RSP Controller Application to the 
    EdgeX [Command](https://docs.edgexfoundry.org/Ch-Command.html) service
- Translates/sends RFID reads from an Intel® RSP Sensor to EdgeX Core Data service

To accomplish this, modifications were made to:
- Support subscriptions to multiple topics 
- Consume RSP Controller Application-specific messages 
- Register new devices on first-discovery
- Translate RSP Controller Application commands and responses
- Validate incoming messages against expected schemas

## Contents
- [Building and Launching the MQTT Device Service with EdgeX](#Building-and-Launching-the-MQTT-Device-Service-with-EdgeX)
    - [Prerequisites](#prerequisites)
        - [Intel® RSP Controller Application](#Intel®-RSP-Controller-Application)
        - [EdgeX, Edinburgh Release](#EdgeX,-Edinburgh-Release)
    - [Getting the Source Code](#getting-the-source-code)
    - [Building and Creating the Docker Image](#building-and-creating-the-docker-image)
        - [Make Targets](#make-targets)
    - [Adding to EdgeX](#Adding-to-EdgeX)
    - [Starting the Services](#starting-the-services)
- [Sending Commands to RSP Controller Application](#sending-commands-to-rsp-controller-application)
    - [Listing Commands](#listing-commands)
- [Retrieving Raw Sensor Data from EdgeX Core Data](#Retrieving-raw-sensor-data-from-EdgeX-Core-Data)


## Building and Launching the MQTT Device Service with EdgeX
### Prerequisites
You'll need the follow software packages to follow the instructions below; you 
should refer to your distribution's package management documentation for more
specific installation instructions if the `apt`-based commands are not relevant
for your OS. 

- make: `sudo apt -y install make`
- Go: [Install Instructions](https://golang.org/doc/install)
- Intel® RSP Controller Application: instructions [below](#Intel®-RSP-Controller-Application)
- EdgeX, Edinburgh Release: instructions [below](#EdgeX,-Edinburgh-Release)

>   :heavy_check_mark: If you installed the _Docker_ version of the Intel® RSP 
    Controller Application, you already have the following dependencies:

>   :warning: If you installed the _native_ version of the Intel® RSP Controller 
    Application, you will also need these to run EdgeX and RSP MQTT Device Service
    in Docker:

- curl: `sudo apt -y install curl`
- Docker: `sudo apt -y install docker.io`
- docker-compose: 
```bash 
    sudo curl \
    -L "https://github.com/docker/compose/releases/download/1.24.0/docker-compose-$(uname -s)-$(uname -m)" \
    -o /usr/local/bin/docker-compose && \
    sudo chmod a+x /usr/local/bin/docker-compose
```

#### Intel® RSP Controller Application
This service connects the Intel® RSP Controller Application to EdgeX, so you
should follow the [*Getting Started with Intel® RFID Sensor Platform (RSP) on Linux*](https://software.intel.com/en-us/getting-started-with-intel-rfid-sensor-platform-on-linux) 
to ensure it is installed and running. The RSP MQTT Device service connects to 
its MQTT broker, registers the it and its commands with EdgeX, and handles the
communication between the two of them.

#### EdgeX, Edinburgh Release 
The instructions in this `README` expect that you're running EdgeX's Docker services.
If you haven't already, you can download the EdgeX Edinburgh `docker-compose` file 
[here](https://raw.githubusercontent.com/edgexfoundry/developer-scripts/master/releases/edinburgh/compose-files/docker-compose-edinburgh-no-secty-1.0.1.yml).
Save it as `docker-compose.yml`. This file contains the service descriptions needed
to deploy EdgeX with docker; you'll edit it later to add the RSP MQTT Device Service. 
Refer to EdgeX's documentation for more information about running EdgeX and adding
device services.

### Getting the Source Code
Simply clone this repository, preferably to a shorter directory name like `mqtt-device-service`:
```bash
git clone https://github.com/intel/rsp-sw-toolkit-im-suite-mqtt-device-service.git mqtt-device-service
```

### Building and Creating the Docker Image
Go to the directory where you cloned the repo and run `make build image`; you  
may need `sudo` rights if you are not in the `docker` group: 
```bash
cd mqtt-device-service
sudo make build image 
```

#### Make Targets
The included [Makefile](Makefile) has some other useful targets for building and 
testing the service. Here's a quick description of these targets:

- `$(SERVICE_NAME)` (default is `mqtt-device-service`): builds the service 
- `build`: alias for `$(SERVICE_NAME)` 
- `test`: runs the test suite with coverage 
- `clean`: deletes the service executable
- `image`: builds and tags a Docker image
- `clean-img` deletes the Docker image

### Adding to EdgeX
1. To use this service with Docker, go to the directory with the EdgeX 
    `docker-compose.yml` file you downloaded in the [EdgeX prerequisites section](#EdgeX,-Edinburgh-Release). 
2. Add the following code snippet to the DEVICE SERVICES section of the EdgeX 
    `docker-compose.yml`.  This snippet also gives it network access to the 
    EdgeX services and the MQTT broker. If the EdgeX services are reachable on a 
    network named `edgex-network` (this is the default name in the EdgeX 
    Edinburgh docker-compose.yml) and the MQTT broker is reachable via 
    `172.17.0.1`. 

Section to add to the `docker-compose.yml` (remember spacing and alignment is 
important!):

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

### Starting the Services
Use `docker-compose` to launch the services. This command must be run within the
directory of your `docker-compose.yml` file; you may need `sudo` rights if your
user is not part of the `docker` group:
```bash
sudo docker-compose up -d
```


## Sending Commands to RSP Controller Application
You can receive data and send commands to RSP Controller Application via EdgeX.
The following demonstrates this using a web browser and `curl`, though you can 
use any tool capable of sending HTTP requests. 
 
> :heavy_exclamation_mark: The following examples use `localhost`; if your EdgeX 
> instance is running elsewhere, replace it with the relevant IP address.
 
> :heavy_exclamation_mark: The following examples use the default EdgeX ports; 
> if your EdgeX instance is using non-standard ports, replace them with the 
> relevant ones.


### Listing Commands
This API is used to find all the executable commands for a particular device;
`rsp-controller` is the default name of the RSP Controller, so we'll use it to 
get the available `rsp-controller` commands. Because it's a `GET` request, you 
can [view it in your browser](http://localhost:48082/api/v1/device/name/rsp-controller), 
or use `curl` to retrieve the output:
    
    curl -o- http://localhost:48082/api/v1/device/name/rsp-controller

If the request is successful, you'll get a JSON response listing the commands. 

> :heavy_check_mark: you may find it helpful to use tools like `jq`, Firefox,
> Postman, or Chrome's DevTools to format the JSON output. The images below show
> the formatted JSON as rendered by Postman.

![GET device](docs/Command_list.png)

The response includes the URLs of the available commands. You can make `GET` 
requests to these endpoints to execute the commands. For example, 
[this API](http://localhost:48082/api/v1/device/name/rsp-controller/command/behavior_get_all)
sends the command `behavior_get_all`:

    curl -o- http://localhost:48082/api/v1/device/name/rsp-controller/command/behavior_get_all

> remember that you may need to modify the host to match your Docker host's IP address

The output from EdgeX represents the `readings` generated as a result of the command;
see EdgeX's documentation for more information about `readings` and `events`, but
note that the RSP Controller's response is encoded in the `value` field of the
first `reading`.

![GET command](docs/Response.png)

## Retrieving raw sensor data from EdgeX Core Data
You can get the raw sensor data from EdgeX's Core Data Service by using those APIs.
For example, [this endpoint](http://localhost:48080/api/v1/reading/device/rsp-controller/1)
returns the most recent data sent by an RSP Sensor, encoded in the `value`: 

    curl -o- http://localhost:48080/api/v1/reading/device/rsp-controller/1

The response is an array of `readings` (in this case, the array has only 1 value):
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

  
