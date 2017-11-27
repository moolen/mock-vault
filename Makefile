.PHONY: default run build
build: clean
	docker build -t moolen/mock-vault .
run: build
	docker run --rm --name mock-vault-local -d -e VAULT_PORT=3000 -e VAULT_PATH=/data -p 3000:3000 moolen/mock-vault
clean:
	rm -f ./mock-vault
	docker rm -f mock-vault-local
publish: build
	docker push moolen/mock-vault
test: run
	sh test.sh