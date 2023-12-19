# makefile variables
reponame = $(shell basename `git rev-parse --show-toplevel`)
branch = $(shell git rev-parse --abbrev-ref HEAD)
service_name=product-es-migration

.PHONY: help
help: # print all available make commands and their usages
	@printf "\e[32m\nMakefile help\n\e[0m"	
	@printf "\e[32mUsage: \e[0m\e[36m make <command>\e[0m\n\n"
	@grep -E '^[a-zA-Z_-]+:.*?# .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?# "}; {printf "\033[36m%-30s\033[0m \033[32m%s\033[0m\n", $$1, $$2}'
	@printf "\e[32m\nMore info: https://jamtangan.atlassian.net/l/cp/8G0KWMM1\e[0m\n\n"	

.PHONY: setup
setup: # set up the necessary dependencies and configurations required for development
	@echo "\n🔎 Setting up repo requirements...\n"
	@.dev/setup.sh

.PHONY: env-local
env-local: # generate local .env files from vault 
	@echo "Generating local .env file.."
	@.dev/generate-env.sh $(GITHUB_TOKEN) $(service_name) local
	
.PHONY: env-dev
env-dev: # generate development .env files from vault 
	@echo "Generating development .env file.."
	@.dev/generate-env.sh $(GITHUB_TOKEN) $(service_name) development

.PHONY: env-staging
env-staging: # generate staging .env files from vault 
	@echo "Generating staging .env file.."
	@.dev/generate-env.sh $(GITHUB_TOKEN) $(service_name) staging

.PHONY: create-index
create-index: # create index base on environment ELASTICSEARCH_PRODUCT_INDEX
	go run main.go create-index

.PHONY: delete-index
delete-index: # delete index base on environment ELASTICSEARCH_PRODUCT_INDEX
	go run main.go delete-index

.PHONY: check-index
check-index: # check index base on environment ELASTICSEARCH_PRODUCT_INDEX is exist or not
	go run main.go check-index

.PHONY: migrate
migrate: # migrate product from postgre db to elastic db
	go run main.go migrate
