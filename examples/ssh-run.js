import ssh from 'k6/x/ssh';

export default function () {
  ssh.connect({
    username: 'helphub',
	  host: "18.198.123.166",
	  port: 22,
    // rsa_key: "/Users/lxkuz/.ssh/id_rsa"
  })
  console.log(ssh.run('pwd'))
  console.log(ssh.run('ls -la'))
}