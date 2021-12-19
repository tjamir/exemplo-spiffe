# exemplo-spiffe


## Comandos

### 1- Iniciar o Spire Server

``` bash

docker-compose -f docker-compose-spireserver.yaml  up

```



### 2- Iniciar o Container do echo-server

``` bash
docker-compose -f docker-compose-echoserver.yaml run echo-server
```

### 3- Gerar um token para o Agent do echo-server

```bash
docker exec exemplo-spiffe_spire-server_1 /opt/spire/bin/spire-server token generate -spiffeID spiffe://example.org/server
```

### 4- Iniciar o agent (dentro do container do echo-server)

```bash
cd /opt/spire
 ./bin/spire-agent run -joinToken [Token gerado no passo 4] &
 ```

### 5- Criar um SVID para o echo-server

```bash
docker exec exemplo-spiffe_spire-server_1 /opt/spire/bin/spire-server entry create \
 -parentID spiffe://example.org/server \
 -spiffeID spiffe://example.org/server/echo-server \
 -selector docker:label:com.docker.compose.service:echo-server
```

### 6- Iniciar o echo-server (dentro do container do echo-server)

```bash
cd /opt/echo-server
 ./bin/echo-server
 ```

### 7- Criar token para os agents dos clientes

```bash
docker exec exemplo-spiffe_spire-server_1 /opt/spire/bin/spire-server token generate -spiffeID spiffe://example.org/client
docker exec exemplo-spiffe_spire-server_1 /opt/spire/bin/spire-server token generate -spiffeID spiffe://example.org/client
```

(guardar os dois tokens)

### 8- Iniciar o container do trusted-echo-client

```bash
docker-compose -f docker-compose-trusted-echoclient.yaml run trusted-client
```

### 9 - Iniciar o agent do trusted-echo-client (dentro do container trusted-echo-client)

```bash
cd /opt/spire
 ./bin/spire-agent run -joinToken [Primeiro token gerado no passo 7] &
 ```

 ### 10 Criar um SVID para o trusted-echo-client

 ```bash
docker exec exemplo-spiffe_spire-server_1 /opt/spire/bin/spire-server entry create \
 -parentID spiffe://example.org/client \
 -spiffeID spiffe://example.org/client/trusted-echo-client \
 -selector docker:label:com.docker.compose.service:trusted-client
```

### 11 Iniciar o trusted-echo-client (dentro do container trusted-echo-client)

```bash
cd /opt/echo-client
./bin/echo-client
```

### 12 Observe a saída

```text
/opt/echo-client # ./bin/echo-client 
2021/12/19 04:41:34 Establecendo conexão tls com o servidor echo-server:8888
DEBU[0019] PID attested to have selectors                pid=20 selectors="[type:\"unix\" value:\"uid:0\" type:\"unix\" value:\"user:root\" type:\"unix\" value:\"gid:0\" type:\"unix\" value:\"group:root\" type:\"unix\" value:\"supplementary_gid:0\" type:\"unix\" value:\"supplementary_group:root\" type:\"unix\" value:\"supplementary_gid:1\" type:\"unix\" value:\"supplementary_group:bin\" type:\"unix\" value:\"supplementary_gid:2\" type:\"unix\" value:\"supplementary_group:daemon\" type:\"unix\" value:\"supplementary_gid:3\" type:\"unix\" value:\"supplementary_group:sys\" type:\"unix\" value:\"supplementary_gid:4\" type:\"unix\" value:\"supplementary_group:adm\" type:\"unix\" value:\"supplementary_gid:6\" type:\"unix\" value:\"supplementary_group:disk\" type:\"unix\" value:\"supplementary_gid:10\" type:\"unix\" value:\"supplementary_group:wheel\" type:\"unix\" value:\"supplementary_gid:11\" type:\"unix\" value:\"supplementary_group:floppy\" type:\"unix\" value:\"supplementary_gid:20\" type:\"unix\" value:\"supplementary_group:dialout\" type:\"unix\" value:\"supplementary_gid:26\" type:\"unix\" value:\"supplementary_group:tape\" type:\"unix\" value:\"supplementary_gid:27\" type:\"unix\" value:\"supplementary_group:video\" type:\"docker\" value:\"label:com.docker.compose.service:trusted-client\" type:\"docker\" value:\"label:com.docker.compose.slug:a99b6daba068d7563e4b1bd71ace5e3256a7160e7735aed6c9194273b1a55da\" type:\"docker\" value:\"label:com.docker.compose.version:1.29.2\" type:\"docker\" value:\"label:com.docker.compose.oneoff:True\" type:\"docker\" value:\"label:com.docker.compose.project:exemplo-spiffe\" type:\"docker\" value:\"label:com.docker.compose.project.config_files:docker-compose-trusted-echoclient.yaml\" type:\"docker\" value:\"label:com.docker.compose.project.working_dir:/home/tjs/workspace/exemplo-spiffe\" type:\"docker\" value:\"env:PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\" type:\"docker\" value:\"env:APP_NAME_ENV=echo-client\" type:\"docker\" value:\"image_id:exemplo-spiffe_trusted-client\"]" subsystem_name=workload_attestor
DEBU[0019] Fetched X.509 SVID                            count=1 method=FetchX509SVID pid=20 registered=true service=WorkloadAPI spiffe_id="spiffe://example.org/client/trusted-echo-client" subsystem_name=endpoints ttl=3580.611556278
2021/12/19 04:41:34 Conexão estabelecida, enviando mensagens echo-server:8888
Digite uma mensagem:
olá
2021/12/19 04:41:35 Mensagem enviada, aguardando resposta
Resposta recebida: olá
```

### 13- Iniciar o container do untrusted-echo-client

```bash
docker-compose -f docker-compose-untrusted-echoclient.yaml run untrusted-client
```

### 14 - Iniciar o agent do untrusted-echo-client (dentro do container trusted-echo-client)

```bash
cd /opt/spire
 ./bin/spire-agent run -joinToken [Segundo token gerado no passo 7] &
 ```

### 15 Criar um SVID para o untrusted-echo-client

 ```bash
docker exec exemplo-spiffe_spire-server_1 /opt/spire/bin/spire-server entry create \
 -parentID spiffe://example.org/client \
 -spiffeID spiffe://example.org/client/untrusted-echo-client \
 -selector docker:label:com.docker.compose.service:untrusted-client
```

### 16 Iniciar o untrusted-echo-client (dentro do container trusted-echo-client)

```bash
cd /opt/echo-client
./bin/echo-client
```

### 17 Observe a saída

```text
2021/12/19 04:52:52 Establecendo conexão tls com o servidor echo-server:8888
DEBU[0013] PID attested to have selectors                pid=22 selectors="[type:\"unix\" value:\"uid:0\" type:\"unix\" value:\"user:root\" type:\"unix\" value:\"gid:0\" type:\"unix\" value:\"group:root\" type:\"unix\" value:\"supplementary_gid:0\" type:\"unix\" value:\"supplementary_group:root\" type:\"unix\" value:\"supplementary_gid:1\" type:\"unix\" value:\"supplementary_group:bin\" type:\"unix\" value:\"supplementary_gid:2\" type:\"unix\" value:\"supplementary_group:daemon\" type:\"unix\" value:\"supplementary_gid:3\" type:\"unix\" value:\"supplementary_group:sys\" type:\"unix\" value:\"supplementary_gid:4\" type:\"unix\" value:\"supplementary_group:adm\" type:\"unix\" value:\"supplementary_gid:6\" type:\"unix\" value:\"supplementary_group:disk\" type:\"unix\" value:\"supplementary_gid:10\" type:\"unix\" value:\"supplementary_group:wheel\" type:\"unix\" value:\"supplementary_gid:11\" type:\"unix\" value:\"supplementary_group:floppy\" type:\"unix\" value:\"supplementary_gid:20\" type:\"unix\" value:\"supplementary_group:dialout\" type:\"unix\" value:\"supplementary_gid:26\" type:\"unix\" value:\"supplementary_group:tape\" type:\"unix\" value:\"supplementary_gid:27\" type:\"unix\" value:\"supplementary_group:video\" type:\"docker\" value:\"label:com.docker.compose.project:exemplo-spiffe\" type:\"docker\" value:\"label:com.docker.compose.project.config_files:docker-compose-untrusted-echoclient.yaml\" type:\"docker\" value:\"label:com.docker.compose.project.working_dir:/home/tjs/workspace/exemplo-spiffe\" type:\"docker\" value:\"label:com.docker.compose.service:untrusted-client\" type:\"docker\" value:\"label:com.docker.compose.slug:33abc8e8ed57db7812746d3cfa52ce272e914ed78ed6136cde26d86dabe662\" type:\"docker\" value:\"label:com.docker.compose.version:1.29.2\" type:\"docker\" value:\"label:com.docker.compose.oneoff:True\" type:\"docker\" value:\"env:PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\" type:\"docker\" value:\"env:APP_NAME_ENV=echo-client\" type:\"docker\" value:\"image_id:exemplo-spiffe_untrusted-client\"]" subsystem_name=workload_attestor
DEBU[0013] Fetched X.509 SVID                            count=1 method=FetchX509SVID pid=22 registered=true service=WorkloadAPI spiffe_id="spiffe://example.org/client/untrusted-echo-client" subsystem_name=endpoints ttl=3586.6794877
2021/12/19 04:52:52 Conexão estabelecida, enviando mensagens echo-server:8888
Digite uma mensagem:
olá
2021/12/19 04:53:04 Mensagem enviada, aguardando resposta
2021/12/19 04:53:04 Erro enviando mensagem remote error: tls: bad certificate
```
