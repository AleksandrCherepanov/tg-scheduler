build-image:
	ansible-playbook build.yaml -e @secrets.yaml --ask-vault-pass
deploy-nginx: 
	ansible-playbook nginx.yaml --ask-become-pass
deploy-app:
	ansible-playbook deploy.yaml -e @secrets.yaml --ask-vault-pass

