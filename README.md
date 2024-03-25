# cms_template_go_v2
 A template for creating a CMS from general Go + React code (Post Nextsense and ODG)

## Setting Up a new CMS
1. Create a new repo on Github using this repo as a template
2. Find and replace all instances of "cms_template_go_v2" with the name of this repo
3. Run `go mod init`, to initalise this repo as a go module
4. Copy `ip_whitelist_example` -> `ip_whitelist` (check if it needs modification)
5. Copy `config_example.json` -> `config.json` AND `config-dev.json`
	- You'll likely need to change "bucket" and "region" properties
	- Other properties ABOVE "special_types" may also need to be changed
6. Create a `.env` file, copy the section below to it:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=postgres
DB_NAME=postgres
DB_SSLMODE=disable

BUCKET=
REGION=
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=

REACT_APP_BACKEND_URL=

IS_DEV=true
```
NOTE: Only the "DB_" and "IS_DEV" properties are relevant for local testing

7. To build, install npm modules (just the first time) and run the project run:
`cd react_frontend && npm i && npm run build && cd .. && GOOS=linux GOARCH=amd64 go build -o ./bin/application`
8. To setup the eb instance (once you've setup the AWS CLI tool), run `eb init --profile [PROFILE_NAME]`
9. To deploy to app to elastic beanstalk run `eb deploy` (you may also need the `--profile` argument)
