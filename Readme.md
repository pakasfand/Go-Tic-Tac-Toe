To update main.wasm run the following in Powershell:
$env:GOOS="js"
$env:GOARCH="wasm"
cd client
go build -o main.wasm .

To build new docker image:
docker build --no-cache -t pakasfand/go-tic-tac-toe:latest .
docker push pakasfand/go-tic-tac-toe:latest

To deploy build and pull latest docker image:
sudo docker pull pakasfand/go-tic-tac-toe:latest

# Quick Update Workflow

# Local: Build and push
docker build --no-cache -t pakasfand/go-tic-tac-toe:latest .
docker push pakasfand/go-tic-tac-toe:latest

# Lightsail: Pull and deploy  
sudo docker pull pakasfand/go-tic-tac-toe:latest
sudo docker stop go-tic-tac-toe 2>/dev/null || true
sudo docker rm go-tic-tac-toe 2>/dev/null || true
sudo docker run -d --name go-tic-tac-toe -p 80:80 -p 8080:8080 pakasfand/go-tic-tac-toe:latest
