> ### ⚠️ This is a proof of concept
>
> As this is a proof of concept,  it won't be supported by the k6 team.
> It may also break in the future as xk6 evolves. USE AT YOUR OWN RISK!
> Any issues with the tool should be raised [here](https://github.com/k6io/xk6-ssh/issues).

</br>
</br>

<div align="center">

# xk6-ssh
A k6 extension for using of SSH in testing. Built for [k6](https://github.com/loadimpact/k6) using [xk6](https://github.com/k6io/xk6).

</div>

## Build

To build a `k6` binary with this extension, first ensure you have the prerequisites:

- [Go toolchain](https://go101.org/article/go-toolchain.html)
- Git

Then:

1. Download `xk6`:
  ```bash
  $ go get -u github.com/k6io/xk6
  ```

2. Build the binary:
  ```bash
  $ xk6 build --with github.com/k6io/xk6-ssh
  ```

## Example

```javascript
import ssh from 'k6/x/ssh';

export default function () {
  ssh.connect({
    username: 'USERNAME',
	  host: "HOST_ADDRESS",
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

Inspect examples folder for more details.