import ssh from 'k6/x/ssh';

export default function () {
  ssh.connect({
    username: 'USERNAME',
	  host: "HOST_ADDRESS",
	  port: 22,
    // rsa_key: "PRIVATE_KEY_PATH" ~/.ssh/id_rsa by default
  })
  console.log(ssh.run('pwd'))
  console.log(ssh.run('ls -la'))
}