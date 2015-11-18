# Buttons APP

Make sure you have [Docker](https://www.docker.com/) installed and build the
app with:

```sh
git clone https://github.com/fgrehm/pucrs-iot-experiment.git
cd pucrs-iot-experiment
make build
```

To deploy to the device:

```sh
# Change the Makefile with the IP of the Raspberry PI
make deploy.rpi
```
