#!/usr/bin/env bash
#run.sh -a AWS_ACCOUNT_ID -k AWS_ACCESS_KEY_ID -s AWS_SECRET_ACCESS_KEY -r AWS_DEFAULT_REGION  -i APP_ID -p PORT

POSITIONAL=()
while [[ $# -gt 0 ]]; do
  key="$1"

  case $key in
    -a|--account_id)
      AWS_ACCOUNT_ID="$2"
      shift 
      shift 
      ;;
    -k|--access_key_id)
      AWS_ACCESS_KEY_ID="$2"
      shift 
      shift 
      ;;
    -s|--secret_access_key)
      AWS_SECRET_ACCESS_KEY="$2"
      shift 
      shift 
      ;;
    -r|--default_region)
      AWS_DEFAULT_REGION="$2"
      shift 
      shift 
      ;;
    -i|--app_id)
      APP_ID="$2"
      shift 
      shift 
      ;;
    -p|--port)
      PORT="$2"
      shift 
      shift 
      ;;
    *)
      POSITIONAL+=("$1")
      shift 
      ;;
  esac
done

set -- "${POSITIONAL[@]}" # restore positional parameters

docker_login () {
  echo "docker login"
  export AWS_ACCESS_KEY_ID
  export AWS_SECRET_ACCESS_KEY
  export AWS_DEFAULT_REGION
  aws ecr get-login-password --region ap-northeast-2 | docker login --username AWS --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.ap-northeast-2.amazonaws.com
}

docker_stop () {
  echo "docker stop"
  docker ps|grep "${APP_ID}"|awk '{ if ( $1 != "" ) system("docker stop " $1) }'
}

docker_prune () {
  echo "docker prune"
  docker system prune -af --volumes
}

docker_run () {
  echo "docker run"
  docker run --rm -d \
  --name "${APP_ID}" \
  -p "${PORT}:${PORT}" \
  "${AWS_ACCOUNT_ID}".dkr.ecr.ap-northeast-2.amazonaws.com/"${APP_ID}":latest
}

main() {
  docker_login
  docker_stop
  docker_prune
  docker_run
}

main "$@"
