# Go Distributed Store
[![Go Report Card](https://goreportcard.com/badge/github.com/Dhruv-mak/godiststore)](https://goreportcard.com/report/github.com/Dhruv-mak/godiststore)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A high-performance, secure, and scalable peer-to-peer distributed file storage system built in Go. This project implements a content-addressable storage (CAS) architecture with built-in encryption, making it suitable for sensitive data storage across distributed networks.

## 🚀 Key Features

- **Distributed P2P Architecture**: Seamlessly distribute and retrieve files across multiple nodes in a peer-to-peer network
- **Content-Addressable Storage**: Implement efficient file deduplication and retrieval using SHA-1 content hashing
- **Military-Grade Security**: AES-256 encryption in CTR mode for secure file transfer and storage
- **High Performance**: Optimized TCP transport layer with custom handshake mechanisms
- **Horizontal Scalability**: Easily add new nodes to increase storage capacity and redundancy
- **Flexible Configuration**: Customizable path transformation and storage strategies

## 🛠️ Technical Implementation

### Architecture Overview
The system is built on a modular architecture with several key components:

```
├── crypto/      # Encryption and hashing implementations
├── p2p/         # Network transport and peer management
├── store/       # Core storage logic and CAS implementation
└── server/      # API and service coordination
```

### Key Components

#### Cryptographic Layer
- AES-256 encryption in CTR mode for file security
- SHA-1 based content addressing
- Secure random ID generation for unique file identification

#### P2P Network
- Custom TCP transport implementation
- Robust peer discovery and handshake mechanism
- Message encoding with GOB for efficient data transfer

#### Storage Engine
- Content-addressable storage (CAS) with customizable path transformation
- Concurrent read/write operations
- Automatic file deduplication

## 🔧 Installation

```bash
# Clone the repository
git clone https://github.com/username/go-distributed-store

# Build the project
make build

# Run tests
make test
```

## 📝 Usage Example

```go
// Initialize a new storage node
server := NewFileServer(FileServerOpts{
    ListenAddr: ":3000",
    StorageRoot: "./data",
})

// Start the server
server.Start()

// Store a file with encryption
server.Store("myfile.txt", data, true)

// Retrieve a file
data, err := server.Get("myfile.txt")
```


## 🔍 Technical Challenges Solved

1. **Distributed Consensus**: Implemented a custom protocol for maintaining consistency across nodes
2. **Security**: Developed a robust encryption system that maintains performance
3. **Network Resilience**: Built-in retry mechanisms and fault tolerance
4. **Data Integrity**: SHA-1 verification ensures file consistency

## 🛣️ Roadmap

- [ ] Implementation of Reed-Solomon error correction
- [ ] DHT-based peer discovery
- [ ] Blockchain-based file tracking
- [ ] Multi-region support

## 🤝 Contributing

Contributions are welcome! Make a Pull Request with clear Description on what functionality it adds. And I will merge it.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.