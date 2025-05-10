# xk6-ssh
A k6 extension for using of SSH in testing. Built for [k6](https://github.com/grafana/k6) using [xk6](https://github.com/grafana/xk6).

## Build

To build a `k6` binary with this extension, first ensure you have the prerequisites:

- [Go toolchain](https://go101.org/article/go-toolchain.html)
- Git

Then:

1. Download `xk6`:
  ```bash
  go install github.com/grafana/xk6/cmd/xk6@latest
  ```

2. Build the binary:
  ```bash
  xk6 build --with github.com/grafana/xk6-ssh@latest
  ```

This will result in a `k6` binary in the current directory.

## Example

```javascript
import ssh from 'k6/x/ssh';

export default function () {
  ssh.connect({
    username: `${__ENV.K6_USERNAME}`,
    password: `${__ENV.K6_PASSWORD}`,
    host: [HOSTNAME],
	port: 22
  })
  console.log(ssh.run('pwd'))
}
```

Result output:

```plain
$ ./k6 run script.js

          /\      |‾‾| /‾‾/   /‾‾/   
     /\  /  \     |  |/  /   /  /    
    /  \/    \    |     (   /   ‾‾\  
   /          \   |  |\  \ |  (‾)  | 
  / __________ \  |__| \__\ \_____/ .io

  execution: local
     script: ../xk6-ssh/script.js
     output: -

  scenarios: (100.00%) 1 scenario, 1 max VUs, 10m30s max duration (incl. graceful stop):
           * default: 1 iterations for each of 1 VUs (maxDuration: 10m0s, gracefulStop: 30s)

INFO[0001] /home/userfolder                                 source=console

running (00m01.4s), 0/1 VUs, 1 complete and 0 interrupted iterations
default ✓ [======================================] 1 VUs  00m01.4s/10m0s  1/1 iters, 1 per VU

     data_received........: 0 B 0 B/s
     data_sent............: 0 B 0 B/s
     iteration_duration...: avg=1.41s min=1.41s med=1.41s max=1.41s p(90)=1.41s p(95)=1.41s
     iterations...........: 1   0.706079/s
     vus..................: 1   min=1 max=1
     vus_max..............: 1   min=1 max=1

```

## Testing Locally
This repo includes a [docker-compose.yml](docker-compose.yml) file that starts an [OpenSSH Server](https://docs.linuxserver.io/images/docker-openssh-server) from [LinuxServer.io](https://www.linuxserver.io/).
The `examples` directory contains scripts that are configured to work with this environment out of the box.

> :warning: Be sure that you've already compiled your custom `k6` binary as described in the [Build](#build) section!

We'll use this environment to run some examples.

1. Start the docker compose environment.

   ```shell
   docker compose up -d
   ```
   Once you see the following, you should be ready.
   ```shell
   [+] Running 2/2
    ⠿ Network xk6-ssh_default             Created
    ⠿ Container xk6-ssh-openssh-server-1  Started
   ```
   Next, we'll use the `k6` binary we compiled in the [Build section](#build) above.

1. Using our custom `k6` binary, we can execute our [example scripts](examples/).
   ```shell
   ./k6 run examples/connect-by-rsa-key.js
   ``` 
   The RSA example will then connect to the local SSH server using the `example_rsa` private key.

## FAQ

### How to start `sudo` commands?

Basically we don't provide sudo password autofill. We suggest to use `/etc/sudoers` for this purpose. Please checkout this [article](https://www.cyberciti.biz/faq/linux-unix-running-sudo-command-without-a-password/) for more details.
