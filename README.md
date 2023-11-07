![GitHub go.mod Go version (branch & subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/wreckitkenny/vngitsub)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-v1.2x-blue)](https://kubernetes.io/)
# vngitsub
vngitSub - A VNGITBOTV2 subscriber

#
```bash
vngitsub
├── CHANGELOG
├── Dockerfile
├── go.mod
├── go.sum
├── LICENSE
├── main.go
├── model
│   ├── message.go
│   └── publish.go
├── pkg
│   ├── controller
│   │   ├── changeImage.go
│   │   ├── gitlab.go
│   │   ├── mongodb.go
│   │   ├── rabbitmq.go
│   │   ├── request.go
│   │   └── telegram.go
│   ├── ping.go
│   └── utils
│       ├── imageChange.go
│       └── logger.go
├── README.md
└── VERSION

4 directories, 19 files
```
#
## Requirements
```bash
go==1.18 or newer
```
#
## Usage
```golang
go run main.go
```
#
## License
[MIT](https://choosealicense.com/licenses/mit/)
