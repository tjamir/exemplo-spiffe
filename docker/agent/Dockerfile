FROM spire-agent/latest-local

ADD start_agent.sh /tmp/start_agent.sh
ADD token_env.sh /tmp/token_env.sh

ENTRYPOINT ["/usr/bin/dumb-init", "/opt/spire/bin/spire-agent", "run", "-joinToken"]