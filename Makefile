compose:
	@echo "Starting Family Tree with Docker Compose"
	docker-compose up

docker:
	@echo "Starting Family Tree with Docker"
	cd backend; docker build -t family-tree .
	docker run --publish=7474:7474 --publish=7687:7687 --volume=$$HOME/neo4j/data:/data --env=NEO4J_AUTH=neo4j/1234 --rm -d neo4j
	docker run -p 8989:8989 --net=host --rm -d family-tree