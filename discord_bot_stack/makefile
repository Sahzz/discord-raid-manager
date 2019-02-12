environment=discord-bot

.PHONY: create-stack update-stack delete-stack

create-stack:
	aws cloudformation create-stack --capabilities CAPABILITY_IAM CAPABILITY_NAMED_IAM --stack-name ${environment} --template-body file://stack.yaml
update-stack:
	aws cloudformation update-stack --capabilities CAPABILITY_IAM CAPABILITY_NAMED_IAM --stack-name ${environment} --template-body file://stack.yaml
delete-stack:
	aws cloudformation delete-stack --stack-name ${environment}
recreate-stack: delete-stack create-stack
