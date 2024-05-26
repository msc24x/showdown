![](banner.png)
# Showdown
### Portable server to execute and judge code for multiple languages


## Run
```bash
docker stop showdown
docker remove showdown
docker build . -t showdown
docker run --privileged --name showdown --env-file .config -p 7070:7070 showdown
```
