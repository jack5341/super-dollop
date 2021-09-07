<p align="center">
  <img height="300" src="https://user-images.githubusercontent.com/53150440/132265128-ae52428b-ace5-423e-b726-5d6099d00206.png"/><br/>
  <a>
        <img src="https://img.shields.io/github/v/release/jack5341/super-dollop?style=flat&labelColor=1C2C2E&color=abc3d6&logo=GitHub&logoColor=white">
  </a>
  <a>
        <img src="https://img.shields.io/github/license/jack5341/super-dollop?style=flat&labelColor=1C2C2E&color=abc3d6&logoColor=white">
  </a>
  <a>
        <img src="https://img.shields.io/github/stars/jack5341/super-dollop?style=flat&labelColor=1C2C2E&color=abc3d6&logoColor=white">
  </a>
  
</p>

# Super Dollop
**Super Dollop** can encrypt your files and notes by your own [GPG](https://docs.github.com/en/github/authenticating-to-github/connecting-to-github-with-ssh/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent) key and save them in [S3](https://docs.aws.amazon.com/sdk-for-go/api/service/s3/) or [minIO](https://docs.min.io/docs/golang-client-api-reference.html) to keep them safe and portability, also you can use **Super Dollop** for encrypt your file quickly to print it. So with **Super Dollop** you'll solve your keep your notes with security problem easily with Gopher.

![list-command](https://user-images.githubusercontent.com/53150440/132383823-3970f586-1281-4bde-8f9d-6e6463263b48.gif)

```sh
dollop list
```

### Dollop can also do print your encrypted file directly!
With `-p` flag dollop can also print it directly without saving.

![print](https://user-images.githubusercontent.com/53150440/132387521-c3e0d0b5-4d98-4b87-b3dc-fa42644594db.gif)

# Requirements
- [Go](https://golang.org/) `>= 1.16
- [MinIO](https://docs.min.io/docs/minio-quickstart-guide.html)
- Core dependencies: `gnugpg`, `gpgme>=1.7.0`, `libgpg-error`

# Installation

Install [gnupg](https://www.gnupg.org/) 

```sh
sudo apt-get install gnupg
```

Set your environments to your terminal

> .zshrc
```sh
# Environment variables for MinIO
export MINIO_ENDPOINT=127.0.0.1:9000
export MINIO_ACCESS_KEY=admin
export MINIO_SECRET_KEY=secretadmin
export MINIO_GPG_ID=GPG-ID
export MINIO_BUCKET_NAME=dollop-files
```

Get your [MinIO](https://docs.min.io/docs/minio-quickstart-guide.html) container.

> docker-compose.yml
```yaml
version: "3"
services:
  s3:
    image: "minio/minio"
    hostname: "storage"
    restart: "no"
    volumes:
      - data:/data
    ports:
      - "9000:9000"
      - "9001:9001"
    entrypoint: ["minio", "server", "/data","--console-address",":9001"]
    networks:
      - local
volumes:
  data:

networks:
  local:
```

```
// Pull minio/mc
docker pull minio/mc

// Run pulled image with docker-compose.yml file
docker-compose up
```

Give first gas to **Super Dollop**

```sh
// clone the super-dollop repository
git clone https://github.com/jack5341/super-dollop && cd super-dollop

// try to run list command
go run . list
```

# Usage
```sh
dollop [FLAGS] [OPTIONS]
```

```sh
COMMANDS:
    completion  generate the autocompletion script for the specified shell
    dec         List your all encrypted files and notes.
    enc         A brief description of your command
    help        Help about any command
    list        List your all encrypted files and notes.
```
