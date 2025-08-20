sudo docker rm -f $(sudo docker ps -aq) && \
sudo docker rmi -f $(sudo docker images -aq) && \
sudo docker volume rm $(sudo docker volume ls -q) && \
sudo docker network prune -f
