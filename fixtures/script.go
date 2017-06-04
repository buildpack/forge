package fixtures

import "fmt"

func RunScript() string {
	return fmt.Sprintf(runnerScript, forwardScript, "./app")
}

func CommitScript() string {
	return fmt.Sprintf(runnerScript, "", "")
}

func StageSyncScript() string {
	return fmt.Sprintf(stageScript, "", "rsync -a /tmp/app/ /tmp/local/")
}

const forwardScript = `
	echo 'Forwarding: some-name some-other-name'
	sshpass -p 'some-code' ssh -f -N \
		-o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no \
		-o LogLevel=ERROR -o ExitOnForwardFailure=yes \
		-o ServerAliveInterval=10 -o ServerAliveCountMax=60 \
		-p 'some-port' 'some-user@some-ssh-host' \
		-L 'some-from:some-to' \
		-L 'some-other-from:some-other-to'`

const runnerScript = `
	set -e%s
	tar --exclude=%s -C /home/vcap -xzf /tmp/droplet
	chown -R vcap:vcap /home/vcap
	command=$1
	if [[ -z $command ]]; then
		command=$(jq -r .start_command /home/vcap/staging_info.yml)
	fi
	exec /tmp/lifecycle/launcher /home/vcap/app "$command" ''
`

const stageScript = `
	set -e
	chown -R vcap:vcap /tmp/app /tmp/cache
	%ssu vcap -p -c "PATH=$PATH exec /tmp/lifecycle/builder -buildpackOrder $0 -skipDetect=$1"
	%s
`
