import ssh from 'k6/x/ssh';

export default function () {
  ssh.connect({
    username: 'YOUR_USERNAME',
	  host: "YOUR_HOST",
    sudo_password: "YOUR_SUDO_PASSWORD",
	  port: 22
  })
  console.log(ssh.run(['sudo pwd'], { sudo: true }))
}