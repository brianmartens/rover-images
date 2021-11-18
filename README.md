# rover-images

CLI for obtaining mars rover images

## How To Use

Clone this repository and install the `rover-images` CLI

```shell
git clone https://github.com/brianmartens/rover-images
cd rover-images
go install rover-images
```

**OPTIONAL: customize your experience with config.yaml**
```yaml
cache_file: /path/to/my/.cache #default is the user's $HOME/.rover-images.cache
```

Start getting Mars images!!!
```shell
rover-images get -r curiosity -C NAVCAM # curiosity NAVCAM images will be returned by default 
```

check `rover-images -h` for additional help