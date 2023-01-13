build-image:
	ansible-playbook app/build.yaml -e @secrets.yaml --ask-vault-pass
deploy-nginx: 
	ansible-playbook nginx.yaml --ask-become-pass
deploy-app:
	ansible-playbook app/deploy.yaml -e @secrets.yaml --ask-vault-pass
deploy-redis:
	ansible-playbook redis.yaml
deploy-postgres:
	ansible-playbook postgresql.yaml -e @secrets.yaml --ask-vault-pass
deploy-migrations:
	ansible-playbook migrations.yaml -e @secrets.yaml --ask-vault-pass

.PHONY: build-image deploy-nginx deploy-app deploy-redis deploy-postgres deploy-migrations
