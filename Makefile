user-rpc-dev:
	@make -f deploy/make/user_rpc.mk release-test

user-api-dev:
	@make -f deploy/make/user_api.mk release-test

social-rpc-dev:
	@make -f deploy/make/social_rpc.mk release-test

social-api-dev:
	@make -f deploy/make/social_api.mk release-test


release-test: user-rpc-dev user-api-dev social-rpc-dev social-api-dev


install-server:
	cd ./deploy/script && chmod +x release-test.sh && ./release-test.sh