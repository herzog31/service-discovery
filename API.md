**Request**
```
GET /api/containers
```

**Example**
```json
[{
    "Id": "fbc48f38ca8bb95ac9ab2e1763f189c7d6b4ae7d01530b9cf237d5f5ac6dfd34",
    "Image": "production_node",
    "Command": "npm start",
    "Created": 1440679438,
    "Status": "Up 16 minutes",
    "Ports": [{
        "PrivatePort": 80,
        "PublicPort": 80,
        "Type": "tcp",
        "IP": "0.0.0.0"
    }],
    "Names": ["/production_node_1"],
    "Labels": {
        "com.docker.compose.config-hash": "12c986558113a024735d8a9caf2b5fa6b21a2cd92beadb71bcb1f5d84c4519ed",
        "com.docker.compose.container-number": "1",
        "com.docker.compose.oneoff": "False",
        "com.docker.compose.project": "production",
        "com.docker.compose.service": "node",
        "com.docker.compose.version": "1.4.0"
    }
}, {
    "Id": "21dcf1f9167d9dc1fc616c13509a6442fa28fe3ef7fc5a9590666d18c4b21b71",
    "Image": "mongo",
    "Command": "/entrypoint.sh mongod",
    "Created": 1440679427,
    "Status": "Up 16 minutes",
    "Ports": [{
        "PrivatePort": 27017,
        "Type": "tcp"
    }],
    "Names": ["/production_mongo_1", "/production_node_1/db", "/production_node_1/mongo_1", "/production_node_1/production_mongo_1"],
    "Labels": {
        "com.docker.compose.config-hash": "02137310822926a39cea3dfb151661b3b7d5e0d78ce2c5ccc515d7646bc57f7c",
        "com.docker.compose.container-number": "1",
        "com.docker.compose.oneoff": "False",
        "com.docker.compose.project": "production",
        "com.docker.compose.service": "mongo",
        "com.docker.compose.version": "1.4.0"
    }
}]
```
***

**Request**
```
GET /api/containersFull
```

**Example**
```json
[{
    "Id": "fbc48f38ca8bb95ac9ab2e1763f189c7d6b4ae7d01530b9cf237d5f5ac6dfd34",
    "Created": "2015-08-27T12:43:58.19794574Z",
    "Path": "npm",
    "Args": ["start"],
    "Config": {
        "Hostname": "fbc48f38ca8b",
        "ExposedPorts": {
            "80/tcp": {}
        },
        "Env": ["PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin", "NODE_VERSION=0.12.7", "NPM_VERSION=2.11.3"],
        "Cmd": ["npm", "start"],
        "Image": "production_node",
        "WorkingDir": "/usr/src/app",
        "Entrypoint": null,
        "Labels": {
            "com.docker.compose.config-hash": "12c986558113a024735d8a9caf2b5fa6b21a2cd92beadb71bcb1f5d84c4519ed",
            "com.docker.compose.container-number": "1",
            "com.docker.compose.oneoff": "False",
            "com.docker.compose.project": "production",
            "com.docker.compose.service": "node",
            "com.docker.compose.version": "1.4.0"
        }
    },
    "State": {
        "Running": true,
        "Pid": 2558,
        "StartedAt": "2015-10-06T12:02:29.031779213Z",
        "FinishedAt": "2015-09-22T15:07:01.355938493Z"
    },
    "Image": "ddd89f2732a6a27e68ce37c3a77b8fb27dcda42f29443dee5b4d23a05717d730",
    "NetworkSettings": {
        "IPAddress": "172.17.0.4",
        "IPPrefixLen": 16,
        "MacAddress": "02:42:ac:11:00:04",
        "Gateway": "172.17.42.1",
        "Ports": {
            "80/tcp": [{
                "HostIP": "0.0.0.0",
                "HostPort": "80"
            }]
        },
        "NetworkID": "aba8fe8d8fa1aa44880e3d49aabba71a2a81afc4cc2fc2283ac7a32bc34b4ac2",
        "EndpointID": "2364d0209debebdeb6ef4de0b768f6b4dcdcf4741011b785bde8e46dfae105ae",
        "SandboxKey": "/var/run/docker/netns/fbc48f38ca8b"
    },
    "ResolvConfPath": "/var/lib/docker/containers/fbc48f38ca8bb95ac9ab2e1763f189c7d6b4ae7d01530b9cf237d5f5ac6dfd34/resolv.conf",
    "HostnamePath": "/var/lib/docker/containers/fbc48f38ca8bb95ac9ab2e1763f189c7d6b4ae7d01530b9cf237d5f5ac6dfd34/hostname",
    "HostsPath": "/var/lib/docker/containers/fbc48f38ca8bb95ac9ab2e1763f189c7d6b4ae7d01530b9cf237d5f5ac6dfd34/hosts",
    "LogPath": "/var/lib/docker/containers/fbc48f38ca8bb95ac9ab2e1763f189c7d6b4ae7d01530b9cf237d5f5ac6dfd34/fbc48f38ca8bb95ac9ab2e1763f189c7d6b4ae7d01530b9cf237d5f5ac6dfd34-json.log",
    "Name": "node",
    "Driver": "aufs",
    "HostConfig": {
        "PortBindings": {
            "80/tcp": [{
                "HostPort": "80"
            }]
        },
        "Links": ["/production_mongo_1:/production_node_1/db", "/production_mongo_1:/production_node_1/mongo_1", "/production_mongo_1:/production_node_1/production_mongo_1"],
        "NetworkMode": "default",
        "RestartPolicy": {},
        "LogConfig": {
            "Type": "json-file"
        }
    },
    "FullName": "production_node_1",
    "Project": "production",
    "Number": 1
}, {
    "Id": "21dcf1f9167d9dc1fc616c13509a6442fa28fe3ef7fc5a9590666d18c4b21b71",
    "Created": "2015-08-27T12:43:47.093764903Z",
    "Path": "/entrypoint.sh",
    "Args": ["mongod"],
    "Config": {
        "Hostname": "21dcf1f9167d",
        "ExposedPorts": {
            "27017/tcp": {}
        },
        "Env": ["affinity:container==02b4ada668c57b855dd32cf420ecb7dc2a1125137171f789a2d1caf89fc5d62c", "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin", "MONGO_MAJOR=3.0", "MONGO_VERSION=3.0.6"],
        "Cmd": ["mongod"],
        "Image": "mongo",
        "Volumes": {
            "/data/db": {}
        },
        "Entrypoint": ["/entrypoint.sh"],
        "Labels": {
            "com.docker.compose.config-hash": "02137310822926a39cea3dfb151661b3b7d5e0d78ce2c5ccc515d7646bc57f7c",
            "com.docker.compose.container-number": "1",
            "com.docker.compose.oneoff": "False",
            "com.docker.compose.project": "production",
            "com.docker.compose.service": "mongo",
            "com.docker.compose.version": "1.4.0"
        }
    },
    "State": {
        "Running": true,
        "Pid": 2437,
        "StartedAt": "2015-10-06T12:02:25.354508566Z",
        "FinishedAt": "2015-09-22T15:06:58.928111527Z"
    },
    "Image": "5c9464760d54612edf1df762d13207117aa4480b2174d9c23962c44afaa4d808",
    "NetworkSettings": {
        "IPAddress": "172.17.0.2",
        "IPPrefixLen": 16,
        "MacAddress": "02:42:ac:11:00:02",
        "Gateway": "172.17.42.1",
        "Ports": {
            "27017/tcp": null
        },
        "NetworkID": "aba8fe8d8fa1aa44880e3d49aabba71a2a81afc4cc2fc2283ac7a32bc34b4ac2",
        "EndpointID": "20b64787abdf1ccb7c46dcb6a17bb20dcdbb9888e4281072c9c15ec3dc25424c",
        "SandboxKey": "/var/run/docker/netns/21dcf1f9167d"
    },
    "ResolvConfPath": "/var/lib/docker/containers/21dcf1f9167d9dc1fc616c13509a6442fa28fe3ef7fc5a9590666d18c4b21b71/resolv.conf",
    "HostnamePath": "/var/lib/docker/containers/21dcf1f9167d9dc1fc616c13509a6442fa28fe3ef7fc5a9590666d18c4b21b71/hostname",
    "HostsPath": "/var/lib/docker/containers/21dcf1f9167d9dc1fc616c13509a6442fa28fe3ef7fc5a9590666d18c4b21b71/hosts",
    "LogPath": "/var/lib/docker/containers/21dcf1f9167d9dc1fc616c13509a6442fa28fe3ef7fc5a9590666d18c4b21b71/21dcf1f9167d9dc1fc616c13509a6442fa28fe3ef7fc5a9590666d18c4b21b71-json.log",
    "Name": "mongo",
    "Driver": "aufs",
    "Mounts": [{
        "Source": "/var/lib/docker/volumes/929e8e87e3572825562cb7e10942105fd8c46a9724cf22bb35c42bda515cdcb4/_data",
        "Destination": "/data/db",
        "Mode": "rw",
        "RW": true
    }],
    "HostConfig": {
        "Binds": ["/var/lib/docker/volumes/929e8e87e3572825562cb7e10942105fd8c46a9724cf22bb35c42bda515cdcb4/_data:/data/db:rw"],
        "NetworkMode": "default",
        "RestartPolicy": {},
        "LogConfig": {
            "Type": "json-file"
        }
    },
    "FullName": "production_mongo_1",
    "Project": "production",
    "Number": 1
}]
```
***

**Request**
```
GET /api/container/:name
```

**Example**
```
GET /api/container/production_node_1
```
```json
{
    "Id": "fbc48f38ca8bb95ac9ab2e1763f189c7d6b4ae7d01530b9cf237d5f5ac6dfd34",
    "Created": "2015-08-27T12:43:58.19794574Z",
    "Path": "npm",
    "Args": ["start"],
    "Config": {
        "Hostname": "fbc48f38ca8b",
        "ExposedPorts": {
            "80/tcp": {}
        },
        "Env": ["PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin", "NODE_VERSION=0.12.7", "NPM_VERSION=2.11.3"],
        "Cmd": ["npm", "start"],
        "Image": "production_node",
        "WorkingDir": "/usr/src/app",
        "Entrypoint": null,
        "Labels": {
            "com.docker.compose.config-hash": "12c986558113a024735d8a9caf2b5fa6b21a2cd92beadb71bcb1f5d84c4519ed",
            "com.docker.compose.container-number": "1",
            "com.docker.compose.oneoff": "False",
            "com.docker.compose.project": "production",
            "com.docker.compose.service": "node",
            "com.docker.compose.version": "1.4.0"
        }
    },
    "State": {
        "Running": true,
        "Pid": 2558,
        "StartedAt": "2015-10-06T12:02:29.031779213Z",
        "FinishedAt": "2015-09-22T15:07:01.355938493Z"
    },
    "Image": "ddd89f2732a6a27e68ce37c3a77b8fb27dcda42f29443dee5b4d23a05717d730",
    "NetworkSettings": {
        "IPAddress": "172.17.0.4",
        "IPPrefixLen": 16,
        "MacAddress": "02:42:ac:11:00:04",
        "Gateway": "172.17.42.1",
        "Ports": {
            "80/tcp": [{
                "HostIP": "0.0.0.0",
                "HostPort": "80"
            }]
        },
        "NetworkID": "aba8fe8d8fa1aa44880e3d49aabba71a2a81afc4cc2fc2283ac7a32bc34b4ac2",
        "EndpointID": "2364d0209debebdeb6ef4de0b768f6b4dcdcf4741011b785bde8e46dfae105ae",
        "SandboxKey": "/var/run/docker/netns/fbc48f38ca8b"
    },
    "ResolvConfPath": "/var/lib/docker/containers/fbc48f38ca8bb95ac9ab2e1763f189c7d6b4ae7d01530b9cf237d5f5ac6dfd34/resolv.conf",
    "HostnamePath": "/var/lib/docker/containers/fbc48f38ca8bb95ac9ab2e1763f189c7d6b4ae7d01530b9cf237d5f5ac6dfd34/hostname",
    "HostsPath": "/var/lib/docker/containers/fbc48f38ca8bb95ac9ab2e1763f189c7d6b4ae7d01530b9cf237d5f5ac6dfd34/hosts",
    "LogPath": "/var/lib/docker/containers/fbc48f38ca8bb95ac9ab2e1763f189c7d6b4ae7d01530b9cf237d5f5ac6dfd34/fbc48f38ca8bb95ac9ab2e1763f189c7d6b4ae7d01530b9cf237d5f5ac6dfd34-json.log",
    "Name": "node",
    "Driver": "aufs",
    "HostConfig": {
        "PortBindings": {
            "80/tcp": [{
                "HostPort": "80"
            }]
        },
        "Links": ["/production_mongo_1:/production_node_1/db", "/production_mongo_1:/production_node_1/mongo_1", "/production_mongo_1:/production_node_1/production_mongo_1"],
        "NetworkMode": "default",
        "RestartPolicy": {},
        "LogConfig": {
            "Type": "json-file"
        }
    },
    "FullName": "production_node_1",
    "Project": "production",
    "Number": 1
}
```
***

**Request**
```
GET /api/container/:name/logs
```

**Example**
```
GET /api/container/production_node_1/logs
```
```
2015-08-27T12:43:58.876618867Z > nodeDemo@1.0.0 start /usr/src/app
2015-08-27T12:43:58.876625572Z > node app.js
2015-08-27T12:43:58.876629763Z 
2015-08-27T12:43:59.045902673Z Example app listening at http://:::80
```
***

**Request**
```
GET /api/container/:name/mappings
```

**Example**
```
GET /api/container/production_node_1/mappings
```
```json
[{
    "Container": {
        "Port": 80,
        "Protocol": "tcp"
    },
    "Host": {
        "Port": 80,
        "Protocol": "tcp"
    },
    "Hostname": "192.168.178.27"
}]
```
***

**Request**
```
GET /api/containers/:name/mapping/:port
```

**Parameter**
- text — Converts JSON response to plaintext

**Example**
```
GET /api/container/production_node_1/mapping/80
```
```json
{
    "Container": {
        "Port": 80,
        "Protocol": "tcp"
    },
    "Host": {
        "Port": 80,
        "Protocol": "tcp"
    },
    "Hostname": "192.168.178.27"
}
```

```
GET /api/container/production_node_1/mapping/80?text
```
```
192.168.178.27:80
```
***

**Request**
```
GET /api/containers/:name/mapping/:port/:protocol
```

**Parameter**
- text — Converts JSON response to plaintext

**Example**
```
GET /api/container/production_node_1/mapping/80/tcp
```
```json
{
    "Container": {
        "Port": 80,
        "Protocol": "tcp"
    },
    "Host": {
        "Port": 80,
        "Protocol": "tcp"
    },
    "Hostname": "192.168.178.27"
}
```
```
GET /api/container/production_node_1/mapping/80?text
```
```
192.168.178.27:80
```
***

**Request**
```
GET /api/projectUp/:project
```

**Example**
```
GET /api/projectUp/production
```
```
true
```
***