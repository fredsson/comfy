# Comfy
This is a small application that can enable/disable smartmode in sensibo depending on the current electricity prices


## Environment
The application looks for a `.env` file in the root directory with the following content
```
TIBBER_API_KEY=API Key from Tibber
SENSIBO_API_KEY=API Key from Sensibo
```

## Deploying with local portainer

Create a tarball for building the docker image

```
$  tar -czvf prod.tar.gz .
```

Once this has been created, go to the local portainer environment and click images, once inside images build a new image with the upload tarball option (dont forget to set path to dockerfile)


After the image has been sucessfully built go to containers and click comfy, click duplicate/edit and change the image to the newly created image tag. Finally click deploy the container and wait for it to start.

To ensure it's running correctly click comfy again and check the logs link, this should show that an inital check has been made and how long it will wait until next one.
