#!/bin/bash


root_dir=/opt/restaurant
net_name=$(ip -o -4 route show to default | awk '{print $5}')
listen_addr=$(ifconfig ${net_name} | grep -E 'inet\W' | grep -o -E [0-9]+.[0-9]+.[0-9]+.[0-9]+ | head -n 1)
if [ -z "${LOG_LEVEL}" ]; then
 export LOG_LEVEL="DEBUG"
fi
writeConfig(){
echo "write template config..."
cat <<EOM > ${root_dir}/conf/chassis.yaml
cse:
  protocols:
    rest:
      listenAddress: ${listen_addr}:30110
  handler:
    chain:
      Provider:
        default: bizkeeper-provider,ratelimiter-provider
EOM
cat <<EOM > ${root_dir}/conf/lager.yaml
logger_level: ${LOG_LEVEL}

logger_file: log/chassis.log

log_format_text: false

rollingPolicy: size

log_rotate_date: 1

log_rotate_size: 10

log_backup_count: 7
EOM

if [ ! -z "$VERSION" ]; then
    /bin/echo "Version ENV: $VERSION"
    cat << EOF > ${root_dir}/conf/microservice.yaml
service_description:
  name: restaurant
  version: $VERSION
EOF
fi
}


echo "prepare config file...."
writeConfig

echo "start service"
/opt/restaurant/main