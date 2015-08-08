test:
	docker run -it -v $(pwd):/home/ trapeze /home/integration/balancerTest

clean:
	docker rm $(docker ps -a -q)
	docker rmi trapeze
